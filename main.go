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
          return m, tea.Quit

        case "up", "k":
          if m.cursor > 0 {
            m.cursor--
          }

        case "down", "j":
          if m.cursor < len(m.workspaces)-1 {
            m.cursor++
          }

          case "enter", " ":
            // TODO: make this enter a workspace
      }
  }

  return m, nil
}

func (m model) View() string {
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
    workspaces: workspaces,
  }

  p := tea.NewProgram(initialModel)

  if err := p.Start(); err != nil {
    fmt.Printf("Alas, there's been an error: %v", err)
    os.Exit(1)
  }
}
