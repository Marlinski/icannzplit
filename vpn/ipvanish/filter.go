package ipvanish

import (
	"icannzplit/util"
	"math/rand"
	"time"
)

// FilterSettings for IPvanish
type FilterSettings struct {
	FilterType    ExitFilterType `json:"filter-type"`              // set how to pick up exits
	FilterRandom  int            `json:"filter-random,omitempty"`  // select a fix amount of exits randomly
	FilterCountry map[string]int `json:"filter-country,omitempty"` // select a fix amount of exits per country randomly
	FilterStatic  []string       `json:"filter-static,omitempty"`  // select exits statically
}

// ExitFilterType how to filter config
type ExitFilterType int

// FilterType
const (
	FilterTypeRandom  ExitFilterType = 0
	FilterTypeCountry                = 1
	FilterTypeStatic                 = 2
)

// DefaultFilterSettings return default ipvanish Filter settings
func DefaultFilterSettings() FilterSettings {
	return FilterSettings{
		FilterType:   FilterTypeRandom,
		FilterRandom: 5,
	}
}

// FilterConfigs filters out ipanish configurations
func FilterConfigs(s Settings) {
	switch s.FilterSettings.FilterType {
	case FilterTypeRandom:
		FilterRandom(s.FilterSettings.FilterRandom)
	case FilterTypeCountry:
		FilterCountry(s.FilterSettings.FilterCountry)
	case FilterTypeStatic:
		FilterStatic(s.FilterSettings.FilterStatic)
	}
}

// FilterRandom randomly picks a fixed amount of config
func FilterRandom(amount int) {
	// shuffle the ipvanish configs
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(configs), func(i, j int) { configs[i], configs[j] = configs[j], configs[i] })

	// pick up only the amount of exit point as configured
	amount = util.Min(amount, len(configs))
	configs = configs[:amount]
}

// FilterCountry randomly picks a fixed amount of config per country
func FilterCountry(countryMap map[string]int) {
	// shuffle the ipvanish configs
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(configs), func(i, j int) { configs[i], configs[j] = configs[j], configs[i] })

	filtered := make([]TunnelConfig, 0)

	// copy the map because we will modify the counters
	// and we don't want to modify the settings
	copyCountryMap := make(map[string]int)
	for k, v := range countryMap {
		copyCountryMap[k] = v
	}

	// select configurations per country
	for _, ipvnConfig := range configs {
		nb, ok := copyCountryMap[ipvnConfig.CountryCode]
		if !ok {
			// this country does not belong to the filter map
			continue
		}

		if nb > 0 {
			// we add this configuration if we still need
			filtered = append(filtered, ipvnConfig)
			copyCountryMap[ipvnConfig.CountryCode] = nb - 1
		}
	}
	configs = filtered
}

// FilterStatic statically choose configs
func FilterStatic(list []string) {
	filtered := make([]TunnelConfig, 0)
	for _, ipvnConfig := range configs {
		_, ok := util.Find(list, ipvnConfig.CountryID)
		if !ok {
			continue
		}

		filtered = append(filtered, ipvnConfig)

		// small optimization
		if len(filtered) == len(list) {
			break
		}
	}
	configs = filtered
}

// ListInterface returns all the interface for the ipvanish provider
func ListInterface() []string {
	ret := make([]string, len(configs))
	for i, config := range configs {
		ret[i] = config.DeviceName
	}
	return ret
}
