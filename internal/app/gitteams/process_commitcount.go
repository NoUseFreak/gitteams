package gitteams

import (
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
)

func init() {
	p := new(CommitCountProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "commits",
		Short:        "Count commits in repository",
		Long:         "Count the numbers of commits in the repository.",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

// CommitCountProcessor counts all commits in a repository.
type CommitCountProcessor struct{}

// GetReportColumn defines how the information should be shown in the report.
func (p *CommitCountProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "commits",
		Name:      "commit count",
		Sort:      table.DscNumeric,
		ValueType: "int",
		Value:     func(r *Repo) interface{} { return r.Data["commitcount"] },
	}
}

// Process collects the commit count.
func (p *CommitCountProcessor) Process(repo Repo) Repo {
	repo.Data["commitcount"] = 0

	out, err := repoExec(repo, "git", "rev-list", "--all", "--count")
	if err != nil {
		logrus.Warnf("Failed to count commits for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

	if number, err := strconv.Atoi(out); err == nil {
		repo.Data["commitcount"] = number
	}

	return repo
}
