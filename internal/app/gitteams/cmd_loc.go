package gitteams

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var locCmd = &cobra.Command{
	Use:   "loc",
	Short: "Get LOC count",
	Long:  `Get Lines of Code count`,
	Run:   executeLoc,
}

func init() {
	rootCmd.AddCommand(locCmd)
}

func executeLoc(cmd *cobra.Command, args []string) {
	logrus.Info("Collecting repos")
	repos := CollectRepos()

	logrus.Info("Processing")
	result := Process(repos, 50, []Processor{
		CountLoc,
	})

	logrus.Info("Report")
	Report(result, &ReportOptions{
		Sort: "loc",
		Columns: []ReportColumn{
			repositoryColumn,
			locColumn,
		},
	})
}
