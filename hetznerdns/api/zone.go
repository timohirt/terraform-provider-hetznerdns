package api

// Zone represents a DNS Zone
type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
}

// CreateZoneRequest represents the body of a POST Zone request
type CreateZoneRequest struct {
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
}

// CreateZoneResponse represents the content of a POST Zone response
type CreateZoneResponse struct {
	Zone Zone `json:"zone"`
}

// GetZoneResponse represents the content of a GET Zone request
type GetZoneResponse struct {
	Zone Zone `json:"zone"`
}

// ZoneResponse represents the content of response containing a Zone
type ZoneResponse struct {
	Zone Zone `json:"zone"`
}

// GetZonesByNameResponse represents the content of a GET Zones response
type GetZonesByNameResponse struct {
	Zones []Zone `json:"zones"`
}
