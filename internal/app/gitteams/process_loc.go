package gitteams

import (
	"fmt"

	"github.com/hhatto/gocloc"
	"github.com/jedib0t/go-pretty/table"
)

var locColumn = ReportColumn{
	ID:    "loc",
	Name:  "Lines of code",
	Sort:  table.DscNumeric,
	Value: func(r *Repo) interface{} { return r.Data["loc"] },
}

func CountLoc(repo Repo) Repo {
	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()
	paths := []string{repo.TmpDir}

	processor := gocloc.NewProcessor(languages, options)
	result, err := processor.Analyze(paths)
	if err != nil {
		fmt.Printf("gocloc fail. error: %v\n", err)
		return repo
	}

	repo.Data["loc"] = result.Total.Total

	return repo
}
