package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/go-github/v73/github"
	"golang.org/x/oauth2"

	"github-profiler/internal/models"
)

// GitHubService handles all GitHub API interactions
type GitHubService struct {
	client *github.Client
	ctx    context.Context
}

// NewGitHubService creates a new GitHub service instance
func NewGitHubService(token string) *GitHubService {
	ctx := context.Background()
	var client *github.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	return &GitHubService{
		client: client,
		ctx:    ctx,
	}
}

// GetUserProfile fetches comprehensive user profile data
func (s *GitHubService) GetUserProfile(username string) (*models.UserProfile, error) {
	// Fetch user basic info
	user, _, err := s.client.Users.Get(s.ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Fetch repositories
	repos, err := s.fetchAllRepositories(username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}

	// Calculate statistics
	languages := s.calculateLanguageStats(repos, username)
	stats := s.calculateProfileStats(repos)
	activity := s.calculateActivityStats(repos)
	ranking := s.calculateRanking(user, repos, stats, activity)

	return &models.UserProfile{
		User:         user,
		Repositories: repos,
		Languages:    languages,
		Stats:        stats,
		Activity:     activity,
		Ranking:      ranking,
	}, nil
}

// fetchAllRepositories gets all repositories for a user
func (s *GitHubService) fetchAllRepositories(username string) ([]*github.Repository, error) {
	var allRepos []*github.Repository

	opts := &github.RepositoryListOptions{
		Type:        "owner",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := s.client.Repositories.List(s.ctx, username, opts)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

// calculateLanguageStats analyzes programming language usage
func (s *GitHubService) calculateLanguageStats(repos []*github.Repository, username string) models.LanguageStats {
	languageBytes := make(map[string]int)
	languageRepos := make(map[string]int)
	totalBytes := 0

	for _, repo := range repos {
		if repo.GetFork() || repo.GetPrivate() {
			continue
		}

		languages, _, err := s.client.Repositories.ListLanguages(s.ctx, username, repo.GetName())
		if err == nil {
			for lang, bytes := range languages {
				languageBytes[lang] += bytes
				languageRepos[lang]++
				totalBytes += bytes
			}
		}

		// Rate limiting protection
		time.Sleep(50 * time.Millisecond)
	}

	// Convert to structured format
	languageStats := make(map[string]models.LanguageInfo)
	for lang, bytes := range languageBytes {
		percentage := 0.0
		if totalBytes > 0 {
			percentage = (float64(bytes) / float64(totalBytes)) * 100
		}

		languageStats[lang] = models.LanguageInfo{
			Name:       lang,
			Bytes:      bytes,
			Percentage: percentage,
			RepoCount:  languageRepos[lang],
		}
	}

	return models.LanguageStats{
		TotalBytes: totalBytes,
		Languages:  languageStats,
	}
}

// calculateProfileStats computes repository statistics
func (s *GitHubService) calculateProfileStats(repos []*github.Repository) models.ProfileStats {
	stats := models.ProfileStats{
		RepoTypes:       make(map[string]int),
		UpdateFrequency: make(map[string]int),
	}

	yearCounts := make(map[int]int)
	totalRepos := 0

	for _, repo := range repos {
		if repo.GetFork() {
			stats.RepoTypes["forks"]++
			continue
		}

		if repo.GetPrivate() {
			stats.RepoTypes["private"]++
		} else {
			stats.RepoTypes["public"]++
		}

		stats.TotalStars += repo.GetStargazersCount()
		stats.TotalForks += repo.GetForksCount()
		stats.TotalSize += int64(repo.GetSize())

		// Timeline analysis
		if repo.CreatedAt != nil {
			year := repo.CreatedAt.Year()
			yearCounts[year]++
		}

		// Update frequency analysis
		if repo.UpdatedAt != nil {
			daysSinceUpdate := int(time.Since(repo.UpdatedAt.Time).Hours() / 24)
			switch {
			case daysSinceUpdate <= 7:
				stats.UpdateFrequency["weekly"]++
			case daysSinceUpdate <= 30:
				stats.UpdateFrequency["monthly"]++
			case daysSinceUpdate <= 90:
				stats.UpdateFrequency["quarterly"]++
			case daysSinceUpdate <= 365:
				stats.UpdateFrequency["yearly"]++
			default:
				stats.UpdateFrequency["stale"]++
			}
		}

		totalRepos++
	}

	// Calculate averages
	if totalRepos > 0 {
		stats.AvgStarsPerRepo = float64(stats.TotalStars) / float64(totalRepos)
	}

	// Convert year counts to timeline
	var years []int
	for year := range yearCounts {
		years = append(years, year)
	}
	sort.Ints(years)

	for _, year := range years {
		stats.CreationTimeline = append(stats.CreationTimeline, models.TimelineEntry{
			Year:  year,
			Count: yearCounts[year],
		})
	}

	return stats
}

// calculateActivityStats computes user activity patterns
func (s *GitHubService) calculateActivityStats(repos []*github.Repository) models.ActivityStats {
	activity := models.ActivityStats{
		CommitFrequency: make(map[string]int),
		ProductiveHours: make(map[string]int),
	}

	// Calculate contribution score based on repository activity
	for _, repo := range repos {
		if !repo.GetFork() {
			activity.ContributionScore += float64(repo.GetStargazersCount()) * 0.5
			activity.ContributionScore += float64(repo.GetForksCount()) * 0.3
			activity.ContributionScore += float64(repo.GetWatchersCount()) * 0.2
		}
	}

	return activity
}

// calculateRanking determines the user's developer ranking
func (s *GitHubService) calculateRanking(user *github.User, repos []*github.Repository, stats models.ProfileStats, activity models.ActivityStats) models.RankingInfo {
	// Social Score (0-25 points)
	socialScore := 0.0
	followers := user.GetFollowers()
	switch {
	case followers >= 10000:
		socialScore = 25.0
	case followers >= 1000:
		socialScore = 20.0
	case followers >= 500:
		socialScore = 15.0
	case followers >= 100:
		socialScore = 10.0
	case followers >= 50:
		socialScore = 7.5
	case followers >= 10:
		socialScore = 5.0
	default:
		socialScore = 2.5
	}

	// Code Score (0-30 points)
	codeScore := 0.0
	publicRepos := stats.RepoTypes["public"]
	if publicRepos > 0 {
		codeScore += float64(publicRepos) * 0.5
		if codeScore > 15 {
			codeScore = 15
		}
	}

	if stats.TotalStars > 0 {
		starScore := float64(stats.TotalStars) * 0.1
		if starScore > 15 {
			starScore = 15
		}
		codeScore += starScore
	}

	// Activity Score (0-25 points)
	activityScore := activity.ContributionScore * 0.01
	if activityScore > 25 {
		activityScore = 25
	}

	// Innovation Score (0-20 points)
	innovationScore := 0.0
	if stats.AvgStarsPerRepo > 0 {
		innovationScore = stats.AvgStarsPerRepo * 0.5
		if innovationScore > 20 {
			innovationScore = 20
		}
	}

	totalScore := socialScore + codeScore + activityScore + innovationScore
	percentile := (totalScore / 100.0) * 100

	rank := models.GetRankByScore(totalScore)

	return models.RankingInfo{
		OverallRank:     rank.Name,
		Badge:           rank.Badge,
		TotalScore:      totalScore,
		Percentile:      percentile,
		SocialScore:     socialScore,
		CodeScore:       codeScore,
		ActivityScore:   activityScore,
		InnovationScore: innovationScore,
	}
}

// GetDemoProfile returns mock data for demo purposes
func (s *GitHubService) GetDemoProfile() (*models.UserProfile, error) {
	return CreateMockProfile(), nil
}
