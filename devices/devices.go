// devices/devices.go
package devices

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"os"
)

type Device interface {
	ID() string
	TurnOn() error
	TurnOff() error
	Status() string
}

type ThermostatConfig struct {
	ID      string `yaml:"id"`
	Model   string `yaml:"model"`
	MACID   string `yaml:"mac_id"`
	Address string `yaml:"address"`
}

type Config struct {
	Thermostats []ThermostatConfig `yaml:"thermostats"`
}

var (
	deviceRegistry = map[string]Device{}
	registryMu     sync.RWMutex
)

func Init(ctx context.Context) {
	log.Info().Msg("Initializing devices...")

	data, err := os.ReadFile("../config/config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config.yaml")
		return
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config.yaml")
		return
	}

	for _, t := range cfg.Thermostats {
		ip, err := FindIPByMAC(t.MACID)
		if err != nil {
			log.Error().Str("device_id", t.ID).Str("mac", t.MACID).Err(err).Msg("Could not find IP for thermostat")
			continue
		}

		t.Address = ip
		dev := NewMockThermostat(t.ID, t.Address)
		RegisterDevice(dev)
	}
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

func FindIPByMAC(mac string) (string, error) {
	out, err := exec.Command("arp", "-a").Output()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), strings.ToLower(mac)) {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				ip := strings.Trim(parts[1], "()")
				return ip, nil
			}
		}
	}
	return "", fmt.Errorf("MAC %s not found", mac)
}
