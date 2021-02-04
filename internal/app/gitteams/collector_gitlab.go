package gitteams

import (
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

func init() {
	setRootFlag("gitlab-token", "", "", "Gitlab token")
	setRootFlag("gitlab-team", "", "", "Gitlab team")
	setRootFlagBool("gitlab-include-personal", "", true, "Gitlab include personal repositories")
}

// GitlabCollector collects repositories hosted on github.com.
type GitlabCollector struct{}

// GetName return the name of the collector.
func (gl *GitlabCollector) GetName() string {
	return "gitlab"
}

// IsAvailable checks if the required config is present to collect the data.
func (gl *GitlabCollector) IsAvailable() bool {
	return viper.GetString("gitlab-group") != "" || viper.GetBool("gitlab-include-personal")
}

// Collect get the data from the origins api.
func (gl *GitlabCollector) Collect() []Repo {
	return gl.collectGitlab(
		viper.GetString("gitlab-token"),
		viper.GetString("gitlab-group"),
		viper.GetBool("gitlab-include-personal"),
	)
}

func (gl *GitlabCollector) collectGitlab(token, group string, personal bool) []Repo {
	origin := RepoOrigin{
		Name:  "gitlab",
		Short: "gl",
	}

	result := []Repo{}

	for _, data := range gl.fetchGitlabRepos(token, group, personal) {
		repo := NewRepo("git", &origin)

		repo.MainBranch = data.DefaultBranch
		repo.URL = data.SSHURLToRepo
		repo.Name = data.PathWithNamespace

		result = append(result, repo)
	}

	return result
}

func (gl *GitlabCollector) fetchGitlabRepos(token, group string, personal bool) []*gitlab.Project {
	client, err := gitlab.NewClient(token)
	if err != nil {
		panic(err)
	}

	all := []*gitlab.Project{}

	if group != "" {
		projects, _, err := client.Groups.ListGroupProjects(group, &gitlab.ListGroupProjectsOptions{})
		if err != nil {
			panic(err)
		}
		all = append(all, projects...)
	}

	if personal {
		projects, _, err := client.Projects.ListProjects(&gitlab.ListProjectsOptions{
			Owned: gitlab.Bool(true),
		})
		if err != nil {
			panic(err)
		}
		all = append(all, projects...)
	}

	return all
}
