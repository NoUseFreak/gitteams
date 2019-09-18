package gitteams

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var branchesCmd = &cobra.Command{
	Use:   "branches",
	Short: "Count number of branches",
	Long:  `Count number of branches`,
	Run:   executeBranches,
}

func init() {
	rootCmd.AddCommand(branchesCmd)
}

func executeBranches(cmd *cobra.Command, args []string) {
	logrus.Info("Collecting repos")
	repos := CollectRepos()

	logrus.Info("Processing")
	result := Process(repos, 50, []Processor{
		GetBranches,
	})

	logrus.Info("Report")
	Report(result, &ReportOptions{
		Sort: "branches",
		Columns: []ReportColumn{
			repositoryColumn,
			branchCountColumn,
		},
	})
}
