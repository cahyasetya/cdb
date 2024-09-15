package main

import (
	"github.com/cahyasetya/cdb/file_updater/atomic_update"
	"github.com/cahyasetya/cdb/file_updater/in_place"
)

// SaveData1 saves the given data to the specified file path.
//
// Parameters:
// - path: The path of the file to save the data to.
// - data: The byte array containing the data to be saved.
//
// Returns:
// - error: An error if there was a problem opening or writing to the file.

func main() {
	atomic_update.Save("test.txt", []byte("Hello, World!"))
	inplace.Save("test_inplace.txt", []byte("Hello, World!"))
}
