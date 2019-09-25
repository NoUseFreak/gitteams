package gitteams

import (
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
)

func init() {
	p := new(AuthorCountProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "author",
		Short:        "Count number of authors",
		Long:         "Count the number of authors in each repository.",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

// AuthorCountProcessor counts the authors in a repository.
type AuthorCountProcessor struct{}

// GetReportColumn defines how the information should be shown in the report.
func (p *AuthorCountProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "author",
		Name:      "author count",
		Weight:    25,
		Sort:      table.DscNumeric,
		ValueType: "int",
		Value:     func(r *Repo) interface{} { return r.Data["authorcount"] },
	}
}

// Process collects the commit count.
func (p *AuthorCountProcessor) Process(repo Repo) Repo {
	repo.Data["authorcount"] = 0

	out, err := repoExec(repo, "bash", "-c", "git shortlog -s -n --all --no-merges | wc -l")
	if err != nil {
		logrus.Warnf("Failed to count commits for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

	if number, err := strconv.Atoi(out); err == nil {
		repo.Data["authorcount"] = number
	}

	return repo
}
