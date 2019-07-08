package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/tylergeery/trash_hunt/src/game"
	"github.com/tylergeery/trash_hunt/src/test"
)

// TestPlayerCreateFailures
// Expect that we cannot create playes for expected reasons
func TestPlayerCreateFailures(t *testing.T) {
	type testCase struct {
		body        map[string]string
		status_code int
		expected    string
	}
	failure_cases := []testCase{
		testCase{map[string]string{}, http.StatusBadRequest, "Password must be at least 8 characters\n"},
		testCase{
			map[string]string{
				"pw":    "!2312125sd-test",
				"email": "tt",
			},
			http.StatusBadRequest,
			"Invalid email format: tt\n",
		},
		testCase{
			map[string]string{
				"pw":    "!2312125sd-test",
				"email": "tyler@test.com",
			},
			http.StatusBadRequest,
			"Invalid name: \n",
		},
	}

	for _, c := range failure_cases {
		player_encoded, _ := json.Marshal(c.body)
		body := strings.NewReader(string(player_encoded))
		resp := GetControllerResponse(t, "POST", "/v1/player/", body, map[string]string{"Content-Type": "application/json"})
		playerResponse, _ := ioutil.ReadAll(resp.Result().Body)

		test.ExpectEqualInt64s(t, int64(c.status_code), int64(resp.Result().StatusCode))
		test.ExpectEqualString(t, c.expected, string(playerResponse))
	}
}

// TestAttemptCreatePlayerAndCannotReuseEmail
// Tries to create player with re-used email
func TestAttemptCreatePlayerAndCannotReuseEmail(t *testing.T) {
	player := map[string]string{
		"name":  "Test Player",
		"pw":    "1234213kdsl;kdg",
		"email": game.GetTestEmail("email-reuse"),
	}
	player_encoded, _ := json.Marshal(player)
	body := strings.NewReader(string(player_encoded))
	resp := GetControllerResponse(t, "POST", "/v1/player/", body, map[string]string{"Content-Type": "application/json"})

	test.ExpectEqualInt64s(t, int64(http.StatusOK), int64(resp.Result().StatusCode))

	player["name"] = "Different Test Name"
	player_encoded, _ = json.Marshal(player)
	body = strings.NewReader(string(player_encoded))
	resp = GetControllerResponse(t, "POST", "/v1/player/", body, map[string]string{"Content-Type": "application/json"})
	playerResponse, _ := ioutil.ReadAll(resp.Result().Body)

	test.ExpectEqualInt64s(t, int64(http.StatusBadRequest), int64(resp.Result().StatusCode))
	test.ExpectEqualString(t, fmt.Sprintf("Email %s belongs to an existing user\n", player["email"]), string(playerResponse))
}

func TestCreateUpdateAndDeletePlayer(t *testing.T) {
	player := map[string]string{
		"name":  "Test Player",
		"pw":    "1234213kdsl;kdg",
		"email": game.GetTestEmail("email-reuse"),
	}
	player_encoded, _ := json.Marshal(player)
	body := strings.NewReader(string(player_encoded))
	resp := GetControllerResponse(t, "POST", "/v1/player/", body, map[string]string{"Content-Type": "application/json"})
	createResponse, _ := ioutil.ReadAll(resp.Result().Body)

	var createdPlayerResponse map[string]interface{}
	_ = json.Unmarshal(createResponse, &createdPlayerResponse)
	fmt.Printf("createResponse: %s", createResponse)
	fmt.Printf("createdPlayerResponse: %s", createdPlayerResponse)

	test.ExpectEqualInt64s(t, int64(http.StatusOK), int64(resp.Result().StatusCode))

	player["name"] = "Different Test Name"
	player_encoded, _ = json.Marshal(player)
	body = strings.NewReader(string(player_encoded))
	headers := map[string]string{"Content-Type": "application/json"}
	headers["Authorization"] = "Bearer " + createdPlayerResponse["token"].(string)
	resp = GetControllerResponse(t, "PUT", "/v1/player/", body, headers)
	playerResponse, _ := ioutil.ReadAll(resp.Result().Body)

	var updatedPlayer game.Player
	_ = json.Unmarshal(playerResponse, &updatedPlayer)
	fmt.Printf("playerResponse: %s", playerResponse)
	fmt.Printf("updatedPlayer: %v", updatedPlayer)

	test.ExpectEqualInt64s(t, int64(http.StatusOK), int64(resp.Result().StatusCode))
	test.ExpectEqualString(t, "Different Test Name", updatedPlayer.Name)
}
