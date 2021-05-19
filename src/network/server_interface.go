package network

import "github.com/lgrossi/go-scaffold/src/configs"

type ServerInterface interface {
	Run(globalConfigs configs.GlobalConfigs) error
	GetName() string
}
