package cmd

import (
	"fmt"

	"github.com/pedAltuve/deckr/internal/tools"
	"github.com/spf13/cobra"
)

func newCurrentCmd(svc *tools.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "current <tool>",
		Short:   "Show the current deck for a tool",
		Example: "deckr current nvim",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := svc.Current(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			fmt.Fprintf(
				cmd.OutOrStdout(),
				"current deck for %s: %s\n",
				args[0],
				result,
			)

			return nil
		},
	}

	return cmd
}
