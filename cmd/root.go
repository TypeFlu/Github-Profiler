package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v56/github"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	githubToken  string
	outputFormat string
	version      = "1.0.0"
	author       = "github@Tyeflu"
)

var rootCmd = &cobra.Command{
	Use:   "github-profiler [username]",
	Short: "A beautiful CLI tool to analyze GitHub user profiles",
	Long: `GitHub Profiler is a CLI tool that provides comprehensive analysis
of any GitHub user's public profile with beautiful visualizations,
advanced statistics, and ranking information.

Use 'github-profiler demo' to see the tool in action with sample data.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runProfiler,
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a demo with sample data",
	Long:  `Demonstrates the GitHub Profiler output using mock data.`,
	Run: func(cmd *cobra.Command, args []string) {
		DemoWithMockData(outputFormat)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&githubToken, "token", "t", "", "GitHub personal access token (optional for public data)")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "terminal", "Output format: terminal, json, html")
	
	rootCmd.AddCommand(demoCmd)
	demoCmd.Flags().StringVarP(&outputFormat, "format", "f", "terminal", "Output format: terminal, json, html")
	
	if githubToken == "" {
		githubToken = os.Getenv("GITHUB_TOKEN")
	}
}

func runProfiler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}
	
	username := args[0]
	var client *github.Client
	
	if githubToken != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}
	
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Fetching GitHub data..."
	s.Start()
	
	userData, err := fetchComprehensiveUserData(client, username)
	s.Stop()
	
	if err != nil {
		fmt.Printf("❌ Error fetching data for '%s': %v\n", username, err)
		fmt.Println("\n💡 Tips:")
		fmt.Println("   • Make sure the username is correct")
		fmt.Println("   • Check your internet connection") 
		fmt.Println("   • Set GITHUB_TOKEN for higher rate limits")
		fmt.Println("   • Try 'github-profiler demo' to see example output")
		os.Exit(1)
	}
	
	switch outputFormat {
	case "json":
		displayJSON(userData)
	case "html":
		displayHTML(userData)
	default:
		displayTerminalRich(userData)
	}
}

type UserData struct {
	User         *github.User         `json:"user"`
	Repositories []*github.Repository `json:"repositories"`
	Languages    map[string]int       `json:"languages"`
	Stats        *UserStats           `json:"stats"`
	Activity     *ActivityStats       `json:"activity"`
	Ranking      *RankingData         `json:"ranking"`
}

type UserStats struct {
	TotalStars        int                    `json:"total_stars"`
	TotalForks        int                    `json:"total_forks"`
	TotalSize         int64                  `json:"total_size_kb"`
	AvgStarsPerRepo   float64                `json:"avg_stars_per_repo"`
	MostUsedLanguages []LanguageStats        `json:"most_used_languages"`
	RepoTypes         map[string]int         `json:"repo_types"`
	CreationTimeline  []RepoTimelineEntry    `json:"creation_timeline"`
	UpdateFrequency   map[string]int         `json:"update_frequency"`
}

type ActivityStats struct {
	RecentCommits     []*github.RepositoryCommit `json:"recent_commits,omitempty"`
	CommitFrequency   map[string]int             `json:"commit_frequency"`
	ActiveDays        []string                   `json:"active_days"`
	ProductiveHours   map[int]int                `json:"productive_hours"`
	ContributionScore float64                    `json:"contribution_score"`
}

type LanguageStats struct {
	Name       string  `json:"name"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
	Bytes      int     `json:"bytes"`
}

type RepoTimelineEntry struct {
	Year  int `json:"year"`
	Count int `json:"count"`
}

type RankingData struct {
	OverallScore    float64 `json:"overall_score"`
	SocialScore     float64 `json:"social_score"`
	CodeScore       float64 `json:"code_score"`
	ActivityScore   float64 `json:"activity_score"`
	QualityScore    float64 `json:"quality_score"`
	Rank            string  `json:"rank"`
	Badge           string  `json:"badge"`
	Percentile      float64 `json:"percentile"`
}

