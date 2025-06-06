package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/logging"
	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

var (
	Settings       model.Settings
	SettingsConfig string
	OutputFormat   string
	Verbose        bool
)

func init() {
	Settings = *model.NewSettings()
	
	// Initialize viper
	viper.SetEnvPrefix("PAC")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	
	// Add the global flags to the root command
	RootCmd.PersistentFlags().StringVar(&SettingsConfig, "settings", "", "Path to the settings configuration file")
	RootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "", "Output format for configuration (yaml|json)")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "V", false, "Enable verbose output for debugging")
	
	// VCS flags
	RootCmd.PersistentFlags().String("repo", "", "Path or URL to the git repository")
	RootCmd.PersistentFlags().String("gh-token", "", "GitHub token to use for authentication")
	
	// General flags
	RootCmd.PersistentFlags().Bool("strict", false, "Enable strict mode")
	
	// Template flags
	RootCmd.PersistentFlags().StringSlice("template", []string{}, "Template location (can be used multiple times)")
	RootCmd.PersistentFlags().StringSlice("template-output", []string{}, "Template output file (must match number of templates)")
	
	// Properties flags - multiple ways to specify
	RootCmd.PersistentFlags().StringToString("property", map[string]string{}, "Properties to pass to templates (key=value)")
	RootCmd.PersistentFlags().StringToString("set", map[string]string{}, "Set property values (alias for --property)")
	
	// Task system flags
	RootCmd.PersistentFlags().StringSlice("task-system", []string{}, "Task system name (e.g., jira, github, none)")
	RootCmd.PersistentFlags().StringToString("task-regex", map[string]string{}, "Task regex patterns (name=pattern, can be used multiple times per system)")
	RootCmd.PersistentFlags().StringToString("task-query", map[string]string{}, "Task query strings (name=query)")
	RootCmd.PersistentFlags().String("task-system-json", "", "Task systems configuration as JSON")
	
	// Bind flags to viper
	viper.BindPFlag("vcs.repo", RootCmd.PersistentFlags().Lookup("repo"))
	viper.BindPFlag("vcs.token", RootCmd.PersistentFlags().Lookup("gh-token"))
	viper.BindPFlag("general.strict", RootCmd.PersistentFlags().Lookup("strict"))
	viper.BindPFlag("templates", RootCmd.PersistentFlags().Lookup("template"))

	// Add root command flags
	RootCmd.Flags().BoolP("version", "v", false, "Display the version of the CLI application")
	
	// Add commands to the root command
	RootCmd.AddCommand(FromCmd)
	RootCmd.AddCommand(ConfigCmd)
	RootCmd.AddCommand(VersionCmd)
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
		// If output format is specified, generate config
		if OutputFormat != "" {
			ConfigCmd.Run(cmd, args)
			return
		}
		fmt.Println("Please use the help command to see available options.")
	},

	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	// Set up a PersistentPreRun function to run after flags are parsed but before any command
	RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Enable verbose logging if flag is set
		logging.SetVerbose(Verbose)
		logging.VerboseSection("PAC Starting")
		logging.Verbose("Command: %s", cmd.Name())
		
		// Skip for config command
		if cmd.Name() == "config" {
			return
		}
		
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			VersionCmd.Run(cmd, args)
			os.Exit(0)
		}

		// Load settings from file if specified
		if SettingsConfig != "" {
			logging.Verbose("Loading settings from: %s", SettingsConfig)
			viper.SetConfigFile(SettingsConfig)
			if err := viper.ReadInConfig(); err != nil {
				fmt.Printf("Error reading config file: %v\n", err)
				os.Exit(1)
			}
			logging.Verbose("Settings file loaded successfully")
		}

		// Merge settings from viper
		logging.Verbose("Merging settings from configuration")
		if err := viper.Unmarshal(&Settings); err != nil {
			fmt.Printf("Error merging settings: %v\n", err)
			os.Exit(1)
		}

		// Apply command line flags using settings builder
		logging.Verbose("Applying command line flags to settings")
		builder := NewSettingsBuilder(cmd)
		builder.ApplyFlagsToSettings(&Settings)
		
		// Handle environment variable for GitHub token if not set
		if Settings.VCS.Token == "" {
			if token := os.Getenv("GITHUB_TOKEN"); token != "" {
				logging.Verbose("Using GitHub token from environment variable")
				Settings.VCS.Token = token
			}
		}
		
		logging.Verbose("Settings initialized successfully")
	}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
