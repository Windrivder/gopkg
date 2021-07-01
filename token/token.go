package token

import (
	"context"
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/windrivder/gopkg/typex"
)

// Context is the context of the JSON web token.
type Context struct {
	ID    string            `json:"id"`
	Value typex.GenericType `json:"value,omitempty"`
}

// secretFunc validates the secret format.
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Make sure the `alg` is what we except.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(http.StatusText(http.StatusBadRequest))
		}

		return []byte(secret), nil
	}
}

// Parse validates the token with the specified secret,
// and returns the context if the token was valid.
func Parse(secret, tokenString string) (*Context, error) {
	ctx := &Context{}

	// Parse the token.
	token, err := jwt.Parse(tokenString, secretFunc(secret))

	// Parse error.
	if err != nil {
		return ctx, err

		// Read the token if it's valid.
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.ID = claims["id"].(string)
		if m, ok := claims["value"]; ok {
			ctx.Value = m.(typex.GenericType)
		}
		return ctx, nil

		// Other errors.
	} else {
		return ctx, err
	}
}

// Sign signs the context with the specified secret.
func Sign(ctx context.Context, secret string, exp time.Duration, c Context) (tokenString string, err error) {
	// The token content.
	// iss: （Issuer）签发者
	// iat: （Issued At）签发时间，用Unix时间戳表示
	// exp: （Expiration Time）过期时间，用Unix时间戳表示
	// aud: （Audience）接收该JWT的一方
	// sub: （Subject）该JWT的主题
	// nbf: （Not Before）不要早于这个时间
	// jti: （JWT ID）用于标识JWT的唯一ID
	now := time.Now()

	claims := jwt.MapClaims{
		"id":  c.ID,
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"exp": now.Add(exp * time.Second).Unix(),
	}

	if c.Value != nil {
		claims["value"] = c.Value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the specified secret.
	tokenString, err = token.SignedString([]byte(secret))

	return
}
