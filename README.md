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
├── tests/                 # Test suite
│   ├── integration/      # Integration tests
│   ├── benchmark/        # Performance benchmarks
│   └── fuzz/             # Fuzz testing
├── internal/              # Private application code
│   └── testutil/         # Test utilities and helpers
├── docker/               # Docker configurations
├── .github/              # GitHub Actions workflows
│   ├── workflows/        # CI/CD pipelines
│   ├── reports/          # Test reports and coverage
│   ├── ISSUE_TEMPLATE/   # Issue templates
│   ├── CODEOWNERS        # Code owners definition
│   ├── dependabot.yml    # Dependency auto-update configuration
│   ├── labeler.yml       # Auto-labeling configuration
│   └── ...              # Other GitHub configurations
├── .dockerignore         # Docker ignore file
├── .gitignore           # Git ignore file
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
- ✅ **Comprehensive testing suite**: Unit tests, integration tests, benchmarks, and fuzz tests
- ✅ **Test automation**: Automated testing with coverage reports and CI/CD integration

## Quick Start

### Prerequisites

- Go 1.23.0 or later (currently supports up to Go 1.24.3)
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

## Testing

This template includes a comprehensive testing framework with multiple test types to ensure code quality and reliability.

### Test Architecture

The project includes four types of tests:

1. **Unit Tests** (`core/config/`, `cmd/app/`, `cmd/cli/`):
   - Test individual functions and components in isolation
   - 21 tests covering configuration management, main application logic, and CLI functionality
   - Use testify framework for assertions and test utilities

2. **Integration Tests** (`tests/integration/`):
   - Test complete workflows and component interactions
   - 8 tests covering end-to-end scenarios with real binaries
   - Validate cross-component functionality

3. **Benchmark Tests** (`tests/benchmark/`):
   - Performance testing and memory profiling
   - 10 benchmark functions measuring execution time and memory usage
   - Help identify performance bottlenecks

4. **Fuzz Tests** (`tests/fuzz/`):
   - Automated testing with random inputs to discover edge cases
   - 4 fuzz functions targeting JSON parsing and configuration handling
   - Help improve code robustness

### Test Reports and Coverage

All test outputs are organized in `.github/reports/`:
- `test-results.xml`: JUnit format XML test reports for CI/CD integration
- `coverage.out`: Raw coverage data for analysis
- `coverage.html`: Visual coverage reports with line-by-line analysis

### Test Dependencies

The testing framework uses:
- **testify/v1.9.0**: Assertions and testing utilities
- **go-junit-report/v2**: XML report generation in JUnit format
- **Built-in Go tools**: Coverage analysis and benchmarking

### Available Make Commands

```bash
make help              # Show available commands
make all               # Build all commands locally
make build-all         # Build for all platforms
make clean             # Remove build artifacts
make fmt               # Format Go code
make test              # Run comprehensive test suite (unit, integration, benchmark, fuzz tests + coverage + linting)
make run               # Build and run the application
make install           # Install binaries to system (/usr/local/bin)
```

The `make test` command automatically:
- 🧪 Runs all unit tests with coverage analysis
- 🔗 Executes integration tests 
- ⚡ Performs benchmark testing
- 🎯 Conducts fuzz testing for edge case discovery
- 🔍 Runs code linting and static analysis
- 📊 Generates XML test reports (JUnit format)
- 📈 Creates visual coverage reports
- 📁 Saves all outputs to `.github/reports/`

## Testing

This project includes a comprehensive testing suite with multiple types of tests and XML reporting capabilities.

### Test Types

1. **Unit Tests**: Test individual functions and methods in isolation
   - Location: `*_test.go` files alongside source code
   - Coverage: `cmd/app` (83.3%), `cmd/cli` (94.1%), `core/config` (94.7%)

2. **Integration Tests**: Test end-to-end functionality with real binaries
   - Location: `tests/integration/`
   - Tests actual binary execution and command-line interfaces

3. **Benchmark Tests**: Performance testing and memory allocation analysis
   - Location: `tests/benchmark/`
   - Monitors CPU time, memory usage, and allocations per operation

4. **Fuzz Tests**: Discover edge cases and potential panics with random input
   - Location: `tests/fuzz/`
   - Tests JSON parsing, config loading, and path handling

### Running Tests

#### Quick test run:
```bash
make test
```

#### Run all tests with coverage and XML report:
```bash
make test-all
```

#### Run tests with XML output only:
```bash
make test-xml
```

#### Run tests with both coverage and XML reports:
```bash
make test-xml-coverage
```

#### Run specific test types:
```bash
make test-integration    # Integration tests only
make test-benchmark      # Benchmark tests only
make test-race          # Race condition detection
make test-fuzz          # Fuzz tests (30s duration)
```

### Test Reports and XML Output

All test reports are generated in the `.github/reports/` directory:

- **`.github/reports/test-results.xml`** - JUnit XML format test results
  - Compatible with Jenkins, GitHub Actions, Azure DevOps, and other CI/CD systems
  - Contains test results, execution times, and error details
  - Includes coverage information when using `test-xml-coverage`

- **`.github/reports/coverage.html`** - Visual coverage report
- **`.github/reports/coverage.out`** - Raw coverage data

### Test Coverage Goals

- **Overall Coverage**: > 80% for core packages
- **Core Package (`core/config`)**: 94.7% 
- **CLI Package (`cmd/cli`)**: 94.1%
- **App Package (`cmd/app`)**: 83.3%

### CI/CD Integration

GitHub Actions automatically:
- Runs all test types on push/PR
- Generates XML reports and coverage data
- Uploads reports as artifacts
- Displays test results in GitHub interface
- Tests cross-platform compilation

### Test Statistics

- **Total Tests**: 67 test cases
- **Test Suites**: 7 packages
- **All Passing**: ✅ Zero failures
- **Execution Time**: < 30 seconds

## GitHub Actions Workflows

This template includes several pre-configured workflows:

- **`build_package.yml`**: Builds and releases packages for all platforms
- **`build_image.yml`**: Builds and pushes Docker images
- **`auto_labeler.yml`**: Automatically labels pull requests
- **`jira.yml`**: JIRA integration for issue tracking
- **`updater.yml`**: Automated dependency updates
- **`secret_scan.yml`**: Security scanning workflow

### Customizing Workflows

1. **Update repository references**:
   - Replace example repository URLs with your actual repository URL
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