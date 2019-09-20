package gitteams

import (
	"context"

	"github.com/google/go-github/v28/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func init() {
	setRootFlag("github-token", "", "", "Github token")
	setRootFlag("github-team", "", "", "Github team")
	setRootFlag("github-username", "", "", "Github username")
}

type GithubCollector struct{}

func (b *GithubCollector) GetName() string {
	return "github"
}

func (gh *GithubCollector) IsAvailable() bool {
	return viper.GetString("github-team") != "" ||
		viper.GetString("github-username") != ""
}

func (gh *GithubCollector) Collect() []Repo {
	return gh.collectGithub(
		viper.GetString("github-token"),
		viper.GetString("github-team"),
		viper.GetString("github-username"),
	)
}

func (gh *GithubCollector) collectGithub(token, team, username string) []Repo {
	origin := RepoOrigin{
		Name:  "github",
		Short: "gh",
	}

	result := []Repo{}

	for _, ghrepo := range gh.fetchGithubRepos(token, team, username) {
		repo := NewRepo("git", &origin)

		repo.MainBranch = ghrepo.GetDefaultBranch()
		repo.Name = ghrepo.GetFullName()
		repo.URL = ghrepo.GetGitURL()

		result = append(result, repo)
	}

	return result
}

func (gh *GithubCollector) fetchGithubRepos(token, team, username string) []*github.Repository {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	ghrepos := []*github.Repository{}
	if team != "" {
		teamrepos, _, err := client.Repositories.ListByOrg(ctx, team, &github.RepositoryListByOrgOptions{
			Type: "sources",
		})
		if err != nil {
			logrus.Error(err)
		}
		ghrepos = append(ghrepos, teamrepos...)
	}
	if username != "" {
		userrepos, _, err := client.Repositories.List(ctx, username, &github.RepositoryListOptions{
			Type: "sources",
		})
		if err != nil {
			logrus.Error(err)
		}
		ghrepos = append(ghrepos, userrepos...)
	}

	return ghrepos
}
