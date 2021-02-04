package gitteams

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get team project stats",
	Long:  `Get team project stats`,
	Run:   executeStats,
}

func init() {
	statsCmd.Flags().String("columns", "", "column names seperated by comma")
	rootCmd.AddCommand(statsCmd)
}

func executeStats(cmd *cobra.Command, args []string) {
	processors := []Processor{}
	reportColumns := []ReportColumn{repositoryColumn}
	if columnString, err := cmd.Flags().GetString("columns"); err == nil && columnString != "" {
		columns := strings.Split(columnString, ",")
		for _, c := range commands {
			if idx := getStringSliceIndex(c.Name, columns); idx >= 0 {
				processors = append(processors, c.Processor)
				c.ReportColumn.Weight = idx
				reportColumns = append(reportColumns, c.ReportColumn)
			}
		}

	} else {
		for _, c := range commands {
			processors = append(processors, c.Processor)
			reportColumns = append(reportColumns, c.ReportColumn)
		}
	}

	logrus.Info("Collecting repos")
	repos := CollectRepos()

	logrus.Info("Processing")
	result := Process(repos, runtime.NumCPU(), processors)

	logrus.Info("Report")
	Report(result, &ReportOptions{
		Sort:    "name",
		Columns: reportColumns,
	})
}

func getStringSliceIndex(v string, s []string) int {
	for si, sv := range s {
		if sv == v {
			return si
		}
	}

	return -1
}
