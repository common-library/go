// Package arguments provides command line arguments
package arguments

import "os"

// Get is get the command line arguments.
func Get(index int) string {
	return os.Args[index]
}

// GetAll is gets all command line arguments.
func GetAll() []string {
	return os.Args
}
