package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// NewCmdCompletion returns the completion command
func NewCmdCompletion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

	$ source <(tfbackend completion bash)

	# To load completions for each session, execute once:
	# Linux:
	$ tfbackend completion bash > /etc/bash_completion.d/tfbackend
	# macOS:
	$ tfbackend completion bash > /usr/local/etc/bash_completion.d/tfbackend

Zsh:

	# If shell completion is not already enabled in your environment,
	# you will need to enable it.  You can execute the following once:

	$ echo "autoload -U compinit; compinit" >> ~/.zshrc

	# To load completions for each session, execute once:
	$ tfbackend completion zsh > "${fpath[1]}/_tfbackend"

	# You will need to start a new shell for this setup to take effect.

fish:

	$ tfbackend completion fish | source

	# To load completions for each session, execute once:
	$ tfbackend completion fish > ~/.config/fish/completions/tfbackend.fish

PowerShell:

	PS> tfbackend completion powershell | Out-String | Invoke-Expression

	# To load completions for every new session, run:
	PS> tfbackend completion powershell > tfbackend.ps1
	# and source this file from your PowerShell profile.
`,
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}

	return cmd
}
