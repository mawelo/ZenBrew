package main

import (
	"encoding/json"
	log "log/slog"
	"os"
)

func main() {
	log.Info("hello slog")
}

type Settings struct {
	AutoCleanup bool   `json:"cleanup"`
	RootDir     string `json:"dir"`
	Repos        []string `json:"repos"`
}

var settings Settings

func get_settings() Settings {
	// Read the JSON file
	file, err := os.ReadFile("settings.json")
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	// Parse the JSON data into the settings structure
	var settings Settings
	err = json.Unmarshal(file, &settings)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	return settings
}