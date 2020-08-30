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
	client := createTestClient(config)

	opts := CreateZoneOpts{Name: "mydomain.com", TTL: 3600}
	zone, err := client.CreateZone(opts)

	assert.NoError(t, err)
	assert.Equal(t, Zone{ID: "12345", Name: "mydomain.com", TTL: 3600}, *zone)
	assert.NotNil(t, requestBodyReader, "The request body should not be nil")
	jsonRequestBody, _ := ioutil.ReadAll(requestBodyReader)
	assert.Equal(t, `{"name":"mydomain.com","ttl":3600}`, string(jsonRequestBody))
}

func TestClientCreateZoneInvalidDomain(t *testing.T) {
	responseBody := []byte(`{"zone":{"id":"","name":"","ttl":0,"registrar":"","legacy_dns_host":"","legacy_ns":null,"ns":null,"created":"","verified":"","modified":"","project":"","owner":"","permission":"","zone_type":{"id":"","name":"","description":"","prices":null},"status":"","paused":false,"is_secondary_dns":false,"txt_verification":{"name":"","token":""},"records_count":0},"error":{"message":"422 : invalid TLD","code":422}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusUnprocessableEntity, responseBodyJSON: responseBody}

	client := createTestClient(config)
	opts := CreateZoneOpts{Name: "this.is.invalid", TTL: 3600}
	_, err := client.CreateZone(opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API returned HTTP 422 Unprocessable Entity error with message: '422 : invalid TLD'")
}

func TestClientCreateZoneInvalidTLD(t *testing.T) {
	var irrelevantConfig RequestConfig
	client := createTestClient(irrelevantConfig)
	opts := CreateZoneOpts{Name: "thisisinvalid", TTL: 3600}
	_, err := client.CreateZone(opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'thisisinvalid' is not a valid domain")
}

func TestClientUpdateZoneSuccess(t *testing.T) {
	zoneWithUpdates := Zone{ID: "12345678", Name: "zone1.online", TTL: 3600}
	zoneWithUpdatesJSON := `{"id":"12345678","name":"zone1.online","ttl":3600}`
	var requestBodyReader io.Reader
	responseBody := []byte(`{"zone":{"id":"12345678","name":"zone1.online","ttl":3600}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, requestBodyReader: &requestBodyReader, responseBodyJSON: responseBody}
	client := createTestClient(config)

	updatedZone, err := client.UpdateZone(zoneWithUpdates)

	assert.NoError(t, err)
	assert.Equal(t, zoneWithUpdates, *updatedZone)
	assert.NotNil(t, requestBodyReader, "The request body should not be nil")
	jsonRequestBody, _ := ioutil.ReadAll(requestBodyReader)
	assert.Equal(t, zoneWithUpdatesJSON, string(jsonRequestBody))
}

func TestClientGetAllZonesWhenThereAreNoZonesCreatedYet(t *testing.T) {
	responseBody := []byte(`{"zones":[]}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, responseBodyJSON: responseBody}
	client := createTestClient(config)

	allZones, err := client.GetAllZones()

	assert.NoError(t, err)
	assert.Empty(t, allZones)
}

func TestClientGetAllZones(t *testing.T) {
	responseBody := []byte(`{"zones":[{"id":"12345678","name":"zone1.online","ttl":3600},{"id":"98765432","name":"zone2.online","ttl":3600}]}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, responseBodyJSON: responseBody}
	client := createTestClient(config)

	allZones, err := client.GetAllZones()

	assert.NoError(t, err)
	assert.Contains(t, allZones, Zone{ID: "12345678", Name: "zone1.online", TTL: 3600})
	assert.Contains(t, allZones, Zone{ID: "98765432", Name: "zone2.online", TTL: 3600})
}

func TestClientGetZone(t *testing.T) {
	responseBody := []byte(`{"zone":{"id":"12345678","name":"zone1.online","ttl":3600}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, responseBodyJSON: responseBody}
	client := createTestClient(config)

	zone, err := client.GetZone("12345678")

	assert.NoError(t, err)
	assert.Equal(t, Zone{ID: "12345678", Name: "zone1.online", TTL: 3600}, *zone)
}

func TestClientGetZoneReturnNilIfNotFound(t *testing.T) {
	config := RequestConfig{responseHTTPStatus: http.StatusNotFound}
	client := createTestClient(config)

	zone, err := client.GetZone("12345678")

	assert.NoError(t, err)
	assert.Nil(t, zone)
}

func TestClientGetZoneByName(t *testing.T) {
	responseBody := []byte(`{"zones":[{"id":"12345678","name":"zone1.online","ttl":3600}]}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, responseBodyJSON: responseBody}
	client := createTestClient(config)

	zone, err := client.GetZoneByName("zone1.online")

	assert.NoError(t, err)
	assert.Equal(t, Zone{ID: "12345678", Name: "zone1.online", TTL: 3600}, *zone)
}

func TestClientGetZoneByNameReturnNilIfnotFound(t *testing.T) {
	config := RequestConfig{responseHTTPStatus: http.StatusNotFound}
	client := createTestClient(config)

	zone, err := client.GetZoneByName("zone1.online")

	assert.NoError(t, err)
	assert.Nil(t, zone)
}

func TestClientDeleteZone(t *testing.T) {
	config := RequestConfig{responseHTTPStatus: http.StatusOK}
	client := createTestClient(config)

	err := client.DeleteZone("irrelevant")

	assert.NoError(t, err)
}

func TestClientGetRecord(t *testing.T) {
	aTTL := 3600
	responseBody := []byte(`{"record":{"zone_id":"wwwlsksjjenm","id":"12345678","name":"zone1.online","ttl":3600,"type":"A","value":"192.168.1.1"}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, responseBodyJSON: responseBody}
	client := createTestClient(config)

	record, err := client.GetRecord("12345678")

	assert.NoError(t, err)
	assert.Equal(t, Record{ZoneID: "wwwlsksjjenm", ID: "12345678", Name: "zone1.online", TTL: &aTTL, Type: "A", Value: "192.168.1.1"}, *record)
}

func TestClientGetRecordWithUndefinedTTL(t *testing.T) {
	responseBody := []byte(`{"record":{"zone_id":"wwwlsksjjenm","id":"12345678","name":"zone1.online","type":"A","value":"192.168.1.1"}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, responseBodyJSON: responseBody}
	client := createTestClient(config)

	record, err := client.GetRecord("12345678")

	assert.NoError(t, err)
	assert.Equal(t, Record{ZoneID: "wwwlsksjjenm", ID: "12345678", Name: "zone1.online", TTL: nil, Type: "A", Value: "192.168.1.1"}, *record)
}

func TestClientGetRecordReturnNilIfNotFound(t *testing.T) {
	config := RequestConfig{responseHTTPStatus: http.StatusNotFound}
	client := createTestClient(config)

	record, err := client.GetRecord("irrelevant")

	assert.NoError(t, err)
	assert.Nil(t, record)
}

func TestClientCreateRecordSuccess(t *testing.T) {
	var requestBodyReader io.Reader
	responseBody := []byte(`{"record":{"zone_id":"wwwlsksjjenm","id":"12345678","name":"zone1.online","ttl":3600,"type":"A","value":"192.168.1.1"}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, requestBodyReader: &requestBodyReader, responseBodyJSON: responseBody}
	client := createTestClient(config)

	aTTL := 3600
	opts := CreateRecordOpts{ZoneID: "wwwlsksjjenm", Name: "zone1.online", TTL: &aTTL, Type: "A", Value: "192.168.1.1"}
	record, err := client.CreateRecord(opts)

	assert.NoError(t, err)
	assert.Equal(t, Record{ZoneID: "wwwlsksjjenm", ID: "12345678", Name: "zone1.online", TTL: &aTTL, Type: "A", Value: "192.168.1.1"}, *record)
	assert.NotNil(t, requestBodyReader, "The request body should not be nil")
	jsonRequestBody, _ := ioutil.ReadAll(requestBodyReader)
	assert.Equal(t, `{"zone_id":"wwwlsksjjenm","type":"A","name":"zone1.online","value":"192.168.1.1","ttl":3600}`, string(jsonRequestBody))
}

func TestClientRecordZone(t *testing.T) {
	config := RequestConfig{responseHTTPStatus: http.StatusOK}
	client := createTestClient(config)

	err := client.DeleteRecord("irrelevant")

	assert.NoError(t, err)
}

func TestClientUpdateRecordSuccess(t *testing.T) {
	aTTL := 3600
	recordWithUpdates := Record{ZoneID: "wwwlsksjjenm", ID: "12345678", Name: "zone2.online", TTL: &aTTL, Type: "A", Value: "192.168.1.1"}
	recordWithUpdatesJSON := `{"zone_id":"wwwlsksjjenm","id":"12345678","type":"A","name":"zone2.online","value":"192.168.1.1","ttl":3600}`
	var requestBodyReader io.Reader
	responseBody := []byte(`{"record":{"zone_id":"wwwlsksjjenm","id":"12345678","type":"A","name":"zone2.online","value":"192.168.1.1","ttl":3600}}`)
	config := RequestConfig{responseHTTPStatus: http.StatusOK, requestBodyReader: &requestBodyReader, responseBodyJSON: responseBody}
	client := createTestClient(config)

	updatedRecord, err := client.UpdateRecord(recordWithUpdates)

	assert.NoError(t, err)
	assert.Equal(t, recordWithUpdates, *updatedRecord)
	assert.NotNil(t, requestBodyReader, "The request body should not be nil")
	jsonRequestBody, _ := ioutil.ReadAll(requestBodyReader)
	assert.Equal(t, recordWithUpdatesJSON, string(jsonRequestBody))
}

func TestClientHandleUnauthorizedRequest(t *testing.T) {
	responseBody := []byte(`{"message":"Invalid API key"}`)
	config := RequestConfig{responseHTTPStatus: http.StatusUnauthorized, responseBodyJSON: responseBody}
	client := createTestClient(config)

	opts := CreateZoneOpts{Name: "mydomain.com", TTL: 3600}
	_, err := client.CreateZone(opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'Invalid API key'", "Error message didn't contain error message from API.")
}

type RequestConfig struct {
	responseHTTPStatus int
	responseBodyJSON   []byte
	requestBodyReader  *io.Reader
}

func createTestClient(config RequestConfig) Client {
	fakeHTTPClient := TestClient{config: config}
	createFakeHTTPClient := func() *http.Client {
		return &http.Client{Transport: fakeHTTPClient}
	}
	return Client{apiToken: "irrelevant", createHTTPClient: createFakeHTTPClient}
}

type TestClient struct {
	config RequestConfig
}

// See https://golang.org/pkg/net/http/#RoundTripper
func (f TestClient) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil && f.config.requestBodyReader != nil {
		*f.config.requestBodyReader = req.Body
	}

	var jsonBody io.ReadCloser = nil
	if f.config.responseBodyJSON != nil {
		jsonBody = ioutil.NopCloser(bytes.NewReader(f.config.responseBodyJSON))
	}
	resp := http.Response{StatusCode: f.config.responseHTTPStatus, Body: jsonBody}
	return &resp, nil
}
