# SQLite3 Command Reference

## Opening the Database

```bash
sqlite3 termsheet.db
```

## Essential Commands

### Formatting Output

```sql
.mode column        -- Column-aligned output
.headers on         -- Show column headers
.width 12 20 30     -- Set column widths (adjust as needed)
```

For best results, run all three commands together:

```sql
.mode column
.headers on
.width 40 30 30 30 20
```

### Other Display Modes

```sql
.mode list          -- Values delimited by separator (default |)
.mode line          -- One value per line
.mode csv           -- Comma-separated values
.mode markdown      -- Markdown table format
.mode table         -- ASCII table (default in newer sqlite3)
.mode box           -- Box-drawing characters
```

## Viewing Tables

### List All Tables

```sql
.tables
```

### Show Table Schema

```sql
.schema              -- All tables
.schema provider     -- Specific table
```

### Query Data

```sql
-- All providers
SELECT * FROM provider;

-- All clients
SELECT * FROM client;

-- All invoices with details
SELECT
    i.id,
    p.name as provider,
    c.name as client,
    i.paid,
    i.date_created
FROM invoice i
LEFT JOIN provider p ON i.provider_id = p.id
LEFT JOIN client c ON i.client_id = c.id;

-- Invoice items
SELECT * FROM invoice_item;
```

## Quick One-Liners (From Terminal)

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

## Useful Queries

### Count Records

```sql
SELECT COUNT(*) as total_providers FROM provider;
SELECT COUNT(*) as total_clients FROM client;
SELECT COUNT(*) as total_invoices FROM invoice;
```

### Recent Invoices

```sql
SELECT * FROM invoice
ORDER BY date_created DESC
LIMIT 10;
```

### Unpaid Invoices

```sql
SELECT * FROM invoice WHERE paid = 0;
```

### Invoice with Items

```sql
SELECT
    i.id,
    ii.item_name,
    ii.amount,
    ii.cost_per_unit,
    (ii.amount * ii.cost_per_unit) as total
FROM invoice i
JOIN invoice_item ii ON i.id = ii.invoice_id
WHERE i.id = 1;  -- Replace 1 with actual invoice ID
```

## Exit SQLite3

```sql
.quit
-- or --
.exit
-- or press --
Ctrl+D
```

## Pro Tips

1. **Set defaults in ~/.sqliterc**

   ```sql
   .mode column
   .headers on
   .prompt '> '
   ```

2. **Use `.width` for custom column widths**

   ```sql
   .width 10 30 50 20
   SELECT id, name, address, email FROM provider;
   ```

3. **Execute SQL from file**

   ```bash
   sqlite3 termsheet.db < queries.sql
   ```

4. **Enable timing for performance**

   ```sql
   .timer on
   SELECT * FROM invoice;
   ```
