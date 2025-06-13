package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/git"
	"github.com/Praqma/Praqmatic-Automated-Changelog/src/logging"
	"github.com/Praqma/Praqmatic-Automated-Changelog/src/report"
	"github.com/Praqma/Praqmatic-Automated-Changelog/src/task"
)

var fromHelpShortText = "Specify a git commit SHA or tag name to state the starting point for the changelog generation."

var fromHelpLongText = "Generate changelog from <oldest-ref> to an optional <newest-ref>. For git this takes anything that rev-parse accepts, such as HEAD~3, Git SHA, or tag name."

var (
	fromSha string
	toSha   string
)

func init() {
	FromCmd.Flags().BoolP("analyze", "a", false, "Analyze the task Json file output and print the results to stdout.")
	FromCmd.Flags().StringVar(&toSha, "to", "", "Specify where to start searching for commit <newest-ref>. For git this takes anything that rev-parse accepts. Such as HEAD / Git sha or tag name.")
}

var FromCmd = &cobra.Command{
	Use:   "from <oldest-ref> [--to=<newest-ref>]",
	Short: fromHelpShortText,
	Long:  fromHelpLongText,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fromSha = args[0]

		logging.VerboseSection("Changelog Generation")
		logging.Verbose("Repository: %s", Settings.VCS.Repo)
		logging.Verbose("From SHA: %s", fromSha)

		if toSha != "" {
			logging.Verbose("To SHA: %s", toSha)
		} else {
			toSha = "HEAD"
			logging.Verbose("To SHA: HEAD (default)")
		}
		fmt.Println()

		logging.VerboseSection("Git VCS Initialization")
		gitVCS, err := git.NewGitVCS(Settings.VCS)
		if err != nil {
			fmt.Printf("Error initializing Git VCS: %v\n", err)
			return
		}
		logging.Verbose("Git VCS initialized successfully")

		logging.VerboseSection("Fetching Commits")
		logging.Verbose("Getting commits between %s and %s", fromSha, toSha)
		commits, err := gitVCS.GetCommitsBetween(fromSha, toSha)
		if err != nil {
			fmt.Printf("Error getting commits: %v\n", err)
			return
		}
		logging.Verbose("Found %d commits", commits.Count())

		logging.VerboseSection("Processing Tasks")
		taskCollection, err := task.TaskIDList(Settings, commits)
		if err != nil {
			fmt.Printf("Error processing tasks: %v\n", err)
			return
		}
		logging.Verbose("Processed %d tasks", len(taskCollection.Tasks))
		logging.Verbose("Found %d unreferenced commits", len(taskCollection.GetUnreferencedCommits()))

		if cmd.Flags().Changed("analyze") {
			logging.Verbose("Analyze mode - outputting JSON")
			fmt.Println(string(taskCollection.ToJSON()))
			return
		}

		logging.VerboseSection("Report Generation")
		logging.Verbose("Number of templates configured: %d", len(Settings.Templates))

		generator := report.NewGenerator(taskCollection)
		if err := generator.Generate(&Settings); err != nil {
			fmt.Printf("Error generating report: %v\n", err)
			return
		}

		logging.Verbose("Report generation completed successfully")
	},
}
