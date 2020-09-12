package auth

import "gopkg.in/square/go-jose.v2/jwt"

// User represents a authed user
type User string

// NewUser creates a new user from the claim's subject
func NewUser(claims *jwt.Claims) User {
	return User(claims.Subject)
}
