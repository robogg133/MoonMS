package app

import (
	"time"

	"github.com/dgraph-io/badger/v4"
)

func dbGC(db *badger.DB, interval time.Duration, discardRatio float64) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		for {
			if err := db.RunValueLogGC(discardRatio); err != nil {
				break
			}
		}
	}
}
