package game

import (
	"errors"
	"fmt"
	"testing"

	_ "github.com/tylergeery/trash_hunt/src/test"
)

func TestPlayerRegisterFailures(t *testing.T) {
	type TestCase struct {
		args [4]string
		err  error
	}
	testCases := []TestCase{
		TestCase{
			args: [4]string{"", "asdffdssadf", "", ""},
			err:  errors.New("Invalid email format: "),
		},
		TestCase{
			args: [4]string{"test", "asdffdssadf", "", ""},
			err:  errors.New("Invalid email format: test"),
		},
		TestCase{
			args: [4]string{"test@yahoo.com", "1234", "", ""},
			err:  errors.New("Password must be at least 8 characters"),
		},
		TestCase{
			args: [4]string{"tyger@geerydev.com", "test1234", "", ""},
			err:  errors.New("Invalid name: "),
		},
	}

	for _, test := range testCases {
		_, err := PlayerRegister(test.args[0], test.args[1], test.args[2], test.args[3])

		if err == nil || fmt.Sprintf("%s", err) != fmt.Sprintf("%s", test.err) {
			t.Fatalf("Expected err: %s, received: %s", test.err, err)
		}
	}
}

func TestPlayerRegisterSuccess(t *testing.T) {
	type TestCase struct {
		args   [4]string
		player *Player
	}
	testEmail := GetTestEmail("success")
	testCases := []TestCase{
		TestCase{
			args:   [4]string{testEmail, "asdffdssadf", "jk", ""},
			player: PlayerNew(0, testEmail, "asdffdssadf", "jk", "", "", PlayerStatusActive, "", ""),
		},
	}

	for _, test := range testCases {
		p, err := PlayerRegister(test.args[0], test.args[1], test.args[2], test.args[3])

		if p.ID <= 0 {
			t.Fatalf("Received invalid player ID: %d, err: %s", p.ID, err)
		}

		if p.Email != test.player.Email {
			t.Fatalf("Expected player email: %s, received: %s", p.Email, test.player.Email)
		}
		if p.pw == test.player.pw {
			t.Fatalf("Expected hashed player pw for: %s, received: %s", test.player.pw, p.pw)
		}
		if p.Name != test.player.Name {
			t.Fatalf("Expected player name: %s, received: %s", p.Name, test.player.Name)
		}
		if p.FacebookID != test.player.FacebookID {
			t.Fatalf("Expected player FacebookID: %s, received: %s", p.FacebookID, test.player.FacebookID)
		}
		if p.Status != test.player.Status {
			t.Fatalf("Expected player status: %d, received: %d", p.Status, test.player.Status)
		}

		playerByEmail := PlayerGetByEmail(p.Email)
		if playerByEmail.ID != p.ID {
			t.Fatalf("PlayerByEmail does not have the correct ID: %d", playerByEmail.ID)
		}
	}
}

func TestPlayerLogin(t *testing.T) {
	email := GetTestEmail("login")
	password := "saklfsdlkfsa"
	p, _ := PlayerRegister(email, password, "asdflksas TLkdlsff", "")

	p1, err := PlayerLogin(email, password)
	if err != nil {
		t.Fatalf("Unexpected login err; %s", err)
	}

	if p1.ID != p.ID {
		t.Fatalf("Unexpected user returned from login")
	}

	_, err = PlayerLogin(email, password+"ah")
	if err == nil {
		t.Fatalf("Did not received expected login error for user with bad pass")
	}
}

func TestPlayerUpdateError(t *testing.T) {
	p := PlayerNew(0, "test@test.com", "", "", "", "", PlayerStatusActive, "", "")

	p.Name = "Tester"
	err := p.Update()

	if err == nil || fmt.Sprintf("%s", err) != "Could not update non-existent player" {
		t.Fatalf("Unexpected update err; %s", err)
	}
}

func TestPlayerUpdate(t *testing.T) {
	p, _ := PlayerRegister(GetTestEmail("update"), "saklfsdlkfsa", "asdflksas TLkdlsff", "")

	p.Name = "Tester"
	err := p.Update()

	if err != nil {
		t.Fatalf("Unexpected register err; %s", err)
	}
}
