package auth

import "github.com/dgrijalva/jwt-go"

type MultiAgentsAuthenticate interface {
	RegisterLogin(agentID, token string) (err error)
	DecrLoginCount(agentID string) error
	RemTokenFromList(agentID, token string) error
	RemAllTokens(agentID string) error
	GetLoginCount(agentID string) (int, error)
}

type Claims struct {
	EntID  string `json:"ent_id"`
	UserID string `json:"user_id"`
	*jwt.StandardClaims
}
