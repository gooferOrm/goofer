# Goofer ORM Examples

Goofer ORM includes comprehensive examples to help you understand its features and capabilities. These examples demonstrate various ways to utilize the ORM for different scenarios and database types.

## 🚀 Quick Navigation

- **[Complete Tutorial](../getting-started/complete-tutorial)** - Build a blog app from scratch
- **[Comprehensive Guide](../../COMPREHENSIVE_GUIDE)** - Deep dive into all features  
- **[Repository Examples](https://github.com/gooferOrm/goofer/tree/main/examples)** - All working examples on GitHub

## Core Examples

### Basic CRUD Operations

The **[Basic Example](./basic)** demonstrates fundamental database operations:

```bash
cd examples/basic
go run main.go
```

**What you'll learn:**
- Entity definition with struct tags
- Database connection and setup
- Table creation from entities
- Complete CRUD operations (Create, Read, Update, Delete)
- Basic filtering and querying
- Transaction usage

**Key concepts covered:**
- Primary keys and auto-increment
- Relationships (One-to-Many)
- Repository pattern usage
- SQL generation and execution

### Client Usage Pattern

The **[Client Example](https://github.com/gooferOrm/goofer/tree/main/examples/client)** shows simplified high-level usage:

```bash
cd examples/client
go run main.go
```

**What you'll learn:**
- High-level client interface
- Simplified setup and configuration
- Automatic table management
- Clean, minimal code patterns

**Perfect for:** Quick prototypes and simple applications

### Advanced Querying

The **[Custom Queries Example](https://github.com/gooferOrm/goofer/tree/main/examples/custom_queries)** demonstrates complex database operations:

```bash
cd examples/custom_queries
go run main.go
```

**What you'll learn:**
- Complex WHERE conditions (IN, BETWEEN, LIKE)
- Aggregation and grouping (COUNT, GROUP BY)
- Joins across multiple tables
- Custom SQL with struct mapping
- Subqueries and advanced patterns
- Transaction usage for data consistency

**Features demonstrated:**
- Query builder advanced usage
- Raw SQL integration
- Performance optimization techniques
- Data transformation patterns

## Relationship Examples

### Entity Relationships

The **[Relationships Example](./relationships)** covers all relationship types:

**Relationship types:**
- **One-to-One**: User ↔ Profile
- **One-to-Many**: User → Posts
- **Many-to-One**: Post → User
- **Many-to-Many**: User ↔ Roles (with join table)

**Key concepts:**
- Foreign key configuration
- Join table management
- Eager vs lazy loading
- Relationship querying

## Lifecycle and Validation

### Hooks and Events

The **[Hooks Example](https://github.com/gooferOrm/goofer/tree/main/examples/hooks)** demonstrates lifecycle management:

```bash
cd examples/hooks
go run main.go
```

**What you'll learn:**
- Before/After hooks for Create, Update, Delete, Save
- Data normalization and transformation
- Audit logging and tracking
- Business logic enforcement
- Error handling in hooks

**Real-world scenarios:**
- Password hashing
- Email normalization
- Timestamps management
- Activity logging

### Data Validation

The **[Validation Example](https://github.com/gooferOrm/goofer/tree/main/examples/validation)** shows data integrity:

```bash
cd examples/validation
go run main.go
```

**What you'll learn:**
- Struct tag validation
- Custom validation logic
- Error handling and messages
- Integration with go-playground/validator
- Business rule enforcement

## Database Dialects

### SQLite (Default)
Most examples use SQLite for simplicity and portability:
- No external dependencies
- Perfect for development and testing
- File-based or in-memory databases

### MySQL Integration
**[MySQL Example](https://github.com/gooferOrm/goofer/tree/main/examples/mysql)** shows MySQL-specific features:

```bash
cd examples/mysql
# Requires MySQL server running
go run main.go
```

**MySQL-specific features:**
- AUTO_INCREMENT handling
- Engine and charset configuration
- MySQL-specific data types
- Connection string formats

### PostgreSQL Integration
**[PostgreSQL Example](https://github.com/gooferOrm/goofer/tree/main/examples/postgres)** demonstrates PostgreSQL usage:

```bash
cd examples/postgres  
# Requires PostgreSQL server running
go run main.go
```

**PostgreSQL-specific features:**
- SERIAL type handling
- PostgreSQL-specific data types
- Schema and namespace support
- Advanced PostgreSQL features

## Command Line Applications

### CLI Application

The **[CLI App Example](https://github.com/gooferOrm/goofer/tree/main/examples/cli_app)** builds a complete command-line interface:

```bash
cd examples/cli_app
go run main.go --help
```

**Features:**
- Complete blog management CLI
- User and post management
- Database statistics
- Cobra CLI framework integration
- Production-ready command structure

**Commands available:**
- `init` - Initialize database
- `user create/list/delete` - User management
- `post create/list/publish` - Post management
- `stats` - Database statistics

### Simple CLI

The **[Simple CLI Example](https://github.com/gooferOrm/goofer/tree/main/examples/simple_cli)** demonstrates interactive applications:

```bash
cd examples/simple_cli
go run main.go
```

**What you'll learn:**
- Interactive command processing
- Real-time database operations
- User input handling
- Session management

## Real-World Applications

### Todo Application

The **[Todo Example](https://github.com/gooferOrm/goofer/tree/main/examples/todo)** shows a practical application:

```bash
cd examples/todo
go run main.go
```

**Features:**
- Complete task management
- Status tracking
- Data persistence
- Clean architecture patterns

### Advanced Features Examples

Additional examples in the repository demonstrate:

- **[Bulk Operations](https://github.com/gooferOrm/goofer/tree/main/examples/bulk_operations)** - Efficient batch processing
- **[Introspection](https://github.com/gooferOrm/goofer/tree/main/examples/introspection)** - Schema analysis and code generation
- **[Migrations](https://github.com/gooferOrm/goofer/tree/main/examples/migrations)** - Database schema evolution
- **[Soft Delete](https://github.com/gooferOrm/goofer/tree/main/examples/soft_delete)** - Logical deletion patterns

## Running Examples

### Prerequisites

```bash
# Install Go 1.21 or later
go version

# Clone the repository
git clone https://github.com/gooferOrm/goofer.git
cd goofer
```

### Running Individual Examples

```bash
# Navigate to any example
cd examples/[example-name]

# Install dependencies (if needed)
go mod tidy

# Run the example
go run main.go
```

### Example Status

| Example | Status | Database | Key Features |
|---------|--------|----------|--------------|
| basic | ✅ Working | SQLite | CRUD, Relationships, Transactions |
| client | ✅ Working | SQLite | High-level API, Simplified usage |
| custom_queries | ✅ Working | SQLite | Advanced queries, Aggregation |
| hooks | ✅ Working | SQLite | Lifecycle events, Data transformation |
| validation | ✅ Working | SQLite | Data validation, Error handling |
| cli_app | ✅ Working | SQLite | CLI interface, Commands |
| todo | ✅ Working | SQLite | Complete application, Task management |
| mysql | ⚠️ Requires MySQL | MySQL | MySQL-specific features |
| postgres | ⚠️ Requires PostgreSQL | PostgreSQL | PostgreSQL-specific features |

## Learning Path

### Beginner Path
1. **[Basic Example](./basic)** - Learn fundamentals
2. **[Client Example](https://github.com/gooferOrm/goofer/tree/main/examples/client)** - Simplified usage
3. **[Validation Example](https://github.com/gooferOrm/goofer/tree/main/examples/validation)** - Data integrity

### Intermediate Path  
1. **[Relationships Example](./relationships)** - Entity relationships
2. **[Hooks Example](https://github.com/gooferOrm/goofer/tree/main/examples/hooks)** - Lifecycle management
3. **[Custom Queries Example](https://github.com/gooferOrm/goofer/tree/main/examples/custom_queries)** - Advanced querying

### Advanced Path
1. **[Complete Tutorial](../getting-started/complete-tutorial)** - Real-world application
2. **[Comprehensive Guide](../../COMPREHENSIVE_GUIDE)** - Master all features
3. **[Migration Guide](../../MIGRATION_GUIDE)** - Production deployments

## Contributing Examples

Found an issue with an example? Want to contribute a new one?

1. **Report Issues**: [GitHub Issues](https://github.com/gooferOrm/goofer/issues)
2. **Submit Examples**: Create pull requests with new examples
3. **Improve Documentation**: Help us make examples clearer

## Next Steps

- Explore the **[Complete Tutorial](../getting-started/complete-tutorial)** for a comprehensive walkthrough
- Read the **[Comprehensive Guide](../../COMPREHENSIVE_GUIDE)** for deep technical details
- Check out **[Best Practices](../reference/best-practices)** for production usage
- Join our **[Community](../community)** for questions and discussions

All examples are self-contained and include the necessary code to demonstrate featured functionality. Start with the basic example and work your way up to more complex scenarios!

## Basic Example

The basic example shows simple CRUD operations with an SQLite database:

- Entity definition and registration
- Connection setup
- Table creation
- Basic CRUD operations (Create, Read, Update, Delete)
- Simple filtering and querying

See the [Basic Example](./basic) for more details.

## Database Dialects

Goofer ORM supports multiple database dialects:

- [SQLite](./basic) (covered in the basic example)
- [MySQL](./mysql) - Working with MySQL databases
- [PostgreSQL](./postgres) - PostgreSQL integration

## Advanced Features

Explore these examples to learn about more advanced features:

- [Relationships](./relationships) - One-to-one, one-to-many, and many-to-many relationships
- [Migrations](./migrations) - Evolving your database schema over time
- [Validation](./validation) - Using validation tags and custom validation
- [Hooks](./hooks) - Lifecycle hooks for automated tasks
- [Custom Queries](./custom-queries) - Advanced querying with raw SQL and aggregate functions
- [Soft Delete](./soft-delete) - Implementing soft delete functionality
- [Bulk Operations](./bulk-operations) - Efficient bulk database operations

## Complete Applications

For more complex examples showing how to put it all together:

- [RESTful API](./restful-api) - Building a RESTful API with Goofer ORM
- [Web Application](./web-application) - Integrating Goofer in a web application

## Running Examples

All examples are located in the `examples` directory of the Goofer ORM repository. To run an example:

```bash
cd examples/[example-name]
go run main.go
```

For instance, to run the basic example:

```bash
cd examples/basic
go run main.go
```

The examples are self-contained and include all the necessary code to demonstrate the featured functionality.