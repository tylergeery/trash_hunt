package requests

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	routing "github.com/go-ozzo/ozzo-routing"
)

// CreateAuthTokenRequest - request body for CreateAuthToken request
type CreateAuthTokenRequest struct {
	Key string `json:"key"`
}

// NewCreateAuthTokenRequest - Create req from context
func NewCreateAuthTokenRequest(c *routing.Context) *CreateAuthTokenRequest {
	var v CreateAuthTokenRequest

	if c.Request.Body == nil {
		return &v
	}

	req, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal(req, &v)

	return &v
}

// Validate request
func (r *CreateAuthTokenRequest) Validate() error {
	if r.Key == "" {
		return errors.New("Empty auth token")
	}

	return nil
}
