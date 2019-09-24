package gitteams

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/sirupsen/logrus"
)

func init() {
	p := new(MergedProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "merged",
		Short:        "Count merged branches",
		Long:         "Count the number of branches merged into the main branch",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

// MergedProcessor retrieves a count for the number of branches fully merged
// into the main branch.
type MergedProcessor struct{}

// GetReportColumn defines how the information should be shown in the report.
func (p *MergedProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "merged",
		Name:      "Merged",
		Sort:      table.DscNumeric,
		ValueType: "int32",
		Value:     func(r *Repo) interface{} { return r.Data["merged"] },
	}
}

// Process collects the branch count.
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

		out, err := repoExec(repo, "git", "rev-list",
			fmt.Sprintf("origin/%s..origin/%s", repo.MainBranch, branch),
			"--count")
		if err != nil {
			logrus.Warnf("Failed to fetch merged branches for %s in %s", repo.URL, repo.TmpDir)
			return repo
		}

		if strings.TrimSpace(out) == "0" {
			logrus.Debugf("Found fully merged branch - %s", branch)
			repo.Data["merged"] = repo.Data["merged"].(int32) + 1
		}
	}

	return repo
}
