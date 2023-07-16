package shelly

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/federicoleon/go-httpclient/core"
	"github.com/federicoleon/go-httpclient/gohttp"
)

// Client is a client for communication with a Shelly-powered device.
type Client struct {
	device Device
}

// NewClient returns a new Client that is configured to be talking to the given Device.
func NewClient(device Device) *Client {
	return &Client{
		device: device,
	}
}

// SetPowerOn turns the power on for the configured device.
func (x *Client) SetPowerOn() error {
	return x.SetPower(true)
}

// SetPowerOff turns the power off for the configured device.
func (x *Client) SetPowerOff() error {
	return x.SetPower(false)
}

// SetPower sets the power state for the configured device.
// See documentation at https://shelly-api-docs.shelly.cloud/gen1/#shelly-plug-plugs-relay-0
func (x *Client) SetPower(on bool) error {
	val := "on"
	if !on {
		val = "off"
	}
	response, err := x.post("relay/0?turn=" + val)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(response.String()), &data)
	if err != nil {
		return err
	}
	newState, err := x.getBoolByPath(data, "ison")
	if err != nil {
		return err
	}
	if newState != on {
		return fmt.Errorf("invalid")
	}
	return nil
}

// IsPowerOn retrieves if the power is on for the configured device.
func (x *Client) IsPowerOn() (bool, error) {
	response, err := x.get("relay/0")
	if err != nil {
		return false, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(response.String()), &data)
	if err != nil {
		return false, err
	}
	return x.getBoolByPath(data, "ison") // = isOn!
}

// GetCurrentPower returns the current power consumption of the attached devices in Watts.
func (x *Client) GetCurrentPower() (float64, error) {
	response, err := x.get("meter/0")
	if err != nil {
		return 0, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(response.String()), &data)
	if err != nil {
		return 0, err
	}

	return x.getFloatByPath(data, "power")
}

func (x *Client) httpClient() gohttp.Client {
	return gohttp.NewBuilder().
		SetConnectionTimeout(10 * time.Second).
		SetResponseTimeout(20 * time.Second).
		SetUserAgent("go-shelly").
		Build()
}

func (x *Client) get(cmd string) (*core.Response, error) {
	apiURL := x.device.GetAPIBaseURL() + cmd
	return x.httpClient().Get(apiURL)
}

func (x *Client) post(cmd string) (*core.Response, error) {
	apiURL := x.device.GetAPIBaseURL() + cmd
	return x.httpClient().Post(apiURL, "")
}

func (x *Client) getFloatByPath(data map[string]interface{}, path string) (float64, error) {
	keys := strings.Split(path, ".")

	var (
		value any
		ok    bool
	)
	for _, key := range keys {
		value, ok = data[key]
		if !ok {
			return 0, fmt.Errorf("key '%s' not found", key)
		}

		data, ok = value.(map[string]interface{})
		if !ok {
			break
		}
	}

	vInt, ok := value.(int)
	if ok {
		return float64(vInt), nil
	}

	vFloat, ok := value.(float64)
	if ok {
		return vFloat, nil
	}

	return 0.0, fmt.Errorf("value is not a float")
}

func (x *Client) getBoolByPath(data map[string]interface{}, path string) (bool, error) {
	keys := strings.Split(path, ".")

	var (
		value any
		ok    bool
	)
	for _, key := range keys {
		value, ok = data[key]
		if !ok {
			return false, fmt.Errorf("key '%s' not found", key)
		}

		data, ok = value.(map[string]interface{})
		if !ok {
			break
		}
	}

	vBool, ok := value.(bool)
	if ok {
		return vBool, nil
	}

	return false, fmt.Errorf("value is not a bool")
}
