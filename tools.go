package main

func formatDuration(seconds int64) (days, hours, minutes, remainingSeconds int64) {
	minutesPerHour := int64(60)
	minutesPerDay := int64(24 * 60)
	secondsPerMinute := int64(60)

	days = seconds / (minutesPerDay * secondsPerMinute)
	seconds %= minutesPerDay * secondsPerMinute

	hours = seconds / (minutesPerHour * secondsPerMinute)
	seconds %= minutesPerHour * secondsPerMinute

	minutes = seconds / secondsPerMinute
	remainingSeconds = seconds % secondsPerMinute

	return days, hours, minutes, remainingSeconds
}
