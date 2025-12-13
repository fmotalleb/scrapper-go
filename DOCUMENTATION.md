# Scrapper-Go Engine Documentation

This document provides detailed information about the core components of the Scrapper-Go engine: **Pipeline Configuration**, **Middlewares**, and **Steps**. These components work together to execute web scraping pipelines defined in YAML configuration files.

## Pipeline Configuration

At the root of your YAML file, you can specify several high-level configurations that control the browser environment and define variables.

```yaml
pipeline:
  browser: chromium
  browser_params:
    headless: false
  browser_page_options:
    screen:
      width: 1280
      height: 763
  vars:
    - name: my_variable
      value: "static value"
    - name: user_agent
      value: "MyScraper/1.0"
  steps:
    # ... your steps here
```

- **`browser`**: Specifies the browser to use. Can be `chromium`, `firefox`, or `webkit`.
- **`browser_params`**: Parameters passed to the browser instance. Any valid Playwright `BrowserTypeLaunchOptions` can be used here (e.g., `headless`, `slow_mo`).
- **`browser_page_options`**: Parameters passed when a new page is created. Any valid Playwright `BrowserNewPageOptions` can be used (e.g., `screen`, `user_agent`).
- **`vars`**: A list of variables to be made available to the steps via templating.
- **`steps`**: The list of actions to be performed in the pipeline.

### The `vars` Block

The `vars` block allows you to pre-define variables. These can be static values or dynamically generated.

**Static Variables:**
```yaml
vars:
  - name: loginBtn
    value: "#login-button"
```

**Dynamic (Random) Variables:**
You can generate random strings, which is useful for creating unique test data.

```yaml
vars:
  - name: random_user_id
    random: always # 'always' generates a new value each time the var is used
    random_chars: "abcdefghijklmnopqrstuvwxyz"
    random_length: 10
  - name: session_email
    random: once # 'once' generates the value once and reuses it
    random_chars: "abcdef123456789"
    random_length: 8
    postfix: "@example.com"
```

These variables can be accessed in your steps using `{{ .variable_name }}`.

## Middlewares

Middlewares are the backbone of the scraping process. They form a chain of functions that wrap around each `Step`'s execution, allowing for cross-cutting concerns like error handling, conditional execution, and looping.

---

### `mid_00_err_handler.go`

Provides robust error handling. The behavior is controlled by the `on-error` key.

**YAML Configuration & Modes:**

- **`ignore`**: Suppresses the error and continues execution. Useful for optional elements.
  ```yaml
  - click: ".optional-popup-close"
    on-error: "ignore"
  ```
- **`print`**: Logs the error but continues execution.
- **`panic`**: Stops execution immediately.

---

### `mid_10_if.go`

Enables conditional execution of a step.

**YAML Configuration:**

```yaml
- click: ".a"
  if: "'.some-selector' | query.exists"
```
The `if` condition is a powerful expression that can use templates and the custom `query` language to check for the existence of elements, their attributes, and more. Execution proceeds only if the condition evaluates to `true`.

---

### `mid_11_loop.go`

Executes a set of nested steps multiple times.

**YAML Configuration:**

**1. Fixed Number of Iterations:**
The `loop-key` (`index` here) holds the current iteration number.
```yaml
- loop: "3"
  loop-key: "index" # Optional, defaults to "item"
  steps:
    - debug: "This is iteration number {{ .index }}"
```

**2. Iterating Over a List (from a variable):**
```yaml
- loop: "{{ .my_links_variable }}"
  loop-key: "link" # Optional, defaults to "item"
  steps:
    - goto: "{{ .link }}"
```

**3. Dynamic Loop from JavaScript Evaluation:**
This advanced example gets a list of option values from a dropdown and then loops over them.
```yaml
- loop: '{{ eval "JSON.stringify([...document.querySelectorAll(''#my-select > option'')].map(o => o.value))" }}'
  on-error: ignore
  steps:
    - select: "#my-select"
      value: "{{ item }}"
    - debug: "Selected option {{ item }}"
```

---

### `mid_zz_execute.go`

The final middleware that executes the `Step` and optionally stores its result in a variable using `set-var`.

**YAML Configuration:**

```yaml
- element: "h1"
  mode: "text"
  set-var: "page_title"
```
If you use `set-var` with the same variable name multiple times, it will create a list and append the new values. This is useful for collecting data in a loop.

---

## Steps

Steps are the individual actions in a pipeline.

---
### `click.go`

