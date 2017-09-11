// Package jwtutil holds JWT related funcs
package jwtutil

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/LeKovr/go-base/logger"
)

const (
	// Bearer is a prefix of JWT header value
	Bearer = "Bearer"
)

// -----------------------------------------------------------------------------

// Flags is a package flags sample
// in form ready for use with github.com/jessevdk/go-flags
type Flags struct {
	Key      string   `long:"jwt_key" description:"Key to sign JWT results"`
	Age      int      `long:"jwt_age" default:"96" description:"JWT age to expire token (hours)"`
	Producer string   `long:"jwt_producer" default:"dbrpc" description:"JWT producer name"`
	Issuers  []string `long:"jwt_issuer" default:"*" description:"required JWT issuer name(s)"`
}

// -----------------------------------------------------------------------------

// App - Класс сервера API
type App struct {
	Log    *logger.Log
	Config *Flags
}

// -----------------------------------------------------------------------------
// Functional options

// Config sets config struct
func Config(c *Flags) func(a *App) error {
	return func(a *App) error {
		return a.setConfig(c)
	}
}

// -----------------------------------------------------------------------------
// Internal setters

func (a *App) setConfig(c *Flags) error {
	a.Config = c
	return nil
}

// -----------------------------------------------------------------------------

// New - Конструктор сервера API
func New(log *logger.Log, options ...func(a *App) error) (a *App, err error) {

	a = &App{Log: log.WithField("in", "jwt")}

	for _, option := range options {
		err := option(a)
		if err != nil {
			return nil, err
		}
	}
	return

}

// CustomClaims holds JWT fields
type CustomClaims struct {
	Data *json.RawMessage `json:"data"`
	jwt.StandardClaims
}

// CustomRes holds JWT creation result
type CustomRes struct {
	Token string `json:"token"`
}

// Session holds custom JWT fields
type Session map[string]interface{}

// SessionSlice holds parsed JWT data
type SessionSlice struct {
	Data []Session `json:"data"`
	jwt.StandardClaims
}

// -----------------------------------------------------------------------------

// Create - creates JWT
func (a *App) Create(method string, s *json.RawMessage) (*json.RawMessage, error) {

	a.Log.Debugf("Got JWT data: %s", s)

	// Create the Claims
	claims := CustomClaims{
		s,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(a.Config.Age)).Unix(),
			Issuer:    fmt.Sprintf("%s:%s", a.Config.Producer, method),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(a.Config.Key))

	a.Log.Debugf("JWT: %v", t)
	ts := CustomRes{Token: t}
	tjs, err := json.Marshal(ts)
	traw := json.RawMessage(tjs)
	return &traw, err

}

// -----------------------------------------------------------------------------

// Parse - parses JWT
func (a *App) Parse(s string) (*Session, error) {

	// to the callback, providing flexibility.
	token, err := jwt.ParseWithClaims(s, &SessionSlice{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.Config.Key), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SessionSlice); ok && token.Valid {
		a.Log.Debugf("Got JWT session: %+v (%+v)", claims, token.Header)

		if !stringExists(a.Config.Issuers, claims.Issuer, "*") {
			return nil, fmt.Errorf("Uncorrect JWT Issuer %s", claims.Issuer)
		}

		return &claims.Data[0], nil
	}
	return nil, fmt.Errorf("Invalid JWT")

}

// -----------------------------------------------------------------------------

// Check if str or any exists in strings slice
func stringExists(strings []string, str string, any string) bool {
	if len(strings) > 0 { // lookup if host is allowed
		for _, s := range strings {
			if str == s || (any != "" && s == any) {
				return true
			}
		}
	}
	return false
}
