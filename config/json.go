package config

func DefaultConfig() string {
	var configData string = `{
		"name": "YouTube Video to Audio Converter",
		"version": "1.0.0",
		"useSixel": true
	}`
	return configData
}
