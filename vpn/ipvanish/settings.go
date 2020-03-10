package ipvanish

// Settings for IPvanish
type Settings struct {
	FilterSettings
}

// DefaultSettings return default ipvvanish settings
func DefaultSettings() Settings {
	return Settings{
		FilterSettings: DefaultFilterSettings(),
	}
}
