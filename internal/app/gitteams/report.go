package gitteams

import (
	"fmt"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.PersistentFlags().StringP("sort", "s", "", "Sort by column (name, branch, loc)")
	rootCmd.PersistentFlags().StringP("format", "o", "table", "Output format (table, html, csv)")
	viper.BindPFlag("sort", rootCmd.PersistentFlags().Lookup("sort"))
	viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
}

var repositoryColumn = ReportColumn{
	ID:    "name",
	Name:  "Repo",
	Sort:  table.Asc,
	Value: func(r *Repo) interface{} { return fmt.Sprintf("%s:%s", r.Origin.Short, r.Name) },
}

type ReportOptions struct {
	Sort    string
	Format  string
	Columns []ReportColumn
}

type ReportColumn struct {
	ID        string
	Name      string
	Sort      table.SortMode
	ValueType string
	Value     func(*Repo) interface{}
}

type ReportModel struct {
	Headers []string
	Data    [][]interface{}
}

func Report(repos []Repo, options *ReportOptions) {
	options = applyCommandArgs(options)
	model := buildModel(repos, options.Columns)

	tw := table.NewWriter()
	headerRow := table.Row{}
	for _, v := range model.Headers {
		headerRow = append(headerRow, v)
	}
	tw.AppendHeader(headerRow)

	for _, line := range model.Data {
		row := table.Row{}
		for _, v := range line {
			row = append(row, v)
		}
		tw.AppendRow(row)
	}

	var sortColumn ReportColumn
	for _, c := range options.Columns {
		if c.ID == options.Sort {
			sortColumn = c
		}
	}
	tw.SortBy([]table.SortBy{{
		Name: sortColumn.Name,
		Mode: sortColumn.Sort,
	}})

	switch options.Format {
	case "csv":
		fmt.Println(tw.RenderCSV())
	case "html":
		fmt.Println(tw.RenderHTML())
	default:
		tw.SetStyle(table.StyleLight)
		tw.Style().Options.SeparateColumns = false
		fmt.Println(tw.Render())
	}
}

func applyCommandArgs(opts *ReportOptions) *ReportOptions {
	if sort := viper.GetString("sort"); sort != "" {
		opts.Sort = sort
	}
	if format := viper.GetString("format"); format != "" {
		opts.Format = format
	}
	return opts
}

func buildModel(repos []Repo, columns []ReportColumn) ReportModel {
	model := ReportModel{
		Headers: []string{},
		Data:    [][]interface{}{},
	}
	for _, v := range columns {
		model.Headers = append(model.Headers, v.Name)
	}
	for _, repo := range repos {
		row := []interface{}{}
		for _, v := range columns {
			row = append(row, v.Value(&repo))
		}
		model.Data = append(model.Data, row)
	}

	return model
}
