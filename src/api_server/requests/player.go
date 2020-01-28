package requests

// PlayerCreateRequest - Request structure for PlayerCreate request
type PlayerCreateRequest struct {
	Email      string `json:"email"`
	Pw         string `json:"pw"`
	Username       string `json:"username"`
}

// PlayerUpdateRequest - Request structure for PlayerUpdate request
type PlayerUpdateRequest struct {
	Email      string `json:"email"`
	Username       string `json:"username"`
}

// PlayerLoginRequest - Request structure for PlayerLogin request
type PlayerLoginRequest struct {
	Email string `json:"email"`
	Pw    string `json:"pw"`
}
