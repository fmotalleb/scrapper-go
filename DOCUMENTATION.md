# Scrapper-Go Engine Documentation

This document provides detailed information about the core components of the Scrapper-Go engine: **Middlewares** and **Steps**. These components work together to execute web scraping pipelines defined in YAML configuration files.

## Middlewares

Middlewares are the backbone of the scraping process. They form a chain of functions that wrap around each `Step`'s execution, allowing for cross-cutting concerns like error handling, conditional execution, and looping. The middlewares are executed in a specific order, determined by their filenames.

---

### `middleware.go`

This file contains the core logic for managing and executing the middleware chain.

- **`registerMiddleware(m middleware)`**: Adds a new middleware to the global list of middlewares.
- **`HandleStep(...)`**: Initiates the execution of the middleware chain for a given `Step`. It calls the first middleware in the chain.
- **`middlewareExec(index, ...)`**: A recursive function that executes the middleware at the current `index` and passes a `next` function that, when called, will execute the next middleware in the chain. The last middleware receives `nil` as the `next` function.

---

### `mid_00_err_handler.go`

This is the first middleware in the chain and provides robust error handling for each step.

**YAML Configuration:**

The behavior of this middleware is controlled by the `on-error` key within a step's configuration.

```yaml
- click: ".some-button"
  on-error: "ignore" # Can be "ignore", "print", or "panic"
```

**Modes:**

- **`ignore`** (Default if `on-error` is not specified and an error occurs): The error is suppressed, and the pipeline execution continues to the next step.
- **`print`**: The error is logged to the console, but the execution continues.
- **`panic`**: The application panics and stops execution immediately.

---

### `mid_10_if.go`

This middleware enables conditional execution of a step based on a specified condition.

**YAML Configuration:**

```yaml
- click: ".a"
  if: "'.some-selector' | query.exists"
```

**Logic:**

1.  It looks for an `if` key in the step's configuration.
2.  If `if` is not present, it proceeds to the next middleware in the chain.
3.  If `if` is present, it evaluates the condition string. The string can contain templates and uses the `query` language to evaluate conditions on the page.
4.  If the condition evaluates to `true`, it proceeds with the execution of the step.
5.  If the condition evaluates to `false`, it returns an `errTestFailed`, and the step (and its subsequent middlewares) are not executed.

---

### `mid_11_loop.go`

This middleware allows a set of nested steps to be executed multiple times, either for a fixed number of iterations or over a list of items.

**YAML Configuration:**

**1. Fixed Number of Iterations:**

```yaml
- loop: "3"
  loop-key: "index" # Optional, defaults to "item"
  steps:
    - debug: "This is iteration number {{ .index }}"
```

**2. Iterating Over a List (from a variable):**

```yaml
- set-var: "my_links"
  element: "a"
  mode: "html"
- loop: "{{ .my_links }}"
  loop-key: "link"
  steps:
    - debug: "Found link: {{ .link }}"
```

**Logic:**

1.  It checks for a `loop` key in the step's configuration.
2.  The value of `loop` is evaluated. It can be a number (for a fixed loop) or a template that resolves to a JSON array.
3.  It iterates based on the evaluated result.
4.  In each iteration, it sets a variable in the context (default key is `item`, configurable with `loop-key`).
5.  It then executes the `steps` defined within the loop configuration for each iteration.

---

### `mid_zz_execute.go`

This is the final middleware in the chain. Its responsibility is to execute the actual `Step` and handle the result.

**YAML Configuration:**

```yaml
- element: "h1"
  mode: "text"
  set-var: "page_title"
```

**Logic:**

1.  It calls the `Execute` method on the `Step` itself.
2.  If the execution is successful and a `set-var` key is present in the step's configuration, it stores the result of the `Execute` method in the results map (`r`).
3.  The `set-var` logic is smart: if you set the same variable multiple times, it will automatically convert the variable into a list and append the new values.

---

## Steps

Steps are the individual actions that make up a scraping pipeline. Each step is a struct that implements the `Step` interface, which has two methods: `Execute` and `GetConfig`.

---

### `step.go`

This file defines the `Step` interface and contains the logic for building a list of `Step` objects from a YAML configuration.

