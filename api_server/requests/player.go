package requests

// PlayerCreateRequest - Request structure for PlayerCreate request
type PlayerCreateRequest struct {
	Email      string `json:"email"`
	Pw         string `json:"pw"`
	Name       string `json:"name"`
	FacebookID string `json:"facebookID"`
}

// PlayerLoginRequest - Request structure for PlayerLogin request
type PlayerLoginRequest struct {
	Email string `json:"id"`
	Pw    string `json:"pw"`
}
