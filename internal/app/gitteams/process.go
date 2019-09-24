package gitteams

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/exec"
	"path"

	"github.com/NoUseFreak/go-parallel"
	"github.com/sirupsen/logrus"
)

// Processor is an interface that performs an action on a repository.
// Collected data can be stored on repo.Data.
type Processor func(Repo) Repo

// Process triggers every Processor. It does this in multiple threads.
func Process(repos []Repo, threads int, actions []Processor) []Repo {
	logrus.Debugf("Processing with %d threads", threads)
	input := parallel.Input{}
	for _, repo := range repos {
		input = append(input, repo)
	}
	p := parallel.Processor{Threads: threads}
	result := p.Process(input, func(i interface{}) interface{} {
		repo, err := clone(i.(Repo))
		if err != nil {
			return repo
		}
		logrus.Tracef("Processing %s", repo.URL)
		for _, a := range actions {
			repo = a(repo)
		}
		return repo
	})

	res := []Repo{}
	for _, r := range result {
		res = append(res, r.(Repo))
	}

	return res
}

func clone(repo Repo) (Repo, error) {
	repo.TmpDir = getTmpDir(repo.URL)

	if _, err := os.Stat(repo.TmpDir); os.IsNotExist(err) {
		logrus.Tracef("Cloning %s", repo.URL)
		cloneCmd := exec.Command("git", "clone", repo.URL, repo.TmpDir)
		err := cloneCmd.Run()
		if err != nil {
			logrus.Errorf("Failed to clone %s (%s)", repo.URL, err)
			return repo, err
		}
	} else {
		logrus.Tracef("Fetching %s", repo.URL)
		fetchCmd := exec.Command("git", "fetch", "-p")
		fetchCmd.Dir = repo.TmpDir
		err := fetchCmd.Run()
		if err != nil {
			logrus.Errorf("Failed to fetch %s (%s)", repo.URL, err)
			return repo, err
		}
	}

	return repo, nil
}

func getTmpDir(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return path.Join(os.TempDir(), "gitcleanup", hash)
}
