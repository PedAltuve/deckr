package cmd

import (
	"fmt"

	"github.com/pedAltuve/deckr/internal/tools"
	"github.com/spf13/cobra"
)

func newSwitchCmd(svc *tools.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "switch <tool> <deck>",
		Short:   "Switch the active deck for a tool",
		Example: `deckr switch nvim lazy`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := svc.Switch(cmd.Context(), tools.SwitchInput{
				Tool:   args[0],
				ToDeck: args[1],
			})

			if err != nil {
				return err
			}

			fmt.Fprintf(
				cmd.OutOrStdout(),
				"switched deck\nactive deck: %s\n",
				result.ActiveDeck,
			)
			return nil

		},
	}
	return cmd
}
