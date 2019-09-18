package gitteams

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/jedib0t/go-pretty/table"
)

var branchCountColumn = ReportColumn{
	ID:    "branches",
	Name:  "branch count",
	Sort:  table.DscNumeric,
	Value: func(r *Repo) interface{} { return len(r.Branches) },
}

func GetBranches(repo Repo) Repo {
	cmd := exec.Command("git", "branch", "-r")
	cmd.Dir = repo.TmpDir
	var outb bytes.Buffer
	cmd.Stderr = &outb
	cmd.Stdout = &outb
	cmd.Run()

	repo.Branches = []string{}
	out := strings.TrimSuffix(outb.String(), "\n")
	for _, s := range strings.Split(out, "\n") {
		s = strings.TrimSpace(s)
		branch := strings.Split(s, " ")[0][7:]
		if branch != "HEAD" {
			repo.Branches = append(repo.Branches, branch)
		}
	}

	return repo
}
