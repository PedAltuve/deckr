package cmd

import (
	"fmt"

	"github.com/pedAltuve/deckr/internal/tools"
	"github.com/spf13/cobra"
)

func newCreateCmd(svc *tools.Service) *cobra.Command {
	var fromDeck string
	var empty bool
	cmd := &cobra.Command{
		Use:     "create <tool> <deck>",
		Short:   "Create a new deck for a tool",
		Example: `deckr create nvim new`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := svc.Create(cmd.Context(), tools.CreateInput{
				Tool:     args[0],
				NewDeck:  args[1],
				FromDeck: fromDeck,
				Empty:    empty,
			})

			if err != nil {
				return err
			}

			fmt.Fprintf(
				cmd.OutOrStdout(),
				"created new deck %s\nactive deck: %s\n",
				result.Deck,
				result.ActiveDeck,
			)
			return nil
		},
	}
	cmd.Flags().StringVar(&fromDeck, "from", "", "source deck to clone")
	cmd.Flags().BoolVar(&empty, "empty", false, "create an empty deck")
	cmd.MarkFlagsMutuallyExclusive("from", "empty")

	return cmd
}
