package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/minhajuddin/github-personal-stats/pkg/github"
)

func main() {
	// Parse command line flags
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub API token")
	user := flag.String("user", "", "GitHub username")
	org := flag.String("org", "", "GitHub organization")
	flag.Parse()

	// Validate required flags
	if *token == "" {
		fmt.Println("Error: GitHub token is required. Set GITHUB_TOKEN environment variable or use -token flag.")
		os.Exit(1)
	}

	if *user == "" {
		fmt.Println("Error: GitHub username is required. Use -user flag.")
		os.Exit(1)
	}

	if *org == "" {
		fmt.Println("Error: GitHub organization is required. Use -org flag.")
		os.Exit(1)
	}

	// Create GitHub client
	client := github.NewClient(*token, *user)

	// Get time periods
	periods := github.GetTimePeriods()

	// Create tabwriter for formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	fmt.Fprintln(w, "Period\tPRs Merged\tLines Added\tLines Deleted\tPRs Reviewed\t")
	fmt.Fprintln(w, "------\t---------\t----------\t------------\t-----------\t")

	// Get stats for each time period
	for _, period := range periods {
		stats, err := client.GetStats(*org, period)
		if err != nil {
			fmt.Printf("Error getting stats for %s: %v\n", period.Name, err)
			continue
		}

		// Print stats
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t\n",
			stats.Period.Name,
			stats.PRsMerged,
			stats.LinesAdded,
			stats.LinesDeleted,
			stats.PRsReviewed,
		)
	}

	// Print summary
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Organization: %s\n", *org)
	fmt.Fprintf(w, "User: %s\n", *user)
}
