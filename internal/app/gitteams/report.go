package gitteams

import (
	"fmt"
	"math"
	"sort"

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
	ID:     "name",
	Name:   "Repo",
	Weight: -100,
	Sort:   table.Asc,
	Value:  func(r *Repo) interface{} { return fmt.Sprintf("%s:%s", r.Origin.Short, r.Name) },
}

// ReportOptions represents how the reports should be shown.
type ReportOptions struct {
	Sort    string
	Format  string
	Columns []ReportColumn
}

// ReportModel is an intermediate data object to allow more dynamic reporting.
type ReportModel struct {
	Headers []string
	Data    [][]interface{}
	Totals  []interface{}
}

// Report generates the actual report.
func Report(repos []Repo, options *ReportOptions) {
	options = applyCommandArgs(options)

	sort.Sort(byWeight(options.Columns))
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

	footerRow := table.Row{}
	for _, v := range model.Totals {
		switch v.(type) {
		case float64, float32:
			footerRow = append(footerRow, fmt.Sprintf("%.0f", v))
		default:
			footerRow = append(footerRow, v)
		}
	}
	tw.AppendFooter(footerRow)

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
		Totals:  []interface{}{},
	}
	for _, v := range columns {
		model.Headers = append(model.Headers, v.Name)
		if v.ValueType == "string" {
			model.Totals = append(model.Totals, "")
		} else {
			model.Totals = append(model.Totals, 0)
		}
	}
	model.Totals[0] = "Total"

	for _, repo := range repos {
		row := []interface{}{}
		for i, v := range columns {
			val := v.Value(&repo)

			if val == nil {
				row = append(row, 0)
			} else {
				row = append(row, v.Value(&repo))
			}
			totNum, err := findNumber(model.Totals[i])
			if err == nil {
				valNum, _ := findNumber(v.Value(&repo))
				model.Totals[i] = totNum + valNum
			}
		}
		model.Data = append(model.Data, row)
	}

	return model
}

func findNumber(num interface{}) (float64, error) {
	switch i := num.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int:
		return float64(i), nil
	case int32:
		return float64(i), nil
	default:
		return math.NaN(), fmt.Errorf("getFloat: unknown value is of incompatible type %T", num)
	}
}