func fetchComprehensiveUserData(client *github.Client, username string) (*UserData, error) {
	ctx := context.Background()
	
	user, _, err := client.Users.Get(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	
	repoOpts := &github.RepositoryListOptions{
		Type:        "owner",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, username, repoOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch repositories: %v", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		repoOpts.Page = resp.NextPage
	}
	
	stats := calculateComprehensiveStats(allRepos)
	
	// Get language statistics for all repos
	languages := make(map[string]int)
	languageBytes := make(map[string]int)
	
	for _, repo := range allRepos {
		if repo.GetFork() || repo.GetPrivate() {
			continue
		}
		
		repoLangs, _, err := client.Repositories.ListLanguages(ctx, username, repo.GetName())
		if err == nil {
			for lang, bytes := range repoLangs {
				languages[lang]++
				languageBytes[lang] += bytes
			}
		}
		
		time.Sleep(50 * time.Millisecond)
	}
	
	activity := calculateActivityStats(client, username, allRepos)
	ranking := calculateAdvancedRanking(user, allRepos, stats, activity)
	
	return &UserData{
		User:         user,
		Repositories: allRepos,
		Languages:    languages,
		Stats:        stats,
		Activity:     activity,
		Ranking:      ranking,
	}, nil
}

func calculateComprehensiveStats(repos []*github.Repository) *UserStats {
	stats := &UserStats{
		RepoTypes:        make(map[string]int),
		CreationTimeline: make([]RepoTimelineEntry, 0),
		UpdateFrequency:  make(map[string]int),
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
	
	if totalRepos > 0 {
		stats.AvgStarsPerRepo = float64(stats.TotalStars) / float64(totalRepos)
	}
	
	var years []int
	for year := range yearCounts {
		years = append(years, year)
	}
	sort.Ints(years)
	
	for _, year := range years {
		stats.CreationTimeline = append(stats.CreationTimeline, RepoTimelineEntry{
			Year:  year,
			Count: yearCounts[year],
		})
	}
	
	return stats
}

func calculateActivityStats(client *github.Client, username string, repos []*github.Repository) *ActivityStats {
	ctx := context.Background()
	activity := &ActivityStats{
		CommitFrequency: make(map[string]int),
		ActiveDays:      make([]string, 0),
		ProductiveHours: make(map[int]int),
	}
	
	// Get recent commits from top repositories
	repoCount := 0
	for _, repo := range repos {
		if repoCount >= 5 || repo.GetFork() || repo.GetPrivate() {
			continue
		}
		
		// Get recent commits
		commitOpts := &github.CommitsListOptions{
			Author: username,
			Since:  time.Now().AddDate(0, -3, 0), // Last 3 months
			ListOptions: github.ListOptions{PerPage: 10},
		}
		
		commits, _, err := client.Repositories.ListCommits(ctx, username, repo.GetName(), commitOpts)
		if err == nil {
			activity.RecentCommits = append(activity.RecentCommits, commits...)
			
			// Analyze commit patterns
			for _, commit := range commits {
				if commit.Commit != nil && commit.Commit.Author != nil && commit.Commit.Author.Date != nil {
					date := commit.Commit.Author.Date.Time
					hour := date.Hour()
					activity.ProductiveHours[hour]++
					
					dayOfWeek := date.Weekday().String()
					activity.CommitFrequency[dayOfWeek]++
				}
			}
		}
		
		repoCount++
		time.Sleep(100 * time.Millisecond) // Rate limiting
	}
	
	// Calculate contribution score based on activity
	totalCommits := len(activity.RecentCommits)
	activity.ContributionScore = float64(totalCommits) * 2.5
	
	return activity
}

func calculateAdvancedRanking(user *github.User, repos []*github.Repository, stats *UserStats, activity *ActivityStats) *RankingData {
	// Social influence score (0-25 points)
	followers := user.GetFollowers()
	socialScore := calculateSocialScore(followers)
	
	// Code quality score (0-30 points)  
	codeScore := calculateCodeQualityScore(stats, repos)
	
	// Activity score (0-25 points)
	activityScore := calculateActivityScore(activity, stats)
	
	// Innovation score (0-20 points)
	qualityScore := calculateQualityScore(user, repos)
	
	totalScore := socialScore + codeScore + activityScore + qualityScore
	
	rank, badge := getDeveloperRank(totalScore)
	percentile := calculatePercentile(totalScore)
	
	return &RankingData{
		OverallScore:  totalScore,
		SocialScore:   socialScore,
		CodeScore:     codeScore,
		ActivityScore: activityScore,
		QualityScore:  qualityScore,
		Rank:         rank,
		Badge:        badge,
		Percentile:   percentile,
	}
}

func displayTerminalRich(data *UserData) {
	user := data.User
	
	// Define enhanced styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 2).
		Margin(1, 0)
	
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		Margin(1, 0)
	
	statStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true)
	
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#626262")).
		Width(20)
	
	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)
	
	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")).
		Bold(true)
	
	fmt.Println(titleStyle.Render("🚀 GitHub Profile Analysis"))
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("   v" + version + " • " + author))
	
	userInfo := fmt.Sprintf(`
👤 %s (%s)
📝 %s
🏢 %s
📍 %s
🌐 %s
🐦 %s
📅 Joined: %s
📊 Public Repos: %d | Followers: %d | Following: %d
`,
		getStringValue(user.Name, user.GetLogin()),
		user.GetLogin(),
		getStringValue(user.Bio, "No bio available"),
		getStringValue(user.Company, "Not specified"),
		getStringValue(user.Location, "Unknown"),
		getStringValue(user.Blog, "Not specified"),
		getStringValue(user.TwitterUsername, "Not specified"),
		user.CreatedAt.Format("January 2006"),
		user.GetPublicRepos(),
		user.GetFollowers(),
		user.GetFollowing(),
	)
	
	fmt.Println(headerStyle.Render(userInfo))
	
	// Repository Statistics
	fmt.Println(titleStyle.Render("📊 Repository Analytics"))
	
	repoStats := [][]string{
		{"⭐ Total Stars", fmt.Sprintf("%d", data.Stats.TotalStars)},
		{"🍴 Total Forks", fmt.Sprintf("%d", data.Stats.TotalForks)},
		{"📦 Repository Size", fmt.Sprintf("%.1f MB", float64(data.Stats.TotalSize)/1024)},
		{"📈 Avg Stars/Repo", fmt.Sprintf("%.1f", data.Stats.AvgStarsPerRepo)},
		{"🔒 Private Repos", fmt.Sprintf("%d", data.Stats.RepoTypes["private"])},
		{"🌍 Public Repos", fmt.Sprintf("%d", data.Stats.RepoTypes["public"])},
		{"� Forked Repos", fmt.Sprintf("%d", data.Stats.RepoTypes["forks"])},
	}
	
	for _, stat := range repoStats {
		fmt.Printf("%s %s\n", labelStyle.Render(stat[0]+":"), statStyle.Render(stat[1]))
	}
	
	// Developer Ranking
	fmt.Println()
	fmt.Println(titleStyle.Render("� Developer Ranking System"))
	
	ranking := data.Ranking
	fmt.Printf("🎖️  %s %s\n", labelStyle.Render("Overall Rank:"), 
		successStyle.Render(ranking.Rank))
	fmt.Printf("🎯 %s %s\n", labelStyle.Render("Badge:"), 
		warningStyle.Render(ranking.Badge))
	fmt.Printf("📊 %s %s\n", labelStyle.Render("Total Score:"), 
		statStyle.Render(fmt.Sprintf("%.1f/100", ranking.OverallScore)))
	fmt.Printf("� %s %s\n", labelStyle.Render("Percentile:"), 
		statStyle.Render(fmt.Sprintf("%.1f%%", ranking.Percentile)))
	
	fmt.Println()
	
	// Score breakdown
	scoreTable := tablewriter.NewWriter(os.Stdout)
	scoreTable.SetHeader([]string{"Category", "Score", "Max", "Percentage"})
	scoreTable.SetBorder(true)
	scoreTable.SetAlignment(tablewriter.ALIGN_CENTER)
	
	scores := [][]string{
		{"Social Influence", fmt.Sprintf("%.1f", ranking.SocialScore), "25", fmt.Sprintf("%.1f%%", ranking.SocialScore/25*100)},
		{"Code Quality", fmt.Sprintf("%.1f", ranking.CodeScore), "30", fmt.Sprintf("%.1f%%", ranking.CodeScore/30*100)},
		{"Activity Level", fmt.Sprintf("%.1f", ranking.ActivityScore), "25", fmt.Sprintf("%.1f%%", ranking.ActivityScore/25*100)},
		{"Innovation", fmt.Sprintf("%.1f", ranking.QualityScore), "20", fmt.Sprintf("%.1f%%", ranking.QualityScore/20*100)},
	}
	
	for _, score := range scores {
		scoreTable.Append(score)
	}
	scoreTable.Render()
	
	// Language Analysis
	fmt.Println()
	fmt.Println(titleStyle.Render("💻 Language Proficiency"))
	
	displayAdvancedLanguageBreakdown(data.Languages)
	
	// Repository Timeline
	fmt.Println()
	fmt.Println(titleStyle.Render("📅 Repository Creation Timeline"))
	
	displayRepositoryTimeline(data.Stats.CreationTimeline)
	
	// Activity Analysis
	fmt.Println()
	fmt.Println(titleStyle.Render("🎯 Activity Patterns"))
	
	displayActivityAnalysis(data.Activity)
	
	// Top Repositories with enhanced details
	fmt.Println()
	fmt.Println(titleStyle.Render("⭐ Top Repositories (Detailed)"))
	
	displayTopRepositoriesEnhanced(data.Repositories)
	
	// Repository Health Analysis
	fmt.Println()
	fmt.Println(titleStyle.Render("� Repository Health"))
	
	displayRepositoryHealth(data.Stats)
}

