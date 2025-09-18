## 12. Coding Standards

### 12.1. Core Standards

-   **Languages & Runtimes:** Go 1.24.6
-   **Style & Linting:**
    -   `gofmt` will be used to format all Go code.
    -   `golangci-lint` will be used for static analysis with the default set of linters.
    -   Lines should be no longer than 120 characters
-   **Test Organization:** Test files will be named `_test.go` and will be located in the same package as the code they are testing.

### 12.2. Naming Conventions

-   Standard Go naming conventions will be followed (e.g., `PascalCase` for exported identifiers, `camelCase` for unexported identifiers).

### 12.3. Critical Rules

-   **No `init()` functions:** The `init()` function is forbidden. All initialization should be done explicitly in the `main()` function or in component constructors.
-   **No global variables:** Global variables are forbidden. All state should be managed within structs and passed as dependencies.
-   **Use `slog` for all logging:** The `log` package is forbidden. All logging must be done through the configured `slog` logger.
-   **Errors must be wrapped:** All errors returned from external libraries or other components must be wrapped with context.
-   **Use `context.Context`:** The `context.Context` must be passed as the first argument to all functions that perform I/O or interact with external systems.
