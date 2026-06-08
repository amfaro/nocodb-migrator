# NocoDB Migration Tool

A tool for managing database migrations in NocoDB via Meta API v3.

## Installation

```bash
go build -o nocodb-migrate
```

## Configuration

Create a `.env` file in the project root:

```env
NOCODB_URL=http://localhost:8080
NOCODB_API_TOKEN=your_api_token_here
NOCODB_BASE_ID=your_base_id_here
NOCODB_MIGRATIONS_DIR=./migrations
```

Or export environment variables:

```bash
export NOCODB_URL=http://localhost:8080
export NOCODB_API_TOKEN=your_api_token_here
export NOCODB_BASE_ID=your_base_id_here
export NOCODB_MIGRATIONS_DIR=./migrations
```

### Environment Variables

- **NOCODB_URL** (required) - URL of your NocoDB instance
- **NOCODB_API_TOKEN** (required) - API token for accessing NocoDB
- **NOCODB_BASE_ID** (required) - Database ID in NocoDB
- **NOCODB_MIGRATIONS_DIR** (optional) - Directory with migration files (default: `./migrations`)

## Usage

### Creating a Migration

```bash
nocodb-migrate create add_users_table
```

This will create two files in the migrations directory (default `migrations/`, can be changed via `NOCODB_MIGRATIONS_DIR` environment variable):
- `{timestamp}-add_users_table.up.json` - migration for applying
- `{timestamp}-add_users_table.down.json` - migration for rollback

### Applying Migrations

```bash
# Apply all pending migrations
nocodb-migrate up

# Apply only N migrations
nocodb-migrate up 3
```

### Rolling Back Migrations

```bash
# Rollback the last migration
nocodb-migrate down 1

# Rollback N migrations
nocodb-migrate down 3
```

### Viewing Status

```bash
nocodb-migrate info
```

Shows the current version and list of applied migrations.

## Migration Format

Migrations are JSON files with an array of operations.

### Supported Operations

1. **create_table** - create a table
2. **alter_table** - modify a table
3. **drop_table** - delete a table
4. **create_field** - create a field
5. **alter_field** - modify a field
6. **drop_field** - delete a field
7. **insert_row** - insert data
8. **delete_row** - delete data

### Example Migration (up.json)

```json
{
  "operations": [
    {
      "type": "create_table",
      "table": "Users",
      "columns": [
        {
          "name": "Id",
          "type": "ID",
          "required": true
        },
        {
          "name": "Name",
          "type": "SingleLineText",
          "required": true
        },
        {
          "name": "Email",
          "type": "Email",
          "required": true,
          "unique": true
        },
        {
          "name": "Age",
          "type": "Number",
          "required": false,
          "default_value": 0
        },
        {
          "name": "CreatedAt",
          "type": "DateTime",
          "required": true
        }
      ]
    },
    {
      "type": "create_field",
      "table": "Users",
      "column": {
        "name": "Bio",
        "type": "LongText",
        "required": false
      }
    }
  ]
}
```

### Example Migration (down.json)

```json
{
  "operations": [
    {
      "type": "drop_table",
      "table": "Users"
    }
  ]
}
```

### NocoDB Data Types

The following data types are supported:

- **SingleLineText** - single-line text
- **LongText** - multi-line text
- **Number** - integer
- **Decimal** - decimal number
- **Currency** - currency
- **Percent** - percentage
- **DateTime** - date and time
- **Date** - date
- **Email** - email address
- **PhoneNumber** - phone number
- **URL** - URL address
- **SingleSelect** - single value selection
- **MultiSelect** - multiple value selection
- **Checkbox** - checkbox
- **Rating** - rating
- **Attachment** - attachment
- **JSON** - JSON object
- **LinkToAnotherRecord** - link to another record
- **User** - user
- **CreatedTime** - creation time
- **CreatedBy** - creator
- **LastModifiedTime** - last modification time
- **LastModifiedBy** - last modifier
- **ID** - identifier (primary key)

