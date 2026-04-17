/*
Copyright © 2026 Pedro Altuve <pedaltuve@protonmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/pedAltuve/deckr/internal/tools"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func NewRootCmd() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "deckr",
		Short: "Git-backed config deck manager",
		Long:  "deckr manages config decks for tools like Neovim, tmux, Ghostty and so on",
	}
	dataDir, err := deckrDataDir()
	if err != nil {
		return nil, fmt.Errorf("resolve deckr data directory: %w", err)
	}
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
	return rootCmd, nil
}

func deckrDataDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config directory: %w", err)
	}
	return filepath.Join(configDir, "deckr"), nil
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
