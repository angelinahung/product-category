package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"

	"github.com/angelinahung/product-category/pkg/config"
)

var (
	// onceInit guarantee only once logger initialization
	onceInit sync.Once
)

// Init initializes logger
func Init() {
	onceInit.Do(func() {
		var logger *zap.Logger
		var err error
		if config.Options.DevMode {
			logger, err = zap.NewDevelopment()
		} else {
			logger, err = zap.NewProduction()
		}
		if err != nil {
			log.Fatal("failed to initialize zap logger:", err)
		}
		zap.RedirectStdLog(logger)
		zap.ReplaceGlobals(logger)
	})
}
