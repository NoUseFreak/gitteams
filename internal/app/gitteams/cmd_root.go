package gitteams

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "gitteams",
	Short: "Git Teams",
	Long:  `Git Teams helps you manage all project at once.`,
}

// Execute runs the cli application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	viper.SetEnvPrefix("GITTEAMS")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	logrus.SetOutput(os.Stderr)

	rootCmd.PersistentFlags().String("loglevel", "info", "Log level")
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		lvl, err := logrus.ParseLevel(viper.GetString("loglevel"))
		if err != nil {
			return err
		}
		logrus.SetLevel(lvl)
		return nil
	}
}

func setRootFlag(name, shorthand, value, usage string) {
	rootCmd.PersistentFlags().StringP(name, shorthand, value, usage)
	viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
	viper.BindEnv(name)
}

func setRootFlagBool(name, shorthand string, value bool, usage string) {
	rootCmd.PersistentFlags().BoolP(name, shorthand, value, usage)
	viper.BindPFlag(name, rootCmd.PersistentFlags().Lookup(name))
	viper.BindEnv(name)
}
