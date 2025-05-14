package main

import "github.com/gooferOrm/goofer/engine"

// Connect initializes a new database connection with the specified driver and DSN
// It's the main entry point for the Goofer ORM
//
// Example:
//   db, err := goofer.Connect("sqlite3", "test.db")
//   if err != nil {
//       log.Fatal(err)
//   }
//   defer db.Close()
func Connect(driver, dsn string) (*engine.Client, error) {
	return engine.Connect(driver, dsn)
}

// Config creates a new database configuration with the specified driver and DSN
// This allows for more advanced configuration before connecting
//
// Example:
//   db, err := goofer.Config("postgres", "user=postgres dbname=mydb").
//       WithLogLevel("debug").
//       Connect()
func Config(driver, dsn string) *engine.Config {
	return engine.NewConfig(driver, dsn)
}
