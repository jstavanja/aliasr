package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/charm/kv"
)

type model struct {
  screen string
  currentWorkspace string
  workspaces []string
  cursor int
  commandsInCurrentWorkspace []string
  db *kv.KV
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
    case tea.KeyMsg:
      switch msg.String() {
        case "ctrl+c", "q":
          if m.screen == "workspace-list" {
            return m, tea.Quit
          }

          if m.screen == "workspace-detail" {
            m.cursor = 0
            m.screen = "workspace-list"
          }

          if m.screen == "add-workspace" {
            m.screen = "workspace-list"
          }

          if m.screen == "add-command" {
            m.screen = "workspace-detail"
          }

        case "up", "k":
          if m.cursor > 0 {
            m.cursor--
          }

        case "down", "j":
          if m.cursor < len(m.workspaces)-1 {
            m.cursor++
          }

        case "enter", " ":
          if m.screen == "workspace-list" {
            m.currentWorkspace = m.workspaces[m.cursor]
            m.commandsInCurrentWorkspace = getCommmandsForCurrentWorkspace(m)
            m.cursor = 0
            m.screen = "workspace-detail"
          }
          
          if m.screen == "workspace-detail" {
            // TODO: get the command under the cursor and execute it, like:
            // exec.Command(m.commandsInCurrentWorkspace[m.cursor])
          }

          if m.screen == "add-workspace" {
            // TODO: add an input, read it and save the new workspace name
            m.screen = "workspace-list"
          }

          if m.screen == "add-command" {
            // TODO: add an input, read it and save the new command under the current workspace
            m.screen = "workspace-detail"
          }

        case "n":
          if m.screen == "workspace-list" {
            m.screen = "add-workspace"
          }
          if m.screen == "workspace-detail" {
            m.screen = "add-command"
          }

        case "x":
          if m.screen == "workspace-list" {
            // TODO: delete the workspace m.currentWorkspace
          }
          if m.screen == "workspace-detail" {
            // TODO: delete the command under m.commandsInCurrentWorkspace[m.cursor]
          }
          // then, decrement the m.cursor if possible
      }
  }

  return m, nil
}

func getCommmandsForCurrentWorkspace(m model) []string {
  // Get the persisted workspace commands
  v, err := m.db.Get([]byte("workspace-commands")) // format: workspace_name:command,command|workspace_name_2:command,command

  if err != nil || len(v) == 0 {
    return []string{}
  }

  workspaces_with_commands := strings.Split(string(v), "|")

  for _, workspace_with_commands := range workspaces_with_commands {
    workspace_name_and_commands_list := strings.Split(workspace_with_commands, ":")

    workspace_name := workspace_name_and_commands_list[0]
    commands := workspace_name_and_commands_list[1]

    if workspace_name == m.currentWorkspace {
      return strings.Split(commands, ",")
    }
  }

  return []string{}
}

func renderWorkspaceListView(m model) string {
  s := "Choose your workspace:\n\n"

  for i, choice := range m.workspaces {
    cursor := " "
    if m.cursor == i {
      cursor = ">"
    }
    s += fmt.Sprintf("%s %s\n", cursor, choice)
  }

  s += "\nPress n to add a new workspace."
  // s += "\nPress x to delete the selected workspace."
  s += "\nPress q to quit."

  return s
}

func renderWorkspaceDetailView(m model) string {
  s := "[Workspace]: " + m.currentWorkspace + "\n"
  s += "\nChoose your command:\n"

  commands := m.commandsInCurrentWorkspace

  for i, choice := range commands {
    cursor := " "
    if m.cursor == i {
      cursor = ">"
    }
    s += fmt.Sprintf("%s %s\n", cursor, choice)
  }

  s += "\nPress n to add a new command."
  // s += "\nPress x to delete the selected command."
  s += "\nPress q to see all workspaces."

  return s
}

func renderAddWorkspaceView(m model) string {
  s := "Name your new workspace:\n\n"

  s += "Some workspace"

  s += "\n\nPress enter to add the workspace."
  s += "\nPress q to abort adding a workspace."

  return s
}

func renderAddCommandView(m model) string {
  s := "Input your new command:\n\n"

  s += "npm run todo"

  s += "\n\nPress enter to add the command."
  s += "\nPress q to abort adding a command."

  return s
}

func (m model) View() string {
  if m.screen == "workspace-list" {
    return renderWorkspaceListView(m)
  }

  if m.screen == "workspace-detail" {
    return renderWorkspaceDetailView(m)
  }

  if m.screen == "add-workspace" {
    return renderAddWorkspaceView(m)
  }

  if m.screen == "add-command" {
    return renderAddCommandView(m)
  }

  return "ERROR: Application in unknown state."
}

func main() {
  // Open or create the aliasr database that will persist the user's data locally 
  db, err := kv.OpenWithDefaults("aliasr-db")

  if err != nil {
    log.Fatal(err)
  }

  defer db.Close()

  // Sync the data from the database
  if err := db.Sync(); err != nil {
    log.Fatal(err)
  }

  // Set some test workspaces as an example
  if err := db.Set([]byte("workspaces"), []byte("test,bro,ski")); err != nil {
    log.Fatal(err)
  }

  // Set some test commands as an example
  if err := db.Set([]byte("workspace-commands"), []byte("test:npm run dev,npm run test|bro:nx serve app|ski:yarn dev,npm run e2e")); err != nil {
    log.Fatal(err)
  }

  // Get the persisted workspaces
  v, err := db.Get([]byte("workspaces"))

  // If the workspaces weren't ever created, initialize them
  if err != nil {
    db.Set([]byte("workspaces"), []byte(""))
    v = []byte("")
  } 
  
  // Create a slice with workspaces from the db string value
  workspaces := strings.Split(string(v), ",")

  initialModel := model{
    screen: "workspace-list",
    workspaces: workspaces,
    commandsInCurrentWorkspace: []string{},
    db: db,
  }

  p := tea.NewProgram(initialModel)

  if err := p.Start(); err != nil {
    fmt.Printf("Alas, there's been an error: %v", err)
    os.Exit(1)
  }
}
