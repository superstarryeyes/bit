// ABOUTME: File I/O operations for favorites persistence.
// ABOUTME: Handles reading/writing favorites.json in ~/.config/bit/.

package favorites

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	configDirName     = "bit"
	favoritesFileName = "favorites.json"
)

// GetConfigDir returns the config directory path, creating it if needed.
// Uses ~/.config/bit/ following XDG conventions.
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ".config", configDirName)
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", err
	}

	return configDir, nil
}

// GetFavoritesFilePath returns the full path to favorites.json
func GetFavoritesFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, favoritesFileName), nil
}

// Load reads favorites from disk. Returns empty store if file doesn't exist.
func Load() (*FavoritesStore, error) {
	filePath, err := GetFavoritesFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &FavoritesStore{Favorites: []Favorite{}}, nil
		}
		return nil, err
	}

	var store FavoritesStore
	err = json.Unmarshal(data, &store)
	if err != nil {
		return nil, err
	}

	// Ensure Favorites is never nil
	if store.Favorites == nil {
		store.Favorites = []Favorite{}
	}

	return &store, nil
}

// Save writes favorites to disk
func Save(store *FavoritesStore) error {
	filePath, err := GetFavoritesFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
