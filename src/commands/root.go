package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

var (
	Settings       model.Settings
	SettingsConfig string
	Template       string
	Repo           string
)

func init() {
	Settings = *model.NewSettings()
	
	// Add the global flags to the root command
	RootCmd.PersistentFlags().StringVar(&Repo, "repo", "", "Path or URL to the git repository")
	RootCmd.PersistentFlags().StringVar(&SettingsConfig, "settings", "", "Path to the settings configuration file")
	RootCmd.PersistentFlags().StringVar(&Template, "template", "", "Path to the template file(s) to use for changelog generation")
	
	// Add root command flags
	RootCmd.Flags().BoolP("version", "v", false, "Display the version of the CLI application")
	
	// Add commands to the root command
	RootCmd.AddCommand(FromCmd)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of the CLI application",
	Long:  `Display the version of the CLI application`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PAC CLI version 1.0.0")
	},
}

var RootCmd = &cobra.Command{
	Use:   "go-pac",
	Short: "Praqmatic Automation Changelog (PAC) - Command Line Interface",
	Long:  `Praqmatic Automation Changelog (PAC) - Command Line Interface`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please use the help command to see available options.")
	},

	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {

	// Set up a PersistentPreRun function to run after flags are parsed but before any command
	RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			VersionCmd.Run(cmd, args)
			os.Exit(0)
		}

		if SettingsConfig != "" {
			s, err := config.GenerateSettings(SettingsConfig)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			Settings = *s
		}

		if Repo != "" {
			Settings.VCS.Repo = Repo
		}
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
