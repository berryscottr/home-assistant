package devices

import "github.com/rs/zerolog/log"

type mockThermostat struct {
	id      string
	address string
	on      bool
}

func NewMockThermostat(id, address string) Device {
	return &mockThermostat{
		id:      id,
		address: address,
		on:      false,
	}
}

func (t *mockThermostat) ID() string {
	return t.id
}

func (t *mockThermostat) TurnOn() error {
	log.Info().Str("device_id", t.id).Str("address", t.address).Msg("Turning thermostat ON")
	t.on = true
	return nil
}

func (t *mockThermostat) TurnOff() error {
	log.Info().Str("device_id", t.id).Str("address", t.address).Msg("Turning thermostat OFF")
	t.on = false
	return nil
}

func (t *mockThermostat) Status() string {
	if t.on {
		return "on"
	}
	return "off"
}
