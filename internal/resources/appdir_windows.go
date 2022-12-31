package resources

import (
	"log"
	"os"
	"path/filepath"
)

func getApplicationDirectory() string {

	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(dir, "AppData", appName)
}
