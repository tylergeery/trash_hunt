package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/game"
	"github.com/tylergeery/trash_hunt/test"
)

// TestCreateWithInvalidKey
// Expect that we cannot create a temporary auth token without a valid user token
func TestCreateWithInvalidKey(t *testing.T) {
	player := game.GetTestPlayer()
	token, _ := auth.CreateToken(player)

	resp := GetControllerResponse(t, "POST", "/v1/auth", nil, map[string]string{"Authorization": "Bearer " + token, "Content-Type": "application/json"})

	test.ExpectEqualInt64s(t, int64(http.StatusBadRequest), int64(resp.Result().StatusCode))
}

// TestCreateAuthToken
// Test that we can retrieve a tempory auth token from valid user key
func TestCreateAuthToken(t *testing.T) {
	player := game.GetTestPlayer()
	token, _ := auth.CreateToken(player)
	req := map[string]string{
		"key": token,
	}
	content, _ := json.Marshal(req)
	body := strings.NewReader(string(content))

	resp := GetControllerResponse(t, "POST", "/v1/auth", body, map[string]string{"Authorization": "Bearer " + token, "Content-Type": "application/json"})

	test.ExpectEqualInt64s(t, int64(http.StatusOK), int64(resp.Result().StatusCode))

	var m map[string]interface{}
	str, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(str, &m)

	test.ExpectNotEmptyString(t, m["token"].(string))
}
