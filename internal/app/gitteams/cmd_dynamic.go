package gitteams

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DynamicCommand represents an command that can be registered.
type DynamicCommand struct {
	Name         string
	Short        string
	Long         string
	Processor    Processor
	ReportColumn ReportColumn
}

var commands = []DynamicCommand{}

// CreateDynamicCommands registers any DynamicCommand instances added in their
// init function.
func CreateDynamicCommands() {
	for _, c := range commands {
		createCommand(c)
	}
}

func createCommand(c DynamicCommand) {
	rootCmd.AddCommand(&cobra.Command{
		Use:   c.Name,
		Short: c.Short,
		Long:  c.Long,
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Info("Collecting repos")
			repos := CollectRepos()

			logrus.Info("Processing")
			result := Process(repos, 50, []Processor{
				c.Processor,
			})

			logrus.Info("Report")
			Report(result, &ReportOptions{
				Sort: c.ReportColumn.ID,
				Columns: []ReportColumn{
					repositoryColumn,
					c.ReportColumn,
				},
			})
		},
	})
}
