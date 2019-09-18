package gitteams

import (
	"context"

	"github.com/google/go-github/v28/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func init() {
	setRootFlag("github-token", "", "", "Github token")
	setRootFlag("github-team", "", "", "Github team")
}

type GithubCollector struct{}

func (b *GithubCollector) GetName() string {
	return "github"
}

func (gh *GithubCollector) IsAvailable() bool {
	return viper.GetString("github-team") != ""
}

func (gh *GithubCollector) Collect() []Repo {
	return gh.collectGithub(
		viper.GetString("github-token"),
		viper.GetString("github-team"),
	)
}

func (gh *GithubCollector) collectGithub(token, team string) []Repo {
	origin := RepoOrigin{
		Name:  "github",
		Short: "gh",
	}

	result := []Repo{}

	for _, ghrepo := range gh.fetchGithubRepos(token, team) {
		repo := NewRepo("git", &origin)

		repo.MainBranch = ghrepo.GetDefaultBranch()
		repo.Name = ghrepo.GetFullName()
		repo.URL = ghrepo.GetGitURL()

		result = append(result, repo)
	}

	return result
}

func (gh *GithubCollector) fetchGithubRepos(token, team string) []*github.Repository {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	opt := &github.RepositoryListByOrgOptions{Type: "sources"}
	ghrepos, _, _ := client.Repositories.ListByOrg(ctx, team, opt)

	return ghrepos
}
