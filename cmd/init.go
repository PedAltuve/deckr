package cmd

import (
	"fmt"

	"github.com/pedAltuve/deckr/internal/tools"
	"github.com/spf13/cobra"
)

func newInitCmd(svc *tools.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <tool> <target-path>",
		Short: "Initialize a managed tool",
		Long:  "Initialize a tool so deckr can manage its config decks.",
		Example: `deckr init nvim ~/.config/nvim
		deckr init tmux ~/.config/tmux`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := svc.Init(cmd.Context(), tools.InitInput{
				Name:       args[0],
				TargetPath: args[1],
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(
				cmd.OutOrStdout(),
				"initialized %s\nactive deck: %s\ntarget: %s\n",
				result.Name,
				result.ActiveDeck,
				result.TargetPath,
			)

			return nil
		},
	}

	return cmd
}
