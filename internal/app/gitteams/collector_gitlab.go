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

type GitlabCollector struct{}

func (b *GitlabCollector) GetName() string {
	return "gitlab"
}

func (gh *GitlabCollector) IsAvailable() bool {
	return viper.GetString("gitlab-group") != "" || viper.GetBool("gitlab-include-personal")
}

func (gh *GitlabCollector) Collect() []Repo {
	return gh.collectGitlab(
		viper.GetString("gitlab-token"),
		viper.GetString("gitlab-group"),
		viper.GetBool("gitlab-include-personal"),
	)
}

func (gh *GitlabCollector) collectGitlab(token, group string, personal bool) []Repo {
	origin := RepoOrigin{
		Name:  "gitlab",
		Short: "gl",
	}

	result := []Repo{}

	for _, data := range gh.fetchGitlabRepos(token, group, personal) {
		repo := NewRepo("git", &origin)

		repo.MainBranch = data.DefaultBranch
		repo.URL = data.SSHURLToRepo
		repo.Name = data.PathWithNamespace

		result = append(result, repo)
	}

	return result
}

func (gh *GitlabCollector) fetchGitlabRepos(token, group string, personal bool) []*gitlab.Project {
	client := gitlab.NewClient(nil, token)

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
