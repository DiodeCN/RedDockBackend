package CanSendVerificationCode

import "time"

var ipPhoneTimers = make(map[string]time.Time)

func CanSendVerificationCode(ip string, phoneNumber string) bool {
	key := ip + ":" + phoneNumber
	lastSendTime, exists := ipPhoneTimers[key]

	if exists && time.Since(lastSendTime) < 60*time.Second {
		return false
	}

	ipPhoneTimers[key] = time.Now()
	return true
}
