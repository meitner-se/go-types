# Go Types Package

A comprehensive Go package providing custom types for the Meitner project with support for JSON serialization/deserialization and SQL database operations.

## Overview

This package provides custom types that distinguish between `NULL` and `UNDEFINED` values coming from HTTP clients. It offers a clear API contract for types that can be used in APIs and how they should be used. For example, it differentiates between `DATE` and `TIMESTAMP` types with different formats.

## Key Features

- **Three-state values**: Each type can be `defined`, `nil`, or `undefined`
- **JSON support**: Full JSON marshaling/unmarshaling with proper null handling
- **SQL support**: Database driver interface implementation for seamless database operations
- **Type safety**: Strong typing with clear contracts
- **Rich text processing**: HTML to plain text conversion for rich content
- **Time handling**: Separate `Date` and `Timestamp` types with different formats

## Supported Types

| Type | Description | JSON Format | SQL Type |
|------|-------------|-------------|----------|
| `Bool` | Boolean values | `true`/`false`/`null` | `BOOLEAN` |
| `Date` | Date only (no time) | `"2023-12-25"` | `DATE` |
| `Float64` | 64-bit floating point | `123.45`/`null` | `DOUBLE PRECISION` |
| `Int` | 32-bit integer | `123`/`null` | `INTEGER` |
| `Int16` | 16-bit integer | `123`/`null` | `SMALLINT` |
| `Int64` | 64-bit integer | `123`/`null` | `BIGINT` |
| `JSON` | JSON raw message | `{"key": "value"}`/`null` | `JSONB` |
| `RichText` | HTML content | `"<p>content</p>"`/`null` | `TEXT` |
| `String` | Plain text | `"text"`/`null` | `VARCHAR` |
| `Time` | Hour and minute | `"15:04"` | `TIME` |
| `Timestamp` | Date and time without timezone | `"2023-12-25T15:04:05Z"` | `TIMESTAMP` |
| `UUID` | UUID/GUID | `"123e4567-e89b-12d3-a456-426614174000"` | `UUID` |

## Installation

```bash
go get github.com/meitner-se/types
```

## Usage

### Basic Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/meitner-se/types"
)

type Person struct {
    FirstName types.String
    LastName  types.String
    Age       types.Int
    IsActive  types.Bool
    BirthDate types.Date
}

func main() {
    // Create a person with some defined and some nil values
    person := Person{
        FirstName: types.NewString("John"),
        LastName:  types.NewStringFromPtr(nil), // nil value
        Age:       types.NewInt(30),
        IsActive:  types.NewBool(true),
        BirthDate: types.NewDate(time.Date(1993, 12, 25, 0, 0, 0, 0, time.UTC)),
    }

    // Marshal to JSON
    jsonBytes, _ := json.Marshal(person)
    fmt.Println(string(jsonBytes))
    // Output: {"FirstName":"John","LastName":null,"Age":30,"IsActive":true,"BirthDate":"1993-12-25"}

    // Unmarshal from JSON
    jsonStr := `{"FirstName": "Jane", "LastName": null, "Age": 25}`
    var newPerson Person
    json.Unmarshal([]byte(jsonStr), &newPerson)

    // Check states
    fmt.Printf("FirstName: defined=%t, nil=%t, value='%s'\n",
        newPerson.FirstName.IsDefined(),
        newPerson.FirstName.IsNil(),
        newPerson.FirstName.String())
}
```

### Three-State Values

Each type supports three states:

1. **Defined**: Value was explicitly set (including null)
2. **Nil**: Value is null (explicitly set to null)
3. **Undefined**: Value was not provided in the input

```go
// Defined with value
name := types.NewString("John")
fmt.Println(name.IsDefined()) // true
fmt.Println(name.IsNil())     // false

// Defined as null
nullName := types.NewStringFromPtr(nil)
fmt.Println(nullName.IsDefined()) // true
fmt.Println(nullName.IsNil())     // true

// Undefined (not provided)
undefinedName := types.NewStringUndefined()
fmt.Println(undefinedName.IsDefined()) // false
fmt.Println(undefinedName.IsNil())     // true
```

### Rich Text Processing

The `RichText` type can convert HTML content to plain text:

```go
htmlContent := `<p>Hello <strong>world</strong>!</p><ul><li>Item 1</li><li>Item 2</li></ul>`
richText := types.NewRichText(htmlContent)

plainText, err := richText.Text()
if err == nil {
    fmt.Println(plainText)
    // Output: Hello world!
    // 
    // Item 1
    // 
    // Item 2
}
```

### Time Handling

Different time types for different use cases:

```go
// Date only (no time component)
date := types.NewDate(time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC))
fmt.Println(date.String()) // "2023-12-25"

// Timestamp with time and timezone
timestamp := types.NewTimestamp(time.Now())
fmt.Println(timestamp.String()) // "2023-12-25T15:04:05Z"

// Timestamp utilities
startOfDay := timestamp.StartOfDay(time.UTC)
endOfDay := timestamp.EndOfDay(time.UTC)
```

### UUID Handling

```go
// Generate random UUID
randomUUID := types.NewRandomUUID()

// Parse from string
uuid, err := types.UUIDFromString("123e4567-e89b-12d3-a456-426614174000")
if err == nil {
    fmt.Println(uuid.String())
}

// Convert arrays
strings := []string{"uuid1", "uuid2"}
uuids := types.UUIDsFromStrings(strings)
backToStrings := types.UUIDsToStrings(uuids)
```

### Database Operations

All types implement the `sql.Scanner` and `driver.Valuer` interfaces:

```go
// Scanning from database
var name types.String
err := db.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&name)

// Storing to database
user := types.NewString("John Doe")
_, err = db.Exec("INSERT INTO users (name) VALUES (?)", user)
```

### Type Conversion

Parse values from strings:

```go
// Parse different types from strings
boolVal, _ := types.BoolFromString("true")
intVal, _ := types.IntFromString("123")
floatVal, _ := types.Float64FromString("123.45")
dateVal, _ := types.DateFromString("2023-12-25")
timestampVal, _ := types.TimestampFromString("2023-12-25T15:04:05Z")
uuidVal, _ := types.UUIDFromString("123e4567-e89b-12d3-a456-426614174000")
```

## Requirements

- Go 1.23.3 or later
- Dependencies (automatically managed):
  - `github.com/aarondl/null/v8`
  - `github.com/friendsofgo/errors`
  - `github.com/google/uuid`
  - `github.com/stretchr/testify` (for testing)
  - `golang.org/x/net`

## Testing

Run the test suite:

```bash
go test ./...
```

Run with coverage:

```bash
go test -cover ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is part of the Meitner ecosystem. Please refer to the project's license documentation for details.
