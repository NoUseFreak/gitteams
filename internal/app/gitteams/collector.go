package gitteams

import "github.com/sirupsen/logrus"

type Collector interface {
	GetName() string
	IsAvailable() bool
	Collect() []Repo
}

func GetCollectors() []Collector {
	return []Collector{
		new(BitbucketCollector),
		new(GithubCollector),
		new(GitlabCollector),
	}
}

func CollectRepos() []Repo {
	repos := []Repo{}

	for _, c := range GetCollectors() {
		if c.IsAvailable() {
			logrus.Debugf("Collecting %s", c.GetName())
			repos = append(repos, c.Collect()...)
		}
	}

	return repos
}
