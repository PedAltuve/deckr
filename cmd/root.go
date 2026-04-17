/*
Copyright © 2026 Pedro Altuve <pedaltuve@protonmail.com>
*/
package cmd

import (
	"path/filepath"

	"github.com/pedAltuve/deckr/internal/tools"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "deckr",
		Short: "Git-backed config deck manager",
		Long:  "deckr manages config decks for tools like Neovim, tmux, Ghostty and so on",
	}

	dataDir := filepath.Join(".", ".deckr-dev")

	svc := &tools.Service{
		Paths:    tools.OSPaths{},
		Registry: tools.NewFileRegistry(filepath.Join(dataDir, "registry")),
		Backend:  tools.NewLocalBackend(filepath.Join(dataDir, "backend")),
	}

	rootCmd.AddCommand(newInitCmd(svc))
	rootCmd.AddCommand(newCreateCmd(svc))
	rootCmd.AddCommand(newSwitchCmd(svc))
	rootCmd.AddCommand(newDeleteCmd())
	rootCmd.AddCommand(newCurrentCmd(svc))
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newPushCmd())
	rootCmd.AddCommand(newPullCmd())

	return rootCmd

}

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <tool> <deck>",
		Short: "Delete a deck",
	}
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [tool]",
		Short: "List tools or decks",
	}
}

func newPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push <tool> [deck]",
		Short: "Push a deck to its remote repository",
	}
}

func newPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull <tool> [deck]",
		Short: "Pull a deck from its remote repository",
	}
}
