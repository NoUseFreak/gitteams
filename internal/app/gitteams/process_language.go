package gitteams

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
	"github.com/src-d/enry/v2"
)

func init() {
	p := new(LanguageProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "language",
		Short:        "Show main language in repository",
		Long:         "Show main language in each repository",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

// LanguageProcessor detects what languages are used in a repository.
type LanguageProcessor struct{}

// GetReportColumn defines how the information should be shown in the report.
func (p *LanguageProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "language",
		Name:      "Language",
		Sort:      table.Asc,
		ValueType: "string",
		Value:     func(r *Repo) interface{} { return r.Data["language"] },
	}
}

// Process collects the languages used in the repository.
func (p *LanguageProcessor) Process(repo Repo) Repo {
	out, err := repoExec(repo, "git", "ls-tree", "-r", "HEAD", "--name-only")
	if err != nil {
		logrus.Warnf("Failed to fetch languages for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

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
	repo.Data["language"] = p.formatLanguageResult(languages)

	return repo
}

func (p *LanguageProcessor) formatLanguageResult(data map[string]int) string {
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
