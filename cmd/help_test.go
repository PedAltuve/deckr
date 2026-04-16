package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelpShowsAvailableCommands(t *testing.T) {
	rootCmd := NewRootCmd()

	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := out.String()

	expected := []string{
		"deckr",
		"Available Commands:",
		"init",
		"create",
		"switch",
		"delete",
		"current",
		"list",
		"push",
		"pull",
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

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := out.String()

	expected := []string{
		"Initialize a managed tool",
		"deckr init",
	}

	for _, want := range expected {
		if !strings.Contains(output, want) {
			t.Errorf("expected help output to contain %q\noutput:\n%s", want, output)
		}
	}
}
