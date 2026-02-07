package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	//I get the directory from the path
	dir := filepath.Dir(path)

	//I Create the temp file used to save the data
	tmp, err := os.CreateTemp(dir, "atomic-*")
	if err != nil {
		return fmt.Errorf("Error creating the temp file: %w", err)
	}
	tmpName := tmp.Name()

	//write data to the temp file
	if _, err = tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("Error writing data to temp file: %w", err)
	}

	//flush the disck
	if err = tmp.Sync(); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return fmt.Errorf("Error syncing temp file: %w", err)
	}

	//closing the file
	if err = tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("Error closing temp file: %w", err)
	}

	//set permissions
	if err = os.Chmod(tmpName, perm); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("Error setting permissions: %w", err)
	}

	//atomic rename of the file
	if err = os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("Erorr replacing file: %w", err)
	}

	return nil
}
