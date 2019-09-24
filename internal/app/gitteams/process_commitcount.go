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
	p := new(CommitCountProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "commits",
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
	cmd := exec.Command("git", "rev-list", "--all", "--count")
	cmd.Dir = repo.TmpDir
	var outb bytes.Buffer
	cmd.Stderr = &outb
	cmd.Stdout = &outb
	cmd.Run()

	if cmd.ProcessState.ExitCode() != 0 {
		logrus.Warnf("Failed to count commits for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

	out := strings.TrimSuffix(outb.String(), "\n")
	if number, err := strconv.Atoi(out); err == nil {
		repo.Data["commitcount"] = number
	}

	return repo
}
