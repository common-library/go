// Package json provides utilities for converting between JSON and Go structs.
//
// This package simplifies JSON marshaling and unmarshaling operations with convenient
// wrapper functions around Go's encoding/json package.
//
// Features:
//   - Struct to JSON string conversion
//   - Formatted JSON output with indentation
//   - Generic type conversion from JSON files
//   - Generic type conversion from JSON strings
//
// Example:
//
//	type User struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	user := User{Name: "Alice", Age: 30}
//	jsonStr, _ := json.ToString(user)
package json

import (
	"encoding/json"

	"github.com/common-library/go/file"
)

// ToString converts a Go struct to a JSON string.
//
// This function marshals the input struct into compact JSON format without
// any formatting or indentation.
//
// Parameters:
//   - input: Any Go value to convert to JSON (typically a struct)
//
// Returns:
//   - string: JSON representation of the input
//   - error: Error if marshaling fails (e.g., unsupported type, circular reference)
//
// The function uses json.Marshal internally, so the same rules apply:
//   - Struct fields must be exported (capitalized) to be included
//   - Use struct tags to control JSON field names: `json:"fieldName"`
//   - Unexported fields are ignored
//   - Nil pointers, empty slices, and zero values are included
//
// Example:
//
//	type User struct {
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	    Age   int    `json:"age"`
//	}
//
//	user := User{
//	    Name:  "Alice",
//	    Email: "alice@example.com",
//	    Age:   30,
//	}
//
//	jsonStr, err := ToString(user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Println(jsonStr)
//	// Output: {"name":"Alice","email":"alice@example.com","age":30}
//
// Example with slice:
//
//	users := []User{
//	    {Name: "Alice", Email: "alice@example.com", Age: 30},
//	    {Name: "Bob", Email: "bob@example.com", Age: 25},
//	}
//
//	jsonStr, err := ToString(users)
//	// Output: [{"name":"Alice",...},{"name":"Bob",...}]
func ToString(input any) (string, error) {
	output, err := json.Marshal(input)
	return string(output), err
}

// ToStringIndent converts a Go struct to a formatted JSON string with indentation.
//
// This function marshals the input struct into human-readable JSON format with
// custom prefix and indentation strings for each nesting level.
//
// Parameters:
//   - input: Any Go value to convert to JSON (typically a struct)
//   - prefix: String to prepend to each line (typically empty "")
//   - indent: String to use for each indentation level (e.g., "\t" or "  ")
//
// Returns:
//   - string: Formatted JSON representation of the input
//   - error: Error if marshaling fails
//
// The prefix is added to the beginning of each line, and indent is added for
// each level of nesting in the JSON structure.
//
// Example with tab indentation:
//
//	type Address struct {
//	    City    string `json:"city"`
//	    Country string `json:"country"`
//	}
//
//	type User struct {
//	    Name    string  `json:"name"`
//	    Age     int     `json:"age"`
//	    Address Address `json:"address"`
//	}
//
//	user := User{
//	    Name: "Alice",
//	    Age:  30,
//	    Address: Address{
//	        City:    "New York",
//	        Country: "USA",
//	    },
//	}
//
//	jsonStr, err := ToStringIndent(user, "", "\t")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Println(jsonStr)
//	// Output:
//	// {
//	//     "name": "Alice",
//	//     "age": 30,
//	//     "address": {
//	//         "city": "New York",
//	//         "country": "USA"
//	//     }
//	// }
//
// Example with 2-space indentation:
//
//	jsonStr, err := ToStringIndent(user, "", "  ")
//	// Uses 2 spaces per indentation level
//
// Example with prefix:
//
//	jsonStr, err := ToStringIndent(user, "> ", "  ")
//	// Each line starts with "> "
func ToStringIndent(input any, prefix string, indent string) (string, error) {
	output, err := json.MarshalIndent(input, prefix, indent)
	return string(output), err
}

// ConvertFromFile reads a JSON file and converts it to the specified type T.
//
// This generic function reads JSON data from a file and unmarshals it into
// the specified Go type. The type parameter T must be provided explicitly
// at the call site.
//
// Parameters:
//   - fileName: Path to the JSON file to read
//
// Returns:
//   - T: The unmarshaled data as type T
//   - error: Error if file reading fails or JSON unmarshaling fails
//
// The function performs two operations:
//  1. Reads the file contents using file.Read
//  2. Unmarshals the JSON data into type T
//
// Type parameter T can be any valid Go type including structs, slices, maps,
// or primitive types.
//
// Example with struct:
//
//	type Config struct {
//	    Host     string `json:"host"`
//	    Port     int    `json:"port"`
//	    Database string `json:"database"`
//	}
//
//	// config.json contains:
//	// {
//	//   "host": "localhost",
//	//   "port": 5432,
//	//   "database": "mydb"
//	// }
//
//	config, err := json.ConvertFromFile[Config]("config.json")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Host: %s, Port: %d\n", config.Host, config.Port)
//
// Example with slice:
//
//	type User struct {
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	}
//
//	// users.json contains:
//	// [
//	//   {"name": "Alice", "email": "alice@example.com"},
//	//   {"name": "Bob", "email": "bob@example.com"}
//	// ]
//
//	users, err := json.ConvertFromFile[[]User]("users.json")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Loaded %d users\n", len(users))
//
// Example with map:
//
//	settings, err := json.ConvertFromFile[map[string]interface{}]("settings.json")
//	// Unmarshals JSON object to map[string]interface{}
func ConvertFromFile[T any](fileName string) (T, error) {
	if data, err := file.Read(fileName); err != nil {
		var t T
		return t, err
	} else {
		return ConvertFromString[T](data)
	}
}

// ConvertFromString converts a JSON string to the specified type T.
//
// This generic function unmarshals JSON data from a string into the specified
// Go type. The type parameter T must be provided explicitly at the call site.
//
// Parameters:
//   - data: JSON string to parse
//
// Returns:
//   - T: The unmarshaled data as type T
//   - error: Error if JSON parsing or unmarshaling fails
//
// The function uses json.Unmarshal internally, so the same rules apply:
//   - JSON field names are matched to struct field names (case-insensitive)
//   - Struct tags control how JSON fields map to struct fields
//   - Extra JSON fields are ignored
//   - Missing JSON fields result in zero values
//
// Example with struct:
//
//	type User struct {
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	    Age   int    `json:"age"`
//	}
//
//	jsonData := `{
//	    "name": "Alice",
//	    "email": "alice@example.com",
//	    "age": 30
//	}`
//
//	user, err := json.ConvertFromString[User](jsonData)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("User: %+v\n", user)
//	// Output: User: {Name:Alice Email:alice@example.com Age:30}
//
// Example with slice:
//
//	jsonData := `[
//	    {"name": "Alice", "age": 30},
//	    {"name": "Bob", "age": 25}
//	]`
//
//	users, err := json.ConvertFromString[[]User](jsonData)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Loaded %d users\n", len(users))
//
// Example with map:
//
//	jsonData := `{"key1": "value1", "key2": "value2"}`
//
//	data, err := json.ConvertFromString[map[string]string](jsonData)
//	// Unmarshals to map[string]string
//
// Example with primitive type:
//
//	jsonData := `[1, 2, 3, 4, 5]`
//
//	numbers, err := json.ConvertFromString[[]int](jsonData)
//	// Unmarshals to []int
func ConvertFromString[T any](data string) (T, error) {
	var t T

	err := json.Unmarshal([]byte(data), &t)

	return t, err
}
