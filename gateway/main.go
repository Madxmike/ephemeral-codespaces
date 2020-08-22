package main

import "github.com/tierzer0/gateway/auth"

func main() {
	authenticator := auth.Authenticator{
		ProjectID: "test",
	}
	authenticator.RetrievePublicKeys()

}
