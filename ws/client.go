package ws

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zmoog/ws/v2/ws/identity"
)

type Client struct {
	client          *http.Client
	endpoint        string
	endpointVersion string
	identity        identity.Manager
}

func NewClient(identity identity.Manager, endpoint string, endpointVersion string) *Client {
	return &Client{
		client:          http.DefaultClient,
		identity:        identity,
		endpoint:        endpoint,
		endpointVersion: endpointVersion,
	}
}

func (c *Client) ListLocations() ([]Location, error) {
	token, err := c.identity.GetToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.endpoint+"/"+c.endpointVersion+"/locations", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token.ID)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var locations []Location
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, err
	}

	return locations, nil
}

func (c *Client) ListRooms(location string) ([]Room, error) {
	token, err := c.identity.GetToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.endpoint+"/"+c.endpointVersion+"/rooms", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token.ID)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rooms []Room
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}
