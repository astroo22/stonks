package utils

import (
	"time"
)

func ScheduleTask(interval time.Duration, task func()) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				task()
			}
		}
	}()
}
