package gitteams

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
	"github.com/src-d/enry/v2"
)

var languageColumn = ReportColumn{
	ID:        "language",
	Name:      "Language",
	Sort:      table.Asc,
	ValueType: "string",
	Value:     func(r *Repo) interface{} { return r.Data["language"] },
}

func GetLanguage(repo Repo) Repo {
	cmd := exec.Command("git", "ls-tree", "-r", "HEAD", "--name-only")
	cmd.Dir = repo.TmpDir
	var outb bytes.Buffer
	cmd.Stderr = &outb
	cmd.Stdout = &outb
	cmd.Run()

	out := strings.TrimSuffix(outb.String(), "\n")
	languages := map[string]int{}
	for _, s := range strings.Split(out, "\n") {
		p := path.Join(repo.TmpDir, s)
		dat, err := ioutil.ReadFile(p)
		if err != nil {
			logrus.Tracef("Failed to read %s, skipping", p)
			continue
		}
		lang := enry.GetLanguage(p, dat)
		if _, ok := languages[lang]; !ok {
			languages[lang] = 1
		} else {
			languages[lang]++
		}
	}
	delete(languages, "")
	repo.Data["languages"] = languages
	repo.Data["language"] = formatLanguageResult(languages)

	return repo
}

func formatLanguageResult(data map[string]int) string {
	key := "unknown"
	total := 0
	high := -1

	for k, v := range data {
		total += v
		if v > high {
			key = k
			high = v
		}
	}

	return fmt.Sprintf("%s (%2.0f%%)", key, float64(data[key])/float64(total)*100)
}
