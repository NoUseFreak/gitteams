package gitteams

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var languageCmd = &cobra.Command{
	Use:   "language",
	Short: "List main language",
	Long:  `Get main language.`,
	Run:   executeLanguage,
}

func init() {
	rootCmd.AddCommand(languageCmd)
}

func executeLanguage(cmd *cobra.Command, args []string) {
	logrus.Info("Collecting repos")
	repos := CollectRepos()

	logrus.Info("Processing")
	result := Process(repos, 10, []Processor{
		GetLanguage,
	})

	logrus.Info("Report")
	Report(result, &ReportOptions{
		Sort: "language",
		Columns: []ReportColumn{
			repositoryColumn,
			languageColumn,
		},
	})
}
