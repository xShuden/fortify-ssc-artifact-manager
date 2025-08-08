# SSC Approver CLI 🔍

A command-line tool to manage and list artifacts requiring approval in Fortify Software Security Center (SSC).

## Features ✨

- 📋 List all artifacts with "Requires Approval" status across all projects
- 🔎 View detailed processing messages (code line changes, etc.)
- 🏷️ Filter by project name
- 📊 Multiple output formats (Table, CSV, JSON)
- 🗂️ List all projects and their versions
- 📦 View artifacts for specific project versions
- 🔐 Flexible authentication (CLI flags, environment variables, config files)

## Installation 📦

### Prerequisites
- Go 1.21 or higher
- Access to Fortify SSC instance
- Valid SSC API token

### Build from source

```bash
git clone https://github.com/xShuden/ssc-artifact-manager.git
cd ssc-artifact-manager
go mod download
go build -o ssc-approver
```

### Download binary

Download the latest release from the [releases page](https://github.com/xShuden/ssc-artifact-manager/releases).

## Configuration 🔧

The tool supports multiple configuration methods (in order of precedence):

### 1. Command-line flags (highest priority)
```bash
./ssc-approver list --url https://ssc.example.com --token your-token
./ssc-approver list -u https://ssc.example.com -t your-token
```

### 2. Environment variables
```bash
export FORTIFY_SSC_URL=https://ssc.example.com
export FORTIFY_SSC_TOKEN=your-api-token
```

### 3. Configuration file (.env)
Create a `.env` file in the project directory:
```env
FORTIFY_SSC_URL=https://ssc.example.com
FORTIFY_SSC_TOKEN=your-api-token
```

## Usage 🚀

### List artifacts requiring approval

```bash
# List all artifacts requiring approval
./ssc-approver list

# With detailed processing messages
./ssc-approver list -d
./ssc-approver list --details

# Filter by project name
./ssc-approver list -p "finance"
./ssc-approver list --project "backend-service"

# Export to CSV
./ssc-approver list -o csv > pending_approvals.csv

# Combined: detailed view + project filter + custom server
./ssc-approver list -d -p "finance" -u https://ssc.example.com -t your-token
```

### List all projects

```bash
# Using environment variables
./ssc-approver projects

# With custom server
./ssc-approver projects -u https://ssc.example.com -t your-token
```

### List artifacts for a specific project version

```bash
# Using project version ID
./ssc-approver artifacts 123

# With custom server
./ssc-approver artifacts 123 -u https://ssc.example.com -t your-token
```

## Output Examples 📊

### Standard output
```
Found 4 artifacts requiring approval:

+--------------------+---------------------+----------+-----------+----------------+-------------------+
|      PROJECT       |    UPLOAD DATE      | FILE NAME| SIZE (MB) |   UPLOAD IP    |      STATUS       |
+--------------------+---------------------+----------+-----------+----------------+-------------------+
| finance-service    | 2025-08-08 17:05:21 | scan.fpr |    2.10   | 10.185.27.21   | Requires Approval |
| backend-service    | 2025-08-08 16:52:19 | scan.fpr |    1.85   | 10.185.27.21   | Requires Approval |
+--------------------+---------------------+----------+-----------+----------------+-------------------+
```

### Detailed view (with -d flag)
```
Found 2 artifacts requiring approval:

+--------------------+---------------------+----------+-----------+----------------+--------------------------------------------------+
|      PROJECT       |    UPLOAD DATE      | FILE NAME| SIZE (MB) |   UPLOAD IP    |                PROCESSING MESSAGES               |
+--------------------+---------------------+----------+-----------+----------------+--------------------------------------------------+
| finance-service    | 2025-08-08 17:05:21 | scan.fpr |    2.10   | 10.185.27.21   | The amount of executable lines of code in the   |
|                    |                     |          |           |                | new scan is higher by more than 10% (1,193      |
|                    |                     |          |           |                | lines in old scan, 4,876 lines in new scan).    |
|                    |                     |          |           |                | This could be due to major code changes, or     |
|                    |                     |          |           |                | this Analysis Result may be from a different    |
|                    |                     |          |           |                | codebase.                                        |
+--------------------+---------------------+----------+-----------+----------------+--------------------------------------------------+
```

## Creating an SSC API Token 🔑

1. Log in to SSC web interface
2. Navigate to Administration > Users
3. Click on Token Management tab
4. Click "New Token"
5. Provide a name and select required permissions:
   - "View application versions"
   - "View artifacts"
   - "Approve artifacts" (for future approve functionality)

## Project Structure 📁

```
ssc-approver/
├── main.go              # CLI commands and main logic
├── client/
│   └── client.go        # SSC API client implementation
├── models/
│   └── models.go        # Data models for API responses
├── config/
│   └── config.go        # Configuration management
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
├── .env                 # Environment variables (create this)
├── .gitignore           # Git ignore rules
└── README.md            # This file
```

## Contributing 🤝

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Security 🔒

- Never commit API tokens or credentials to the repository
- Keep your `.env` file local and never push it to git
- Use environment variables or secure secret management in production
- The tool accepts self-signed certificates (can be configured)

## License 📄

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author 👨‍💻

Created by [xShuden](https://github.com/xShuden)

## Acknowledgments 🙏

- Built with [Cobra](https://github.com/spf13/cobra) for CLI interface
- [Tablewriter](https://github.com/olekukonko/tablewriter) for formatted output
- [Color](https://github.com/fatih/color) for colored terminal output

## Support 💬

For issues, questions, or suggestions, please [open an issue](https://github.com/xShuden/ssc-artifact-manager/issues) on GitHub.