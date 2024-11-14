package ws

import (
	"encoding/json"
	"net/http"

	"github.com/zmoog/ws/ws/identity"
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

func (c *Client) ListLocations() ([]Location, error) {
	token, err := c.identity.GetToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.endpoint+"/v2.6/locations", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

	req, err := http.NewRequest("GET", c.endpoint+"/v2.6/rooms", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rooms []Room
	if err := json.NewDecoder(resp.Body).Decode(&rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}
