// devices/devices.go
package devices

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

type Device interface {
	ID() string
	TurnOn() error
	TurnOff() error
	Status() string
}

var (
	deviceRegistry = map[string]Device{}
	registryMu     sync.RWMutex
)

func Init(ctx context.Context) {
	log.Info().Msg("Initializing devices...")
	// In a full version, load from config/config.yaml and initialize devices
	// Example: load MQTT settings and Tasmota devices
}

func RegisterDevice(dev Device) {
	registryMu.Lock()
	defer registryMu.Unlock()
	deviceRegistry[dev.ID()] = dev
	log.Info().Str("device_id", dev.ID()).Msg("Device registered")
}

func GetDevice(id string) (Device, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	dev, ok := deviceRegistry[id]
	return dev, ok
}