// Package file provides utilities for file and directory operations.
//
// This package offers simplified functions for common file system operations including
// reading, writing, listing, creating, and removing files and directories. It wraps
// Go's standard os and path/filepath packages with convenient interfaces.
//
// Features:
//   - Read and write file contents
//   - List files and directories (recursive and non-recursive)
//   - Create directories (single and nested)
//   - Remove files and directories
//   - File permission management
//
// Example:
//
//	data, _ := file.Read("config.txt")
//	file.Write("output.txt", data, 0644)
//	files, _ := file.List("./data", true)
package file

import (
	"os"
	"path/filepath"
)

// Read reads and returns the entire contents of a file as a string.
//
// Parameters:
//   - fileName: Path to the file to read (absolute or relative)
//
// Returns:
//   - string: Contents of the file as a string
//   - error: Error if file doesn't exist, cannot be read, or permission denied
//
// The function reads the entire file into memory, so it's best suited for small to
// medium-sized text files. For large files, consider using streaming approaches.
//
// Example:
//
//	data, err := file.Read("config.json")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(data)
func Read(fileName string) (string, error) {
	if data, err := os.ReadFile(fileName); err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

// Write writes data to a file, creating it if it doesn't exist or truncating it if it does.
//
// Parameters:
//   - fileName: Path to the file to write (absolute or relative)
//   - data: String content to write to the file
//   - fileMode: File permissions (e.g., 0644, 0600, os.ModePerm)
//
// Returns:
//   - error: Error if file cannot be created or written, nil on success
//
// If the file already exists, it will be truncated before writing. The fileMode parameter
// sets the Unix file permissions. Common modes: 0644 (rw-r--r--), 0600 (rw-------),
// os.ModePerm (0777).
//
// Example:
//
//	// Write with owner read/write only
//	err := file.Write("secret.txt", "password123", 0600)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Write with full permissions
//	err = file.Write("public.txt", "hello world", os.ModePerm)
func Write(fileName string, data string, fileMode os.FileMode) error {
	return os.WriteFile(fileName, []byte(data), fileMode)
}

// List returns a list of files and directories in the specified path.
//
// Parameters:
//   - path: Directory path to list (absolute or relative)
//   - recursive: If true, recursively lists all subdirectories; if false, lists only immediate children
//
// Returns:
//   - []string: Slice of file and directory paths. Directories end with filepath.Separator
//   - error: Error if path doesn't exist, cannot be read, or permission denied
//
// For non-recursive listing, directory entries end with the path separator (/ or \).
// For recursive listing, all files are included with full paths, and only empty
// subdirectories are listed.
//
// Example:
//
//	// List immediate children only
//	files, err := file.List("./data", false)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, f := range files {
//	    if strings.HasSuffix(f, string(filepath.Separator)) {
//	        fmt.Printf("Directory: %s\n", f)
//	    } else {
//	        fmt.Printf("File: %s\n", f)
//	    }
//	}
//
//	// List all files recursively
//	allFiles, err := file.List("./project", true)
//	for _, f := range allFiles {
//	    fmt.Println(f)
//	}
func List(path string, recursive bool) ([]string, error) {
	var list func(*[]string, string, bool) error
	list = func(result *[]string, path string, recursive bool) error {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if fileInfo, err := file.Stat(); err != nil {
			return err
		} else if !fileInfo.IsDir() {
			*result = append(*result, path)
			return nil
		} else if fileInfos, err := file.Readdir(0); err != nil {
			return err
		} else if len(fileInfos) == 0 {
			*result = append(*result, file.Name()+string(filepath.Separator))
			return nil
		} else {
			for _, fileInfo := range fileInfos {
				dir := filepath.Dir(file.Name() + string(filepath.Separator))
				name := dir + string(filepath.Separator) + fileInfo.Name()

				if !fileInfo.IsDir() {
					*result = append(*result, name)
				} else if !recursive {
					*result = append(*result, name+string(filepath.Separator))
				} else {
					list(result, name, recursive)
				}
			}

			return nil
		}
	}

	result := []string{}
	if err := list(&result, path, recursive); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

// CreateDirectory creates a single directory.
//
// Parameters:
//   - name: Directory path to create
//   - fileMode: Directory permissions (e.g., 0755, 0700, os.ModePerm)
//
// Returns:
//   - error: Error if directory already exists, parent doesn't exist, or permission denied
//
// This function creates only a single directory. All parent directories must already exist.
// To create nested directories, use CreateDirectoryAll instead. Common modes: 0755 (rwxr-xr-x),
// 0700 (rwx------), os.ModePerm (0777).
//
// Example:
//
//	// Create directory with standard permissions
//	err := file.CreateDirectory("data", 0755)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Create directory with owner-only permissions
//	err = file.CreateDirectory("private", 0700)
func CreateDirectory(name string, fileMode os.FileMode) error {
	return os.Mkdir(name, fileMode)
}

// CreateDirectoryAll creates a directory along with all necessary parent directories.
//
// Parameters:
//   - path: Directory path to create (can include multiple levels)
//   - fileMode: Directory permissions applied to all created directories
//
// Returns:
//   - error: Error if creation fails or permission denied, nil if successful or already exists
//
// This function creates all directories in the path that don't exist. If the directory
// already exists, it returns nil (no error). All created directories receive the same
// fileMode permissions.
//
// Example:
//
//	// Create nested directories
//	err := file.CreateDirectoryAll("data/2024/January", 0755)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Create with full permissions
//	err = file.CreateDirectoryAll("/tmp/project/build/output", os.ModePerm)
func CreateDirectoryAll(path string, fileMode os.FileMode) error {
	return os.MkdirAll(path, fileMode)
}

// Remove deletes a single file or empty directory.
//
// Parameters:
//   - name: Path to the file or empty directory to remove
//
// Returns:
//   - error: Error if file/directory doesn't exist, is not empty (for directories), or permission denied
//
// This function only removes a single file or an empty directory. To remove directories
// with contents, use RemoveAll instead. If the target doesn't exist, an error is returned.
//
// Example:
//
//	// Remove a file
//	err := file.Remove("temp.txt")
//	if err != nil {
//	    log.Printf("Failed to remove: %v", err)
//	}
//
//	// Remove an empty directory
//	err = file.Remove("empty_dir")
func Remove(name string) error {
	return os.Remove(name)
}

// RemoveAll recursively removes a path and all its contents.
//
// Parameters:
//   - path: Path to the file or directory to remove
//
// Returns:
//   - error: Error if removal fails or permission denied, nil if successful or path doesn't exist
//
// This function removes the specified path and all its contents recursively. If the path
// doesn't exist, it returns nil (no error). Use with caution as this operation cannot
// be undone.
//
// Example:
//
//	// Remove directory and all contents
//	err := file.RemoveAll("build/output")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Remove project artifacts
//	file.RemoveAll("node_modules")
//	file.RemoveAll("dist")
func RemoveAll(path string) error {
	return os.RemoveAll(path)
}
