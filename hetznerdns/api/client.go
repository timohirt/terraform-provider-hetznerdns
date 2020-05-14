package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type doHTTPRequests func(apiToken string, method string, url string, body io.Reader) (*http.Response, error)

func defaultDoHTTPRequest(apiToken string, method string, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	log.Printf("[DEBUG] HTTP request to API %s %s", method, url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("[DEBUG] Error while creating HTTP request to API %s", err)
		return nil, err
	}

	req.Header.Add("Auth-API-Token", apiToken)
	req.Header.Add("Accept", "application/json; charset=utf-8")
	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("[DEBUG] Error while sending HTTP request to API %s", err)
		return nil, err
	}
	return resp, nil
}

// Client for the Hetzner DNS API.
type Client struct {
	apiToken      string
	doHTTPRequest doHTTPRequests
}

// NewClient creates a new API Client using a given api token.
func NewClient(apiToken string) (*Client, error) {
	return &Client{apiToken: apiToken, doHTTPRequest: defaultDoHTTPRequest}, nil
}

func (c *Client) doGetRequest(url string) (*http.Response, error) {
	return c.doHTTPRequest(c.apiToken, http.MethodGet, url, nil)
}

func (c *Client) doDeleteRequest(url string) (*http.Response, error) {
	return c.doHTTPRequest(c.apiToken, http.MethodDelete, url, nil)
}

func (c *Client) doPostRequest(url string, bodyJSON interface{}) (*http.Response, error) {
	reqJSON, err := json.Marshal(bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error serializing JSON body %s", err)
	}
	body := bytes.NewReader(reqJSON)

	return c.doHTTPRequest(c.apiToken, http.MethodPost, url, body)
}

func (c *Client) doPutRequest(url string, bodyJSON interface{}) (*http.Response, error) {
	reqJSON, err := json.Marshal(bodyJSON)
	if err != nil {
		return nil, fmt.Errorf("Error serializing JSON body %s", err)
	}
	body := bytes.NewReader(reqJSON)

	return c.doHTTPRequest(c.apiToken, http.MethodPut, url, body)
}

func readAndParseJSONBody(resp *http.Response, respType interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading HTTP response body %s", err)
	}

	err = json.Unmarshal(body, &respType)
	if err != nil {
		return fmt.Errorf("Error parsing JSON %s: %s", err, body)
	}
	return nil
}

// GetZone reads the current state of a DNS zone
func (c *Client) GetZone(id string) (*Zone, error) {
	resp, err := c.doGetRequest(fmt.Sprintf("https://dns.hetzner.com/api/v1/zones/%s", id))
	if err != nil {
		return nil, fmt.Errorf("Error getting zone %s: %s", id, err)
	}

	if resp.StatusCode == http.StatusOK {
		var response GetZoneResponse
		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}
		return &response.Zone, nil
	}

	return nil, fmt.Errorf("Error getting Zone. HTTP status %d unhandeled", resp.StatusCode)
}

// UpdateZone takes the passed state and updates the respective Zone
func (c *Client) UpdateZone(zone Zone) (*Zone, error) {
	resp, err := c.doPutRequest(fmt.Sprintf("https://dns.hetzner.com/api/v1/zones/%s", zone.ID), zone)
	if err != nil {
		return nil, fmt.Errorf("Error updating zone %s: %s", zone.ID, err)
	}

	if resp.StatusCode == http.StatusOK {
		var response ZoneResponse
		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}
		return &response.Zone, nil
	}

	return nil, fmt.Errorf("Error updating Zone. HTTP status %d unhandeled", resp.StatusCode)
}

// DeleteZone deletes a given DNS zone
func (c *Client) DeleteZone(id string) error {
	resp, err := c.doDeleteRequest(fmt.Sprintf("https://dns.hetzner.com/api/v1/zones/%s", id))
	if err != nil {
		return fmt.Errorf("Error deleting zone %s: %s", id, err)
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("Error deleting Zone. HTTP status %d unhandeled", 400)
}

// GetZoneByName reads the current state of a DNS zone with a given name
func (c *Client) GetZoneByName(name string) (*Zone, error) {
	resp, err := c.doGetRequest(fmt.Sprintf("https://dns.hetzner.com/api/v1/zones?name=%s", name))
	if err != nil {
		return nil, fmt.Errorf("Error getting zone %s: %s", name, err)
	}

	if resp.StatusCode == http.StatusOK {
		var response *GetZonesByNameResponse
		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}

		if len(response.Zones) != 1 {
			return nil, fmt.Errorf("Error getting zone '%s'. No matching zone or multiple matching zones found", name)
		}

		return &response.Zones[0], nil
	}

	return nil, fmt.Errorf("Error getting Zone. HTTP status %d unhandeled", resp.StatusCode)
}

// CreateZoneOpts covers all parameters used to create a new DNS zone
type CreateZoneOpts struct {
	Name string
	TTL  int
}

// CreateZone creates a new DNS zone
func (c *Client) CreateZone(opts CreateZoneOpts) (*Zone, error) {

	if !strings.Contains(opts.Name, ".") {
		return nil, fmt.Errorf("Error creating zone. The name '%s' is not a valid domain name with top level domain", opts.Name)
	}

	reqBody := CreateZoneRequest{Name: opts.Name, TTL: opts.TTL}
	resp, err := c.doPostRequest("https://dns.hetzner.com/api/v1/zones", reqBody)
	if err != nil {
		return nil, fmt.Errorf("Error creating zone %s", err)
	}

	if resp.StatusCode == http.StatusOK {
		var response CreateZoneResponse
		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}

		return &response.Zone, nil
	}

	return nil, fmt.Errorf("Error creating Zone. HTTP status %d unhandeled", resp.StatusCode)
}
