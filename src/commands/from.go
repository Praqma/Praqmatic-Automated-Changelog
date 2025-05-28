package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/git"
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

		fmt.Println("repo:", Settings.VCS.Repo)
		fmt.Printf("Generating changelog from %s", fromSha)
		if toSha != "" {
			fmt.Printf(" to %s", toSha)
		} else {
			toSha = "HEAD"
		}
		fmt.Println()

		gitVCS, err := git.NewGitVCS(Settings.VCS)
		if err != nil {
			fmt.Printf("Error initializing Git VCS: %v\n", err)
			return
		}

		commits, err := gitVCS.GetCommitsBetween(fromSha, toSha)
		if err != nil {
			fmt.Printf("Error getting commits: %v\n", err)
			return
		}

		taskCollection, err := task.TaskIDList(Settings.TaskSystems, commits)
		if err != nil {
			fmt.Printf("Error processing tasks: %v\n", err)
			return
		}

		if cmd.Flags().Changed("analyze") {
			fmt.Println(string(taskCollection.ToJSON()))
			return
		}
		fmt.Println("Generating report(s)...")

		generator := report.NewGenerator(taskCollection)
		if err := generator.Generate(&Settings); err != nil {
			fmt.Printf("Error generating report: %v\n", err)
			return
		}
	},
}
