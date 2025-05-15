# Repository Pattern Examples

This directory contains examples demonstrating different ways to use the Goofer ORM with the repository pattern.

## Examples

### 1. With Engine (`with_engine/`)

This example shows how to use the repository pattern with the Goofer ORM engine. The engine provides a higher-level API and handles connection management.

**Key Features:**
- Uses `engine.Connect` to create a client
- Registers entities with the client
- Uses type-safe repositories with `engine.Repo[T]`
- Automatic connection management

### 2. Without Engine (`without_engine/`)

This example demonstrates using the repository pattern directly without the engine. This approach gives you more control but requires manual management of the database connection.

**Key Features:**
- Direct use of `repository.NewRepository`
- Manual database connection management
- More control over the repository configuration
- No dependency on the engine package

### 3. With Config (`with_config/`)

This example shows a structured approach using a configuration package to manage the database connection and repositories.

**Key Features:**
- Centralized configuration
- Clean separation of concerns
- Easy to test and maintain
- Reusable across the application

## Running the Examples

Each example is self-contained. To run an example:

```bash
cd with_engine  # or without_engine or with_config
go run main.go
```

## Choosing an Approach

- **Use With Engine** if you want a simple, high-level API with built-in connection management.
- **Use Without Engine** if you need more control over the database connection or want to avoid the engine's overhead.
- **Use With Config** for larger applications where you want to maintain a clean architecture and separation of concerns.

## Dependencies

All examples require:
- Go 1.16 or higher
- SQLite3 driver (`github.com/mattn/go-sqlite3`)
- Goofer ORM

## License

MIT
