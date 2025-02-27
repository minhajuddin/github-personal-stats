package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub client with additional functionality
type Client struct {
	client *github.Client
	user   string
	ctx    context.Context
}

// NewClient creates a new GitHub client with the provided token
func NewClient(token, user string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &Client{
		client: client,
		user:   user,
		ctx:    ctx,
	}
}

// TimePeriod represents a time period for stats calculation
type TimePeriod struct {
	Name  string
	Since time.Time
	Until time.Time
}

// GetTimePeriods returns the time periods for stats calculation
func GetTimePeriods() []TimePeriod {
	now := time.Now()
	return []TimePeriod{
		{
			Name:  "Last Week",
			Since: now.AddDate(0, 0, -7),
			Until: now,
		},
		{
			Name:  "Last Month",
			Since: now.AddDate(0, -1, 0),
			Until: now,
		},
		{
			Name:  "Last Year",
			Since: now.AddDate(-1, 0, 0),
			Until: now,
		},
	}
}

// Stats represents GitHub statistics for a user
type Stats struct {
	PRsMerged    int
	LinesAdded   int
	LinesDeleted int
	PRsReviewed  int
	Period       TimePeriod
	Organization string
}

// GetStats retrieves GitHub statistics for a user in an organization
func (c *Client) GetStats(org string, period TimePeriod) (*Stats, error) {
	stats := &Stats{
		Period:       period,
		Organization: org,
	}

	// Get merged PRs
	mergedPRs, err := c.getMergedPRs(org, period)
	if err != nil {
		return nil, fmt.Errorf("error getting merged PRs: %w", err)
	}
	stats.PRsMerged = len(mergedPRs)

	// Get lines added/deleted
	linesAdded, linesDeleted, err := c.getLinesChangedInPRs(mergedPRs)
	if err != nil {
		return nil, fmt.Errorf("error getting lines changed: %w", err)
	}
	stats.LinesAdded = linesAdded
	stats.LinesDeleted = linesDeleted

	// Get PRs reviewed
	reviewedPRs, err := c.getReviewedPRs(org, period)
	if err != nil {
		return nil, fmt.Errorf("error getting reviewed PRs: %w", err)
	}
	stats.PRsReviewed = reviewedPRs

	return stats, nil
}

// getMergedPRs gets the list of PRs merged by the user in the given time period
func (c *Client) getMergedPRs(org string, period TimePeriod) ([]*github.PullRequest, error) {
	// Format the date in the GitHub search format (YYYY-MM-DD)
	sinceStr := period.Since.Format("2006-01-02")
	untilStr := period.Until.Format("2006-01-02")

	// Build the search query
	query := fmt.Sprintf("author:%s org:%s is:pr is:merged merged:%s..%s",
		c.user, org, sinceStr, untilStr)

	var result []*github.PullRequest
	page := 1

	for {
		searchOpts := &github.SearchOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		}

		searchResult, _, err := c.client.Search.Issues(c.ctx, query, searchOpts)
		if err != nil {
			return nil, err
		}

		for _, issue := range searchResult.Issues {
			// Convert issue to PR
			pr, _, err := c.client.PullRequests.Get(c.ctx, org, *issue.Repository.Name, *issue.Number)
			if err != nil {
				continue
			}
			result = append(result, pr)
		}

		if len(searchResult.Issues) < 100 {
			break
		}
		page++
	}

	return result, nil
}

// getLinesChangedInPRs calculates the total lines added and deleted in the given PRs
func (c *Client) getLinesChangedInPRs(prs []*github.PullRequest) (int, int, error) {
	var totalAdded, totalDeleted int

	for _, pr := range prs {
		added := pr.GetAdditions()
		deleted := pr.GetDeletions()

		totalAdded += added
		totalDeleted += deleted
	}

	return totalAdded, totalDeleted, nil
}

// getReviewedPRs gets the count of PRs reviewed by the user in the given time period
func (c *Client) getReviewedPRs(org string, period TimePeriod) (int, error) {
	// Format the date in the GitHub search format (YYYY-MM-DD)
	sinceStr := period.Since.Format("2006-01-02")
	untilStr := period.Until.Format("2006-01-02")

	// Build the search query for PRs reviewed by the user
	query := fmt.Sprintf("org:%s reviewed-by:%s created:%s..%s",
		org, c.user, sinceStr, untilStr)

	var count int
	page := 1

	for {
		searchOpts := &github.SearchOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: 100,
			},
		}

		searchResult, _, err := c.client.Search.Issues(c.ctx, query, searchOpts)
		if err != nil {
			return 0, err
		}

		count += len(searchResult.Issues)

		if len(searchResult.Issues) < 100 {
			break
		}
		page++
	}

	return count, nil
}
