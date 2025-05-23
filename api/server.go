// api/server.go
package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/berryscottr/home-assistant/devices"
)

func StartServer(ctx context.Context) {
	http.HandleFunc("/device/on", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if dev, ok := devices.GetDevice(id); ok {
			if err := dev.TurnOn(); err != nil {
				http.Error(w, "Failed to turn on device", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "Device %s turned on", id)
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/device/off", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if dev, ok := devices.GetDevice(id); ok {
			if err := dev.TurnOff(); err != nil {
				http.Error(w, "Failed to turn off device", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "Device %s turned off", id)
		} else {
			http.NotFound(w, r)
		}
	})

	log.Info().Msg("API server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Err(err).Msg("API server failed")
	}
}
