package cmd

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v73/github"
)

func displayJSON(data *UserData) {
	output := map[string]interface{}{
		"user": map[string]interface{}{
			"username":     data.User.GetLogin(),
			"name":         getStringValue(data.User.Name, data.User.GetLogin()),
			"bio":          getStringValue(data.User.Bio, "No bio available"),
			"company":      getStringValue(data.User.Company, "Not specified"),
			"location":     getStringValue(data.User.Location, "Unknown"),
			"website":      getStringValue(data.User.Blog, "Not specified"),
			"twitter":      getStringValue(data.User.TwitterUsername, "Not specified"),
			"created_at":   data.User.GetCreatedAt().Format(time.RFC3339),
			"followers":    data.User.GetFollowers(),
			"following":    data.User.GetFollowing(),
			"public_repos": data.User.GetPublicRepos(),
		},
		"stats": map[string]interface{}{
			"total_stars":       data.Stats.TotalStars,
			"total_forks":       data.Stats.TotalForks,
			"total_size_kb":     data.Stats.TotalSize,
			"avg_stars_per_repo": data.Stats.AvgStarsPerRepo,
			"repo_types":        data.Stats.RepoTypes,
			"update_frequency":  data.Stats.UpdateFrequency,
		},
		"activity": map[string]interface{}{
			"contribution_score": data.Activity.ContributionScore,
			"commit_frequency":   data.Activity.CommitFrequency,
			"productive_hours":   data.Activity.ProductiveHours,
			"recent_commits":     len(data.Activity.RecentCommits),
		},
		"ranking":            calculateRankingData(data),
		"top_repositories":   getTopRepositoriesData(data.Repositories),
		"languages":          getLanguageData(data.Languages),
		"creation_timeline":  data.Stats.CreationTimeline,
	}
	
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error generating JSON: %v\n", err)
		return
	}
	
	fmt.Println(string(jsonData))
}

