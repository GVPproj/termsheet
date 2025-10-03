# termsheet

A TerminalUI application for managing invoices locally.

## sqlite3 db

Running this app generates/reads a local db to store invoicing data, clients etc.

```bash
# List tables
sqlite3 termsheet.db ".tables"

# Show all schemas
sqlite3 termsheet.db ".schema"

# Query with nice formatting
sqlite3 termsheet.db -column -header "SELECT * FROM provider"

# Export to CSV
sqlite3 termsheet.db -csv "SELECT * FROM provider" > providers.csv

# Markdown format
sqlite3 termsheet.db -markdown "SELECT * FROM provider"
```

## Dev Resources

Theming Huh:
<https://raw.githubusercontent.com/charmbracelet/huh/refs/heads/main/theme.go>
