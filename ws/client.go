package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zmoog/ws/v2/ws/identity"
)

type Client struct {
	client   *http.Client
	endpoint string
	identity identity.Manager
}

func NewClient(identity identity.Manager, endpoint string) *Client {
	return &Client{
		client:   http.DefaultClient,
		identity: identity,
		endpoint: endpoint,
	}
}

func (c *Client) ListDevices() ([]Device, error) {
	token, err := c.identity.GetToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint+"/"+"ListDevices", strings.NewReader("{}"))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.ID)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var devices Devices
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, err
	}

	return devices.Devices, nil
}

func (c *Client) GetDevice(deviceName string) (Device, error) {
	token, err := c.identity.GetToken()
	if err != nil {
		return Device{}, err
	}

	r := struct {
		Name string `json:"name"`
	}{
		Name: deviceName,
	}
	jsonReq, err := json.Marshal(r)
	if err != nil {
		return Device{}, err
	}

	req, err := http.NewRequest("POST", c.endpoint+"/GetDevice", strings.NewReader(string(jsonReq)))
	if err != nil {
		return Device{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.ID)

	resp, err := c.client.Do(req)
	if err != nil {
		return Device{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Device{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var device Device
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return Device{}, err
	}

	return device, nil
}
