package cmd

import (
	"fmt"
	"time"

	"github.com/google/go-github/v73/github"
)

func CreateMockUserData() *UserData {
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
		Login:             github.String("demo-user"),
		Name:              &name,
		Bio:               &bio,
		Company:           &company,
		Location:          &location,
		Blog:              &blog,
		TwitterUsername:   &twitterUsername,
		PublicRepos:       &publicRepos,
		Followers:         &followers,
		Following:         &following,
		CreatedAt:         &github.Timestamp{Time: createdAt},
	}

	// Create mock repositories
	repos := []*github.Repository{
		createMockRepo("awesome-web-app", "A modern web application built with React and Node.js", 324, 89, "JavaScript", false, false),
		createMockRepo("ml-algorithms", "Implementation of popular machine learning algorithms in Python", 189, 45, "Python", false, false),
		createMockRepo("go-microservice", "High-performance microservice written in Go", 156, 32, "Go", false, false),
		createMockRepo("rust-cli-tool", "Fast command-line utility built with Rust", 98, 21, "Rust", false, false),
		createMockRepo("typescript-lib", "Useful TypeScript library for web developers", 67, 15, "TypeScript", false, false),
		createMockRepo("data-visualization", "Beautiful charts and graphs using D3.js", 45, 12, "JavaScript", false, false),
		createMockRepo("api-gateway", "Scalable API gateway with rate limiting", 38, 8, "Go", false, false),
		createMockRepo("mobile-app", "Cross-platform mobile app built with Flutter", 29, 6, "Dart", false, false),
	}

	// Create mock languages
	languages := map[string]int{
		"JavaScript": 12,
		"Python":     8,
		"Go":         6,
		"TypeScript": 5,
		"Rust":       3,
		"Dart":       2,
		"Java":       2,
		"C++":        1,
	}

	// Create mock stats
	stats := &UserStats{
		TotalStars:      947,
		TotalForks:      228,
		TotalSize:       15678,
		AvgStarsPerRepo: 21.0,
		RepoTypes: map[string]int{
			"public":  39,
			"private": 6,
			"forks":   8,
		},
		CreationTimeline: []RepoTimelineEntry{
			{2018, 3},
			{2019, 8},
			{2020, 12},
			{2021, 15},
			{2022, 18},
			{2023, 22},
			{2024, 25},
		},
		UpdateFrequency: map[string]int{
			"weekly":    8,
			"monthly":   12,
			"quarterly": 15,
			"yearly":    8,
			"stale":     4,
		},
	}

	// Create mock activity
	activity := &ActivityStats{
		CommitFrequency: map[string]int{
			"Monday":    15,
			"Tuesday":   18,
			"Wednesday": 22,
			"Thursday":  19,
			"Friday":    16,
			"Saturday":  8,
			"Sunday":    5,
		},
		ProductiveHours: map[int]int{
			9:  12,
			10: 18,
			11: 15,
			14: 20,
			15: 22,
			16: 16,
			20: 8,
			21: 6,
		},
		ContributionScore: 287.5,
	}

	// Create mock ranking
	ranking := &RankingData{
		OverallScore:  78.5,
		SocialScore:   18.0,
		CodeScore:     24.5,
		ActivityScore: 19.0,
		QualityScore:  17.0,
		Rank:         "💪 Experienced Developer",
		Badge:        "EXPERIENCED",
		Percentile:   78.5,
	}

	return &UserData{
		User:         user,
		Repositories: repos,
		Languages:    languages,
		Stats:        stats,
		Activity:     activity,
		Ranking:      ranking,
	}
}

func createMockRepo(name, desc string, stars, forks int, language string, private, fork bool) *github.Repository {
	description := desc
	starCount := stars
	forkCount := forks
	lang := language
	isPrivate := private
	isFork := fork
	size := 1024 + (stars * 10) // Simulate size based on popularity
	createdAt := time.Now().AddDate(-2, 0, -(stars/10)) // Simulate creation date
	updatedAt := time.Now().AddDate(0, 0, -(30 - stars/20)) // Simulate last update

	return &github.Repository{
		Name:            &name,
		Description:     &description,
		StargazersCount: &starCount,
		ForksCount:      &forkCount,
		Language:        &lang,
		Private:         &isPrivate,
		Fork:            &isFork,
		Size:            &size,
		CreatedAt:       &github.Timestamp{Time: createdAt},
		UpdatedAt:       &github.Timestamp{Time: updatedAt},
	}
}

func DemoWithMockData(format string) {
	fmt.Println("🎭 GitHub Profiler Demo")
	fmt.Println("=======================")
	fmt.Printf("   %s • %s\n", version, author)
	fmt.Println()

	mockData := CreateMockUserData()

	switch format {
	case "json":
		fmt.Println("📋 JSON Output:")
		fmt.Println("===============")
		displayJSON(mockData)
	case "html":
		fmt.Println("🌐 HTML Output:")
		fmt.Println("===============")
		displayHTML(mockData)
		fmt.Println("HTML file generated successfully!")
	default:
		fmt.Println("🖥️  Terminal Output:")
		fmt.Println("=====================")
		displayTerminalRich(mockData)
	}

	fmt.Println()
	fmt.Println("💡 Try with a real username: github-profiler <username>")
	fmt.Println("   For higher rate limits, set GITHUB_TOKEN environment variable")
}
