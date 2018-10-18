package auth

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/tylergeery/trash_hunt/game"
)

func TestCreatingAndExtractingToken(t *testing.T) {
	var playerID int64 = 296
	player := game.PlayerNew(playerID, "", "", "", "", "", "", "")
	token, err := CreateToken(player)

	if err != nil || token == "" {
		t.Fatalf("Error creating token from player")
	}

	claims, err := ExtractToken(token)
	if err != nil {
		t.Fatalf("Error extracting claims from token")
	}

	val, ok := claims["player_id"]
	if !ok || int64(val.(float64)) != playerID {
		t.Fatalf("Could not get player_id from claims")
	}

	nano, ok := claims["nbf"]
	if !ok || 0 >= int64(nano.(float64)) || int64(nano.(float64)) > time.Now().Unix() {
		t.Fatalf("Unexpected nbf value (%d) in claims, %d", int64(nano.(float64)), time.Now().Unix())
	}

	playerIDFromClaims, err := GetPlayerIDFromAccessToken(token)
	if err != nil || playerIDFromClaims != playerID {
		t.Fatalf("Unexpected err getting playerID from claims: %s, %d", err, playerIDFromClaims)
	}
}

func TestInvalidClaims(t *testing.T) {
	_, err := GetPlayerIDFromAccessToken("invalid")
	if err == nil {
		t.Fatalf("Expected error retrieving playerID from invalid token")
	}

	claims := jwt.MapClaims{
		"nbf": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(signingKey)

	extracted, _ := ExtractToken(tokenString)

	_, ok := extracted["player_id"]
	if ok {
		t.Fatalf("Did not expect player_id from claims")
	}

	_, err = GetPlayerIDFromAccessToken(tokenString)
	if err == nil {
		t.Fatalf("Expected error retrieving playerID from token without playerID")
	}
}
