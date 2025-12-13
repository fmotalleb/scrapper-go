# Scrapper-Go

Scrapper-Go is a powerful and flexible Go application that acts as a wrapper around Playwright, enabling you to define and execute web scraping pipelines using simple YAML configuration files. It provides a robust engine for automating browser interactions, extracting data, and handling various web scenarios.

## Features

- **YAML-driven Scraping**: Define complex scraping workflows using intuitive YAML configurations.
- **Playwright Integration**: Leverages the full power of Playwright for browser automation, supporting Chromium, Firefox, and WebKit.
- **API Server**: Expose your scraping capabilities as a RESTful API endpoint.
- **Interactive Shell**: Interact with the scrapper in a live shell environment for testing and development.
- **Dependency Management**: Easily install Playwright browsers and drivers with a dedicated setup command.

## Installation

### Prerequisites

- Go (1.18 or higher)
- Node.js (for Playwright dependencies)

### Build from Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/fmotalleb/scrapper-go.git
   cd scrapper-go
   ```

2. **Install Playwright dependencies**:
   ```bash
   go run main.go setup
   ```
   You can specify which browsers to install:
   ```bash
   go run main.go setup --browsers chromium,firefox
   ```
   Or skip browser installation:
   ```bash
   go run main.go setup --skip-browsers
   ```

3. **Build the application**:
   ```bash
   go build -o scrapper-go .
   ```

## Usage

### Executing a Scraping Pipeline

You can run a YAML-defined scraping pipeline directly:

```bash
./scrapper-go -c path/to/your/config.yaml
```

Example `config.yaml`:
```yaml
# Your YAML scraping configuration here
```

### Subcommands

#### `serve` - Start the API Server

Run Scrapper-Go as an API service. By default, it listens on `127.0.0.1:8080`.
**Note**: This application does not support authentication. It is recommended to run it behind a reverse proxy for production use.

```bash
./scrapper-go serve
# Or specify address and port
./scrapper-go serve -a 0.0.0.0 -p 8081
```

#### `setup` - Install Playwright Dependencies

As described in the installation section, this command helps manage Playwright's browsers and drivers.

```bash
./scrapper-go setup --browsers webkit
```

#### `shell` - Interactive Shell

Start an interactive shell for direct interaction and testing of scraping steps.

```bash
./scrapper-go shell
```

## Configuration

Scrapper-Go looks for a configuration file named `.scrapper-go.yaml` in your home directory by default. You can specify a different configuration file using the `-c` or `--config` flag.

## Contributing

We welcome contributions! Please see `CONTRIBUTING.md` (if available) for details on how to contribute.

## License

This project is licensed under the GNU General Public License v2.0 - see the [LICENSE](LICENSE) file for details.