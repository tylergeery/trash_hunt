package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/tylergeery/trash_hunt/model"
	"github.com/tylergeery/trash_hunt/test"
)

// TestPlayerCreateFailures
// Expect that we cannot create playes for expected reasons
func TestPlayerCreateFailures(t *testing.T) {
	type testCase struct {
		body       map[string]string
		statusCode int
		expected   string
	}
	failureCases := []testCase{
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
			"Invalid username: \n",
		},
	}

	for _, c := range failureCases {
		playerEncoded, _ := json.Marshal(c.body)
		body := strings.NewReader(string(playerEncoded))
		resp := GetControllerResponse(t, "POST", "/v1/player/", body, map[string]string{"Content-Type": "application/json"})
		playerResponse, _ := ioutil.ReadAll(resp.Result().Body)

		test.ExpectEqualInt64s(t, int64(c.statusCode), int64(resp.Result().StatusCode))
		test.ExpectEqualString(t, c.expected, string(playerResponse))
	}
}

// TestAttemptCreatePlayerAndCannotReuseEmail
// Tries to create player with re-used email
func TestAttemptCreatePlayerAndCannotReuseEmail(t *testing.T) {
	player := map[string]string{
		"username":  model.GetTestUsername("email-reuse"),
		"pw":    "1234213kdsl;kdg",
		"email": model.GetTestEmail("email-reuse"),
	}
	playerEncoded, _ := json.Marshal(player)
	body := strings.NewReader(string(playerEncoded))
	resp := GetControllerResponse(t, "POST", "/v1/player/", body, map[string]string{"Content-Type": "application/json"})

	test.ExpectEqualInt64s(t, int64(http.StatusCreated), int64(resp.Result().StatusCode))

	player["username"] = "difftestname2"
	playerEncoded, _ = json.Marshal(player)
	body = strings.NewReader(string(playerEncoded))
	resp = GetControllerResponse(t, "POST", "/v1/player/", body, map[string]string{"Content-Type": "application/json"})
	playerResponse, _ := ioutil.ReadAll(resp.Result().Body)

	test.ExpectEqualInt64s(t, int64(http.StatusBadRequest), int64(resp.Result().StatusCode))
	test.ExpectEqualString(t, fmt.Sprintf("Email %s belongs to an existing user\n", player["email"]), string(playerResponse))
}

func TestCreateUpdateAndDeletePlayer(t *testing.T) {
	player := map[string]string{
		"username":  model.GetTestUsername("test-crud"),
		"pw":    "1234213kdsl;kdg",
		"email": model.GetTestEmail("test-crud"),
	}
	playerEncoded, _ := json.Marshal(player)
	body := strings.NewReader(string(playerEncoded))
	resp := GetControllerResponse(t, "POST", "/v1/player/", body, map[string]string{"Content-Type": "application/json"})
	createResponse, _ := ioutil.ReadAll(resp.Result().Body)

	var createdPlayerResponse map[string]interface{}
	_ = json.Unmarshal(createResponse, &createdPlayerResponse)

	test.ExpectEqualInt64s(t, int64(http.StatusCreated), int64(resp.Result().StatusCode))

	player["username"] = model.GetTestUsername("other-crud-name")
	playerEncoded, _ = json.Marshal(player)
	body = strings.NewReader(string(playerEncoded))
	headers := map[string]string{"Content-Type": "application/json"}
	headers["Authorization"] = "Bearer " + createdPlayerResponse["token"].(string)
	resp = GetControllerResponse(t, "PUT", "/v1/player/", body, headers)
	playerResponse, _ := ioutil.ReadAll(resp.Result().Body)

	var updatedPlayer model.Player
	_ = json.Unmarshal(playerResponse, &updatedPlayer)

	test.ExpectEqualInt64s(t, int64(http.StatusOK), int64(resp.Result().StatusCode))
	test.ExpectEqualString(t, player["username"], updatedPlayer.Username)

	// Delete a player
	headers = map[string]string{"Content-Type": "application/json"}
	headers["Authorization"] = "Bearer " + createdPlayerResponse["token"].(string)
	resp = GetControllerResponse(t, "DELETE", "/v1/player/", nil, headers)
	removedResponse, _ := ioutil.ReadAll(resp.Result().Body)

	var removedPlayerResponse map[string]interface{}
	_ = json.Unmarshal(removedResponse, &removedPlayerResponse)

	test.ExpectEqualInt64s(t, int64(http.StatusNoContent), int64(resp.Result().StatusCode))
}

func TestResetPassword(t *testing.T) {
	player := map[string]string{
		"username":  model.GetTestUsername("reset-password"),
		"pw":    "1234213kdsl;kdg",
		"email": model.GetTestEmail("reset-password"),
	}
	playerEncoded, _ := json.Marshal(player)
	body := strings.NewReader(string(playerEncoded))
	resp := GetControllerResponse(t, "POST", "/v1/player/", body, map[string]string{"Content-Type": "application/json"})
	createResponse, _ := ioutil.ReadAll(resp.Result().Body)

	var createdPlayerResponse map[string]interface{}
	_ = json.Unmarshal(createResponse, &createdPlayerResponse)

	test.ExpectEqualInt64s(t, int64(http.StatusCreated), int64(resp.Result().StatusCode))

	// TODO: reset password
}
