# deckr

`deckr` is a Git-backed config deck manager for local tools like Neovim, tmux, Ghostty, and similar apps.

The idea is simple: each tool can have multiple named decks, and each deck represents a different config setup. `deckr` uses Git worktrees by default so each deck can be versioned and backed up without forcing branch checkouts every time you switch.

## Current command set

```bash
deckr init <tool> <target-path>
deckr create <tool> <deck> [--from <deck>]
deckr switch <tool> <deck>
deckr delete <tool> <deck>
deckr current <tool>
deckr list [tool]
deckr push <tool> [deck]
deckr pull <tool> [deck]