// Legacy functions for backwards compatibility - now simplified
func displayTopRepositories(repos []*github.Repository) {
	displayTopRepositoriesEnhanced(repos)
}

func displayLanguageBreakdown(languages map[string]int) {
	displayAdvancedLanguageBreakdown(languages)
}

func displayRanking(data *UserData) {
	ranking := data.Ranking
	fmt.Printf("🏆 Developer Rank: %s\n", ranking.Rank)
	fmt.Printf("📊 Total Score: %.0f/100\n", ranking.OverallScore)
}

// Utility functions
func getDescription(desc *string) string {
	if desc != nil && *desc != "" {
		if len(*desc) > 60 {
			return (*desc)[:57] + "..."
		}
		return *desc
	}
	return "No description"
}

// Helper functions for enhanced display
func getStringValue(value *string, fallback string) string {
	if value != nil && *value != "" {
		return *value
	}
	return fallback
}

func calculateSocialScore(followers int) float64 {
	switch {
	case followers >= 50000:
		return 25.0
	case followers >= 10000:
		return 22.0
	case followers >= 5000:
		return 20.0
	case followers >= 1000:
		return 18.0
	case followers >= 500:
		return 15.0
	case followers >= 100:
		return 12.0
	case followers >= 50:
		return 8.0
	case followers >= 10:
		return 5.0
	default:
		return float64(followers) * 0.5
	}
}

