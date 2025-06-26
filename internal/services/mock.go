package services

import (
	"time"

	"github.com/google/go-github/v73/github"

	"github-profiler/internal/models"
)

// CreateMockProfile creates sample data for demo purposes
func CreateMockProfile() *models.UserProfile {
	// Create mock user
	name := "Demo Developer"
	bio := "Full-stack developer passionate about open source and clean code"
	company := "TechCorp Solutions"
	location := "San Francisco, CA"
	blog := "https://demo-developer.dev"
	twitterUsername := "demo_dev"
	publicRepos := 45
	followers := 1250
	following := 180
	createdAt := time.Date(2018, 3, 15, 0, 0, 0, 0, time.UTC)

	user := &github.User{
		Login:           github.String("demo-user"),
		Name:            &name,
		Bio:             &bio,
		Company:         &company,
		Location:        &location,
		Blog:            &blog,
		TwitterUsername: &twitterUsername,
		PublicRepos:     &publicRepos,
		Followers:       &followers,
		Following:       &following,
		CreatedAt:       &github.Timestamp{Time: createdAt},
	}

	// Create mock repositories
	repos := createMockRepositories()

	// Create mock language stats
	languages := models.LanguageStats{
		TotalBytes: 150000,
		Languages: map[string]models.LanguageInfo{
			"JavaScript": {Name: "JavaScript", Bytes: 46200, Percentage: 30.8, RepoCount: 12},
			"Python":     {Name: "Python", Bytes: 30750, Percentage: 20.5, RepoCount: 8},
			"Go":         {Name: "Go", Bytes: 23100, Percentage: 15.4, RepoCount: 6},
			"TypeScript": {Name: "TypeScript", Bytes: 19200, Percentage: 12.8, RepoCount: 5},
			"Rust":       {Name: "Rust", Bytes: 11550, Percentage: 7.7, RepoCount: 3},
			"Java":       {Name: "Java", Bytes: 7650, Percentage: 5.1, RepoCount: 2},
		},
	}

	// Create mock stats
	stats := models.ProfileStats{
		TotalStars:      947,
		TotalForks:      228,
		TotalSize:       15680, // KB
		AvgStarsPerRepo: 21.0,
		RepoTypes: map[string]int{
			"public":  39,
			"private": 6,
			"forks":   8,
		},
		UpdateFrequency: map[string]int{
			"weekly":    8,
			"monthly":   12,
			"quarterly": 15,
			"yearly":    8,
			"stale":     4,
		},
		CreationTimeline: []models.TimelineEntry{
			{Year: 2018, Count: 3},
			{Year: 2019, Count: 8},
			{Year: 2020, Count: 12},
			{Year: 2021, Count: 15},
			{Year: 2022, Count: 18},
			{Year: 2023, Count: 22},
			{Year: 2024, Count: 25},
		},
	}

	// Create mock activity
	activity := models.ActivityStats{
		ContributionScore: 287.5,
		RecentCommits:     0,
		CommitFrequency: map[string]int{
			"Monday":    15,
			"Tuesday":   18,
			"Wednesday": 22,
			"Thursday":  19,
			"Friday":    16,
			"Saturday":  8,
			"Sunday":    5,
		},
		ProductiveHours: map[string]int{
			"9":  12,
			"10": 18,
			"11": 15,
			"14": 20,
			"15": 22,
			"16": 16,
			"20": 8,
			"21": 6,
		},
	}

	// Create mock ranking
	ranking := models.RankingInfo{
		OverallRank:     "Experienced Developer",
		Badge:           "EXPERIENCED",
		TotalScore:      78.5,
		Percentile:      78.5,
		SocialScore:     18.0,
		CodeScore:       24.5,
		ActivityScore:   19.0,
		InnovationScore: 17.0,
	}

	return &models.UserProfile{
		User:         user,
		Repositories: repos,
		Languages:    languages,
		Stats:        stats,
		Activity:     activity,
		Ranking:      ranking,
	}
}

func createMockRepositories() []*github.Repository {
	repos := []*github.Repository{}

	// Sample repository data
	repoData := []struct {
		name        string
		description string
		language    string
		stars       int
		forks       int
		size        int
	}{
		{"awesome-web-app", "A modern web application built with React and Node.js", "JavaScript", 324, 89, 4264},
		{"ml-algorithms", "Collection of machine learning algorithms implemented in Python", "Python", 189, 45, 2914},
		{"go-microservice", "High-performance microservice written in Go", "Go", 156, 32, 2584},
		{"rust-cli-tool", "Fast command-line tool built with Rust", "Rust", 98, 21, 2004},
		{"typescript-lib", "Type-safe utility library for TypeScript projects", "TypeScript", 67, 15, 1694},
		{"data-visualization", "Interactive data visualization dashboard", "JavaScript", 45, 12, 1474},
		{"api-gateway", "Scalable API gateway service", "Go", 38, 8, 1404},
		{"mobile-app", "Cross-platform mobile app built with Flutter", "Dart", 29, 6, 1314},
	}

	for _, data := range repoData {
		name := data.name
		desc := data.description
		lang := data.language
		stars := data.stars
		forks := data.forks
		size := data.size
		isPrivate := false
		isFork := false
		updatedAt := time.Now().AddDate(0, -1, 0) // 1 month ago

		repo := &github.Repository{
			Name:            &name,
			Description:     &desc,
			Language:        &lang,
			StargazersCount: &stars,
			ForksCount:      &forks,
			Size:            &size,
			Private:         &isPrivate,
			Fork:            &isFork,
			UpdatedAt:       &github.Timestamp{Time: updatedAt},
			CreatedAt:       &github.Timestamp{Time: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
		}

		repos = append(repos, repo)
	}

	return repos
}
