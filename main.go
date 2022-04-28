package main

import "log"
import "github.com/charmbracelet/charm/kv"

func main() {
  db, err := kv.OpenWithDefaults("my-db")

  log.Println("hello")

  if err != nil {
    log.Fatal(err)
  }

  log.Println("hello 2")

  defer db.Close()

  if err := db.Sync(); err != nil {
    log.Fatal(err)
  }

  log.Println("hello 3")

  if err := db.Set([]byte("fave-food"), []byte("gherkin")); err != nil {
    log.Fatal(err)
  }

  log.Println("hello 4")

  v, err := db.Get([]byte("fave-food"))

  if err != nil {
    log.Fatal(err)
  }

  log.Println(string(v))

}
