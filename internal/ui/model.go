package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github-profiler/internal/models"
	"github-profiler/internal/services"
)

// AppState represents the different states of the application
type AppState int

const (
	StateInput AppState = iota
	StateLoading
	StateProfileView
	StateError
)

// Model represents the main application model following Elm Architecture
type Model struct {
	// Application state
	state    AppState
	username string
	token    string
	format   string

	// UI components
	spinner spinner.Model
	profile *models.UserProfile
	error   error
	width   int
	height  int

	// Services
	githubService *services.GitHubService

	// Navigation
	activeView ViewType
	views      []ViewType
}

// ViewType represents different profile views
type ViewType int

const (
	ViewOverview ViewType = iota
	ViewRepositories
	ViewLanguages
	ViewActivity
	ViewRanking
)

// ViewInfo holds view metadata
type ViewInfo struct {
	Type  ViewType
	Title string
	Icon  string
}

// GetViews returns all available views
func GetViews() []ViewInfo {
	return []ViewInfo{
		{ViewOverview, "Overview", "[O]"},
		{ViewRepositories, "Repositories", "[R]"},
		{ViewLanguages, "Languages", "[L]"},
		{ViewActivity, "Activity", "[A]"},
		{ViewRanking, "Ranking", "[K]"},
	}
}

// NewModel creates a new application model
func NewModel(username, token, format string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	githubService := services.NewGitHubService(token)

	state := StateInput
	if username != "" {
		state = StateLoading
	}

	return Model{
		state:         state,
		username:      username,
		token:         token,
		format:        format,
		spinner:       s,
		githubService: githubService,
		activeView:    ViewOverview,
		views:         []ViewType{ViewOverview, ViewRepositories, ViewLanguages, ViewActivity, ViewRanking},
	}
}

// Init implements the bubbletea.Model interface
func (m Model) Init() tea.Cmd {
	if m.state == StateLoading {
		return tea.Batch(
			m.spinner.Tick,
			m.fetchProfile,
		)
	}
	return m.spinner.Tick
}

// Update implements the bubbletea.Model interface
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if m.state == StateInput && m.username != "" {
				m.state = StateLoading
				return m, tea.Batch(
					m.spinner.Tick,
					m.fetchProfile,
				)
			}

		case "left", "h":
			if m.state == StateProfileView {
				return m.previousView(), nil
			}

		case "right", "l":
			if m.state == StateProfileView {
				return m.nextView(), nil
			}

		case "r":
			if m.state == StateError || m.state == StateProfileView {
				m.state = StateLoading
				m.error = nil
				return m, tea.Batch(
					m.spinner.Tick,
					m.fetchProfile,
				)
			}

		default:
			if m.state == StateInput {
				if msg.String() == "backspace" {
					if len(m.username) > 0 {
						m.username = m.username[:len(m.username)-1]
					}
				} else if len(msg.String()) == 1 {
					m.username += msg.String()
				}
				return m, nil
			}
		}

	case ProfileFetchedMsg:
		m.profile = msg.Profile
		m.state = StateProfileView
		return m, nil

	case ProfileErrorMsg:
		m.error = msg.Error
		m.state = StateError
		return m, nil

	case spinner.TickMsg:
		if m.state == StateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

// View implements the bubbletea.Model interface
func (m Model) View() string {
	switch m.state {
	case StateInput:
		return m.renderInputView()
	case StateLoading:
		return m.renderLoadingView()
	case StateProfileView:
		return m.renderProfileView()
	case StateError:
		return m.renderErrorView()
	default:
		return "Unknown state"
	}
}

// Message types for Elm Architecture
type ProfileFetchedMsg struct {
	Profile *models.UserProfile
}

type ProfileErrorMsg struct {
	Error error
}

// fetchProfile is a command that fetches user profile data
func (m Model) fetchProfile() tea.Msg {
	var profile *models.UserProfile
	var err error

	if m.username == "demo-user" {
		profile, err = m.githubService.GetDemoProfile()
	} else {
		profile, err = m.githubService.GetUserProfile(m.username)
	}

	if err != nil {
		return ProfileErrorMsg{Error: err}
	}
	return ProfileFetchedMsg{Profile: profile}
}

// Navigation helpers
func (m Model) nextView() Model {
	currentIndex := -1
	for i, view := range m.views {
		if view == m.activeView {
			currentIndex = i
			break
		}
	}

	if currentIndex >= 0 && currentIndex < len(m.views)-1 {
		m.activeView = m.views[currentIndex+1]
	}

	return m
}

