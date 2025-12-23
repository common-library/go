package utility

import "reflect"

// GetTypeName returns the type name of a value as a string.
//
// This function uses reflection to determine the runtime type of any value
// and returns its string representation.
//
// Parameters:
//   - value: Any value to get the type name from
//
// Returns:
//   - string: Type name (e.g., "int", "string", "*MyStruct", "[]int")
//
// The returned string format:
//   - Built-in types: "int", "string", "bool", etc.
//   - Pointers: "*TypeName"
//   - Slices: "[]TypeName"
//   - Maps: "map[KeyType]ValueType"
//   - Structs: "package.StructName" or "main.StructName"
//
// Example with basic types:
//
//	typeName := utility.GetTypeName(42)
//	fmt.Println(typeName) // Output: "int"
//
//	typeName = utility.GetTypeName("hello")
//	fmt.Println(typeName) // Output: "string"
//
// Example with complex types:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	user := User{Name: "Alice", Age: 30}
//	typeName := utility.GetTypeName(user)
//	fmt.Println(typeName) // Output: "main.User"
//
//	typeName = utility.GetTypeName(&user)
//	fmt.Println(typeName) // Output: "*main.User"
//
// Example with collections:
//
//	slice := []int{1, 2, 3}
//	fmt.Println(utility.GetTypeName(slice)) // Output: "[]int"
//
//	myMap := map[string]int{"a": 1}
//	fmt.Println(utility.GetTypeName(myMap)) // Output: "map[string]int"
//
// Example for type checking:
//
//	func processValue(value any) {
//	    typeName := utility.GetTypeName(value)
//	    switch typeName {
//	    case "int":
//	        fmt.Println("Processing integer")
//	    case "string":
//	        fmt.Println("Processing string")
//	    default:
//	        fmt.Printf("Processing %s\n", typeName)
//	    }
//	}
func GetTypeName(value any) string {
	return reflect.TypeOf(value).String()
}
