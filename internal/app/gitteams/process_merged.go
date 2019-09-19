package gitteams

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
)

var mergedColumn = ReportColumn{
	ID:        "merged",
	Name:      "Merged",
	Sort:      table.DscNumeric,
	ValueType: "int32",
	Value:     func(r *Repo) interface{} { return r.Data["merged"] },
}

func GetMerged(repo Repo) Repo {
	if len(repo.Branches) == 0 {
		repo = GetBranches(repo)
	}

	repo.Data["merged"] = int32(0)
	for _, branch := range repo.Branches {
		switch branch {
		case repo.MainBranch, "develop", "master":
			continue
		}
		cmd := exec.Command("git", "rev-list",
			fmt.Sprintf("origin/%s..origin/%s", repo.MainBranch, branch),
			"--count",
		)
		cmd.Dir = repo.TmpDir
		var outb bytes.Buffer
		cmd.Stderr = &outb
		cmd.Stdout = &outb
		cmd.Run()

		if strings.TrimSpace(outb.String()) == "0" {
			logrus.Debugf("Found fully merged branch - %s", branch)
			repo.Data["merged"] = repo.Data["merged"].(int32) + 1
		}
	}

	return repo
}