### Operation Examples

#### Creating a Table with Fields

```json
{
  "type": "create_table",
  "table": "Products",
  "columns": [
    {
      "name": "Id",
      "type": "ID",
      "required": true
    },
    {
      "name": "Title",
      "type": "SingleLineText",
      "required": true
    },
    {
      "name": "Price",
      "type": "Currency",
      "required": true,
      "options": {
        "code": "USD"
      }
    }
  ]
}
```

#### Modifying a Table

```json
{
  "type": "alter_table",
  "table": "Products",
  "data": {
    "title": "NewProducts",
    "description": "Updated description"
  }
}
```

#### Creating a Field

```json
{
  "type": "create_field",
  "table": "Products",
  "column": {
    "name": "Description",
    "type": "LongText",
    "required": false
  }
}
```

With explicit NocoDB display order:

```json
{
  "type": "create_field",
  "table": "Products",
  "column": {
    "name": "SKU",
    "type": "SingleLineText",
    "required": true,
    "order": 3
  }
}
```

#### Modifying a Field

```json
{
  "type": "alter_field",
  "table": "Products",
  "field_id": "field_id_here",
  "column": {
    "name": "Description",
    "type": "LongText",
    "required": true
  }
}
```

Or by field name:

```json
{
  "type": "alter_field",
  "table": "Products",
  "column": {
    "name": "Description",
    "required": true
  }
}
```

Reordering an existing field:

```json
{
  "type": "alter_field",
  "table": "Products",
  "column": {
    "name": "Description",
    "order": 5
  }
}
```

#### Deleting a Field

```json
{
  "type": "drop_field",
  "table": "Products",
  "field_id": "field_id_here"
}
```

Or by name:

```json
{
  "type": "drop_field",
  "table": "Products",
  "column": {
    "name": "Description"
  }
}
```

#### Inserting Data

```json
{
  "type": "insert_row",
  "table": "Products",
  "data": {
    "Title": "Product 1",
    "Price": 99.99
  }
}
```

#### Deleting Data

By record ID:

```json
{
  "type": "delete_row",
  "table": "Products",
  "record_id": "record_id_here"
}
```

By condition:

```json
{
  "type": "delete_row",
  "table": "Products",
  "where": {
    "Title": "Product 1"
  }
}
```

## Migrations Table

The tool automatically creates a `Migrations` table in NocoDB to track applied migrations. Table structure:

- **Id** (Integer, PK) - record identifier
- **Timestamp** (Number) - timestamp from migration file name
- **Name** (SingleLineText) - migration name
- **AppliedAt** (DateTime) - when the migration was applied
- **Direction** (SingleSelect) - direction: "up" or "down"
- **Status** (SingleSelect) - status: "success" or "failed"

## API Endpoints

The tool uses the following NocoDB Meta API v3 endpoints:

### Tables
- `GET /api/v3/meta/bases/{base_id}/tables` - list tables
- `POST /api/v3/meta/bases/{base_id}/tables` - create table
- `GET /api/v3/meta/bases/{baseId}/tables/{tableId}` - get table schema
- `PATCH /api/v3/meta/bases/{baseId}/tables/{tableId}` - update table
- `DELETE /api/v3/meta/bases/{baseId}/tables/{tableId}` - delete table

### Fields
- `POST /api/v3/meta/bases/{baseId}/tables/{tableId}/fields` - create field
- `GET /api/v3/meta/bases/{baseId}/fields/{fieldId}` - get field
- `PATCH /api/v3/meta/bases/{baseId}/fields/{fieldId}` - update field
- `DELETE /api/v3/meta/bases/{baseId}/fields/{fieldId}` - delete field

### Data
- `GET /api/v3/data/{baseId}/{tableId}/records` - list records
- `POST /api/v3/data/{baseId}/{tableId}/records` - insert record
- `DELETE /api/v3/data/{baseId}/{tableId}/records/{recordId}` - delete record

## License

MIT
