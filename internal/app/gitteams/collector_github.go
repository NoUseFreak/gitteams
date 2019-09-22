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
	setRootFlag("github-username", "", "", "Github username, provide to get personal repositories")
	setRootFlagBool("github-include-forks", "", true, "Github include forks")
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
		viper.GetBool("github-include-forks"),
	)
}

func (gh *GithubCollector) collectGithub(token, team, username string, forks bool) []Repo {
	origin := RepoOrigin{
		Name:  "github",
		Short: "gh",
	}

	result := []Repo{}

	for _, ghrepo := range gh.fetchGithubRepos(token, team, username, forks) {
		repo := NewRepo("git", &origin)

		repo.MainBranch = ghrepo.GetDefaultBranch()
		repo.Name = ghrepo.GetFullName()
		repo.URL = ghrepo.GetGitURL()

		result = append(result, repo)
	}

	return result
}

func (gh *GithubCollector) fetchGithubRepos(token, team, username string, forks bool) []*github.Repository {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	ghrepos := []*github.Repository{}
	if team != "" {
		opts := &github.RepositoryListByOrgOptions{
			Type: "sources",
		}
		for {
			repos, resp, err := client.Repositories.ListByOrg(ctx, team, opts)
			if err != nil {
				logrus.Error(err)
			}
			ghrepos = append(ghrepos, repos...)
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}

	}
	if username != "" {
		opts := &github.RepositoryListOptions{
			Type: "sources",
		}
		for {
			repos, resp, err := client.Repositories.List(ctx, username, opts)
			if err != nil {
				logrus.Error(err)
			}
			ghrepos = append(ghrepos, repos...)
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
	}

	all := []*github.Repository{}
	for _, r := range ghrepos {
		if *r.Fork && !forks {
			continue
		}
		all = append(all, r)
	}

	return all
}
