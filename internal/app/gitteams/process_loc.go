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
		Short:        "Get LOC count in repositories",
		Long:         "Get lines of code in each repository",
		Processor:    p.Process,
		ReportColumn: p.GetReportColumn(),
	})
}

// LOCProcessor counts the lines of code in a repository.
type LOCProcessor struct{}

// GetReportColumn defines how the information should be shown in the report.
func (p *LOCProcessor) GetReportColumn() ReportColumn {
	return ReportColumn{
		ID:        "loc",
		Name:      "Lines of code",
		Weight:    20,
		Sort:      table.DscNumeric,
		ValueType: "int32",
		Value:     func(r *Repo) interface{} { return r.Data["loc"] },
	}
}

// Process counts the lines of code in a repository.
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
