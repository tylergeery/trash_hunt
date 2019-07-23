package middleware

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/game"
)

type TestWriter struct {
	code int
}

func (t TestWriter) Write(b []byte) (int, error) { return len(b), nil }
func (t TestWriter) WriteHeader(statusCode int)  { t.code = statusCode }
func (t TestWriter) Header() http.Header {
	return map[string][]string{
		"hello": []string{"world"},
	}
}
func TestLogRequestAndValidateToken(t *testing.T) {
	count := 0
	foundPlayerID := int64(0)
	player := game.PlayerNew(34, "", "", "", "", "", game.PlayerStatusActive, "", "")
	handler := func(c *routing.Context) error {
		count++
		foundPlayerID = c.Get("PlayerID").(int64)

		return nil
	}
	request, _ := http.NewRequest("GET", "/auth", strings.NewReader(""))
	writer := TestWriter{}
	token, _ := auth.CreateToken(player)
	request.Header.Add("Authorization", "Bearer "+token)

	ctx := routing.NewContext(writer, request, LogRequest(), ValidateToken(), handler)
	ctx.Next()

	if count == 0 {
		t.Fatalf("Handler was never called")
	}

	if foundPlayerID != player.ID {
		t.Fatalf("Could not find playerID in auth token: %d", foundPlayerID)
	}
}
