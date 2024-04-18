package durations

import "time"

// GetDurationWithDefault gets either the value configured without changing anything
// OR the default provided in case the value is 0
func GetDurationWithDefault(value time.Duration, defaultToSet time.Duration) time.Duration {
	var durationToSet time.Duration
	if value > 0 {
		durationToSet = value
	} else {
		durationToSet = defaultToSet
	}
	return durationToSet
}
