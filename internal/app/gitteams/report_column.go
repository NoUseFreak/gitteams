package gitteams

import "github.com/jedib0t/go-pretty/table"

// ReportColumn defines how the column should be handled in the report.
type ReportColumn struct {
	ID        string
	Name      string
	Weight    int
	Sort      table.SortMode
	ValueType string
	Value     func(*Repo) interface{}
}

type byWeight []ReportColumn

func (s byWeight) Len() int {
	return len(s)
}
func (s byWeight) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byWeight) Less(i, j int) bool {
	return s[i].Weight < s[j].Weight
}
