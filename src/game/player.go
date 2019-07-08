package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goware/emailx"
	"github.com/tylergeery/trash_hunt/src/storage"
	"golang.org/x/crypto/bcrypt"
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

// Player is a given player in the game
type PlayerPublicProfile struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

var createColumns = []string{
	"email",
	"name",
	"password",
	"facebook_id",
	"token",
}
var queryColumns = []string{
	"id",
	"email",
	"name",
	"facebook_id",
	"token",
	"created_at",
	"updated_at",
}
var updateColumns = []string{
	"email",
	"name",
	"facebook_id",
	"token",
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

	password, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	p := PlayerNew(0, email, string(password), name, facebookID, "", "", "")

	err = p.save()

	return p, err
}

// PlayerLogin authenticates the login credentials for a player
func PlayerLogin(email, pw string) (*Player, error) {
	p := PlayerGetByEmail(email)

	if p.ID == 0 {
		return nil, fmt.Errorf("Player not found: %s", email)
	}

	err := p.ValidatePassword(pw)
	if err != nil {
		return nil, err
	}

	return p, nil
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

	// Validate email is unique
	existing := PlayerGetByEmail(p.Email)
	if existing.ID != 0 && existing.ID != p.ID {
		return fmt.Errorf("Email %s belongs to an existing user", p.Email)
	}

	if p.ID == 0 {
		id, err = storage.Insert(storage.TABLE_PLAYER, p.toCreateValues(), createColumns)
		p.ID = id
	} else {
		err = storage.Update(storage.TABLE_PLAYER, p.toUpdateValues(), updateColumns, p.ID)
	}

	return err
}

// PlayerGetByID retrieves player based on ID
func PlayerGetByID(id int64) *Player {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE id = $1",
		strings.Join(queryColumns, ","),
		storage.TABLE_PLAYER,
	)

	return scanPlayer(storage.FetchRow(query, id), false)
}

// PlayerGetByToken retrieves player based on auth token
func PlayerGetByToken(authToken string) *Player {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE token = $1",
		strings.Join(queryColumns, ","),
		storage.TABLE_PLAYER,
	)

	return scanPlayer(storage.FetchRow(query, authToken), false)
}

// PlayerGetByEmail retrieves player based on email
func PlayerGetByEmail(userEmail string) *Player {
	query := fmt.Sprintf("SELECT id, email, password, name, facebook_id, token, created_at, updated_at FROM %s WHERE email=$1", storage.TABLE_PLAYER)

	return scanPlayer(storage.FetchRow(query, userEmail), true)
}

// ValidatePassword ensures that the provided password is correct for user
func (p *Player) ValidatePassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(p.pw), []byte(pw))
}

// ToPublicProfile returns a public profile version of a User
func (p *Player) ToPublicProfile() *PlayerPublicProfile {
	return &PlayerPublicProfile{
		ID:   p.ID,
		Name: p.Name,
	}
}

func (p *Player) toCreateValues() []interface{} {
	return []interface{}{p.Email, p.Name, p.pw, p.FacebookID, p.Token}
}

func (p *Player) toUpdateValues() []interface{} {
	return []interface{}{p.Email, p.Name, p.FacebookID, p.Token}
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

type playerScanner interface {
	Scan(dest ...interface{}) error
}

func scanPlayer(scanner playerScanner, includePass bool) *Player {
	var ID int64
	var email, pw, name, facebookID, token, createdAt, updatedAt string

	if includePass {
		scanner.Scan(&ID, &email, &pw, &name, &facebookID, &token, &createdAt, &updatedAt)
	} else {
		scanner.Scan(&ID, &email, &name, &facebookID, &token, &createdAt, &updatedAt)
	}

	return PlayerNew(ID, email, pw, name, facebookID, token, createdAt, updatedAt)
}
