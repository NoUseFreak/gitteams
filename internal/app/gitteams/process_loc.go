package gitteams

import (
	"fmt"

	"github.com/hhatto/gocloc"
	"github.com/jedib0t/go-pretty/table"
)

func init() {
	p := new(LOCProcessor)
	commands = append(commands, DynamicCommand{
		Name:         "loc",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

type LOCProcessor struct{}

func (p *LOCProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "loc",
		Name:      "Lines of code",
		Sort:      table.DscNumeric,
		ValueType: "int32",
		Value:     func(r *Repo) interface{} { return r.Data["loc"] },
	}
}

func (p *LOCProcessor) Process(repo Repo) Repo {
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
