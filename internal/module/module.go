/*

Package module is fetching the current module from the script's context.

*/
package module

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Get fetches the current module name from the go.mod file.
func Get() (string, error) {
	// Get the current directory, i.e. the directory where the binary is called,
	// not the directory where the binary is stored.
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Open the go.mod file.
	file, err := os.Open(filepath.Join(dir, "go.mod"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Scan the first line
	scanner.Scan()
	rawLine := scanner.Text()
	module := strings.TrimPrefix(rawLine, "module ")

	return module, nil
}
