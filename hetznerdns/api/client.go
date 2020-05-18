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

	return nil, fmt.Errorf("Error getting Zone. HTTP status %d unhandled", resp.StatusCode)
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
	return fmt.Errorf("Error deleting Zone. HTTP status %d unhandled", resp.StatusCode)
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

	return nil, fmt.Errorf("Error getting Zone. HTTP status %d unhandled", resp.StatusCode)
}

// CreateZoneOpts covers all parameters used to create a new DNS zone
type CreateZoneOpts struct {
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
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

	return nil, fmt.Errorf("Error creating Zone. HTTP status %d unhandled", resp.StatusCode)
}

// GetRecordByName reads the current state of a DNS Record with a given name and zone id
func (c *Client) GetRecordByName(zoneID string, name string) (*Record, error) {
	resp, err := c.doGetRequest(fmt.Sprintf("https://dns.hetzner.com/api/v1/records?zone_id=%s", zoneID))
	if err != nil {
		return nil, fmt.Errorf("Error getting record %s: %s", name, err)
	}

	if resp.StatusCode == http.StatusOK {
		var response *RecordsResponse
		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}

		if len(response.Records) == 0 {
			return nil, fmt.Errorf("Error getting record '%s'. It seems there are no records in zone %s at all", name, zoneID)
		}

		for _, record := range response.Records {
			if record.Name == name {
				return &record, nil
			}
		}

		return nil, fmt.Errorf("Error getting record '%s'. There are records in zone %s, but %s isn't included", name, zoneID, name)
	}

	return nil, fmt.Errorf("Error getting Record. HTTP status %d unhandled", resp.StatusCode)
}

// GetRecord reads the current state of a DNS Record
func (c *Client) GetRecord(recordID string) (*Record, error) {
	resp, err := c.doGetRequest(fmt.Sprintf("https://dns.hetzner.com/api/v1/records/%s", recordID))
	if err != nil {
		return nil, fmt.Errorf("Error getting record %s: %s", recordID, err)
	}

	if resp.StatusCode == http.StatusOK {
		var response *RecordResponse
		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, fmt.Errorf("Error Reading json response of get record %s request: %s", recordID, err)
		}

		return &response.Record, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	return nil, fmt.Errorf("Error getting Record. HTTP status %d unhandled", resp.StatusCode)
}

// CreateRecordOpts covers all parameters used to create a new DNS record
type CreateRecordOpts struct {
	ZoneID string `json:"zone_id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	TTL    int    `json:"ttl"`
}

// CreateRecord create a new DNS records
func (c *Client) CreateRecord(opts CreateRecordOpts) (*Record, error) {
	reqBody := CreateRecordRequest{ZoneID: opts.ZoneID, Name: opts.Name, TTL: opts.TTL, Type: opts.Type, Value: opts.Value}
	resp, err := c.doPostRequest("https://dns.hetzner.com/api/v1/records", reqBody)
	if err != nil {
		return nil, fmt.Errorf("Error creating record %s: %s", opts.Name, err)
	}

	if resp.StatusCode == http.StatusOK {
		var response RecordResponse
		err = readAndParseJSONBody(resp, &response)
		if err != nil {
			return nil, err
		}

		return &response.Record, nil
	}

	return nil, fmt.Errorf("Error creating Record. HTTP status %d unhandled", resp.StatusCode)
}

// DeleteRecord deletes a given record
func (c *Client) DeleteRecord(id string) error {
	resp, err := c.doDeleteRequest(fmt.Sprintf("https://dns.hetzner.com/api/v1/records/%s", id))
	if err != nil {
		return fmt.Errorf("Error deleting zone %s: %s", id, err)
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("Error deleting Record. HTTP status %d unhandled", resp.StatusCode)
}
