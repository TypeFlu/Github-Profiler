# 🚀 GitHub Profiler

A powerful CLI tool for analyzing GitHub user profiles with advanced statistics, visualizations, and developer rankings.

![GitHub Profiler](https://img.shields.io/badge/Status-Ready-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.24+-blue)
![License](https://img.shields.io/badge/License-MIT-green)

*Created by github@Tyeflu*

## ✨ Features

- 🎨 **Beautiful Terminal UI** - Colorful, well-organized output
- 📊 **Advanced Statistics** - Comprehensive user and repository analysis  
- 🏆 **Developer Ranking** - Smart scoring system based on activity
- 📈 **Language Analysis** - Programming language usage breakdown
- ⭐ **Repository Insights** - Top repositories and contribution patterns
- 📋 **Multiple Formats** - Terminal, JSON, and HTML output
- ⚡ **Fast & Lightweight** - Single binary, no dependencies

## 🚀 Quick Start

### Using Pre-built Binaries

1. Download the latest release for your platform
2. Extract and move to your PATH:
   ```bash
   # macOS/Linux
   tar -xzf github-profiler-*-your-platform.tar.gz
   sudo mv github-profiler /usr/local/bin/
   
   # Windows
   # Extract github-profiler.exe and add to PATH
   ```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/Typeflu/Temp.git
cd Temp

# Build and install
make install

# Or just build
make build
./build/github-profiler --help
```

## 💻 Usage

### Basic Usage
```bash
# Analyze any GitHub user
github-profiler octocat

# With custom GitHub token for higher rate limits
GITHUB_TOKEN=your_token github-profiler username

# Output as JSON
github-profiler username --format json

# Generate HTML report
github-profiler username --format html
```

### Command Line Options
```
Flags:
  -f, --format string   Output format: terminal, json, html (default "terminal")
  -h, --help           help for github-profiler
  -t, --token string   GitHub personal access token (optional for public data)
```

### Environment Variables
- `GITHUB_TOKEN` - Your GitHub personal access token (increases rate limits)

## 📊 Sample Output

### Terminal Output
```
🚀 GitHub Profile Analysis

┌─────────────────────────────────────────────┐
│ 👤 Octocat (octocat)                       │
│ 📝 A great place to start                  │
│ 🏢 GitHub                                  │
│ 📍 San Francisco                           │
│ 🌐 https://github.blog                     │
│ 📅 Joined: January 2011                    │
└─────────────────────────────────────────────┘

📊 Statistics
👥 Followers:      4,000
➡️  Following:      9
📦 Public Repos:   8
🔒 Private Repos:  0
📄 Gists:         8
🏢 Organizations:  1

🎯 Contributions (This Year)
💻 Commits:        52
🐛 Issues:         0
🔄 Pull Requests:  0
👀 PR Reviews:     0
📦 Repositories:   0

⭐ Top Repositories
┌─────────────────────────────────────────┐
│ 📦 Hello-World                         │
│ My first repository on GitHub!         │
│ ⭐ 1,500  🍴 1,000  💻 JavaScript      │
│ 📅 Updated: Mar 2023                   │
└─────────────────────────────────────────┘

💻 Language Breakdown
JavaScript      ████████████████████ 45.0% (9 repos)
Python          ████████████ 25.0% (5 repos)
Go              ████████ 20.0% (4 repos)
TypeScript      ████ 10.0% (2 repos)

🏆 Developer Ranking
🏆 Developer Rank: ⚡ Active Developer
📊 Total Score: 75/100
👥 Social Score: 25/30
📦 Repository Score: 30/40
🎯 Activity Score: 20/30
```

### HTML Report
The HTML format generates a beautiful, interactive web report with:
- Responsive design
- Interactive charts
- Professional styling
- Social sharing ready
- Print-friendly layout

## 🏗️ Development

### Prerequisites
- Go 1.21 or higher
- Make (for build automation)

### Building
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

### Available Make Commands
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

### Project Structure
```
github-profiler/
├── cmd/           # CLI commands and logic
│   ├── root.go    # Main command implementation
│   └── output.go  # Output formatting (JSON, HTML)
├── build/         # Build artifacts
├── main.go        # Application entry point
├── go.mod         # Go module definition
├── Makefile       # Build automation
└── README.md      # This file
```

## 🎯 Ranking System

The developer ranking system evaluates users across three dimensions:

### Social Score (30 points)
- **Elite (30pts)**: 10,000+ followers
- **High (25pts)**: 1,000+ followers  
- **Medium (20pts)**: 500+ followers
- **Growing (15pts)**: 100+ followers
- **Active (10pts)**: 50+ followers
- **Starter (5pts)**: 10+ followers

### Repository Score (40 points)
- **Repository Count**: 0.5 points per repo (max 20pts)
- **Star Count**: 0.1 points per star (max 20pts)
- Excludes private repos and forks

### Activity Score (30 points)
- **Contributions**: 0.01 points per contribution
- Includes commits, issues, PRs, reviews, and repos
- Based on current year activity

### Final Rankings
- 🌟 **Elite Developer** (90-100pts)
- 🚀 **Senior Developer** (80-89pts)  
- 💪 **Experienced Developer** (70-79pts)
- ⚡ **Active Developer** (60-69pts)
- 🌱 **Growing Developer** (50-59pts)
- 👶 **Junior Developer** (30-49pts)
- 🥚 **Beginner** (<30pts)

## 🔧 Configuration

### GitHub Token Setup
While not required for public data, a GitHub token provides:
- Higher rate limits (5,000 vs 60 requests/hour)
- Access to additional user data
- Better reliability

```bash
# Set environment variable
export GITHUB_TOKEN="your_personal_access_token"

# Or pass directly
github-profiler username --token your_token
```

To create a token:
1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Generate new token with `public_repo` scope
3. Copy and use the token

## 📝 Output Formats

### Terminal (Default)
- Colorful, emoji-rich display
- Perfect for command-line usage
- Instant visual feedback

### JSON
- Machine-readable format
- Integration with other tools
- API-friendly structure

### HTML
- Beautiful web report
- Interactive charts and graphs
- Professional presentation
- Easy sharing and printing

## 🤝 Contributing

We welcome contributions! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Workflow
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- GitHub GraphQL API for providing comprehensive user data
- [Charm](https://charm.sh/) libraries for beautiful terminal UI
- [Cobra](https://cobra.dev/) for CLI framework
- Go community for excellent tooling

## 🔗 Links

- [GitHub Repository](https://github.com/TypeFlu/Temp)
- [Report Issues](https://github.com/Typeflu/temp/issues)
- [Latest Releases](https://github.com/Typeflu/Temp/releases)

---

Built with ❤️ by [@typeflu](https://github.com/typeflu)
