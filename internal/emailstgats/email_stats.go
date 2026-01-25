package emailstats

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Stats struct {
	Month         int  `json:"month"`
	Count         int  `json:"count"`
	TierAlertSent bool `json:"tierAlertSent"`
}

var (
	stats       Stats
	filePath    string
	mu          sync.Mutex
	initialized bool
)

func Init(path string) error {
	mu.Lock()
	defer mu.Unlock()

	filePath = path

	// Si no existe, creamos uno nuevo
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		stats = Stats{
			Month:         int(time.Now().Month()),
			Count:         0,
			TierAlertSent: false,
		}
		return save()
	}

	// Si existe, lo cargamos
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &stats); err != nil {
		return err
	}

	initialized = true
	return nil
}

func save() error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func Increment() error {
	mu.Lock()
	defer mu.Unlock()

	nowMonth := int(time.Now().Month())

	// Reset mensual automático
	if stats.Month != nowMonth {
		stats.Month = nowMonth
		stats.Count = 0
		stats.TierAlertSent = false
	}

	stats.Count++
	return save()
}

func GetCount() int {
	mu.Lock()
	defer mu.Unlock()
	return stats.Count
}

func HasSentTierAlert() bool {
	mu.Lock()
	defer mu.Unlock()
	return stats.TierAlertSent
}

func MarkTierAlertSent() error {
	mu.Lock()
	defer mu.Unlock()
	stats.TierAlertSent = true
	return save()
}
