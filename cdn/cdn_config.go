package cdn

import (
	"sync"
)

type CDNConfig struct {
	server string
}

var (
	instanceCfg *CDNConfig
	once        sync.Once
)

func (C *CDNConfig) Create(server string) {
	C.server = server
}

func (C CDNConfig) GetURLServer() string {
	return C.server
}

func GetInstanceConfig() *CDNConfig {
	once.Do(func() {
		instanceCfg = new(CDNConfig)
	})
	return instanceCfg
}
