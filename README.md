# SSC Approver CLI 🔍

A command-line tool to manage and list artifacts requiring approval in Fortify Software Security Center (SSC).

## Features ✨

- 📋 List all artifacts with "Requires Approval" status across all projects
- 🔎 View detailed processing messages (code line changes, security warnings)
- 🏷️ Filter artifacts by project name
- 📊 Multiple output formats:
  - **Table**: Human-readable ASCII tables
  - **CSV**: For Excel/spreadsheet import
  - **JSON**: For automation and API integration
- 🗂️ List all projects and their versions
- 📦 View artifacts for specific project versions
- 🔐 Flexible authentication options:
  - Command-line flags
  - Environment variables
  - Configuration files (.env)
- ⚡ Fast parallel scanning of multiple projects
- 🔒 Supports self-signed certificates

## Installation 📦

### Prerequisites
- Go 1.21 or higher
- Access to Fortify SSC instance
- Valid SSC API token

### Build from source

```bash
git clone https://github.com/xShuden/fortify-ssc-artifact-manager.git
cd fortify-ssc-artifact-manager
go mod download
go build -o ssc-approver
```

### Download binary

Download the latest release from the [releases page](https://github.com/xShuden/fortify-ssc-artifact-manager/releases).

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

# Export to JSON
./ssc-approver list -o json > pending_approvals.json

# JSON with jq processing
./ssc-approver list -o json | jq '.[].project'
./ssc-approver list -o json | jq '.[] | select(.file_size_mb | tonumber > 2)'

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

## Output Formats 📊

### Available Formats

The tool supports three output formats:

1. **Table Format** (default) - Human-readable ASCII table
2. **CSV Format** - For spreadsheet applications and data processing
3. **JSON Format** - For programmatic processing and automation

### Format Examples

#### Table Format (Default)
```
Found 4 artifacts requiring approval:

+--------------------+---------------------+----------+-----------+----------------+-------------------+
|      PROJECT       |    UPLOAD DATE      | FILE NAME| SIZE (MB) |   UPLOAD IP    |      STATUS       |
+--------------------+---------------------+----------+-----------+----------------+-------------------+
| finance-service    | 2025-08-08 17:05:21 | scan.fpr |    2.10   | 10.185.27.21   | Requires Approval |
| backend-service    | 2025-08-08 16:52:19 | scan.fpr |    1.85   | 10.185.27.21   | Requires Approval |
+--------------------+---------------------+----------+-----------+----------------+-------------------+
```

#### CSV Format
```csv
Project,Upload Date,File Name,Size,Upload IP,Messages
finance-service - dev,2025-08-08 17:05:21,scan.fpr,2202009,10.185.27.21,"Code lines increased by 10%"
backend-service - dev,2025-08-08 16:52:19,analysis.fpr,1939865,10.185.27.21,"New findings detected"
```

#### JSON Format
```json
[
  {
    "project": "finance-service - dev",
    "upload_date": "2025-08-08 17:05:21",
    "file_name": "scan.fpr",
    "file_size_bytes": 2202009,
    "file_size_mb": "2.10",
    "upload_ip": "10.185.27.21",
    "status": "Requires Approval",
    "messages": "Code lines increased by 10%"
  },
  {
    "project": "backend-service - dev",
    "upload_date": "2025-08-08 16:52:19",
    "file_name": "analysis.fpr",
    "file_size_bytes": 1939865,
    "file_size_mb": "1.85",
    "upload_ip": "10.185.27.21",
    "status": "Requires Approval",
    "messages": "New findings detected"
  }
]
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

## Advanced Usage Examples 🚀

### JSON Processing with jq

```bash
# Get all project names
./ssc-approver list -o json | jq -r '.[].project' | sort -u

# Filter artifacts larger than 2MB
./ssc-approver list -o json | jq '.[] | select(.file_size_mb | tonumber > 2)'

# Count artifacts per project
./ssc-approver list -o json | jq 'group_by(.project) | .[] | {project: .[0].project, count: length}'

# Export specific fields to CSV
./ssc-approver list -o json | jq -r '.[] | [.project, .file_name, .upload_date] | @csv'
```

### Automation Scripts

```bash
#!/bin/bash
# Daily approval check script

ARTIFACTS=$(./ssc-approver list -o json)
COUNT=$(echo "$ARTIFACTS" | jq '. | length')

if [ "$COUNT" -gt 0 ]; then
    echo "Found $COUNT artifacts requiring approval"
    # Send notification or create ticket
    echo "$ARTIFACTS" | jq -r '.[] | "\(.project): \(.file_name) - \(.upload_date)"'
fi
```

### Integration with CI/CD

```yaml
# GitHub Actions example
- name: Check SSC Approvals
  run: |
    ./ssc-approver list -u ${{ secrets.SSC_URL }} -t ${{ secrets.SSC_TOKEN }} -o json > artifacts.json
    if [ $(jq '. | length' artifacts.json) -gt 0 ]; then
      echo "::warning::Artifacts pending approval in SSC"
      jq -r '.[] | "- \(.project): \(.file_name)"' artifacts.json
    fi
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

## Roadmap 🗺️

- [ ] Add `approve` command to approve artifacts directly from CLI
- [ ] Add `reject` command with comment support
- [ ] Support for bulk operations
- [ ] Add filtering by date range
- [ ] Export detailed reports in PDF format
- [ ] Slack/Teams notification integration
- [ ] Configuration profiles for multiple SSC instances
- [ ] Interactive mode with artifact selection

## Support 💬

For issues, questions, or suggestions, please [open an issue](https://github.com/xShuden/fortify-ssc-artifact-manager/issues) on GitHub.

## Star History ⭐

If you find this tool useful, please consider giving it a star on GitHub!

[![Star History Chart](https://api.star-history.com/svg?repos=xShuden/fortify-ssc-artifact-manager&type=Date)](https://star-history.com/#xShuden/fortify-ssc-artifact-manager&Date)