func (m Model) previousView() Model {
	currentIndex := -1
	for i, view := range m.views {
		if view == m.activeView {
			currentIndex = i
			break
		}
	}

	if currentIndex > 0 {
		m.activeView = m.views[currentIndex-1]
	}

	return m
}

// Rendering methods
func (m Model) renderInputView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(1).
		Render("GitHub Profiler")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("v1.0.0 - github@Tyeflu")

	prompt := lipgloss.NewStyle().
		MarginTop(2).
		MarginBottom(1).
		Render("Enter GitHub username:")

	input := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Render(m.username + "|")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1).
		Render("Press Enter to analyze - Ctrl+C to quit")

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s", title, subtitle, prompt, input, instructions)
}

func (m Model) renderLoadingView() string {
	return fmt.Sprintf("\n%s Fetching GitHub data for %s...\n\n",
		m.spinner.View(),
		lipgloss.NewStyle().Bold(true).Render(m.username))
}

func (m Model) renderErrorView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196")).
		Render("ERROR")

	errorMsg := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1).
		Render(m.error.Error())

	instructions := lipgloss.NewStyle().
		MarginTop(2).
		Foreground(lipgloss.Color("241")).
		Render("Press 'r' to retry - Ctrl+C to quit")

	return fmt.Sprintf("%s\n%s\n%s", title, errorMsg, instructions)
}

func (m Model) renderProfileView() string {
	if m.profile == nil {
		return "No profile data available"
	}

	// Header with navigation
	header := m.renderHeader()

	// Active view content
	var content string
	switch m.activeView {
	case ViewOverview:
		content = m.renderOverviewView()
	case ViewRepositories:
		content = m.renderRepositoriesView()
	case ViewLanguages:
		content = m.renderLanguagesView()
	case ViewActivity:
		content = m.renderActivityView()
	case ViewRanking:
		content = m.renderRankingView()
	}

	// Footer with navigation instructions
	footer := m.renderFooter()

	return fmt.Sprintf("%s\n\n%s\n\n%s", header, content, footer)
}

func (m Model) renderHeader() string {
	user := m.profile.User

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Render("GitHub Profile Analysis")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("v1.0.0 - github@Tyeflu")

	userInfo := fmt.Sprintf("%s (%s)",
		getStringValue(user.Name),
		user.GetLogin())

	// Navigation tabs
	views := GetViews()
	var tabs []string
	for i, view := range views {
		viewType := ViewType(i)
		style := lipgloss.NewStyle().Padding(0, 1)

		if viewType == m.activeView {
			style = style.Background(lipgloss.Color("86")).Foreground(lipgloss.Color("0"))
		} else {
			style = style.Foreground(lipgloss.Color("241"))
		}

		tabs = append(tabs, style.Render(fmt.Sprintf("%s %s", view.Icon, view.Title)))
	}

	navigation := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)

	return fmt.Sprintf("%s\n%s\n\n%s\n\n%s", title, subtitle, userInfo, navigation)
}

func (m Model) renderFooter() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render("← → Navigate • r Refresh • q Quit")
}

// View rendering methods for different sections

func (m Model) renderOverviewView() string {
	user := m.profile.User
	stats := m.profile.Stats

	// User basic info
	userBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		MarginBottom(1)

	userInfo := fmt.Sprintf(`User: %s (%s)
Bio: %s
Company: %s
Location: %s
Website: %s
Joined: %s
Stats: Public Repos: %d | Followers: %d | Following: %d`,
		getStringValue(user.Name),
		user.GetLogin(),
		getStringValue(user.Bio),
		getStringValue(user.Company),
		getStringValue(user.Location),
		getStringValue(user.Blog),
		user.GetCreatedAt().Format("January 2006"),
		user.GetPublicRepos(),
		user.GetFollowers(),
		user.GetFollowing())

	// Quick stats
	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		MarginTop(1)

	quickStats := fmt.Sprintf(`Total Stars: %d
Total Forks: %d
Repository Size: %.1f MB
Avg Stars/Repo: %.1f`,
		stats.TotalStars,
		stats.TotalForks,
		float64(stats.TotalSize)/1024,
		stats.AvgStarsPerRepo)

	return userBox.Render(userInfo) + "\n" + statsBox.Render(quickStats)
}

