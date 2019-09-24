package gitteams

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
)

func init() {
	p := new(TagCountProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "tags",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

// TagCountProcessor counts all commits in a repository.
type TagCountProcessor struct{}

// GetReportColumn defines how the information should be shown in the report.
func (p *TagCountProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "tags",
		Name:      "tag count",
		Sort:      table.DscNumeric,
		ValueType: "int",
		Value:     func(r *Repo) interface{} { return r.Data["tagcount"] },
	}
}

// Process collects the commit count.
func (p *TagCountProcessor) Process(repo Repo) Repo {
	cmd := exec.Command("bash", "-c", "git tag | wc -l")
	cmd.Dir = repo.TmpDir
	var outb bytes.Buffer
	cmd.Stderr = &outb
	cmd.Stdout = &outb
	cmd.Run()

	repo.Data["tagcount"] = 0
	if cmd.ProcessState.ExitCode() != 0 {
		logrus.Warnf("Failed to count tags for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

	out := strings.TrimSpace(strings.TrimSuffix(outb.String(), "\n"))
	if number, err := strconv.Atoi(out); err == nil {
		repo.Data["tagcount"] = number
	}

	return repo
}
