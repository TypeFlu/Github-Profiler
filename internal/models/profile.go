package models

import (
	"github.com/google/go-github/v73/github"
)

// UserProfile represents the comprehensive user profile data
type UserProfile struct {
	User         *github.User         `json:"user"`
	Repositories []*github.Repository `json:"repositories"`
	Languages    LanguageStats        `json:"languages"`
	Stats        ProfileStats         `json:"stats"`
	Activity     ActivityStats        `json:"activity"`
	Ranking      RankingInfo          `json:"ranking"`
}

// LanguageStats represents programming language usage statistics
type LanguageStats struct {
	TotalBytes int                     `json:"total_bytes"`
	Languages  map[string]LanguageInfo `json:"languages"`
}

// LanguageInfo contains detailed language information
type LanguageInfo struct {
	Name       string  `json:"name"`
	Bytes      int     `json:"bytes"`
	Percentage float64 `json:"percentage"`
	RepoCount  int     `json:"repo_count"`
}

// ProfileStats contains repository and contribution statistics
type ProfileStats struct {
	TotalStars       int             `json:"total_stars"`
	TotalForks       int             `json:"total_forks"`
	TotalSize        int64           `json:"total_size_kb"`
	AvgStarsPerRepo  float64         `json:"avg_stars_per_repo"`
	RepoTypes        map[string]int  `json:"repo_types"`
	UpdateFrequency  map[string]int  `json:"update_frequency"`
	CreationTimeline []TimelineEntry `json:"creation_timeline"`
}

// TimelineEntry represents repository creation over time
type TimelineEntry struct {
	Year  int `json:"year"`
	Count int `json:"count"`
}

// ActivityStats represents user activity patterns
type ActivityStats struct {
	ContributionScore float64        `json:"contribution_score"`
	RecentCommits     int            `json:"recent_commits"`
	CommitFrequency   map[string]int `json:"commit_frequency"`
	ProductiveHours   map[string]int `json:"productive_hours"`
}

// RankingInfo represents the developer ranking system
type RankingInfo struct {
	OverallRank     string  `json:"overall_rank"`
	Badge           string  `json:"badge"`
	TotalScore      float64 `json:"total_score"`
	Percentile      float64 `json:"percentile"`
	SocialScore     float64 `json:"social_score"`
	CodeScore       float64 `json:"code_score"`
	ActivityScore   float64 `json:"activity_score"`
	InnovationScore float64 `json:"innovation_score"`
}

// RankTier represents different developer levels
type RankTier struct {
	Name        string
	Badge       string
	MinScore    float64
	MaxScore    float64
	Description string
}

// GetRankTiers returns all available rank tiers
func GetRankTiers() []RankTier {
	return []RankTier{
		{"Elite Developer", "ELITE", 90, 100, "Industry leader and innovator"},
		{"Senior Developer", "SENIOR", 80, 89, "Experienced professional"},
		{"Experienced Developer", "EXPERIENCED", 70, 79, "Skilled practitioner"},
		{"Active Developer", "ACTIVE", 60, 69, "Regular contributor"},
		{"Growing Developer", "GROWING", 50, 59, "Developing skills"},
		{"Junior Developer", "JUNIOR", 30, 49, "Learning and improving"},
		{"Beginner", "BEGINNER", 0, 29, "Just starting out"},
	}
}

// GetRankByScore returns the appropriate rank for a given score
func GetRankByScore(score float64) RankTier {
	tiers := GetRankTiers()
	for _, tier := range tiers {
		if score >= tier.MinScore && score <= tier.MaxScore {
			return tier
		}
	}
	return tiers[len(tiers)-1] // Return beginner as fallback
}
