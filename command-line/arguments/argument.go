// Package arguments provides utilities for accessing command-line arguments.
//
// This package offers a simple wrapper around os.Args for retrieving
// command-line arguments by index or as a complete slice.
//
// Features:
//   - Individual argument access by index
//   - Retrieve all arguments as a slice
//   - Simple wrapper around os.Args
//
// Example:
//
//	arg := arguments.Get(1)  // Get first argument
//	all := arguments.GetAll() // Get all arguments
package arguments

import "os"

// Get returns the command-line argument at the specified index.
// This is a direct wrapper around os.Args[index].
//
// Parameters:
//   - index: The position of the argument to retrieve (0 is the program name)
//
// Returns:
//   - The argument string at the specified index
//
// Example:
//
//	// For command: ./program arg1 arg2
//	programName := arguments.Get(0)  // "./program"
//	firstArg := arguments.Get(1)     // "arg1"
//	secondArg := arguments.Get(2)    // "arg2"
//
// Note: Calling Get with an out-of-bounds index will panic.
func Get(index int) string {
	return os.Args[index]
}

// GetAll returns all command-line arguments as a string slice.
// This is a direct wrapper around os.Args.
//
// Returns:
//   - A slice containing all command-line arguments, including the program name at index 0
//
// Example:
//
//	// For command: ./program arg1 arg2 --flag=value
//	args := arguments.GetAll()
//	// args = ["./program", "arg1", "arg2", "--flag=value"]
//	fmt.Printf("Total arguments: %d\n", len(args))
func GetAll() []string {
	return os.Args
}