func calculateCodeQualityScore(stats *UserStats, repos []*github.Repository) float64 {
	score := 0.0
	
	// Repository count contribution (max 10 points)
	repoCount := stats.RepoTypes["public"]
	if repoCount > 50 {
		score += 10.0
	} else {
		score += float64(repoCount) * 0.2
	}
	
	// Star rating contribution (max 15 points)
	if stats.TotalStars >= 10000 {
		score += 15.0
	} else if stats.TotalStars >= 1000 {
		score += 12.0
	} else if stats.TotalStars >= 100 {
		score += 8.0
	} else {
		score += float64(stats.TotalStars) * 0.08
	}
	
	// Average stars per repo (max 5 points)
	if stats.AvgStarsPerRepo >= 50 {
		score += 5.0
	} else {
		score += stats.AvgStarsPerRepo * 0.1
	}
	
	if score > 30 {
		return 30
	}
	return score
}

func calculateActivityScore(activity *ActivityStats, stats *UserStats) float64 {
	score := 0.0
	
	// Recent commit activity (max 15 points)
	recentCommits := len(activity.RecentCommits)
	if recentCommits >= 100 {
		score += 15.0
	} else {
		score += float64(recentCommits) * 0.15
	}
	
	// Repository freshness (max 10 points)
	freshRepos := stats.UpdateFrequency["weekly"] + stats.UpdateFrequency["monthly"]
	totalRepos := stats.RepoTypes["public"] + stats.RepoTypes["private"]
	if totalRepos > 0 {
		freshnessRatio := float64(freshRepos) / float64(totalRepos)
		score += freshnessRatio * 10
	}
	
	if score > 25 {
		return 25
	}
	return score
}