Performs a click action on an element.
**YAML Key:** `click`
```yaml
- click: "#submit-button"
  # Optional Playwright LocatorClickOptions can be added here
  force: true
```

---
### `config.go`

Configures engine-level settings during execution.
**YAML Key:** `config`
```yaml
- config:
    timeout: 30000 # Default timeout in milliseconds
    nav_timeout: 60000 # Default navigation timeout in milliseconds
```

---
### `debug.go`

Logs a message to the console. Excellent for debugging.
**YAML Key:** `debug`
```yaml
- debug: "Current URL is {{ .page.url }} and title is {{ .page.title }}"
```

---
### `eval.go`

Evaluates JavaScript.
**YAML Key:** `eval`
```yaml
# On the page
- eval: "() => document.title"
  set-var: "pageTitle"

# On an element, to get its text content
- locator: "#my-element"
  eval: "el => el.textContent"
  set-var: "elementText"
```

---
### `fill.go`

Fills an input field.
**YAML Key:** `fill`
```yaml
- fill: "#username"
  value: "my_user"
- fill: "#password"
  value: "{{ .env.PASSWORD }}" # Using an environment variable
```

---
### `get_element.go`

Retrieves content from an element.
**YAML Key:** `element`
```yaml
- element: "h1"
  mode: "text" # "text", "html", "value", "table", "table-flat"
  set-var: "page_title"
```
- **Modes**: `text`, `html` (innerHTML), `value` (input value), `table` (parses a `<table>` into a list of lists), `table-flat` (parses a `<table>` into a flat list).

---
### `goto.go`

Navigates to a URL.
**YAML Key:** `goto`
```yaml
- goto: "https://example.com"
  # Optional Playwright PageGotoOptions can be added here
  waitUntil: "networkidle"
```

---
### `mouse.go`

Performs direct mouse actions.
**YAML Key:** `mouse`
```yaml
- mouse: "100,200" # X,Y coordinates
  action: "click" # "click", "double-click", "scroll", "move", "up", "down"
```

---
### `nop.go`

"No operation" step. Useful as a container for logic like loops when no other action is needed.
**YAML Key:** `nop`
```yaml
- loop: "5" # A loop with no primary action
  steps:
    - debug: "Looping..."
```

---
### `omit.go`

Deletes a variable from the context.
**YAML Key:** `omit`
```yaml
- omit: "variable_to_delete"
```

---
### `screenshot.go`

Takes a screenshot of an element.
**YAML Key:** `screenshot`
```yaml
- screenshot: "#my-chart"
  set-var: "chart_image_b64" # Returns as a base64 encoded string
```

---
### `select.go`

Selects options in a `<select>` dropdown.
**YAML Key:** `select`
```yaml
- select: "#my-select"
  value: "option-1" # Can also use 'label' or 'index'
```

---
### `sleep.go`

Pauses execution.
**YAML Key:** `sleep`
```yaml
- sleep: "5s" # 5 seconds
- sleep: "100ms" # 100 milliseconds
```

---

## Full Example

Here is a complete example that demonstrates many of the features described above. It signs up for a service using randomly generated credentials and extracts an API key.

```yaml
pipeline:
  browser: chromium
  browser_params:
    headless: false

  vars:
    - name: username
      random: once
      random_chars: "abcdefghijklmnopqrstuvwxyz"
      random_length: 10

    - name: email
      random: once
      random_chars: "abcdefghijklmnopqrstuvwxyz123456789"
      random_length: 8
      postfix: "@example.com"

    - name: password
      random: once
      random_chars: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
      random_length: 14

  steps:
    # Open demo signup page
    - goto: https://example.testproject.io/web/

    # Start signup flow
    - click: "a[href='/register']"

    # Fill registration form
    - fill: "input[name='username']"
      value: "{{ username }}"

    - fill: "input[name='email']"
      value: "{{ email }}"

    - fill: "input[name='password']"
      value: "{{ password }}"

    - fill: "input[name='confirmPassword']"
      value: "{{ password }}"

    # Select role from dropdown
    - click: "select#role"
    - select: "select#role"
      value: "tester"

    # Submit registration
    - click: "button[type='submit']"

    # Navigate to dashboard page
    - goto: "https://example.testproject.io/dashboard"

    # Open API key modal
    - click: "button#create-api-key"

    # Extract generated API key
    - element: "input#api-key"
      mode: value
      set-var: result

    - screenshot: "#qrcode"
      set-var: qr
```
