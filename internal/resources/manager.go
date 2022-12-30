package resources

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

// Application data directory: https://hub.jmonkeyengine.org/t/appdata-equivalent-on-macos-and-linux/43735

const (
	appName = "mtgstudio"
	sets    = "sets"
	cards   = "cards"
)

// Manager handles the file access for the application
type Manager struct {
	applicationDirectory string
}

func NewManager() *Manager {
	m := &Manager{applicationDirectory: getApplicationDirectory()}
	err := os.Mkdir(m.applicationDirectory, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(m.applicationDirectory, sets), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir(filepath.Join(m.applicationDirectory, cards), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func (m *Manager) LoadSetIcon(path string) ([]byte, error) {

	return nil, nil
}
