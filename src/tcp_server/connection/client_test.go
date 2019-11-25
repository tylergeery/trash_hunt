package connection

import (
	"fmt"
	"testing"

	"github.com/tylergeery/trash_hunt/auth"
	"github.com/tylergeery/trash_hunt/model"
)

// ResponseOnlyMockConnection object used for mocking Connection interface while testing arena
type ResponseOnlyMockConnection struct {
	message GameMessage
}

func (c *ResponseOnlyMockConnection) gatherInput(input []byte) (s string, err error) {
	return "", fmt.Errorf("Response only")
}

func (c *ResponseOnlyMockConnection) respond(message GameMessage) error {
	c.message = message

	return nil
}

func TestSetupUserInvalidPreferences(t *testing.T) {
	type TestCase struct {
		client *Client
		errMsg string
	}
	testCases := []TestCase{
		TestCase{
			client: NewClient(&MockConnection{
				input: "",
			}),
			errMsg: "unexpected end of JSON input",
		},
		TestCase{
			client: NewClient(&MockConnection{
				input: "{\"user_token\": \"\", \"difficulty\": \"easy\"}",
			}),
			errMsg: "token contains an invalid number of segments",
		},
	}

	for _, tc := range testCases {
		err := tc.client.SetUpUser()
		if err.Error() != tc.errMsg {
			t.Fatalf("Expected (%s) to match (%s)", err.Error(), tc.errMsg)
		}
	}
}

func TestSetupUser(t *testing.T) {
	// Given
	player := model.GetTestPlayer("client-setup-user")
	token, _ := auth.CreateToken(player)
	input := fmt.Sprintf("{\"user_token\": \"%s\", \"difficulty\": \"easy\"}", token)
	client := NewClient(&MockConnection{input: input})

	// When/Then
	err := client.SetUpUser()
	if err != nil {
		t.Fatalf("Unexpected error setting up user: %s", err.Error())
	}
}

func TestWaitForStart(t *testing.T) {
	// Given
	connection := &ResponseOnlyMockConnection{}
	client := NewClient(connection)

	// When
	client.notifications <- eventStartGame
	client.WaitForStart()

	// Then
	if connection.message.Event != messageStatusPending {
		t.Fatalf("Expected message status pending, got %d", connection.message.Event)
	}
}
