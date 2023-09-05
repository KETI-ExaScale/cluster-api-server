package util

import (
	"fmt"
	"time"
)

func PrettyDuration(duration time.Duration) string {
	result := ""
	second := int(duration.Seconds())
	minute := second / 60
	second = second % 60
	hour := minute / 60
	day := hour / 24

	if day > 9 {
		result = fmt.Sprintf("%dd", day)
	} else if day > 0 {
		result = fmt.Sprintf("%dd%dh", day, hour)
	} else if hour > 9 {
		result = fmt.Sprintf("%dh", hour)
	} else if hour > 0 {
		result = fmt.Sprintf("%dh%dm", hour, minute)
	} else if minute > 9 {
		result = fmt.Sprintf("%dm", minute)
	} else if minute > 0 {
		result = fmt.Sprintf("%dm%ds", minute, second)
	} else {
		result = fmt.Sprintf("%ds", second)
	}

	return result
}
