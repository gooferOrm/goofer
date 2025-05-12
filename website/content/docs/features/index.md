# Goofer ORM Features

Goofer ORM is a powerful, type-safe ORM for Go that provides an amazing developer experience. It's designed to make working with databases in Go a pleasant experience with zero drama.

## Core Features

### Entity System

The [Entity System](./entity-system) is the foundation of Goofer ORM. It allows you to define your database schema using Go structs with tags for metadata. This approach provides:

- Type safety through Go's type system
- Compile-time checks for your database models
- Clear and concise schema definition

### Schema Parser

The [Schema Parser](./schema-parser) uses Go's reflection capabilities to analyze your entity structs at runtime. It:

- Extracts metadata from struct tags
- Maps Go types to database types
- Builds a complete schema registry for your application

### Relation Mapping

[Relation Mapping](./relation-mapping) in Goofer ORM makes it easy to work with related entities. It supports:

- One-to-One relationships
- One-to-Many relationships
- Many-to-One relationships
- Many-to-Many relationships with join tables
- Eager and lazy loading strategies

### Migration Engine

The [Migration Engine](./migration-engine) helps you evolve your database schema over time. It provides:

- Automatic SQL generation for schema changes
- Versioned migrations
- Up and down migration support
- Migration status tracking

### Repository Pattern

The [Repository Pattern](./repository-pattern) implementation provides a type-safe API for database operations:

- Generic Repository[T] for each entity type
- CRUD operations (Create, Read, Update, Delete)
- Fluent query building
- Filtering, sorting, and pagination

### Validation

[Validation](./validation) ensures your data meets your requirements before it hits the database:

- Integration with go-playground/validator
- Struct tag support for validation rules
- Custom validation hooks

### Hooks

[Hooks](./hooks) allow you to execute code at specific points in an entity's lifecycle:

- BeforeCreate, AfterCreate
- BeforeUpdate, AfterUpdate
- BeforeDelete, AfterDelete
- BeforeSave, AfterSave

### Dialects Support

[Dialects Support](./dialects) allows Goofer ORM to work with multiple database systems:

- SQLite
- MySQL
- PostgreSQL
- Custom dialect support

### Transactions

[Transactions](./transactions) ensure data integrity for complex operations:

- First-class transaction support
- Automatic rollback on error
- Nested transaction support

## Additional Features

- **Type Safety**: Fully leverages Go's type system with generics
- **Zero Drama**: Simple, intuitive API with minimal boilerplate
- **Query Builder**: Fluent API for building complex queries
- **Custom Queries**: Support for raw SQL when needed
- **No Code Generation**: Uses reflection and generics instead of code generation

Explore each feature in depth by clicking on the links above.