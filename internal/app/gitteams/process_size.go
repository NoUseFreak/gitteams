package gitteams

import (
	"math"
	"regexp"

	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
)

var sizeRegex *regexp.Regexp

func init() {
	p := new(SizeProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "size",
		Short:        "Get repository size",
		Long:         "Get the size of the repository for each repository",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})

	sizeRegex = regexp.MustCompile("size-pack: ([0-9\\.]+ .+B)")
}

// SizeProcessor calculates the size of a repository.
type SizeProcessor struct{}

// GetReportColumn defines how the information should be shown in the report.
func (p *SizeProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "size",
		Name:      "size (kb)",
		Sort:      table.DscNumeric,
		ValueType: "float64",
		Value:     func(r *Repo) interface{} { return r.Data["size"] },
	}
}

// Process collects the commit count.
func (p *SizeProcessor) Process(repo Repo) Repo {
	repo.Data["size"] = 0

	repoExec(repo, "git", "gc")
	out, err := repoExec(repo, "git", "count-objects", "-vH")
	if err != nil {
		logrus.Warnf("Failed to calculate size for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

	if !sizeRegex.MatchString(out) {
		logrus.Warnf("Failed to fetch size for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

	matches := sizeRegex.FindStringSubmatch(out)
	if len(matches) != 2 {
		logrus.Warnf("Failed to parse size for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

	size, err := humanize.ParseBytes(matches[1])
	if err != nil {
		logrus.Warnf("Failed to translate size for %s in %s", repo.URL, repo.TmpDir)
	}

	floatSize := float64(size) / 1024
	repo.Data["size"] = math.Round(floatSize*100) / 100

	return repo
}