- **`Step` interface**:
    - `Execute(...)`: Performs the action of the step.
    - `GetConfig()`: Returns the original configuration of the step.
- **`BuildSteps(...)`**: Takes a list of step configurations and uses a series of `stepSelector` functions to determine which step type to instantiate for each configuration.

---

### `click.go`

Performs a click action on a specified element.

**YAML Key:** `click`

**Configuration:**

```yaml
- click: "#submit-button"
  # Optional Playwright LocatorClickOptions can be added here
  force: true
```

---

### `config.go`

Configures engine-level settings.

**YAML Key:** `config`

**Configuration:**

```yaml
- config:
    timeout: 30000 # in milliseconds
    nav_timeout: 60000 # in milliseconds
```

---

### `debug.go`

Logs a message to the console. Useful for debugging pipelines.

**YAML Key:** `debug`

**Configuration:**

```yaml
- debug: "Current URL is {{ .page.url }}"
```

---

### `eval.go`

Evaluates a JavaScript expression on the page or within the context of a specific element.

**YAML Key:** `eval`

**Configuration:**

**On the page:**
```yaml
- eval: "() => document.title"
  set-var: "pageTitle"
```

**On an element:**
```yaml
- locator: "#my-element"
  eval: "el => el.textContent"
  set-var: "elementText"
```

---

### `fill.go`

Fills an input field with a specified value.

**YAML Key:** `fill`

**Configuration:**

```yaml
- fill: "#username"
  value: "my_user"
- fill: "#password"
  value: "{{ .env.PASSWORD }}" # Using a template
```

---

### `get_element.go`

Retrieves content from an element.

**YAML Key:** `element`

**Configuration:**

```yaml
- element: "h1"
  mode: "text" # "text", "html", "value", "table", "table-flat"
  set-var: "page_title"
```

**Modes:**

- `text`: Gets the `textContent`.
- `html`: Gets the `innerHTML`.
- `value`: Gets the `inputValue` (for form elements).
- `table`: Parses a `<table>` into a nested list of lists.
- `table-flat`: Parses a `<table>` into a flat list of strings.

---

### `goto.go`

Navigates the browser to a specified URL.

**YAML Key:** `goto`

**Configuration:**

```yaml
- goto: "https://example.com"
  # Optional Playwright PageGotoOptions can be added here
  waitUntil: "networkidle"
```

---

### `mouse.go`

Performs mouse actions like clicking, double-clicking, moving, and scrolling at specific coordinates.

**YAML Key:** `mouse`

**Configuration:**

```yaml
- mouse: "100,200" # X,Y coordinates
  action: "click" # "click", "double-click", "scroll", "move", "up", "down"
```

---

### `nop.go`

The "no operation" step. It does nothing but can be useful as a container for other logic, especially loops, when no other action is needed.

**YAML Key:** `nop`

**Configuration:**

```yaml
- nop: "Just a placeholder"
- loop: "5" # A loop without any other primary action
  steps:
    - debug: "Looping..."
```

---

### `omit.go`

Deletes a variable from the context.

**YAML Key:** `omit`

**Configuration:**

```yaml
- omit: "variable_to_delete"
```

---

### `screenshot.go`

Takes a screenshot of a specific element and returns it as a base64 encoded string.

**YAML Key:** `screenshot`

**Configuration:**

```yaml
- screenshot: "#my-chart"
  set-var: "chart_image_b64"
  # Optional Playwright LocatorScreenshotOptions can be added here
```

---

### `select.go`

Selects one or more options in a `<select>` dropdown element.

**YAML Key:** `select`

**Configuration:**

You can select by value, label, or index.

```yaml
- select: "#my-select"
  # Select by value
  value: "option-1"
  # Or by label
  label: "Option 2"
  # Or by index
  index: "2"
  # Multiple selections are also possible
  values:
    - "option-3"
    - "option-4"
```

---

### `sleep.go`

Pauses the execution for a specified duration.

**YAML Key:** `sleep`

**Configuration:**

```yaml
- sleep: "5s" # 5 seconds
- sleep: "100ms" # 100 milliseconds
```
