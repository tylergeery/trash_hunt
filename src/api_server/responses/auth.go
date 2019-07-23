package responses

import "encoding/json"

// AuthTokenCreateResponse - response structure for creating new token
type AuthTokenCreateResponse struct {
	Token string `json:"token"`
}

// MarshalJSON - encode response as JSON
func (res AuthTokenCreateResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"token": res.Token})
}
