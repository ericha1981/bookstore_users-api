package date_utils

import "time"

const (
	apiDateLayout = "2006-01-02T15:04:05.000Z"
	apiDBLayout = "2006-01-02 15:04:05"
)

func GetNow() time.Time { // Having this as a function can make it flexible if we want to do a different formatting in the future.
	return time.Now().UTC()
}
func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

func GetNowDBFormat() string {
	return GetNow().Format(apiDBLayout)
}
