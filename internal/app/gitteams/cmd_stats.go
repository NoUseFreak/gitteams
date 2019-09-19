package gitteams

import (
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
	rootCmd.AddCommand(statsCmd)
}

func executeStats(cmd *cobra.Command, args []string) {
	logrus.Info("Collecting repos")
	repos := CollectRepos()

	logrus.Info("Processing")
	result := Process(repos, 50, []Processor{
		GetBranches,
		CountLoc,
		GetMerged,
		GetLanguage,
	})

	logrus.Info("Report")
	Report(result, &ReportOptions{
		Sort: "name",
		Columns: []ReportColumn{
			repositoryColumn,
			branchCountColumn,
			locColumn,
			mergedColumn,
			languageColumn,
		},
	})
}
