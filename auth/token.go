package auth

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/tylergeery/trash_hunt/game"
)

var signingKey []byte

func init() {
	signingKey = []byte(os.Getenv("TOKEN_SIGNING_VALUE"))
}

// CreateToken creates auth token for player
func CreateToken(player *game.Player) (string, error) {
	claims := jwt.MapClaims{
		"player_id": player.ID,
		"nbf":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return tokenString, err
}

// ExtractToken gets claims from auth token
func ExtractToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("Could not extract claims from token")
	}

	return claims, nil
}

// GetPlayerIDFromAccessToken gets the player ID from access token string
func GetPlayerIDFromAccessToken(token string) (int64, error) {
	claims, err := ExtractToken(token)

	if err != nil {
		return 0, err
	}

	if val, ok := claims["player_id"]; ok {
		return int64(val.(float64)), nil
	}

	return 0, fmt.Errorf("Player ID not found in token")
}
