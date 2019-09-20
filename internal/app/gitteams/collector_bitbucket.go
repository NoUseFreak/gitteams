package gitteams

import (
	"github.com/ktrysmt/go-bitbucket"
	"github.com/spf13/viper"
)

func init() {
	setRootFlag("bitbucket-username", "", "", "Bitbucket username")
	setRootFlag("bitbucket-password", "", "", "Bitbucket password")
	setRootFlag("bitbucket-team", "", "", "Bitbucket token")
	setRootFlagBool("bitbucket-include-personal", "", true, "Bitbucket include personal repositories")
}

type BitbucketCollector struct{}

func (b *BitbucketCollector) GetName() string {
	return "bitbucket"
}

func (b *BitbucketCollector) IsAvailable() bool {
	return viper.GetString("bitbucket-username") != "" &&
		viper.GetString("bitbucket-password") != "" &&
		(viper.GetString("bitbucket-team") != "" || viper.GetBool("bitbucket-include-personal"))
}

func (b *BitbucketCollector) Collect() []Repo {
	return b.collectBitbucket(
		viper.GetString("bitbucket-username"),
		viper.GetString("bitbucket-password"),
		viper.GetString("bitbucket-team"),
		viper.GetBool("bitbucket-include-personal"),
	)
}

func (b *BitbucketCollector) collectBitbucket(username, password, team string, personal bool) []Repo {
	c := bitbucket.NewBasicAuth(username, password)

	all := []Repo{}
	if viper.GetString("bitbucket-team") != "" {
		repos, err := c.Repositories.ListForTeam(&bitbucket.RepositoriesOptions{
			Owner: team,
		})
		if err != nil {
			panic(err)
		}
		all = append(all, b.parseResult(repos)...)
	}
	if viper.GetBool("bitbucket-include-personal") {
		repos, err := c.Repositories.ListForAccount(&bitbucket.RepositoriesOptions{
			Owner: username,
		})
		if err != nil {
			panic(err)
		}
		all = append(all, b.parseResult(repos)...)
	}

	return all
}

func (b BitbucketCollector) parseResult(repos interface{}) []Repo {
	origin := RepoOrigin{
		Name:  "bitbucket",
		Short: "bb",
	}

	result := []Repo{}
	for _, v := range repos.(map[string]interface{})["values"].([]interface{}) {
		repo := NewRepo("git", &origin)
		r := v.(map[string]interface{})
		if r["mainbranch"] == nil {
			continue
		}
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
