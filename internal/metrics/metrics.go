package metrics

import "sync/atomic"

type Metrics struct {
	TotalRequests       int64
	SuccessfulRequests  int64
	ValidationErrors    int64
	HoneypotBlocks      int64
	FirewallBlocks      int64
	RateLimitBlocks     int64
	CooldownBlocks      int64
	UserAgentBlocks     int64
	PayloadTooLarge     int64
	EmailSendErrors     int64
	TotalProcessingTime int64 // ms
}

var M = &Metrics{}

func Inc(counter *int64) {
	atomic.AddInt64(counter, 1)
}

func AddTime(ms int64) {
	atomic.AddInt64(&M.TotalProcessingTime, ms)
}
