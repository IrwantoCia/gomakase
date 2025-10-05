# Gomakase 🍣

[![Go Version](https://img.shields.io/badge/Go-1.24.5+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

A CLI tool for generating Go projects with DDD structure. It helps developers quickly scaffold well-organized, scalable Go projects following clean architecture principles.

## 🎯 What is Gomakase?

Gomakase is a CLI tool that generates production-ready Go applications with Domain-Driven Design (DDD) architecture. It eliminates boilerplate setup time and provides a solid foundation for scalable applications.

### 🚀 Key Benefits

- **Instant Setup**: Generate complete Go projects in seconds, not hours
- **Clean Architecture**: Built-in DDD structure with proper separation of concerns
- **Production-Ready**: Includes web server, database setup, logging, and configuration management
- **Best Practices**: Follows Go community standards and conventions
- **Extensible**: Plugin system for adding functionality like authentication

### 🎯 Perfect For

- **Startups** needing to build applications quickly
- **Teams** wanting consistent project structure
- **Developers** learning DDD architecture
- **Companies** establishing Go development standards

## 📦 Installation

### Prerequisites

- **Go 1.24.5 or higher** - [Install Go](https://golang.org/dl/)

### Method 1: Install via Go Install (Recommended)

```bash
go install github.com/IrwantoCia/gomakase@latest
```

### Method 2: Install from Source

```bash
git clone https://github.com/IrwantoCia/gomakase.git
cd gomakase
make install
```

### Verify Installation

```bash
gomakase --help
# Should show the available commands
```

### 🚨 Troubleshooting

**Issue:** `gomakase: command not found`
- **Solution:** Ensure `~/go/bin` is in your PATH
```bash
echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
source ~/.bashrc
```

## 🚀 Quick Start

Generate your first Go project in under 30 seconds:

### Step 1: Create a New Project

```bash
gomakase new myapp
```

### Step 2: See What Was Generated

```bash
cd myapp
find . -type f -name "*.go" | head -20
```

**Expected Structure:**
```
myapp/
├── cmd/
│   └── server/
│       ├── main.go              # Application entry point
│       ├── app.go               # Application setup
│       └── router.go            # HTTP routes
├── internal/
│   └── shared/                  # Shared utilities
│       ├── config/              # Configuration management
│       ├── db/                  # Database connection
│       ├── logger/              # Logging utilities
│       └── middleware/          # HTTP middleware
├── web/                         # Frontend assets
│   ├── static/
│   │   └── css/
│   └── views/
│       ├── layouts/
│       └── index.html
├── go.mod                       # Go module definition
├── Makefile                     # Build commands
├── .env.example                 # Environment variables template
└── gen.yaml                     # Project configuration
```

### Step 3: Run Your New Application

```bash
# Run the application in development mode
make dev
```

**Success!** Your application is now running at `http://localhost:8080`

### What Just Happened?

Gomakase created a complete Go application with:
- ✅ Clean architecture folder structure
- ✅ Web server setup
- ✅ Database connection configuration
- ✅ Logging with structured logging
- ✅ Environment configuration
- ✅ Modern web UI templates
- ✅ Development tools (Makefile, Docker support)

## 📖 Commands Reference

### Available Commands

#### `gomakase new <project_name>`
Creates a new Go project with complete DDD structure.

**Syntax:**
```bash
gomakase new <project_name>
```

**Examples:**
```bash
# Basic usage
gomakase new myproject

# Project will be created in ./myproject directory
gomakase new github.com/username/myproject
```

#### `gomakase context <context_name>`
Generates a new business context in an existing project.

**Syntax:**
```bash
gomakase context <context_name>
```

**Examples:**
```bash
# Create a new user context
gomakase context user

# Create a new product context
gomakase context product
```

**Generated Context Structure:**
```
internal/<context_name>/
├── domain/
│   ├── <context_name>.entity.go     # Business entities
│   └── <context_name>.repository.go # Repository interfaces
├── application/
│   └── <context_name>.service.go    # Business logic
├── delivery/
│   └── <context_name>.handler.go    # HTTP handlers
└── infrastructure/
    ├── <context_name>.repository.go # Repository implementation
    └── <context_name>.schema.go     # Database schema
```

#### `gomakase list`
Lists all available plugins that can be added to your project.

**Syntax:**
```bash
gomakase list
```

**Examples:**
```bash
gomakase list
# Output: Available plugins:
#          • auth
```

#### `gomakase add <plugin_name>`
Adds a specific plugin to your existing project.

**Syntax:**
```bash
gomakase add <plugin_name>
```

**Examples:**
```bash
# Add authentication functionality
gomakase add auth
```

**Note:** You must run this command from within a project generated by Gomakase.

### Global Flags

```bash
--help, -h          # Show help information
--toggle, -t        # Toggle flag (placeholder)
```

## 🏗️ Generated Project Structure

Gomakase generates projects following Clean Architecture and Domain-Driven Design principles.

### Core Architecture

```
<project_name>/
├── cmd/                           # Application entry points
│   └── server/
│       ├── main.go               # HTTP server startup
│       ├── app.go                # Application configuration
│       └── router.go             # HTTP routes definition
├── internal/                     # Private application code
│   ├── <context_name>/           # Business contexts
│   │   ├── domain/               # Business logic and entities
│   │   │   ├── entity.go         # Business entity definition
│   │   │   └── repository.go     # Repository interface
│   │   ├── application/          # Use cases and business rules
│   │   │   └── service.go        # Business service implementation
│   │   ├── delivery/             # External interface implementations
│   │   │   └── handler.go        # HTTP handlers
│   │   └── infrastructure/       # External dependencies
│   │       ├── repository.go     # Database repository implementation
│   │       └── schema.go         # Database schema definition
│   └── shared/                   # Shared utilities and cross-cutting concerns
│       ├── config/               # Configuration management
│       │   └── config.go         # Environment variables and settings
│       ├── db/                   # Database connection and setup
│       │   └── db.go             # Database connection logic
│       ├── logger/               # Logging utilities
│       │   └── logger.go         # Structured logger configuration
│       └── middleware/           # HTTP middleware
│           └── logger.go         # Request logging middleware
├── web/                          # Frontend assets
│   ├── static/                   # Static files (CSS, JS)
│   │   └── css/
│   │       └── app.css           # Main stylesheet
│   │   └── js/
│   │       ├── main.js           # Main JavaScript
│   │       ├── api.js            # API utilities
│   │       ├── component.js      # Component utilities
│   │       └── alpine-mixins.js  # Alpine.js mixins
│   └── views/                    # HTML templates
│       ├── layouts/              # Page layouts
│       │   └── master.html       # Base template
│       └── index.html            # Home page
├── .env.example                  # Environment variables template
├── .air.toml                     # Air hot reload configuration
├── .gitignore                    # Git ignore rules
├── docker-compose.yaml           # Docker development environment
├── Dockerfile                    # Container configuration
├── Makefile                      # Build and development commands
├── go.mod                        # Go module definition
├── gen.yaml                      # Project configuration
└── package.json                  # Node.js dependencies (for frontend)
```

### Architecture Principles

#### **Why cmd/?**
The `cmd/` directory contains application entry points. This follows Go's standard project layout where executable programs are placed in `cmd/`.

#### **Why internal/?**
The `internal/` directory contains code that should not be imported by other projects. This prevents external packages from depending on your internal implementation details.

#### **Why Domain-Driven Design?**
Each business context follows a layered architecture:
- **Domain**: Core business logic and entities (independent of frameworks)
- **Application**: Use cases and business rules (orchestrates domain objects)
- **Delivery**: External interfaces (HTTP handlers)
- **Infrastructure**: External dependencies (database, external APIs)

#### **Why shared/?**
The `shared/` directory contains utilities used across multiple contexts but not part of any specific domain.

## 🔧 Available Plugins

### Authentication Plugin

The `auth` plugin adds complete authentication functionality to your project:

**What it includes:**
- JWT-based authentication system
- User registration and login functionality
- Authentication middleware
- Login and registration web pages
- Database schema for users
- API endpoints for authentication

**Usage:**
```bash
# First create a project
gomakase new myapp
cd myapp

# Add authentication
gomakase add auth
```

**Generated Files:**
- `internal/auth/` - Complete authentication context
- `web/views/login.html` - Login page
- `web/views/register.html` - Registration page
- `web/static/js/src/components/login.js` - Login functionality
- `internal/shared/middleware/auth.go` - Authentication middleware

## 📝 Project Configuration

Each generated project includes a `gen.yaml` file for project configuration:

```yaml
project:
  name: "myproject"
  module: "github.com/username/myproject"
  version: "1.0.0"
```

This configuration file is used by Gomakase when adding new contexts or plugins to ensure they are properly integrated.

## 🛠️ Development Commands

Generated projects include a Makefile with the following commands:

```bash
make dev         # Start the development server with hot reload (uses Air)
make build       # Build the application (builds CSS/JS and Go binary)
make run         # Run the built application
make clean       # Clean build artifacts and binaries
make help        # Display available commands and usage
```

## 🤝 Contributing

We love contributions! Whether you're fixing a bug, adding a feature, or improving documentation, we appreciate your help.

### 🛠️ Development Setup

1. **Fork and Clone**
```bash
git clone https://github.com/yourusername/gomakase.git
cd gomakase
```

2. **Install Dependencies**
```bash
go mod download
```

3. **Install for Development**
```bash
make install
```

### 🧪 Testing

Before submitting a pull request, ensure:
- All tests pass: `make test`
- Build succeeds: `make build`

### 📝 Code Style

- Follow Go conventions and `gofmt` formatting
- Use meaningful variable and function names
- Add comments for public functions and complex logic
- Keep functions small and focused

### 🚀 Submitting Changes

1. **Create a Feature Branch**
```bash
git checkout -b feature/your-feature-name
```

2. **Make Your Changes**
- Write code following the style guide
- Update documentation if needed

3. **Test Your Changes**
```bash
make build
```

4. **Commit Your Changes**
```bash
git add .
git commit -m "feat: add new feature description"
```

5. **Push and Create Pull Request**
```bash
git push origin feature/your-feature-name
```

### 🐛 Reporting Issues

When reporting bugs, please include:
- Go version and OS
- Gomakase version
- Steps to reproduce
- Expected vs actual behavior
- Any error messages or logs

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

### What This Means

✅ **You can:**
- Use this software for commercial purposes
- Modify and distribute the software
- Use this software in private applications
- Sublicense the software

❌ **You must:**
- Include the copyright notice in all copies
- Include the license notice in all copies
- Provide attribution to the original project

❌ **You cannot:**
- Hold the authors liable for damages
- Use the authors' names for endorsement

### Third-Party Licenses

This project includes third-party components with their own licenses:
- [Cobra](https://github.com/spf13/cobra) - Apache 2.0
- [Viper](https://github.com/spf13/viper) - MIT

See `go.mod` for a complete list of dependencies and their licenses.

---

**Happy coding with Gomakase! 🍣✨**