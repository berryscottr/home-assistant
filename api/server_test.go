package api

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/berryscottr/home-assistant/devices"
)

// A mock device implementing the Device interface for testing
type mockDevice struct {
	id     string
	on     bool
	status string
}

func (m *mockDevice) ID() string           { return m.id }
func (m *mockDevice) TurnOn() error        { m.on = true; m.status = "on"; return nil }
func (m *mockDevice) TurnOff() error       { m.on = false; m.status = "off"; return nil }
func (m *mockDevice) Status() string       { return m.status }

func setupMockDevices() {
	devices.RegisterDevice(&mockDevice{id: "dev1"})
}

func TestDeviceOnOffEndpoints(t *testing.T) {
	setupMockDevices()

	handler := NewHandler()

	reqOn := httptest.NewRequest("GET", "/device/on?id=dev1", nil)
	wOn := httptest.NewRecorder()
	handler.ServeHTTP(wOn, reqOn)

	respOn := wOn.Result()
	bodyOn, _ := io.ReadAll(respOn.Body)
	if respOn.StatusCode != 200 {
		t.Fatalf("Expected 200 OK on /device/on, got %d", respOn.StatusCode)
	}
	if !strings.Contains(string(bodyOn), "Device dev1 turned on") {
		t.Fatalf("Unexpected response body on /device/on: %s", bodyOn)
	}

	reqOff := httptest.NewRequest("GET", "/device/off?id=dev1", nil)
	wOff := httptest.NewRecorder()
	handler.ServeHTTP(wOff, reqOff)

	respOff := wOff.Result()
	bodyOff, _ := io.ReadAll(respOff.Body)
	if respOff.StatusCode != 200 {
		t.Fatalf("Expected 200 OK on /device/off, got %d", respOff.StatusCode)
	}
	if !strings.Contains(string(bodyOff), "Device dev1 turned off") {
		t.Fatalf("Unexpected response body on /device/off: %s", bodyOff)
	}
}
