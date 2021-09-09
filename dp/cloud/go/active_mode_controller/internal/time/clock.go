package time

import "time"

type Clock struct{}

func (*Clock) Now() time.Time {
	return time.Now()
}

func (*Clock) Tick(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}
