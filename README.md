# Go Template Project

English | [繁體中文](README.zh-TW.md) | [简体中文](README.zh-CN.md)

## Overview

This is a Go project template with a well-structured directory layout and comprehensive GitHub Actions workflows. It provides a solid foundation for building Go applications with multiple commands, cross-platform builds, and automated CI/CD pipelines.

## Project Structure

```
.
├── cmd/                    # Command-line applications
│   ├── app/               # Main application
│   └── cli/               # CLI tool
├── core/                  # Core business logic
│   └── config/           # Configuration management
├── docker/               # Docker configurations
├── .github/              # GitHub Actions workflows
│   ├── workflows/        # CI/CD pipelines
│   ├── ISSUE_TEMPLATE/   # Issue templates
│   └── ...              # Other GitHub configurations
├── go.mod               # Go module definition
├── Makefile            # Build automation
└── README.md           # This file
```

## Features

- ✅ **Multi-command structure**: Separate `app` and `cli` commands
- ✅ **Cross-platform builds**: Support for Linux, macOS, Windows (AMD64 & ARM64)
- ✅ **GitHub Actions**: Complete CI/CD pipeline with automated builds and releases
- ✅ **Docker support**: Ready-to-use Docker configuration
- ✅ **Configuration management**: Flexible config system with environment support
- ✅ **Makefile automation**: Easy build, test, and package commands
- ✅ **Multi-language README**: English, Traditional Chinese, Simplified Chinese

## Quick Start

### Prerequisites

- Go 1.23.0 or later
- Make (for using Makefile commands)
- Docker (optional, for containerization)

### Installation

1. **Clone or use this template**:
   ```bash
   git clone <your-repo-url>
   cd go-template
   ```

2. **Update module name**:
   ```bash
   # Replace 'go-template' with your actual module name in go.mod
   go mod edit -module your-module-name
   ```

3. **Install dependencies**:
   ```bash
   go mod tidy
   ```

### Building

#### Build all commands locally:
```bash
make all
```

#### Build for specific platforms:
```bash
make build_linux_amd64
make build_windows_amd64
make build_darwin_arm64
```

#### Build for all platforms:
```bash
make build-all
```

#### Create distribution packages:
```bash
make package-all
```

### Running

#### Run the main application:
```bash
./build/app
# or
make run
```

#### Run the CLI tool:
```bash
./build/cli --help
./build/cli --version
```

## Development

### Project Customization

1. **Update application names**:
   - Modify `BIN_NAME` and `CLI_NAME` in `Makefile`
   - Update binary names in GitHub Actions workflows

2. **Add your business logic**:
   - Implement your application logic in `cmd/app/main.go`
   - Add CLI commands and functionality in `cmd/cli/main.go`
   - Create additional packages under `core/` for shared logic

3. **Configuration**:
   - Modify `core/config/config.go` to add your configuration fields
   - Update default configuration values as needed

4. **Docker**:
   - Customize `docker/Dockerfile` for your application needs
   - Update Docker image names in GitHub Actions

### Available Make Commands

```bash
make help              # Show available commands
make all               # Build all commands locally
make build-all         # Build for all platforms
make package-all       # Create distribution packages
make clean             # Remove build artifacts
make fmt               # Format Go code
make test              # Run tests
make run               # Build and run the application
```

## GitHub Actions Workflows

This template includes several pre-configured workflows:

- **`build_package.yml`**: Builds and releases packages for all platforms
- **`build_image.yml`**: Builds and pushes Docker images
- **`auto_labeler.yml`**: Automatically labels pull requests
- **`jira.yml`**: JIRA integration for issue tracking
- **`updater.yml`**: Automated dependency updates

### Customizing Workflows

1. **Update repository references**:
   - Replace `gitea.mediatek.inc/IT-GAIA/go-template` with your repository URL
   - Update Docker registry URLs if needed

2. **Configure secrets**:
   - `GITHUB_TOKEN`: For repository access
   - `GT_TOKEN`: For Docker registry access
   - `JIRA_TOKEN`: For JIRA integration (optional)
   - `SSH_KEY`: For SSH access (if needed)

3. **Modify build targets**:
   - Update platform targets in workflows if you don't need all platforms
   - Adjust build commands as needed

## Configuration

The application supports configuration through:

1. **Configuration file**: `~/.go-template/config.json`
2. **Environment variables**: `CONFIG_PATH` to specify custom config location
3. **Default values**: Built-in defaults for development

Example configuration:
```json
{
  "version": "1.0.0",
  "environment": "production",
  "log_level": "info",
  "debug": false
}
```

## Docker Support

Build Docker image:
```bash
docker build -f docker/Dockerfile -t your-app .
```

Run with Docker:
```bash
docker run --rm your-app
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This template is provided as-is for your use. Add your own license as needed.

## Support

- Create issues for bugs or feature requests
- Check existing documentation and examples
- Review GitHub Actions logs for build issues

---

## Next Steps

After setting up this template:

1. **Customize the application logic** in `cmd/` directories
2. **Add your business logic** in `core/` packages
3. **Update configuration** to match your needs
4. **Modify GitHub Actions** for your CI/CD requirements
5. **Add tests** in appropriate directories
6. **Update documentation** to reflect your application

Happy coding! 🚀