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

type Processor func(Repo) Repo

func Process(repos []Repo, threads int, actions []Processor) []Repo {
	logrus.Debugf("Processing with %d threads", threads)
	input := parallel.Input{}
	for _, repo := range repos {
		input = append(input, repo)
	}
	p := parallel.Processor{Threads: threads}
	result := p.Process(input, func(i interface{}) interface{} {
		repo := clone(i.(Repo))
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

func clone(repo Repo) Repo {
	repo.TmpDir = getTmpDir(repo.URL)

	if _, err := os.Stat(repo.TmpDir); os.IsNotExist(err) {
		logrus.Tracef("Cloning %s", repo.URL)
		cloneCmd := exec.Command("git", "clone", repo.URL, repo.TmpDir)
		err := cloneCmd.Run()
		if err != nil {
			panic(err)
		}
	} else {
		logrus.Tracef("Fetching %s", repo.URL)
		fetchCmd := exec.Command("git", "fetch", "-p")
		fetchCmd.Dir = repo.TmpDir
		fetchCmd.Run()
	}

	return repo
}

func getTmpDir(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	hash := hex.EncodeToString(hasher.Sum(nil))

	return path.Join(os.TempDir(), "gitcleanup", hash)
}
