package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goware/emailx"
	"github.com/tylergeery/trash_hunt/storage"
	"golang.org/x/crypto/bcrypt"
)

const minPasswordLength = 8

// PlayerStatusActive player status
const PlayerStatusActive = 1
// PlayerStatusRemoved player status
const PlayerStatusRemoved = 2

// Player is a given player in the game
type Player struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	pw         string
	Username       string `json:"username"`
	Token      string `json:"token"`
	Status     int    `json:"status"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// PlayerPublicProfile is a given player in the game
type PlayerPublicProfile struct {
	ID   int64  `json:"id"`
	Username string `json:"username"`
}

var createColumns = []string{
	"email",
	"username",
	"password",
	"token",
	"status",
}
var queryColumns = []string{
	"id",
	"email",
	"username",
	"token",
	"status",
	"created_at",
	"updated_at",
}
var updateColumns = []string{
	"email",
	"username",
	"token",
	"status",
}

// PlayerNew - Constructor
func PlayerNew(id int64, email, pw, username, token string, status int, createdAt, updatedAt string) *Player {
	return &Player{
		ID:         id,
		Email:      email,
		pw:         pw,
		Username:       username,
		Token:      token,
		Status:     status,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

// PlayerRegister - register a new player
func PlayerRegister(email, pw, username string) (*Player, error) {
	// Validate and hash password
	if len(pw) < minPasswordLength {
		return nil, fmt.Errorf("Password must be at least %d characters", minPasswordLength)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	p := PlayerNew(0, email, string(password), username, "", PlayerStatusActive, "", "")

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
	columns := append(queryColumns, "password")
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE email = $1",
		strings.Join(columns, ","),
		storage.TABLE_PLAYER,
	)

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
		Username: p.Username,
	}
}

func (p *Player) toCreateValues() []interface{} {
	return []interface{}{p.Email, p.Username, p.pw, p.Token, PlayerStatusActive}
}

func (p *Player) toUpdateValues() []interface{} {
	return []interface{}{p.Email, p.Username, p.Token, p.Status}
}

func (p *Player) validate() error {
	var err error

	if err = emailx.Validate(p.Email); err != nil {
		return fmt.Errorf("Invalid email format: %s", p.Email)
	}

	if p.Username == "" {
		return fmt.Errorf("Invalid username: %s", p.Username)
	}

	// return error
	return err
}

// GetTestEmail returns a unique test email
func GetTestEmail(ident string) string {
	c := strconv.Itoa(int(time.Now().UnixNano()))

	return fmt.Sprintf("test%s%s@geerydev.com", ident, c)
}

// GetTestUsername provides a unique test username
func GetTestUsername(ident string) string {
	c := strconv.Itoa(int(time.Now().UnixNano()))

	return fmt.Sprintf("testuser-%s-%s", ident, c)
}

// GetTestPlayer returns a unique test user
func GetTestPlayer(ident string) *Player {
	p, _ := PlayerRegister(GetTestEmail(ident), "testpasstestpass", GetTestUsername(ident))

	return p
}

type playerScanner interface {
	Scan(dest ...interface{}) error
}

func scanPlayer(scanner playerScanner, includePass bool) *Player {
	var ID int64
	var status int
	var email, pw, username, token, createdAt, updatedAt string

	if includePass {
		scanner.Scan(&ID, &email, &username, &token, &status, &createdAt, &updatedAt, &pw)
	} else {
		scanner.Scan(&ID, &email, &username, &token, &status, &createdAt, &updatedAt)
	}

	return PlayerNew(ID, email, pw, username, token, status, createdAt, updatedAt)
}
