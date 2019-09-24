package gitteams

import "github.com/sirupsen/logrus"

// Collector is an interface for collecting data from different platforms.
type Collector interface {
	GetName() string
	IsAvailable() bool
	Collect() []Repo
}

// GetCollectors returns a list of available collectors.
func GetCollectors() []Collector {
	return []Collector{
		new(BitbucketCollector),
		new(GithubCollector),
		new(GitlabCollector),
	}
}

// CollectRepos return repositories collected by all collectors.
func CollectRepos() []Repo {
	repos := []Repo{}

	for _, c := range GetCollectors() {
		if c.IsAvailable() {
			logrus.Debugf("Collecting - %s", c.GetName())
			repos = append(repos, c.Collect()...)
		} else {
			logrus.Debugf("Collecting skipped - %s", c.GetName())
		}
	}

	return repos
}
