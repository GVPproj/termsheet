# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a terminal-based invoice management application called "invogo" built in Go using the Bubble Tea TUI framework. The application provides a menu-driven interface for managing providers, clients, and invoices.

## Architecture

The codebase follows a simple modular structure:

- `main.go`: Contains the main Bubble Tea application model, view logic, and event handling
- `types/types.go`: Defines the View enum constants used for navigation between different screens

The application uses the Model-View-Update (MVU) pattern via Bubble Tea:
- Model: Holds application state (current view, cursor position, menu choices)
- Update: Handles user input and state transitions
- View: Renders the current screen based on application state

## Common Commands

### Build and Run
```bash
go run main.go
```

### Build Binary
```bash
go build -o invogo
```

### Install Dependencies
```bash
go mod tidy
```

### Test
```bash
go test ./...
```

## Development Notes

- The application currently has placeholder implementations for Providers view, while Clients and Invoices views show "not implemented yet" messages
- Navigation uses vi-style keybindings (j/k for up/down) plus arrow keys
- ESC returns to the main menu from any view
- The window title is set to "invogo" on startup