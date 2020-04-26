package log

import (
	"fmt"
	"testing"
	"time"
)

func TestNewLogging(t *testing.T) {
	NewLogging()
	Logger.Errorf("test error: %v", fmt.Errorf("error happend"))
	Logger.Infof("test info: %s", "some info")
	Logger.Warnf("test warn: %s", "some warning")

	Logger.WithFields(map[string]interface{}{
		"now": time.Now().Format(time.RFC3339),
		"ip":  "127.0.0.1",
	}).Info()
}
