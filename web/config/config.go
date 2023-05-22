package config

import (
	"time"
)

func LoadConfigFromDatabase() {
	var database_config_cache_time time.Time = time.Now()

	_ = database_config_cache_time
}
