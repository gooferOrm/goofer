# Goofer Usage Flows

Goofer is designed to be flexible, offering two primary usage flows to suit your development style: a "Library-First" approach for those who prefer to work directly with Go code, and a "CLI-First" approach that leverages code generation and automation.

## 1. Library-First (Manual) Flow

This approach is ideal for developers who want maximum control and prefer to define everything in their Go code. You'll use Goofer's packages directly to manage your database schema, models, and queries.

### Steps:

1.  **Installation:**
    Add Goofer to your project:
    ```bash
    go get github.com/gooferOrm/goofer
    ```

2.  **Define Your Structs:**
    Create your Go structs and add `orm` tags to define the schema.

    ```go
    package models

    type User struct {
        ID    uint   `orm:"primaryKey;autoIncrement"`
        Name  string `orm:"type:varchar(255);notnull"`
        Email string `orm:"unique;type:varchar(255);notnull"`
    }
    ```

3.  **Establish a Database Connection:**
    Use Go's standard `database/sql` package to connect to your database.

4.  **Perform Operations:**
    Use the `repository` package to perform CRUD operations.

    ```go
    import "github.com/gooferOrm/goofer/repository"

    userRepo := repository.NewRepository[User](db, dialect)
    user := &User{Name: "John Doe", Email: "john.doe@example.com"}
    userRepo.Save(user)
    ```

5.  **Manage Migrations (Optional):**
    Use the `migration` package to manage schema changes over time.

### Best For:

*   Developers who prefer a code-first approach.
*   Projects with complex, custom logic that doesn't fit a code-generation model.
*   Smaller projects where the overhead of a CLI is not necessary.

## 2. CLI-First (Client) Flow

This approach uses the `goofer` command-line tool to automate common tasks like project setup, code generation, and schema management. It's designed to speed up development and reduce boilerplate.

### Steps:

1.  **Install the CLI:**
    ```bash
    go install github.com/gooferOrm/goofer/cmd/goofer@latest
    ```

2.  **Initialize Your Project:**
    Run `goofer init` in your project's root directory. This will create a `goofer.toml` configuration file and a `goofer` directory for your generated code.

    ```bash
    goofer init
    ```

3.  **Define Your Schema:**
    Define your database schema in the `goofer.toml` file.

4.  **Generate Code:**
    Run `goofer generate` to create your Go models and repositories based on the schema defined in `goofer.toml`.

    ```bash
    goofer generate
    ```

5.  **Use the Generated Code:**
    Import the generated repositories into your application code to interact with the database.

    ```go
    import "your/project/path/goofer"

    userRepo := goofer.NewUserRepository(db)
    // ...
    ```

6.  **Manage Migrations:**
    Use `goofer migrate` to create and apply database migrations.

    ```bash
    goofer migrate create "add_users_table"
    goofer migrate up
    ```

### Best For:

*   Rapid application development.
*   Projects that follow a consistent structure.
*   Developers who prefer to define schemas in a configuration file.
*   Teams that want to enforce a standardized workflow.

## Choosing Your Flow

Both flows can be mixed and matched to some extent. For example, you can use the CLI to generate your initial models and then switch to the manual approach for more complex queries.

The best flow depends on your project's needs and your personal preferences. We recommend starting with the **CLI-First** approach for new projects to get up and running quickly.
