# GitHub Personal Stats

A command-line tool to fetch your personal GitHub statistics for an organization.

## Features

This tool provides the following statistics for a GitHub user in an organization:

1. Number of PRs merged (not just closed)
2. Number of lines added and deleted
3. Number of PRs reviewed

Statistics are provided for three time periods:
- Last Week
- Last Month
- Last Year

## Installation

```bash
# Clone the repository
git clone https://github.com/minhajuddin/github-personal-stats.git
cd github-personal-stats

# Build the application
go build -o github-stats ./cmd/github-stats
```

## Usage

```bash
# Using environment variable for GitHub token
export GITHUB_TOKEN=your_github_personal_access_token
./github-stats -user your_github_username -org your_organization

# Or providing the token directly
./github-stats -token your_github_personal_access_token -user your_github_username -org your_organization
```

### Required Parameters

- `-token`: Your GitHub personal access token (can also be set via GITHUB_TOKEN environment variable)
- `-user`: Your GitHub username
- `-org`: The GitHub organization to fetch stats for

## Creating a GitHub Personal Access Token

1. Go to GitHub Settings > Developer settings > Personal access tokens > Tokens (classic)
2. Click "Generate new token" and select "Generate new token (classic)"
3. Give your token a descriptive name
4. Select the following scopes:
   - `repo` (Full control of private repositories)
   - `read:org` (Read organization membership)
5. Click "Generate token"
6. Copy the token and use it with this application

## Example Output

```
Period      PRs Merged  Lines Added  Lines Deleted  PRs Reviewed  
------      ---------   ----------   ------------   -----------   
Last Week   5           1250         450            12            
Last Month  18          4320         1840           45            
Last Year   156         38450        15230          320           

Organization: example-org
User: example-user
```

## License

MIT 