func calculateQualityScore(user *github.User, repos []*github.Repository) float64 {
	score := 0.0
	
	// Account age bonus (max 5 points)
	if user.CreatedAt != nil {
		years := time.Since(user.CreatedAt.Time).Hours() / (24 * 365)
		if years >= 5 {
			score += 5.0
		} else {
			score += years
		}
	}
	
	// Language diversity (max 8 points)
	languages := make(map[string]bool)
	for _, repo := range repos {
		if repo.GetLanguage() != "" && !repo.GetFork() {
			languages[repo.GetLanguage()] = true
		}
	}
	langCount := len(languages)
	if langCount >= 8 {
		score += 8.0
	} else {
		score += float64(langCount)
	}
	
	// Bio and profile completeness (max 7 points)
	if user.Bio != nil && *user.Bio != "" {
		score += 2.0
	}
	if user.Company != nil && *user.Company != "" {
		score += 1.0
	}
	if user.Location != nil && *user.Location != "" {
		score += 1.0
	}
	if user.Blog != nil && *user.Blog != "" {
		score += 2.0
	}
	if user.Email != nil && *user.Email != "" {
		score += 1.0
	}
	
	if score > 20 {
		return 20
	}
	return score
}

func getDeveloperRank(score float64) (string, string) {
	switch {
	case score >= 95:
		return "🌟 Legendary Developer", "LEGENDARY"
	case score >= 90:
		return "🚀 Elite Developer", "ELITE"
	case score >= 80:
		return "💎 Senior Developer", "SENIOR"
	case score >= 70:
		return "💪 Experienced Developer", "EXPERIENCED"
	case score >= 60:
		return "⚡ Active Developer", "ACTIVE"
	case score >= 50:
		return "🌱 Growing Developer", "GROWING"
	case score >= 40:
		return "📚 Learning Developer", "LEARNING"
	case score >= 30:
		return "👶 Junior Developer", "JUNIOR"
	case score >= 20:
		return "🌟 Aspiring Developer", "ASPIRING"
	default:
		return "🥚 Beginner", "BEGINNER"
	}
}

func calculatePercentile(score float64) float64 {
	// Simulated percentile calculation based on score distribution
	switch {
	case score >= 95:
		return 99.5
	case score >= 90:
		return 95.0
	case score >= 80:
		return 85.0
	case score >= 70:
		return 70.0
	case score >= 60:
		return 55.0
	case score >= 50:
		return 40.0
	case score >= 40:
		return 25.0
	case score >= 30:
		return 15.0
	default:
		return score * 0.5
	}
}

func displayAdvancedLanguageBreakdown(languages map[string]int) {
	if len(languages) == 0 {
		fmt.Println("No language data available.")
		return
	}
	
	// Convert to sorted slice
	type langStat struct {
		name  string
		count int
	}
	
	var langStats []langStat
	total := 0
	for lang, count := range languages {
		langStats = append(langStats, langStat{lang, count})
		total += count
	}
	
	// Sort by count
	sort.Slice(langStats, func(i, j int) bool {
		return langStats[i].count > langStats[j].count
	})
	
	// Display with enhanced formatting
	for i, ls := range langStats {
		if i >= 10 { // Show top 10
			break
		}
		
		percentage := float64(ls.count) / float64(total) * 100
		barWidth := int(percentage / 2)
		if barWidth < 1 && percentage > 0 {
			barWidth = 1
		}
		
		bar := strings.Repeat("█", barWidth)
		color := getLanguageColor(ls.name)
		
		fmt.Printf("%-15s %s %5.1f%% (%d repos)\n",
			ls.name,
			lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(bar),
			percentage,
			ls.count)
	}
}

func getLanguageColor(language string) string {
	colors := map[string]string{
		"Go":         "#00ADD8",
		"JavaScript": "#f1e05a",
		"TypeScript": "#2b7489",
		"Python":     "#3572A5",
		"Java":       "#b07219",
		"C++":        "#f34b7d",
		"C":          "#555555",
		"Rust":       "#dea584",
		"Swift":      "#ffac45",
		"Kotlin":     "#F18E33",
		"Ruby":       "#701516",
		"PHP":        "#4F5D95",
		"C#":         "#239120",
		"Dart":       "#00B4AB",
		"Shell":      "#89e051",
	}
	
	if color, exists := colors[language]; exists {
		return color
	}
	return "#626262" // Default gray
}

func displayRepositoryTimeline(timeline []RepoTimelineEntry) {
	if len(timeline) == 0 {
		fmt.Println("No timeline data available.")
		return
	}
	
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Year", "Repositories Created", "Activity Level"})
	table.SetBorder(true)
	
	for _, entry := range timeline {
		activity := "Low"
		if entry.Count >= 10 {
			activity = "High"
		} else if entry.Count >= 5 {
			activity = "Medium"
		}
		
		table.Append([]string{
			fmt.Sprintf("%d", entry.Year),
			fmt.Sprintf("%d", entry.Count),
			activity,
		})
	}
	
	table.Render()
}

