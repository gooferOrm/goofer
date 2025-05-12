# Goofer ORM Examples

Goofer ORM includes several examples to help you understand its features and capabilities. These examples demonstrate various ways to utilize the ORM for different scenarios and database types.

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