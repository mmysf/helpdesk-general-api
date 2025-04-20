package helpers

import (
	"app/domain/model"
	"os"
	"strconv"
)

func GetSubscriptionDuration() int {
	dur, _ := strconv.Atoi(os.Getenv("SUBSCRIPTION_DURATION"))
	if dur == 0 {
		dur = 12 // default 12 month
	}
	return dur
}

func CustomerBalanceFormat(c *model.Customer) *model.Customer {
	if c.Subscription != nil {
		if balance := c.Subscription.Balance; balance != nil {
			remainingSeconds := c.Subscription.Balance.Time.Total - c.Subscription.Balance.Time.Used
			hours, minutes, seconds := ConvertRemainingSeconds(remainingSeconds)
			c.Subscription.Balance.Time.Remaining = model.RemainingNested{
				Total:  remainingSeconds,
				Hour:   hours,
				Minute: minutes,
				Second: seconds,
			}
		}
	}
	return c
}

func ConvertRemainingSeconds(seconds int64) (int64, int64, int64) {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	return hours, minutes, remainingSeconds
}
