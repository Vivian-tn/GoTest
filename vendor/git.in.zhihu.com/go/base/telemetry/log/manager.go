package log

import "sync"

type logManager struct {
	loggerStore sync.Map
}

func (lm *logManager) add(name string, logger *ZLogger) {
	lm.loggerStore.Store(name, logger)
}

func (lm *logManager) getLogger(name string) *ZLogger {
	if v, ok := lm.loggerStore.Load(name); ok {
		return v.(*ZLogger)
	}
	return nil
}

func (lm *logManager) release(name string) {
	lm.loggerStore.Delete(name)
}
