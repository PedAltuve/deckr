package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelpShowsCommandSections(t *testing.T) {
	rootCmd := NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	output := out.String()
	expected := []string{
		"deckr",
		"Available Commands:",
		"init        Initialize a managed tool",
		"current     Show the current deck for a tool",
		"create      Create a new deck for a tool",
		"Additional help topics:",
		"deckr switch",
		"deckr delete",
		"deckr list",
		"deckr push",
		"deckr pull",
	}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Errorf("expected help output to contain %q\noutput:\n%s", want, output)
		}
	}
}

func TestInitHelpShowsUsage(t *testing.T) {
	rootCmd := NewRootCmd()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"init", "--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	output := out.String()
	expected := []string{
		"Initialize a tool so deckr can manage its config decks.",
		"Usage:",
		"deckr init <tool> <target-path> [flags]",
		"Examples:",
		"deckr init nvim ~/.config/nvim",
		"deckr init tmux ~/.config/tmux",
		"Flags:",
		"-h, --help",
	}
	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Errorf("expected help output to contain %q\noutput:\n%s", want, output)
		}
	}
}
