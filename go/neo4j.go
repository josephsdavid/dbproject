package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "strings"

  "github.com/neo4j/neo4j-go-driver/neo4j"
)

var (
  uri      string
  username string
  password string
  query    string
)

// Simple header printing logic, open to improvements
func processHeaders(result neo4j.Result) {
  if keys, err := result.Keys(); err == nil {
    for index, key := range keys {
      if index > 0 {
        fmt.Print("\t")
      }
      fmt.Printf("%-10s", key)
    }
    fmt.Print("\n")

    for index := range keys {
      if index > 0 {
        fmt.Print("\t")
      }
      fmt.Print(strings.Repeat("=", 10))
    }
    fmt.Print("\n")
  }
}

// Simple record values printing logic, open to improvements
func processRecord(record neo4j.Record) {
  for index, value := range record.Values() {
    if index > 0 {
      fmt.Print("\t")
    }
    fmt.Printf("%-10v", value)
  }
  fmt.Print("\n")
}

// Transaction function
func executeQuery(tx neo4j.Transaction) (interface{}, error) {
  var (
    counter int
    result  neo4j.Result
    err     error
  )

  // Execute the query on the provided transaction
  if result, err = tx.Run(query, nil); err != nil {
    return nil, err
  }

  // Print headers
  processHeaders(result)

  // Loop through record stream until EOF or error
  for result.Next() {
    processRecord(result.Record())
    counter++
  }
  // Check if we encountered any error during record streaming
  if err = result.Err(); err != nil {
    return nil, err
  }

  // Return counter
  return counter, nil
}

func run() error {
  var (
    driver           neo4j.Driver
    session          neo4j.Session
    recordsProcessed interface{}
    err              error
  )

  // Construct a new driver
  if driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(config *neo4j.Config) {
    config.Log = neo4j.ConsoleLogger(neo4j.ERROR)
  }); err != nil {
    return err
  }
  defer driver.Close()

  // Acquire a session
  if session, err = driver.Session(neo4j.AccessModeRead); err != nil {
    return err
  }
  defer session.Close()

  // Execute the transaction function
  if recordsProcessed, err = session.ReadTransaction(executeQuery); err != nil {
    return err
  }

  fmt.Printf("\n%d records processed\n", recordsProcessed)

  return nil
}

func parseAndVerifyFlags() bool {
  flag.Parse()

  if uri == "" || username == "" || password == "" || query == "" {
    flag.Usage()
    return false
  }

  return true
}

func main() {
  if !parseAndVerifyFlags() {
    os.Exit(-1)
  }

  if err := run(); err != nil {
    log.Fatal(err)
    os.Exit(1)
  }

  os.Exit(0)
}

func init() {
  flag.StringVar(&uri, "uri", "bolt://localhost", "the bolt uri to connect to")
  flag.StringVar(&username, "username", "neo4j", "the database username")
  flag.StringVar(&password, "password", "", "the database password")
  flag.StringVar(&query, "query", "", "the query to execute")
}
