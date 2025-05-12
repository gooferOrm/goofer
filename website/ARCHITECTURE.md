# Goofer ORM Architecture Documentation

## Overview

Goofer is a modern ORM for Go that provides type-safe database access through code generation. This document explains the architecture and key components of the Goofer codebase.

## Core Components

### 1. Generator System

The generator system is the heart of Goofer, responsible for generating Go code from your database schema.

#### Key Files:
- `generator/run.go`: Main generator logic
- `generator/templates/*.gotpl`: Template files for code generation
- `generator/types/types.go`: Type definitions for the generator

#### Workflow:
1. Reads Prisma schema
2. Generates Go code using templates
3. Creates database client and models
4. Generates query builders and actions

### 2. Runtime System

The runtime system handles database operations at runtime.

#### Key Components:
- `runtime/builder/builder.go`: Query builder system
- `runtime/engine/engine.go`: Database engine interface
- `runtime/types/types.go`: Core type definitions
- `runtime/transaction/transaction.go`: Transaction management

#### Query Building Process:
1. Query builder creates GraphQL queries
2. Engine executes queries against the database
3. Results are transformed and returned

### 3. CLI System

The CLI system handles command-line operations.

#### Key Files:
- `cli/cli.go`: Main CLI implementation
- `main.go`: Entry point for CLI commands

#### Available Commands:
- `generate`: Generate code from schema
- `init`: Initialize project
- `prefetch`: Pre-fetch binaries

### 4. Binary Management

Handles database engine and CLI binaries.

#### Key Files:
- `binaries/version.go`: Version management
- `binaries/bindata/bindata.go`: Binary data management
- `binaries/platform/platform.go`: Platform-specific logic

## Code Generation Process

1. **Schema Parsing**
   - Reads Prisma schema
   - Parses models, fields, and relations

2. **Template Processing**
   - Uses Go template system
   - Generates code for:
     - Models
     - Queries
     - Actions
     - Enums
     - Errors

3. **Code Generation**
   - Creates database client
   - Generates type-safe operations
   - Implements query builders

## Advanced Features

### Query Caching

```go
// QueryCache provides caching for database queries
func NewQueryCache(defaultTTL time.Duration) *QueryCache {
    return &QueryCache{
        cache: make(map[string]cachedResult),
        ttl:   defaultTTL,
    }
}
```

### Database Sharding

```go
// DatabaseShardManager handles database sharding
func NewDatabaseShardManager(shardURLs []string) *DatabaseShardManager {
    return &DatabaseShardManager{
        shards: shardURLs,
    }
}
```

### Query Profiling

```go
// QueryProfiler tracks slow queries
func NewQueryProfiler(slowThreshold time.Duration) *QueryProfiler {
    return &QueryProfiler{
        slowQueryThreshold: slowThreshold,
        queryLogs:          []QueryLog{},
    }
}
```

## Development Workflow

1. **Code Generation**
   ```bash
   go run github.com/gooferOrm/goofer generate
   ```

2. **Testing**
   - Unit tests in `test/` directory
   - Integration tests with Docker
   - E2E tests for compatibility

3. **Building**
   ```bash
   go build ./...
   ```

## Error Handling

Goofer uses a comprehensive error handling system:
- Type-safe errors
- Detailed error messages
- Context-aware error handling

## Type System

Goofer supports various database types:
- `DateTime`: Time handling
- `Decimal`: Decimal numbers
- `Bytes`: Binary data
- `BigInt`: Large integers
- `JSON`: JSON data

## Best Practices

1. **Code Generation**
   - Always use code generation for database operations
   - Never modify generated code directly

2. **Error Handling**
   - Always check for errors
   - Use proper error types
   - Provide meaningful error messages

3. **Performance**
   - Use query caching for repeated queries
   - Implement proper indexing
   - Use transactions for multiple operations

## Future Development

1. **Enhancements**
   - More database types
   - Advanced query optimizations
   - Better error handling

2. **Features**
   - More advanced caching
   - Better query profiling
   - Improved sharding

3. **Performance**
   - Query optimization
   - Connection pooling
   - Better caching strategies

## Contributing

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Run tests
5. Submit a PR

## Security Considerations

1. **Input Validation**
   - Always validate user input
   - Use type-safe operations
   - Prevent SQL injection

2. **Authentication**
   - Proper database credentials
   - Secure connection handling
   - Environment-based configuration

3. **Data Protection**
   - Proper encryption
   - Secure data handling
   - Audit logging
