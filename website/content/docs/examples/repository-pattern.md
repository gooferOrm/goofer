---
title: Repository Pattern
description: Different ways to use the repository pattern with Goofer ORM
---

# Repository Pattern Examples

This guide demonstrates different approaches to using the repository pattern with Goofer ORM. The repository pattern provides a clean way to access data from your database while abstracting the underlying data access implementation.

## Table of Contents

- [With Engine](#with-engine)
- [Without Engine](#without-engine)
- [With Config](#with-config)
- [Choosing the Right Approach](#choosing-the-right-approach)

## With Engine

The engine provides a high-level API for working with repositories. This is the recommended approach for most applications.

```go
package main

import (
	"log"
	"github.com/goferOrm/goofer/engine"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize the client
	client, err := engine.Connect("sqlite3", "./with_engine.db")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Register entities
	err = client.RegisterEntities(&User{}, &Post{})
	if err != nil {
		log.Fatal(err)
	}

	// Get a repository
	repo := engine.Repo[User](client)
	
	// Use the repository
	user := &User{Name: "John Doe", Email: "john@example.com"}
	err = repo.Save(user)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Without Engine

For more control, you can use the repository directly without the engine:

```go
package main

import (
	"database/sql"
	"log"
	"github.com/goferOrm/goofer/repository"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open a database connection
	db, err := sql.Open("sqlite3", "./without_engine.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a repository
	repo := repository.NewRepository[User](db, repository.SQLite)
	
	// Use the repository
	user := &User{Name: "Jane Smith", Email: "jane@example.com"}
	err = repo.Save(user)
	if err != nil {
		log.Fatal(err)
	}
}
```

## With Config

For larger applications, you can use a configuration package to manage your database connection and repositories:

```go
// goofer/config.go
package goofer

import (
	"database/sql"
	"github.com/goferOrm/goofer/engine"
)

var DB *sql.DB

func Init() (*engine.Client, error) {
	db, err := sql.Open("sqlite3", "./with_config.db")
	if err != nil {
		return nil, err
	}

	client, err := engine.Connect("sqlite3", "./with_config.db")
	if err != nil {
		db.Close()
		return nil, err
	}

	err = client.RegisterEntities(&User{}, &Post{})
	if err != nil {
		client.Close()
		db.Close()
		return nil, err
	}

	DB = db
	return client, nil
}

func Close(client *engine.Client) {
	if client != nil {
		client.Close()
	}
	if DB != nil {
		DB.Close()
	}
}
```

```go
// main.go
package main

import (
	"log"
	"./goofer"
)

func main() {
	client, err := goofer.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer goofer.Close(client)

	repo := client.GetRepository(&goofer.User{})
	user := &goofer.User{Name: "Alice", Email: "alice@example.com"}
	err = repo.Create(user)
	if err != nil {
		log.Fatal(err)
	}
}
```

## Choosing the Right Approach

- **Use With Engine** for most applications where you want a balance of convenience and features.
- **Use Without Engine** when you need more control over the database connection or want to understand the internals.
- **Use With Config** for larger applications where you want better organization and separation of concerns.

## Next Steps

- Learn more about [query building](/docs/reference/query-builder)
- Explore [transactions](/docs/reference/transactions)
- Read about [migrations](/docs/reference/migrations)
