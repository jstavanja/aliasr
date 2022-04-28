package main

import (
  "log"
  "strings"

  "github.com/charmbracelet/charm/kv"
)

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

  for _, workspace := range workspaces {
   log.Println(workspace)
  }
}
