package gitteams

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
)

func init() {
	p := new(BranchProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "branch",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

// BranchProcessor retrieves all branches in a repository.
type BranchProcessor struct{}

// GetReportColumn defines how the information should be shown in the report.
func (p *BranchProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "branches",
		Name:      "branch count",
		Sort:      table.DscNumeric,
		ValueType: "int",
		Value:     func(r *Repo) interface{} { return len(r.Branches) },
	}
}

// Process collects the branches.
func (p *BranchProcessor) Process(repo Repo) Repo {
	cmd := exec.Command("git", "branch", "-r")
	cmd.Dir = repo.TmpDir
	var outb bytes.Buffer
	cmd.Stderr = &outb
	cmd.Stdout = &outb
	cmd.Run()

	if cmd.ProcessState.ExitCode() != 0 {
		logrus.Warnf("Failed to fetch branched for %s in %s", repo.URL, repo.TmpDir)
		return repo
	}

	repo.Branches = []string{}
	out := strings.TrimSuffix(outb.String(), "\n")
	for _, s := range strings.Split(out, "\n") {
		s = strings.TrimSpace(s)
		parts := strings.Split(s, " ")
		if s == "" || len(parts) < 1 || len(parts[0]) < 7 {
			continue
		}
		branch := parts[0][7:]
		if branch != "HEAD" {
			repo.Branches = append(repo.Branches, branch)
		}
	}

	return repo
}
