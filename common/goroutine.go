package common

import "go-crontab/log"

func SafelyGo(work func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("Work failed err %v", err)
			}
		}()
		work()
	}()
}
