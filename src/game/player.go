package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goware/emailx"
	"github.com/tylergeery/trash_hunt/src/storage"
)

const minPasswordLength = 8

// Player is a given player in the game
type Player struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	pw         string
	Name       string `json:"name"`
	FacebookID string `json:"facebook_id"`
	Pos        Pos    `json:"pos"`
	Token      string `json:"token"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

var types = map[string]string{
	"id":          "int",
	"email":       "string",
	"name":        "string",
	"facebook_id": "string",
	"token":       "string",
	"created_at":  "string",
	"updated_at":  "string",
}

// PlayerNew - Constructor
func PlayerNew(id int64, email, pw, name, facebookID, token, createdAt, updatedAt string) *Player {
	return &Player{
		ID:         id,
		Email:      email,
		pw:         pw,
		Name:       name,
		FacebookID: facebookID,
		Token:      token,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

// PlayerRegister - register a new player
func PlayerRegister(email, pw, name, facebookID string) (*Player, error) {
	// Validate and hash password, or validate facebookID
	if len(pw) < minPasswordLength {
		return nil, fmt.Errorf("Password must be at least %d characters", minPasswordLength)
	}

	// Validate email is unique
	existing := PlayerGetByEmail(email)
	if existing.ID != 0 {
		return nil, fmt.Errorf("Email %s belongs to an existing user", email)
	}

	p := PlayerNew(0, email, pw, name, facebookID, "", "", "")
	err := p.save()

	return p, err
}

// Update - update an existing player
func (p *Player) Update() error {
	if p.ID == 0 {
		return errors.New("Could not update non-existent player")
	}

	return p.save()
}

func (p *Player) save() error {
	var err error
	var id int64

	err = p.validate()

	if err != nil {
		return err
	}

	if p.ID == 0 {
		id, err = storage.Insert(storage.TABLE_PLAYER, p.toCreateMap(), types)
		p.ID = id
	} else {
		err = storage.Update(storage.TABLE_PLAYER, p.toUpdateMap(), types, p.ID)
	}

	return err
}

// PlayerGetByToken retrieves player based on auth token
func PlayerGetByToken(authToken string) *Player {
	var ID int64
	var email, pw, name, facebookID, token, createdAt, updatedAt string

	query := fmt.Sprintf("SELECT %s FROM %s WHERE token = $1", getColumns(types), storage.TABLE_PLAYER)

	storage.FetchRow(query, authToken).Scan(&ID, &email, &pw, &name, &facebookID, &token, &createdAt, &updatedAt)

	return PlayerNew(ID, email, pw, name, facebookID, token, createdAt, updatedAt)
}

// PlayerGetByEmail retrieves player based on email
func PlayerGetByEmail(userEmail string) *Player {
	var ID int64
	var email, pw, name, facebookID, token, createdAt, updatedAt string

	query := fmt.Sprintf("SELECT %s FROM %s WHERE email = '$1'", getColumns(types), storage.TABLE_PLAYER)

	storage.FetchRow(query, userEmail).Scan(&ID, &email, &pw, &name, &facebookID, &token, &createdAt, &updatedAt)

	return PlayerNew(ID, email, pw, name, facebookID, token, createdAt, updatedAt)
}

func (p *Player) toSaveMap() map[string]interface{} {
	return map[string]interface{}{
		"name":        p.Name,
		"email":       p.Email,
		"facebook_id": p.FacebookID,
	}
}

// ValidatePassword ensures that the provided password is correct for user
func (p *Player) ValidatePassword(pw string) bool {
	// TODO: handle pw encryption
	return pw == p.pw
}

func (p *Player) toCreateMap() map[string]interface{} {
	m := p.toSaveMap()

	// TODO: handle password encryption
	m["password"] = p.pw

	return m
}

func (p *Player) toUpdateMap() map[string]interface{} {
	m := p.toSaveMap()

	return m
}

func (p *Player) validate() error {
	var err error

	if err = emailx.Validate(p.Email); err != nil {
		return fmt.Errorf("Invalid email format: %s", p.Email)
	}

	if p.Name == "" {
		return fmt.Errorf("Invalid name: %s", p.Name)
	}

	// return error
	return err
}

func getColumns(m map[string]string) string {
	cols := make([]string, 0, len(m))

	for k := range m {
		cols = append(cols, k)
	}

	return strings.Join(cols, ",")
}

// GetTestEmail returns a unique test email
func GetTestEmail(ident string) string {
	c := strconv.Itoa(int(time.Now().UnixNano()))

	return fmt.Sprintf("test%s%s@geerydev.com", ident, c)
}

// GetTestPlayer returns a unique test user
func GetTestPlayer(ident string) *Player {
	p, _ := PlayerRegister(GetTestEmail(ident), "saklfsdlkfsa", "asdflksas TLkdlsff", "")

	return p
}
