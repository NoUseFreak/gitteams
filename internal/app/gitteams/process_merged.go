package gitteams

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
)

func init() {
	p := new(MergedProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "merged",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

type MergedProcessor struct{}

func (p *MergedProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "merged",
		Name:      "Merged",
		Sort:      table.DscNumeric,
		ValueType: "int32",
		Value:     func(r *Repo) interface{} { return r.Data["merged"] },
	}
}
func (p *MergedProcessor) Process(repo Repo) Repo {
	if len(repo.Branches) == 0 {
		repo = new(BranchProcessor).Process(repo)
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
