package devices

import (
	"context"
	"os"
	"os/exec"
	"testing"
)

// mockDevice for testing device registry
type mockDevice struct {
	id     string
	on     bool
	status string
}

func (m *mockDevice) ID() string           { return m.id }
func (m *mockDevice) TurnOn() error        { m.on = true; m.status = "on"; return nil }
func (m *mockDevice) TurnOff() error       { m.on = false; m.status = "off"; return nil }
func (m *mockDevice) Status() string       { return m.status }

func TestRegisterAndGetDevice(t *testing.T) {
	mock := &mockDevice{id: "testdev"}
	RegisterDevice(mock)

	dev, ok := GetDevice("testdev")
	if !ok {
		t.Fatalf("Expected to find registered device")
	}
	if dev.ID() != "testdev" {
		t.Errorf("Expected device ID 'testdev', got '%s'", dev.ID())
	}

	// Check behavior of TurnOn/TurnOff
	if err := dev.TurnOn(); err != nil {
		t.Errorf("TurnOn returned error: %v", err)
	}
	if dev.Status() != "on" {
		t.Errorf("Expected status 'on', got '%s'", dev.Status())
	}

	if err := dev.TurnOff(); err != nil {
		t.Errorf("TurnOff returned error: %v", err)
	}
	if dev.Status() != "off" {
		t.Errorf("Expected status 'off', got '%s'", dev.Status())
	}
}

// overrideExecCommand lets us mock exec.Command
var overrideExecCommand = func(name string, args ...string) *exec.Cmd {
	return mockCmd("? (192.168.1.42) at 5c:fc:e1:7c:24:06 on en0 [ethernet]\n")
}



func init() {
	// Redirect exec.Command to our mock if set
	execCommand = func(name string, arg ...string) *exec.Cmd {
		if overrideExecCommand != nil {
			return overrideExecCommand(name, arg...)
		}
		return exec.Command(name, arg...)
	}
}

var execCommand = exec.Command

// TestFindIPByMAC_Match mocks `arp -a` output to simulate IP lookup
func TestFindIPByMAC_Match(t *testing.T) {
	overrideExecCommand = func(name string, args ...string) *exec.Cmd {
		return mockCmd("5c:fc:e1:7c:24:06 192.168.1.42\n")
	}
	defer func() { overrideExecCommand = nil }()

	ip, err := FindIPByMAC("5C:FC:E1:7C:24:06")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if ip != "192.168.1.42" {
		t.Errorf("Expected IP '192.168.1.42', got '%s'", ip)
	}
}

// TestInit_ParsesConfigAndRegistersDevices writes a fake config file and tests Init()
func TestInit_ParsesConfigAndRegistersDevices(t *testing.T) {
	configYAML := `
thermostats:
  - id: "downstairs"
    model: "honeywell_t5+"
    mac_id: "5C:FC:E1:7C:24:06"
`
	err := os.MkdirAll("tmp_config", 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	err = os.WriteFile("tmp_config/config.yaml", []byte(configYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}
	defer os.Remove("tmp_config/config.yaml")

	overrideExecCommand = func(name string, args ...string) *exec.Cmd {
		return mockCmd("5c:fc:e1:7c:24:06 (192.168.1.10) at eth0\n")
	}
	defer func() { overrideExecCommand = nil }()

	Init(context.Background())

	dev, ok := GetDevice("downstairs")
	if !ok {
		t.Fatal("Expected thermostat 'downstairs' to be registered")
	}
	if dev.ID() != "downstairs" {
		t.Errorf("Expected ID 'downstairs', got '%s'", dev.ID())
	}
}

func mockCmd(output string) *exec.Cmd {
	return exec.Command("echo", output)
}
