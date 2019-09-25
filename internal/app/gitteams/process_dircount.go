package gitteams

import (
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
)

func init() {
	p := new(DirCountProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "dirs",
		Short:        "Get directory count in repositories",
		Long:         "Count directories in each repository",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

// DirCountProcessor counts the lines of code in a repository.
type DirCountProcessor struct{}

// GetReportColumn defines how the information should be shown in the report.
func (p *DirCountProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "dirs",
		Name:      "Directories",
		Sort:      table.DscNumeric,
		ValueType: "int",
		Value:     func(r *Repo) interface{} { return r.Data["directories"] },
	}
}

// Process counts the lines of code in a repository.
func (p *DirCountProcessor) Process(repo Repo) Repo {
	repo.Data["directories"] = 0

	out, err := repoExec(repo, "bash", "-c", "git ls-tree -d -r HEAD --name-only | wc -l")
	if err != nil {
		logrus.Warnf("Failed to count directories for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

	if number, err := strconv.Atoi(out); err == nil {
		repo.Data["directories"] = number
	}

	return repo
}
