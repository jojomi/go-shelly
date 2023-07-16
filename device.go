package shelly

import (
	"fmt"
	"strings"
)

// Device is a Shelly-powered device
type Device struct {
	Hostname string
}

// GetAPIBaseURL returns the base URL to be used to send commands to this Shelly device.
func (x Device) GetAPIBaseURL() string {
	hostname := strings.TrimPrefix(x.Hostname, "http://")
	return fmt.Sprintf("http://%s/", hostname)
}