func displayActivityAnalysis(activity *ActivityStats) {
	fmt.Printf("🔥 Contribution Score: %.1f\n", activity.ContributionScore)
	fmt.Printf("📝 Recent Commits: %d (last 3 months)\n", len(activity.RecentCommits))
	
	if len(activity.ProductiveHours) > 0 {
		fmt.Println("\n⏰ Most Productive Hours:")
		
		// Find peak hours
		maxCommits := 0
		peakHour := 0
		for hour, commits := range activity.ProductiveHours {
			if commits > maxCommits {
				maxCommits = commits
				peakHour = hour
			}
		}
		
		fmt.Printf("   Peak: %02d:00 (%d commits)\n", peakHour, maxCommits)
		
		// Show hourly distribution
		for hour := 0; hour < 24; hour += 4 {
			commits := activity.ProductiveHours[hour]
			fmt.Printf("   %02d:00-%02d:00: %s\n", 
				hour, hour+3, 
				strings.Repeat("▓", commits/2))
		}
	}
	
	if len(activity.CommitFrequency) > 0 {
		fmt.Println("\n📅 Weekly Activity Pattern:")
		days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
		for _, day := range days {
			commits := activity.CommitFrequency[day]
			fmt.Printf("   %-9s: %s (%d)\n", 
				day, 
				strings.Repeat("●", commits/2), 
				commits)
		}
	}
}

func displayTopRepositoriesEnhanced(repos []*github.Repository) {
	count := 0
	
	// Sort repositories by stars
	sort.Slice(repos, func(i, j int) bool {
		return repos[i].GetStargazersCount() > repos[j].GetStargazersCount()
	})
	
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Repository", "⭐ Stars", "🍴 Forks", "Language", "Size (KB)", "Last Update"})
	table.SetBorder(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	
	for _, repo := range repos {
		if count >= 10 || repo.GetFork() || repo.GetPrivate() {
			continue
		}
		
		language := repo.GetLanguage()
		if language == "" {
			language = "Unknown"
		}
		
		lastUpdate := "Unknown"
		if repo.UpdatedAt != nil {
			lastUpdate = repo.UpdatedAt.Format("Jan 2006")
		}
		
		table.Append([]string{
			repo.GetName(),
			fmt.Sprintf("%d", repo.GetStargazersCount()),
			fmt.Sprintf("%d", repo.GetForksCount()),
			language,
			fmt.Sprintf("%d", repo.GetSize()),
			lastUpdate,
		})
		
		count++
	}
	
	table.Render()
}

func displayRepositoryHealth(stats *UserStats) {
	total := stats.RepoTypes["public"] + stats.RepoTypes["private"]
	if total == 0 {
		fmt.Println("No repositories found.")
		return
	}
	
	fmt.Println("📊 Update Frequency Analysis:")
	
	updateTable := tablewriter.NewWriter(os.Stdout)
	updateTable.SetHeader([]string{"Category", "Count", "Percentage", "Health"})
	updateTable.SetBorder(true)
	
	categories := []struct {
		name   string
		key    string
		health string
	}{
		{"Active (Weekly)", "weekly", "Excellent"},
		{"Regular (Monthly)", "monthly", "Good"},
		{"Quarterly", "quarterly", "Fair"},
		{"Yearly", "yearly", "Poor"},
		{"Stale (>1 year)", "stale", "Critical"},
	}
	
	for _, cat := range categories {
		count := stats.UpdateFrequency[cat.key]
		percentage := float64(count) / float64(total) * 100
		
		updateTable.Append([]string{
			cat.name,
			fmt.Sprintf("%d", count),
			fmt.Sprintf("%.1f%%", percentage),
			cat.health,
		})
	}
	
	updateTable.Render()
	
	// Health recommendations
	staleRepos := stats.UpdateFrequency["stale"]
	if staleRepos > 0 {
		fmt.Printf("\n⚠️  Health Warning: %d repositories haven't been updated in over a year\n", staleRepos)
		fmt.Println("💡 Consider archiving unused repositories or adding maintenance schedules")
	}
	
	activeRepos := stats.UpdateFrequency["weekly"] + stats.UpdateFrequency["monthly"]
	if activeRepos > total/2 {
		fmt.Println("\n✅ Good News: Most repositories are actively maintained!")
	}
}


