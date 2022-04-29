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
  currentWorkspace int
  workspaces []string
  cursor int
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
    case tea.KeyMsg:
      switch msg.String() {
        case "ctrl+c", "q":
          if m.screen == "workspace-detail" {
            m.cursor = 0
            m.screen = "workspace-list"
          } else {
            return m, tea.Quit
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
            m.currentWorkspace = m.cursor
            m.cursor = 0
            m.screen = "workspace-detail"
          } else {
            // get the command under the cursor and execute it
          }
      }
  }

  return m, nil
}

func RenderWorkspaceListView(m model) string {
  s := "Choose your workspace:\n"

  for i, choice := range m.workspaces {
    cursor := " "
    if m.cursor == i {
      cursor = ">"
    }
    s += fmt.Sprintf("%s %s\n", cursor, choice)
  }

  s += "\nPress q to quit.\n"

  return s
}

func RenderWorkspaceDetailView(m model) string {
  // TODO: add the name of the workspace in the title of the screen
  s := "Choose your command:\n"

  // TODO: fetch commands for current workspace from the kv db
  commands := []string{"npm run dev", "echo \"hello world\""}

  for i, choice := range commands {
    cursor := " "
    if m.cursor == i {
      cursor = ">"
    }
    s += fmt.Sprintf("%s %s\n", cursor, choice)
  }

  s += "\nPress q to see all workspaces.\n"

  return s
}

func (m model) View() string {
  if m.screen == "workspace-list" {
    return RenderWorkspaceListView(m)
  }

  if m.screen == "workspace-detail" {
    return RenderWorkspaceDetailView(m)
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
  // if err := db.Set([]byte("workspaces"), []byte("test,bro,ski")); err != nil {
  //   log.Fatal(err)
  // }

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
  }

  p := tea.NewProgram(initialModel)

  if err := p.Start(); err != nil {
    fmt.Printf("Alas, there's been an error: %v", err)
    os.Exit(1)
  }
}
