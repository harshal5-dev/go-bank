package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	key paseto.V4SymmetricKey
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	key, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, fmt.Errorf("invalid paseto key: %w", err)
	}

	return &PasetoMaker{key: key}, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token := paseto.NewToken()
	token.SetString("id", payload.Id.String())
	token.SetString("username", payload.Username)
	token.SetTime("iat", payload.IssuedAt)
	token.SetTime("exp", payload.ExpiresAt)

	return token.V4Encrypt(maker.key, nil), nil
}

func (maker *PasetoMaker) VerifyToken(tokenString string) (*Payload, error) {
	parser := paseto.NewParserWithoutExpiryCheck()
	pasetoToken, err := parser.ParseV4Local(maker.key, tokenString, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	idStr, err := pasetoToken.GetString("id")
	if err != nil {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	username, err := pasetoToken.GetString("username")
	if err != nil {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	iat, err := pasetoToken.GetTime("iat")
	if err != nil {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	exp, err := pasetoToken.GetTime("exp")
	if err != nil {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid token id: %w", err)
	}

	payload := &Payload{
		Id:        id,
		Username:  username,
		IssuedAt:  iat,
		ExpiresAt: exp,
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