func (m Model) renderRepositoriesView() string {
	repos := m.profile.Repositories

	if len(repos) == 0 {
		return "No repositories found"
	}

	// Sort repos by stars
	topRepos := repos
	if len(repos) > 8 {
		topRepos = repos[:8]
	}

	var repoList []string
	for i, repo := range topRepos {
		if repo.GetFork() {
			continue
		}

		repoInfo := fmt.Sprintf(`%d. %s
   %s
   Stars: %d  Forks: %d  Language: %s
   Updated: %s`,
			i+1,
			repo.GetName(),
			getStringValue(repo.Description),
			repo.GetStargazersCount(),
			repo.GetForksCount(),
			getStringValue(repo.Language),
			repo.GetUpdatedAt().Format("Jan 2006"))

		repoList = append(repoList, repoInfo)
	}

	return lipgloss.JoinVertical(lipgloss.Left, repoList...)
}

func (m Model) renderLanguagesView() string {
	languages := m.profile.Languages

	if len(languages.Languages) == 0 {
		return "No language data available"
	}

	// Convert to sorted slice
	type langItem struct {
		name       string
		percentage float64
		repos      int
	}

	var langList []langItem
	for _, lang := range languages.Languages {
		langList = append(langList, langItem{
			name:       lang.Name,
			percentage: lang.Percentage,
			repos:      lang.RepoCount,
		})
	}

	// Sort by percentage (simple bubble sort for now)
	for i := 0; i < len(langList)-1; i++ {
		for j := 0; j < len(langList)-i-1; j++ {
			if langList[j].percentage < langList[j+1].percentage {
				langList[j], langList[j+1] = langList[j+1], langList[j]
			}
		}
	}

	var langDisplay []string
	for _, lang := range langList {
		if lang.percentage < 1.0 {
			continue // Skip languages with less than 1%
		}

		// Create progress bar
		barWidth := 20
		fillWidth := int((lang.percentage / 100.0) * float64(barWidth))
		bar := strings.Repeat("█", fillWidth) + strings.Repeat("░", barWidth-fillWidth)

		langInfo := fmt.Sprintf("%-12s %s %.1f%% (%d repos)",
			lang.name,
			bar,
			lang.percentage,
			lang.repos)

		langDisplay = append(langDisplay, langInfo)
	}

	return lipgloss.JoinVertical(lipgloss.Left, langDisplay...)
}

func (m Model) renderActivityView() string {
	stats := m.profile.Stats
	activity := m.profile.Activity

	activityInfo := fmt.Sprintf(`Contribution Score: %.1f
Recent Commits: %d

Repository Update Frequency:
   Weekly: %d repositories
   Monthly: %d repositories  
   Quarterly: %d repositories
   Yearly: %d repositories
   Stale (>1 year): %d repositories

Repository Timeline:`,
		activity.ContributionScore,
		activity.RecentCommits,
		stats.UpdateFrequency["weekly"],
		stats.UpdateFrequency["monthly"],
		stats.UpdateFrequency["quarterly"],
		stats.UpdateFrequency["yearly"],
		stats.UpdateFrequency["stale"])

	// Add timeline
	var timeline []string
	for _, entry := range stats.CreationTimeline {
		timeline = append(timeline, fmt.Sprintf("   %d: %d repositories", entry.Year, entry.Count))
	}

	return activityInfo + "\n" + lipgloss.JoinVertical(lipgloss.Left, timeline...)
}

func (m Model) renderRankingView() string {
	ranking := m.profile.Ranking

	// Rank badge
	badgeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("86")).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	badge := badgeStyle.Render(fmt.Sprintf("BADGE: %s", ranking.Badge))

	// Score breakdown
	scoreInfo := fmt.Sprintf(`Overall Rank: %s
Total Score: %.1f/100 (%.1f%% percentile)

Score Breakdown:
Social Score:     %.1f/25  (%.1f%%)
Code Score:       %.1f/30  (%.1f%%)
Activity Score:   %.1f/25  (%.1f%%)
Innovation Score: %.1f/20  (%.1f%%)`,
		ranking.OverallRank,
		ranking.TotalScore,
		ranking.Percentile,
		ranking.SocialScore,
		(ranking.SocialScore/25)*100,
		ranking.CodeScore,
		(ranking.CodeScore/30)*100,
		ranking.ActivityScore,
		(ranking.ActivityScore/25)*100,
		ranking.InnovationScore,
		(ranking.InnovationScore/20)*100)

	return badge + "\n\n" + scoreInfo
}

// Helper function to safely get string values
func getStringValue(s *string) string {
	if s == nil {
		return "Not specified"
	}
	return *s
}
