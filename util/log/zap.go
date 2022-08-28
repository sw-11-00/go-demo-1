package log

import (
	"go.uber.org/zap"
	"time"
)

// LogZapTest https://zhuanlan.zhihu.com/p/371547318
func LogZapTest() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "http://example.com",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", "http://example.com")
}
