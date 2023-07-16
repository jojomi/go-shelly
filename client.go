package shelly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	var cmd string
	switch x.device.Type {
	case PlugS:
		cmd = "relay/0"
	case BulbRGBW:
		cmd = "light/0"
	}
	if cmd == "" {
		return fmt.Errorf("can't change power state, invalid device type %s", x.device.Type)
	}

	val := "on"
	if !on {
		val = "off"
	}
	response, err := x.post(cmd+"?turn="+val, nil)
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
	var cmd string
	switch x.device.Type {
	case PlugS:
		cmd = "relay/0"
	case BulbRGBW:
		cmd = "light/0"
	}
	if cmd == "" {
		return false, fmt.Errorf("can't check for power on, invalid device type %s", x.device.Type)
	}

	response, err := x.get(cmd, nil)
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
	if x.device.Type != PlugS {
		return 0, fmt.Errorf("can't read power for device type %s", x.device.Type)
	}

	response, err := x.get("meter/0", nil)
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

func (x *Client) SetModeWhite(opts ModeWhiteOptions) error {
	if x.device.Type != BulbRGBW {
		return fmt.Errorf("can't set light mode to white for device type %s", x.device.Type)
	}

	cmd := "light/0/set"
	data := map[string]any{
		"mode": "white",
	}

	if opts.Brightness != nil {
		data["brightness"] = *opts.Brightness
	}
	if opts.Transition != nil {
		data["transition"] = *opts.Transition
	}
	if opts.Temp != nil {
		data["temp"] = *opts.Temp
	}

	response, err := x.get(cmd, &data)
	if err != nil {
		return err
	}
	// TODO add verification through checking response.String()
	_ = response

	return nil
}

func (x *Client) IsModeWhite() (bool, error) {
	mode, err := x.getMode()
	if err != nil {
		return false, err
	}
	return mode == "white", nil
}

func (x *Client) IsModeColor() (bool, error) {
	mode, err := x.getMode()
	if err != nil {
		return false, err
	}
	return mode == "color", nil
}

func (x *Client) SetModeColor(opts ModeColorOptions) error {
	if x.device.Type != BulbRGBW {
		return fmt.Errorf("can't set light mode to white for device %s", x.device.Type)
	}

	cmd := "light/0/set"
	data := map[string]any{
		"mode": "color",
	}

	if opts.Color != nil {
		data["red"] = int(opts.Color.R * 255)
		data["green"] = int(opts.Color.G * 255)
		data["blue"] = int(opts.Color.B * 255)
	}
	if opts.Brightness != nil {
		data["brightness"] = *opts.Brightness
	}
	if opts.Gain != nil {
		data["gain"] = *opts.Gain
	}
	if opts.Transition != nil {
		data["transition"] = *opts.Transition
	}

	response, err := x.get(cmd, &data)
	if err != nil {
		return err
	}
	// TODO verify state in response.String()
	_ = response

	return nil
}

func (x *Client) httpClient() gohttp.Client {
	return gohttp.NewBuilder().
		SetConnectionTimeout(10 * time.Second).
		SetResponseTimeout(20 * time.Second).
		SetUserAgent("go-shelly").
		Build()
}

func (x *Client) get(cmd string, data *map[string]any) (*core.Response, error) {
	apiURL := x.device.GetAPIBaseURL() + cmd

	if data != nil {
		u, err := url.Parse(apiURL)
		if err != nil {
			return nil, err
		}

		q := u.Query()
		for key, value := range *data {
			q.Add(key, fmt.Sprintf("%v", value))
		}

		u.RawQuery = q.Encode()

		apiURL = u.String()
	}

	return x.httpClient().Get(apiURL)
}

func (x *Client) post(cmd string, data *map[string]any) (*core.Response, error) {
	body := ""
	headers := http.Header{}
	if data != nil {
		var err error
		bodyBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = string(bodyBytes)
		headers.Add("Content-Type", "application/json")
	}
	apiURL := x.device.GetAPIBaseURL() + cmd
	return x.httpClient().Post(apiURL, body, headers)
}

func (x *Client) getMode() (string, error) {
	if x.device.Type != BulbRGBW {
		return "", fmt.Errorf("can't read mode for device type %s", x.device.Type)
	}

	response, err := x.get("light/0", nil)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(response.String()), &data)
	if err != nil {
		return "", err
	}

	mode, err := x.getStringByPath(data, "mode")
	if err != nil {
		return "", err
	}

	return mode, nil
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

func (x *Client) getStringByPath(data map[string]interface{}, path string) (string, error) {
	keys := strings.Split(path, ".")

	var (
		value any
		ok    bool
	)
	for _, key := range keys {
		value, ok = data[key]
		if !ok {
			return "", fmt.Errorf("key '%s' not found", key)
		}

		data, ok = value.(map[string]interface{})
		if !ok {
			break
		}
	}

	return fmt.Sprintf("%v", value), nil
}
