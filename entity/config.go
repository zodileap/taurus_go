package entity

import "sync"

type Config struct {
	// BatchSize 单条语句，参数的最大数量，以PostgreSQL为标准。
	BatchSize *int
	// SqlConsole 是否打印sql语句。
	SqlConsole *bool
	// SqlLogger sql语句的日志文件，这个是匹配tlog的日志文件名。
	SqlLogger *string
}

var config *Config
var configMu sync.RWMutex // 嵌入一个互斥锁来保护配置

func initConfig() {
	initConfig()
}

// GetConfig 提供了一个全局访问点
func GetConfig() Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return *config
}

func SetConfig(c Config) {
	configMu.Lock()
	defer configMu.Unlock()
	if c.BatchSize != nil {
		config.BatchSize = c.BatchSize
	}
	if c.SqlConsole != nil {
		config.SqlConsole = c.SqlConsole
	}
}
