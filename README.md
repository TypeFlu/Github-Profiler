# GitHub-Profiler

A professional CLI tool for comprehensive GitHub user profile analysis with advanced statistics, visualizations, and developer rankings.

[![Build Status](https://img.shields.io/badge/Status-Production%20Ready-brightgreen)](https://github.com/Typeflu/Github-Profiler)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![GitHub Release](https://img.shields.io/github/release/Typeflu/Github-Profiler.svg)](https://github.com/TypeFlu/Github-Profiler/releases/latest)

## Overview

GitHub Profiler is a modern terminal-based application built with Go that provides comprehensive analysis of GitHub user profiles. It features an interactive Terminal User Interface (TUI) powered by Bubble Tea and follows the Elm Architecture pattern for maintainable, scalable code.

## Features

### Core Functionality
- **Interactive Terminal UI** - Professional, keyboard-driven interface
- **Advanced Statistics** - Comprehensive user and repository analysis  
- **Developer Ranking System** - Intelligent scoring based on contribution patterns
- **Language Analysis** - Programming language usage breakdown with statistics
- **Repository Insights** - Detailed analysis of top repositories and contribution patterns
- **Real-time Data** - Live GitHub API integration with authentication support
- **Demo Mode** - Offline demonstration with mock data

### Technical Features
- **Modern Architecture** - Elm Architecture pattern with proper separation of concerns
- **Professional UI** - Clean, corporate-friendly terminal interface
- **Type Safety** - Full type safety with go-github v73 integration
- **Error Handling** - Robust error handling and user feedback
- **Cross-platform** - Single binary deployment across platforms

## Quick Start

### Installation

#### Option 1: Download Pre-built Binary
```bash
# Download latest release from GitHub
curl -L https://github.com/Typeflu/Temp/releases/latest/download/github-profiler-linux-amd64.tar.gz | tar xz
sudo mv github-profiler /usr/local/bin/
```

#### Option 2: Build from Source
```bash
git clone https://github.com/Typeflu/Temp.git
cd Temp
make install
```

#### Option 3: Go Install
```bash
go install github.com/Typeflu/Temp@latest
```

### Basic Usage

```bash
# Interactive mode - enter username in TUI
github-profiler

# Direct username analysis
github-profiler octocat

# Demo mode with mock data
github-profiler demo

# With GitHub token for higher rate limits
GITHUB_TOKEN=your_token github-profiler username
```
## Interface Navigation

### Keyboard Controls
- **Left/Right Arrow Keys** - Navigate between different views
- **r** - Refresh data from GitHub API
- **q** - Quit application
- **Enter** - Confirm actions in interactive mode

### Available Views
1. **Overview** - User profile summary and key statistics
2. **Repositories** - Top repositories with detailed metrics
3. **Languages** - Programming language breakdown and statistics
4. **Activity** - Contribution patterns and activity metrics
5. **Ranking** - Developer ranking and scoring breakdown

## Architecture

The application follows the Elm Architecture pattern for maintainable, scalable code:

```
github-profiler/
├── cmd/                    # CLI commands and entry points
│   └── root.go            # Main command and TUI initialization
├── internal/              # Internal application code
│   ├── models/            # Domain models and data structures
│   │   └── profile.go     # GitHub profile models
│   ├── services/          # Service layer
│   │   ├── github.go      # GitHub API client
│   │   └── mock.go        # Mock data for demo mode
│   └── ui/                # User interface components
│       └── model.go       # Bubble Tea TUI (Elm Architecture)
├── main.go                # Application entry point
├── go.mod                 # Go module definition
├── Makefile              # Build automation
└── README.md             # Documentation
```

### Key Design Principles
- **Separation of Concerns** - Clear boundaries between data, business logic, and UI
- **Type Safety** - Comprehensive type safety with go-github v73
- **Error Handling** - Graceful error handling with user-friendly messages
- **Testability** - Modular design enabling comprehensive testing
- **Professional Standards** - Clean, maintainable Go code following best practices

## Sample Output

### Terminal Interface
```
GitHub Profiler - Enhanced TUI Version

User: Octocat (octocat)
Bio: A great place to start
Company: GitHub
Location: San Francisco
Website: https://github.blog
Joined: January 2011
Stats: Public Repos: 8 | Followers: 4000 | Following: 9

Total Stars: 2500
Total Forks: 1000
Repository Size: 45.2 MB
Avg Stars/Repo: 312.5

Top Repositories:
1. Hello-World
   My first repository on GitHub!
   Stars: 1500  Forks: 1000  Language: JavaScript
   Updated: Mar 2023

2. Spoon-Knife
   This repo is for demonstration purposes only
   Stars: 800  Forks: 500  Language: HTML
   Updated: Feb 2023

Programming Languages:
JavaScript    45.0% (9 repositories)
Python        25.0% (5 repositories)
Go            20.0% (4 repositories)
TypeScript    10.0% (2 repositories)
```

## Development

### Prerequisites
- Go 1.24 or higher
- Make (for build automation)
- Git

### Dependencies
The project uses the latest versions of key dependencies:
- `github.com/google/go-github/v73` - GitHub API client
- `github.com/charmbracelet/bubbletea` - Terminal UI framework
- `github.com/charmbracelet/bubbles` - TUI components
- `github.com/charmbracelet/lipgloss` - Style and layout
- `github.com/spf13/cobra` - CLI framework
- `golang.org/x/oauth2` - OAuth2 authentication

### Build Commands
```bash
# Install dependencies
make deps

# Build for current platform
make build

# Install to system PATH
make install

# Uninstall from system
make uninstall

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
make lint

# Clean build artifacts
make clean

# Show all available commands
make help
```

### Available Make Targets
| Command | Description |
|---------|-------------|
| `make build` | Build binary for current platform |
| `make install` | Install to system PATH (/usr/local/bin) |
| `make uninstall` | Remove from system PATH |
| `make build-all` | Build for all platforms |
| `make release` | Create release archives |
| `make test` | Run tests |
| `make lint` | Run code linter |
| `make clean` | Clean build artifacts |
| `make deps` | Install dependencies |
| `make help` | Show all commands |

## Developer Ranking System

The application includes a comprehensive developer ranking system that evaluates users across multiple dimensions:

### Scoring Components

#### Social Score (25 points maximum)
- **Elite (25pts)**: 10,000+ followers
- **High (20pts)**: 1,000+ followers  
- **Medium (15pts)**: 500+ followers
- **Growing (10pts)**: 100+ followers
- **Active (5pts)**: 50+ followers
- **Starter (2pts)**: 10+ followers

#### Code Score (30 points maximum)
- **Repository Count**: 0.5 points per repository (max 15pts)
- **Star Count**: 0.02 points per star (max 15pts)
- Excludes forked repositories to focus on original contributions

#### Activity Score (25 points maximum)
- **Recent Contributions**: Based on commit activity, issues, and pull requests
- **Repository Updates**: Frequency of repository maintenance
- **Community Engagement**: Issues opened, PRs created, and reviews

#### Innovation Score (20 points maximum)
- **Repository Diversity**: Variety of programming languages used
- **Project Complexity**: Assessment based on repository size and structure
- **Community Impact**: Stars and forks received on original projects

### Final Rankings
- **Elite Developer** (90-100pts) - Industry leaders and open source maintainers
- **Senior Developer** (80-89pts) - Experienced professionals with strong contributions
- **Experienced Developer** (70-79pts) - Skilled developers with consistent activity
- **Active Developer** (60-69pts) - Regular contributors to the community
- **Growing Developer** (50-59pts) - Emerging developers building their portfolio
- **Junior Developer** (30-49pts) - Early-career developers learning and contributing
- **Beginner** (<30pts) - New to the platform or development in general

## Configuration

### GitHub Authentication
While the application works with public GitHub data without authentication, using a GitHub Personal Access Token provides several benefits:

- **Higher Rate Limits**: 5,000 requests/hour vs 60 requests/hour for unauthenticated requests
- **Access to Additional Data**: Some profile information requires authentication
- **Better Reliability**: Reduced chance of hitting API limits during analysis

#### Setting up Authentication
```bash
# Method 1: Environment variable
export GITHUB_TOKEN="your_personal_access_token"
github-profiler username

# Method 2: Command line flag
github-profiler username --token your_personal_access_token
```

#### Creating a GitHub Token
1. Navigate to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Select scopes: `public_repo` (sufficient for public data analysis)
4. Generate token and copy the value
5. Store securely and use as shown above

### Environment Variables
- `GITHUB_TOKEN` - GitHub Personal Access Token for API authentication

## API Integration

The application integrates with the GitHub REST API v4 using the official go-github library (v73). Key features include:

- **Comprehensive Data Retrieval** - User profiles, repositories, languages, and activity
- **Rate Limit Handling** - Automatic rate limit detection and user feedback
- **Error Handling** - Graceful handling of API errors and network issues
- **Authentication Support** - Optional token authentication for enhanced access
- **Type Safety** - Full type safety with structured API responses

## Contributing

We welcome contributions from the community! Whether you're fixing bugs, adding features, or improving documentation, your help is appreciated.

### Development Workflow
1. **Fork** the repository on GitHub
2. **Clone** your fork locally
3. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
4. **Make** your changes following the project's coding standards
5. **Test** your changes thoroughly
6. **Commit** your changes (`git commit -m 'Add amazing feature'`)
7. **Push** to your branch (`git push origin feature/amazing-feature`)
8. **Create** a Pull Request with a clear description

### Coding Standards
- Follow standard Go formatting (`go fmt`)
- Write comprehensive tests for new features
- Update documentation for API changes
- Maintain the professional, emoji-free style
- Ensure compatibility with the Elm Architecture pattern

### Testing
```bash
# Run all tests
make test

# Run with coverage
go test -v -cover ./...

# Run linter
make lint
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support and Community

- **Issues**: [GitHub Issues](https://github.com/Typeflu/Temp/issues) for bug reports and feature requests
- **Discussions**: [GitHub Discussions](https://github.com/Typeflu/Temp/discussions) for questions and community support
- **Releases**: [GitHub Releases](https://github.com/Typeflu/Temp/releases) for latest versions and changelogs

## Acknowledgments

This project builds upon excellent open-source libraries and tools:

- **[Charm](https://charm.sh/)** - For the beautiful and powerful terminal UI libraries
- **[Cobra](https://cobra.dev/)** - For the CLI framework and command handling
- **[go-github](https://github.com/google/go-github)** - For comprehensive GitHub API integration
- **Go Community** - For excellent tooling and ecosystem support

## Technical Specifications

### System Requirements
- **Operating System**: Linux, macOS, Windows
- **Architecture**: amd64, arm64
- **Memory**: Minimum 64MB RAM
- **Network**: Internet connection for GitHub API access

### Performance Characteristics
- **Startup Time**: < 100ms
- **Memory Usage**: < 50MB typical
- **API Calls**: Optimized to minimize requests
- **Response Time**: 1-3 seconds for typical user analysis

---

**GitHub Profiler** - Professional GitHub user analysis tool  
Built with Go, powered by modern terminal UI technology

For more information, visit the [project repository](https://github.com/Typeflu/Temp).
