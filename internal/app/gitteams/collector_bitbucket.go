package gitteams

import (
	"github.com/ktrysmt/go-bitbucket"
	"github.com/spf13/viper"
)

func init() {
	setRootFlag("bitbucket-username", "", "", "Bitbucket username")
	setRootFlag("bitbucket-password", "", "", "Bitbucket password")
	setRootFlag("bitbucket-team", "", "", "Bitbucket token")
}

type BitbucketCollector struct{}

func (b *BitbucketCollector) GetName() string {
	return "bitbucket"
}

func (b *BitbucketCollector) IsAvailable() bool {
	return viper.GetString("bitbucket-team") != ""
}

func (b *BitbucketCollector) Collect() []Repo {
	return b.collectBitbucket(
		viper.GetString("bitbucket-username"),
		viper.GetString("bitbucket-password"),
		viper.GetString("bitbucket-team"),
	)
}

func (b *BitbucketCollector) collectBitbucket(username, password, team string) []Repo {
	origin := RepoOrigin{
		Name:  "bitbucket",
		Short: "bb",
	}

	result := []Repo{}

	c := bitbucket.NewBasicAuth(username, password)
	repos, err := c.Repositories.ListForTeam(&bitbucket.RepositoriesOptions{
		Owner: team,
	})
	if err != nil {
		panic(err)
	}

	for _, v := range repos.(map[string]interface{})["values"].([]interface{}) {
		repo := NewRepo("git", &origin)
		r := v.(map[string]interface{})
		repo.MainBranch = r["mainbranch"].(map[string]interface{})["name"].(string)
		repo.Name = r["full_name"].(string)
		repo.URL = b.bbGetURL(r)

		result = append(result, repo)
	}

	return result
}

func (b *BitbucketCollector) bbGetURL(r map[string]interface{}) string {
	for _, l := range r["links"].(map[string]interface{})["clone"].([]interface{}) {
		if l.(map[string]interface{})["name"].(string) == "ssh" {
			return l.(map[string]interface{})["href"].(string)
		}
	}
	return ""
}
