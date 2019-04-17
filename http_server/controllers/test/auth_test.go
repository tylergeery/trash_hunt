package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/game"
	"github.com/tylergeery/trash_hunt/test"
)

func TestCreateWithInvalidKey(t *testing.T) {
	player := game.GetTestPlayer()
	token, _ := auth.CreateToken(player)

	resp := GetControllerResponse(t, "POST", "/auth", nil, map[string]string{"Authorization": "Bearer " + token})

	test.ExpectEqualInt64s(t, int64(http.StatusBadRequest), int64(resp.Result().StatusCode))

	body := []byte{}
	resp.Result().Body.Read(body)
	fmt.Println(resp.Result().Body)
	test.ExpectEqualString(t, "Invalid key supplied", string(body))
}

func TestCreateAuthToken(t *testing.T) {

}
