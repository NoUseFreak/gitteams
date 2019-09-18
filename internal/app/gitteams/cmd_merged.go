package gitteams

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var mergedCmd = &cobra.Command{
	Use:   "merged",
	Short: "List merged branches",
	Long:  `Get all branched that have been fully merged into main branch.`,
	Run:   executeMerged,
}

func init() {
	rootCmd.AddCommand(mergedCmd)
}

func executeMerged(cmd *cobra.Command, args []string) {
	logrus.Info("Collecting repos")
	repos := CollectRepos()

	logrus.Info("Processing")
	result := Process(repos, 10, []Processor{
		GetMerged,
	})

	logrus.Info("Report")
	Report(result, &ReportOptions{
		Sort: "merged",
		Columns: []ReportColumn{
			repositoryColumn,
			mergedColumn,
		},
	})
}
