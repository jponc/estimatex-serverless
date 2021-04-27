package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jponc/estimatex-serverless/internal/types"
)

type ParticipantClaims struct {
	RoomID  string `json:"room_id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
	jwt.StandardClaims
}

// Client interface
type Client interface {
	// CreateAccessToken creates a JWT access token
	CreateAccessToken(participant types.Participant) (string, error)
	// GetClaims
	GetClaims(token string) (*ParticipantClaims, error)
}

type client struct {
	jwtSecret string
}

// NewClient instantiates an Auth Client
func NewClient(jwtSecret string) (Client, error) {
	c := &client{
		jwtSecret: jwtSecret,
	}

	return c, nil
}

func (c *client) CreateAccessToken(participant types.Participant) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := ParticipantClaims{
		RoomID:  participant.RoomID,
		Name:    participant.Name,
		IsAdmin: participant.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.jwtSecret))
}

func (c *client) GetClaims(tokenString string) (*ParticipantClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ParticipantClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.jwtSecret), nil
	})

	if claims, ok := token.Claims.(*ParticipantClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