func displayHTML(data *UserData) {
	htmlTemplate := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GitHub Profile: {{.User.Login}}</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            min-height: 100vh;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .header {
            background: white;
            border-radius: 20px;
            padding: 30px;
            margin-bottom: 20px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
        }
        
        .profile-info {
            display: flex;
            align-items: center;
            gap: 20px;
            margin-bottom: 20px;
        }
        
        .avatar {
            width: 100px;
            height: 100px;
            border-radius: 50%;
            background: linear-gradient(45deg, #667eea, #764ba2);
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-size: 36px;
            font-weight: bold;
        }
        
        .user-details h1 {
            font-size: 2.5em;
            color: #333;
            margin-bottom: 10px;
        }
        
        .user-details p {
            color: #666;
            margin-bottom: 5px;
        }
        
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 20px;
        }
        
        .stat-card {
            background: white;
            border-radius: 15px;
            padding: 25px;
            text-align: center;
            box-shadow: 0 5px 15px rgba(0,0,0,0.1);
            transition: transform 0.3s ease;
        }
        
        .stat-card:hover {
            transform: translateY(-5px);
        }
        
        .stat-number {
            font-size: 2.5em;
            font-weight: bold;
            color: #667eea;
            margin-bottom: 10px;
        }
        
        .stat-label {
            color: #666;
            font-weight: 500;
        }
        
        .section {
            background: white;
            border-radius: 20px;
            padding: 30px;
            margin-bottom: 20px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
        }
        
        .section h2 {
            color: #333;
            margin-bottom: 20px;
            font-size: 1.8em;
        }
        
        .repo-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 20px;
        }
        
        .repo-card {
            border: 2px solid #f0f0f0;
            border-radius: 10px;
            padding: 20px;
            transition: border-color 0.3s ease;
        }
        
        .repo-card:hover {
            border-color: #667eea;
        }
        
        .repo-name {
            font-size: 1.2em;
            font-weight: bold;
            color: #333;
            margin-bottom: 10px;
        }
        
        .repo-desc {
            color: #666;
            margin-bottom: 15px;
        }
        
        .repo-stats {
            display: flex;
            gap: 15px;
            font-size: 0.9em;
            color: #888;
        }
        
        .language-bar {
            height: 30px;
            background: #f0f0f0;
            border-radius: 15px;
            margin-bottom: 10px;
            overflow: hidden;
            position: relative;
        }
        
        .language-fill {
            height: 100%;
            background: linear-gradient(90deg, #667eea, #764ba2);
            border-radius: 15px;
        }
        
        .language-label {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 5px;
        }
        
        .rank-badge {
            display: inline-block;
            background: linear-gradient(45deg, #FFD700, #FFA500);
            color: white;
            padding: 10px 20px;
            border-radius: 25px;
            font-weight: bold;
            font-size: 1.2em;
            margin-bottom: 20px;
        }
        
        .chart-container {
            position: relative;
            height: 300px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="profile-info">
                <div class="avatar">{{substr .User.Login 0 1 | upper}}</div>
                <div class="user-details">
                    <h1>{{.UserName}}</h1>
                    <p><strong>@{{.User.Login}}</strong></p>
                    <p>{{.UserBio}}</p>
                    <p>📍 {{.UserLocation}} | 🏢 {{.UserCompany}}</p>
                    <p>📅 Joined {{.User.CreatedAt.Format "January 2006"}}</p>
                </div>
            </div>
        </div>
        
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-number">{{.User.GetFollowers}}</div>
                <div class="stat-label">Followers</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{.User.GetFollowing}}</div>
                <div class="stat-label">Following</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{.User.GetPublicRepos}}</div>
                <div class="stat-label">Public Repos</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">{{.Stats.TotalStars}}</div>
                <div class="stat-label">Total Stars</div>
            </div>
        </div>
        
        <div class="section">
            <h2>🏆 Developer Ranking</h2>
            <div class="rank-badge">{{.Rank}}</div>
            <div class="chart-container">
                <canvas id="rankChart"></canvas>
            </div>
        </div>
        
        <div class="section">
            <h2>⭐ Top Repositories</h2>
            <div class="repo-grid">
                {{range .TopRepos}}
                <div class="repo-card">
                    <div class="repo-name">📦 {{.GetName}}</div>
                    <div class="repo-desc">{{if .Description}}{{.GetDescription}}{{else}}No description{{end}}</div>
                    <div class="repo-stats">
                        <span>⭐ {{.GetStargazersCount}}</span>
                        <span>🍴 {{.GetForksCount}}</span>
                        <span>💻 {{if .Language}}{{.GetLanguage}}{{else}}Unknown{{end}}</span>
                    </div>
                </div>
                {{end}}
            </div>
        </div>
        
        <div class="section">
            <h2>💻 Language Breakdown</h2>
            {{range .Languages}}
            <div style="margin-bottom: 15px;">
                <div class="language-label">
                    <span>{{.Name}}</span>
                    <span>{{printf "%.1f" .Percentage}}% ({{.Count}} repos)</span>
                </div>
                <div class="language-bar">
                    <div class="language-fill" style="width: {{.Percentage}}%;"></div>
                </div>
            </div>
            {{end}}
        </div>
    </div>
    
    <script>
        // Chart.js for ranking visualization
        const ctx = document.getElementById('rankChart').getContext('2d');
        new Chart(ctx, {
            type: 'doughnut',
            data: {
                labels: ['Social Score', 'Repository Score', 'Activity Score'],
                datasets: [{
                    data: [{{.FollowerScore}}, {{.RepoScore}}, {{.ContributionScore}}],
                    backgroundColor: [
                        '#FF6384',
                        '#36A2EB',
                        '#FFCE56'
                    ],
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom'
                    }
                }
            }
        });
    </script>
    
    <footer style="text-align: center; padding: 20px; color: #666; font-size: 12px;">
        Generated by GitHub Profiler v{{.Version}} • {{.Author}}
    </footer>
</body>
</html>`

	// Create template data
	templateData := struct {
		User               *github.User
		UserName           string
		UserBio            string
		UserLocation       string
		UserCompany        string
		Rank               string
		FollowerScore      float64
		RepoScore          float64
		ContributionScore  float64
		TopRepos           []*github.Repository
		Languages          []LanguageData
		Stats              *UserStats
		Ranking            *RankingData
		Version            string
		Author             string
	}{
		User:              data.User,
		UserName:          getStringValue(data.User.Name, getStringValue(data.User.Login, "Unknown")),
		UserBio:           getStringValue(data.User.Bio, "No bio available"),
		UserLocation:      getStringValue(data.User.Location, "Unknown"),
		UserCompany:       getStringValue(data.User.Company, "Not specified"),
		Rank:              data.Ranking.Rank,
		FollowerScore:     data.Ranking.SocialScore,
		RepoScore:         data.Ranking.CodeScore,
		ContributionScore: data.Ranking.ActivityScore,
		TopRepos:          getTopRepos(data.Repositories, 6),
		Languages:         getLanguageDataStructured(data.Languages),
		Stats:             data.Stats,
		Ranking:           data.Ranking,
		Version:           version,
		Author:            author,
	}
	
	// Add template functions
	funcMap := template.FuncMap{
		"substr": func(s string, start, length int) string {
			if start >= len(s) {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
		"upper": strings.ToUpper,
	}
	
	tmpl, err := template.New("profile").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}
	
	// Create output file
	filename := fmt.Sprintf("%s-profile.html", data.User.Login)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating HTML file: %v\n", err)
		return
	}
	defer file.Close()
	
	err = tmpl.Execute(file, templateData)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}
	
	absPath, _ := filepath.Abs(filename)
	fmt.Printf("HTML report generated: %s\n", absPath)
}

type LanguageData struct {
	Name       string  `json:"name"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

func calculateRankingData(data *UserData) map[string]interface{} {
	ranking := data.Ranking
	
	return map[string]interface{}{
		"total_score":        ranking.OverallScore,
		"social_score":       ranking.SocialScore,
		"code_score":         ranking.CodeScore,
		"activity_score":     ranking.ActivityScore,
		"quality_score":      ranking.QualityScore,
		"rank":               ranking.Rank,
		"badge":              ranking.Badge,
		"percentile":         ranking.Percentile,
	}
}

func getTopRepositoriesData(repos []*github.Repository) []map[string]interface{} {
	var result []map[string]interface{}
	count := 0
	
	for _, repo := range repos {
		if count >= 5 || repo.GetPrivate() || repo.GetFork() {
			continue
		}
		
		result = append(result, map[string]interface{}{
			"name":        repo.GetName(),
			"description": getDescription(repo.Description),
			"stars":       repo.GetStargazersCount(),
			"forks":       repo.GetForksCount(),
			"language":    repo.GetLanguage(),
			"updated_at":  repo.GetUpdatedAt().Format("2006-01-02"),
		})
		count++
	}
	
	return result
}

func getLanguageData(languages map[string]int) map[string]interface{} {
	result := make(map[string]interface{})
	total := 0
	
	for _, count := range languages {
		total += count
	}
	
	for lang, count := range languages {
		percentage := float64(count) / float64(total) * 100
		result[lang] = map[string]interface{}{
			"count":      count,
			"percentage": percentage,
		}
	}
	
	return result
}

func getTopRepos(repos []*github.Repository, limit int) []*github.Repository {
	var result []*github.Repository
	count := 0
	
	for _, repo := range repos {
		if count >= limit || repo.GetPrivate() || repo.GetFork() {
			continue
		}
		result = append(result, repo)
		count++
	}
	
	return result
}

func getLanguageDataStructured(languages map[string]int) []LanguageData {
	var result []LanguageData
	total := 0
	
	for _, count := range languages {
		total += count
	}
	
	for lang, count := range languages {
		percentage := float64(count) / float64(total) * 100
		result = append(result, LanguageData{
			Name:       lang,
			Count:      count,
			Percentage: percentage,
		})
	}
	
	// Sort by count descending
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Count > result[i].Count {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	
	// Return top 10
	if len(result) > 10 {
		result = result[:10]
	}
	
	return result
}
