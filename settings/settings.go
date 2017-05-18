package settings

import (
	"encoding/json"
	"fmt"
	"meanbot/constants"
	"os"
)

// Settings stores settings for a particular chat
type Settings struct {
	InsultInterval int64 `json:"insult_interval"`
}

// GetDefaultSettings returns the default settings for a new chat
func GetDefaultSettings() Settings {
	defaultSettings := Settings{
		InsultInterval: 100,
	}
	return defaultSettings
}

// GetSavedSettings Loads settings from disk
func GetSavedSettings() map[string]Settings {
	settingsFile, settingsFileErr := os.Open(constants.BotPath + "/settings.json")
	defer settingsFile.Close()

	if settingsFileErr != nil {
		fmt.Println("Failed to open settings file")
	}

	settings := make(map[string]Settings)
	settingsDecoder := json.NewDecoder(settingsFile)
	settingsDecoderErr := settingsDecoder.Decode(&settings)

	if settingsDecoderErr != nil {
		fmt.Println("Error decoding settings json")
	}

	fmt.Println("Loaded settings!")

	return settings
}

// SaveSettings Saves settings to disk
func SaveSettings(settings map[string]Settings) {
	settingsFile, settingsFileErr := os.Create(constants.BotPath + "settings.json")
	defer settingsFile.Close()

	if settingsFileErr != nil {
		fmt.Println("Failed to open settings file")
	}

	settingsBuf, settingsBufErr := json.Marshal(settings)

	if settingsBufErr != nil {
		fmt.Println("Failed to encode settings in a json string")
	}

	_, writeErr := settingsFile.Write(settingsBuf)

	if writeErr != nil {
		fmt.Println("Failed to write to the settings file:", writeErr)
	}

	fmt.Println("Saved settings!")

}
