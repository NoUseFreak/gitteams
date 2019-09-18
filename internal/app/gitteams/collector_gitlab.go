package gitteams

import (
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

func init() {
	setRootFlag("gitlab-token", "", "", "Gitlab token")
	setRootFlag("gitlab-team", "", "", "Gitlab team")
}

type GitlabCollector struct{}

func (b *GitlabCollector) GetName() string {
	return "gitlab"
}

func (gh *GitlabCollector) IsAvailable() bool {
	return viper.GetString("gitlab-group") != ""
}

func (gh *GitlabCollector) Collect() []Repo {
	return gh.collectGitlab(
		viper.GetString("gitlab-token"),
		viper.GetString("gitlab-group"),
	)
}

func (gh *GitlabCollector) collectGitlab(token, group string) []Repo {
	origin := RepoOrigin{
		Name:  "gitlab",
		Short: "gl",
	}

	result := []Repo{}

	for _, data := range gh.fetchGitlabRepos(token, group) {
		repo := NewRepo("git", &origin)

		repo.MainBranch = data.DefaultBranch
		repo.URL = data.SSHURLToRepo
		repo.Name = data.PathWithNamespace

		result = append(result, repo)
	}

	return result
}

func (gh *GitlabCollector) fetchGitlabRepos(token, group string) []*gitlab.Project {
	client := gitlab.NewClient(nil, token)

	projects, _, err := client.Groups.ListGroupProjects(group, &gitlab.ListGroupProjectsOptions{})
	if err != nil {
		panic(err)
	}

	return projects
}
