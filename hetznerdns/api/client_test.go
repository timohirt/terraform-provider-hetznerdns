package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientCreateZoneSuccess(t *testing.T) {
	var requestBodyReader io.Reader
	responseBody := []byte(`{"zone":{"id":"12345","name":"mydomain.com","ttl":3600}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, requestBodyReader: &requestBodyReader, responseBodyJSON: responseBody}
	client := Client{apiToken: "irrelevant", doHTTPRequest: interceptRequestAndFakeResponse(config)}

	opts := CreateZoneOpts{Name: "mydomain.com", TTL: 3600}
	zone, err := client.CreateZone(opts)

	assert.NoError(t, err)
	assert.Equal(t, Zone{ID: "12345", Name: "mydomain.com", TTL: 3600}, *zone)
	assert.NotNil(t, requestBodyReader, "The request body should not be nil")
	jsonRequestBody, _ := ioutil.ReadAll(requestBodyReader)
	assert.Equal(t, `{"name":"mydomain.com","ttl":3600}`, string(jsonRequestBody))
}

func TestClientCreateZoneInvalidDomainName(t *testing.T) {
	client := Client{apiToken: "irrelevant"}
	opts := CreateZoneOpts{Name: "thisisinvalid", TTL: 3600}
	_, err := client.CreateZone(opts)

	assert.Error(t, err, "A invalid domain name was used. This should result in an error.")
}

func TestClientUpdateZoneSuccess(t *testing.T) {
	zoneWithUpdates := Zone{ID: "12345678", Name: "zone1.online", TTL: 3600}
	zoneWithUpdatesJSON := `{"id":"12345678","name":"zone1.online","ttl":3600}`
	var requestBodyReader io.Reader
	responseBody := []byte(`{"zone":{"id":"12345678","name":"zone1.online","ttl":3600}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, requestBodyReader: &requestBodyReader, responseBodyJSON: responseBody}
	client := Client{apiToken: "irrelevant", doHTTPRequest: interceptRequestAndFakeResponse(config)}

	updatedZone, err := client.UpdateZone(zoneWithUpdates)

	assert.NoError(t, err)
	assert.Equal(t, zoneWithUpdates, *updatedZone)
	assert.NotNil(t, requestBodyReader, "The request body should not be nil")
	jsonRequestBody, _ := ioutil.ReadAll(requestBodyReader)
	assert.Equal(t, zoneWithUpdatesJSON, string(jsonRequestBody))
}

func TestClientGetZone(t *testing.T) {
	responseBody := []byte(`{"zone":{"id":"12345678","name":"zone1.online","ttl":3600}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, responseBodyJSON: responseBody}
	client := Client{apiToken: "irrelevant", doHTTPRequest: interceptRequestAndFakeResponse(config)}

	zone, err := client.GetZone("12345678")

	assert.NoError(t, err)
	assert.Equal(t, Zone{ID: "12345678", Name: "zone1.online", TTL: 3600}, *zone)
}

func TestClientGetZoneByName(t *testing.T) {
	responseBody := []byte(`{"zones":[{"id":"12345678","name":"zone1.online","ttl":3600}]}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, responseBodyJSON: responseBody}
	client := Client{apiToken: "irrelevant", doHTTPRequest: interceptRequestAndFakeResponse(config)}

	zone, err := client.GetZoneByName("zone1.online")

	assert.NoError(t, err)
	assert.Equal(t, Zone{ID: "12345678", Name: "zone1.online", TTL: 3600}, *zone)
}

func TestClientDeleteZone(t *testing.T) {
	config := RequestConfig{responseHTTPStatus: http.StatusOK}
	client := Client{apiToken: "irrelevant", doHTTPRequest: interceptRequestAndFakeResponse(config)}

	err := client.DeleteZone("irrelevant")

	assert.NoError(t, err)
}

func TestClientGetRecord(t *testing.T) {
	responseBody := []byte(`{"record":{"zone_id":"wwwlsksjjenm","id":"12345678","name":"zone1.online","ttl":3600,"type":"A","value":"192.168.1.1"}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, responseBodyJSON: responseBody}
	client := Client{apiToken: "irrelevant", doHTTPRequest: interceptRequestAndFakeResponse(config)}

	zone, err := client.GetRecord("12345678")

	assert.NoError(t, err)
	assert.Equal(t, Record{ZoneID: "wwwlsksjjenm", ID: "12345678", Name: "zone1.online", TTL: 3600, Type: "A", Value: "192.168.1.1"}, *zone)
}

func TestClientCreateRecordSuccess(t *testing.T) {
	var requestBodyReader io.Reader
	responseBody := []byte(`{"record":{"zone_id":"wwwlsksjjenm","id":"12345678","name":"zone1.online","ttl":3600,"type":"A","value":"192.168.1.1"}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, requestBodyReader: &requestBodyReader, responseBodyJSON: responseBody}
	client := Client{apiToken: "irrelevant", doHTTPRequest: interceptRequestAndFakeResponse(config)}

	opts := CreateRecordOpts{ZoneID: "wwwlsksjjenm", Name: "zone1.online", TTL: 3600, Type: "A", Value: "192.168.1.1"}
	record, err := client.CreateRecord(opts)

	assert.NoError(t, err)
	assert.Equal(t, Record{ZoneID: "wwwlsksjjenm", ID: "12345678", Name: "zone1.online", TTL: 3600, Type: "A", Value: "192.168.1.1"}, *record)
	assert.NotNil(t, requestBodyReader, "The request body should not be nil")
	jsonRequestBody, _ := ioutil.ReadAll(requestBodyReader)
	assert.Equal(t, `{"zone_id":"wwwlsksjjenm","type":"A","name":"zone1.online","value":"192.168.1.1","ttl":3600}`, string(jsonRequestBody))
}

type RequestConfig struct {
	responseHTTPStatus int
	responseBodyJSON   []byte
	requestBodyReader  *io.Reader
}

func interceptRequestAndFakeResponse(config RequestConfig) func(apiToken string, method string, url string, body io.Reader) (*http.Response, error) {
	return func(apiToken string, method string, url string, body io.Reader) (*http.Response, error) {
		if body != nil {
			*config.requestBodyReader = body
		}

		var jsonBody io.ReadCloser = nil
		if config.responseBodyJSON != nil {
			jsonBody = ioutil.NopCloser(bytes.NewReader(config.responseBodyJSON))
		}
		resp := http.Response{StatusCode: config.responseHTTPStatus, Body: jsonBody}
		return &resp, nil
	}
}
