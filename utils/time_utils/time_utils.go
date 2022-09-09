package time_utils

import "time"

// function to return current unix timestamp
func Return_current_epoch_timestamp() int64 {
	now := time.Now()
	unix_timestamp := now.Unix()
	return unix_timestamp
}

// function to read unix timestamp and return date and time
func Return_time_date_from_epoch_timestamp(unix_timestamp int64) string {
	t := time.Unix(unix_timestamp, 0)
	date_time := t.Format("15:04:05 2006-01-02")
	return date_time
}

// function to return current date
func Return_current_date() string {
	now := time.Now()
	date := now.Format("02-Jan-2006")
	return date
}
