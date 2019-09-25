package gitteams

import (
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
)

func init() {
	p := new(FileCountProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "files",
		Short:        "Get file count in repositories",
		Long:         "Count files/directories in each repository",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

// FileCountProcessor counts the lines of code in a repository.
type FileCountProcessor struct{}

// GetReportColumn defines how the information should be shown in the report.
func (p *FileCountProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "files",
		Name:      "Files",
		Sort:      table.DscNumeric,
		ValueType: "int",
		Value:     func(r *Repo) interface{} { return r.Data["files"] },
	}
}

// Process counts the lines of code in a repository.
func (p *FileCountProcessor) Process(repo Repo) Repo {
	repo.Data["files"] = 0

	out, err := repoExec(repo, "bash", "-c", "git ls-files | wc -l")
	if err != nil {
		logrus.Warnf("Failed to count files for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

	if number, err := strconv.Atoi(out); err == nil {
		repo.Data["files"] = number
	}

	return repo
}
