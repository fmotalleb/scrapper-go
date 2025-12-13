# Scrapper-Go API Server Documentation

This document provides a detailed overview of the Scrapper-Go API server, its endpoints, and their functionalities. The API allows for both stateless, one-shot scraping tasks and stateful, interactive browser sessions.

## Overview

The API server is built using the [Echo](https://echo.labstack.com/) framework. It exposes several endpoints to control and interact with the scraping engine.

To start the server, use the `serve` command:
```bash
./scrapper-go serve --address 127.0.0.1 --port 8080
```

---

## 1. Stateless Processing

This endpoint is for simple, one-shot scraping tasks. You provide a full pipeline configuration, and the server executes it and returns the final result.

### `POST /process`

Executes a complete scraping pipeline and returns the collected data.

**Request:**

The request body must be a JSON object containing the full pipeline configuration, identical to the structure used in YAML files.

- **`Content-Type: application/json`**

**Example `curl`:**

```bash
curl -X POST http://127.0.0.1:8080/process \
-H "Content-Type: application/json" \
-d '{
  "pipeline": {
    "browser": "chromium",
    "browser_params": {
      "headless": true
    },
    "steps": [
      {
        "goto": "https://example.com"
      },
      {
        "element": "h1",
        "mode": "text",
        "set-var": "title"
      }
    ]
  }
}'
```

**Response:**

A JSON object containing the variables set during the execution (e.g., using `set-var`).

- **`200 OK`**:
  ```json
  {
    "title": "Example Domain"
  }
  ```
- **`400 Bad Request`**: If the configuration is invalid or an error occurs during execution.

---

## 2. Stateful Sessions

The session endpoints allow for creating persistent browser sessions that can be controlled interactively. This is useful for complex scenarios that require multiple, separate steps, such as logging in and then performing various actions.

### `POST /sessions`

Creates a new persistent browser session.

**Query Parameters:**
- `timeout` (optional): A duration string (e.g., `5m`, `1h`) that specifies how long the session should live before being automatically terminated. Defaults to `5m`.

**Request:**

The request body contains the initial browser and page configuration for the session. The `steps` part of the configuration is ignored.

**Example `curl`:**
```bash
curl -X POST "http://127.0.0.1:8080/sessions?timeout=10m" \
-H "Content-Type: application/json" \
-d '{
  "browser": "chromium",
  "browser_params": {
    "headless": false
  }
}'
```

**Response:**

- **`200 OK`**: A JSON object containing the new session's ID and its timeout.
  ```json
  {
    "id": "a1b2c3d4-e5f6-...",
    "timeout": "10m0s"
  }
  ```
- **`400 Bad Request`**: If the configuration is invalid.

---

### `GET /sessions`

Retrieves a list of all currently active session IDs.

**Example `curl`:**

```bash
curl -X GET http://127.0.0.1:8080/sessions
```

**Response:**

- **`200 OK`**: A JSON array of session ID strings.
  ```json
  [
    "a1b2c3d4-e5f6-...",
    "f7g8h9i0-j1k2-..."
  ]
  ```

---

### `POST /sessions/:id`

Executes one or more steps within an existing session.

**Request:**

The request body can be either a single JSON object representing one step, or a JSON array of step objects.

**Example `curl` (single step):**

```bash
curl -X POST http://127.0.0.1:8080/sessions/a1b2c3d4-e5f6-ùi \
-H "Content-Type: application/json" \
-d '{
  "goto": "https://github.com/fmotalleb/scrapper-go"
}'
```

**Example `curl` (multiple steps):**

```bash
curl -X POST http://127.0.0.1:8080/sessions/a1b2c3d4-e5f6-ùi \
-H "Content-Type: application/json" \
-d '[
  {
    "fill": "#username",
    "value": "my-user"
  },
  {
    "fill": "#password",
    "value": "my-secret-pass"
  },
  {
    "click": "button[type=submit]"
  }
]'
```

**Response:**

- **`200 OK`**: A JSON object containing the results from the executed steps.
- **`404 Not Found`**: If the session ID does not exist.
- **`400 Bad Request`**: If the step configuration is invalid or an error occurs during execution.

---

### `DELETE /sessions/:id`

Terminates and cleans up a specific browser session.

**Example `curl`:**

```bash
curl -X DELETE http://127.0.0.1:8080/sessions/a1b2c3d4-e5f6-ùi
```

**Response:**

- **`200 OK`**: Confirms the session was killed.
  ```json
  {
    "id": "a1b2c3d4-e5f6-ùi"
  }
  ```
- **`44 Not Found`**: If the session ID does not exist.

---

## 3. Live Streaming (WebSocket)

For the most interactive experience, the live stream endpoint provides a WebSocket connection for real-time, bidirectional communication with the scraping engine.

### `GET /live-stream`

Upgrades the HTTP connection to a WebSocket connection.

**Connection Flow:**

1.  **Client connects:** A client establishes a WebSocket connection to `ws://127.0.0.1:8080/live-stream`.

2.  **Client sends initial config:** The client sends the initial browser and page configuration as a JSON message.
    ```json
    {
      "browser": "chromium",
      "browser_params": { "headless": false }
    }
    ```

3.  **Server creates session:** The server starts a browser session based on the config.

4.  **Real-time interaction:**
    - The client can now send individual step configurations as JSON messages.
    - The server executes each step and sends the result (or an error) back to the client as a JSON message. This continues until the connection is closed.

**Example Client-Side JavaScript:**

```javascript
const ws = new WebSocket("ws://127.0.0.1:8080/live-stream");

ws.onopen = () => {
  console.log("WebSocket connection established.");

  // 1. Send initial browser configuration
  const initialConfig = {
    browser: "chromium",
    browser_params: { headless: true },
    steps: [] // Steps can be sent here or later
  };
  ws.send(JSON.stringify(initialConfig));

  // 2. Send a step to execute
  const gotoStep = { goto: "https://example.com" };
  ws.send(JSON.stringify(gotoStep));

  // 3. Send another step to get the title
  const getTitleStep = { element: 'h1', mode: 'text', 'set-var': 'title' };
  ws.send(JSON.stringify(getTitleStep));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log("Received from server:", message);
  // Example response: { title: 'Example Domain' }
};

ws.onclose = () => {
  console.log("WebSocket connection closed.");
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};
```