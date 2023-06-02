package delegation

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"auth-proxy/platform/authenticator"
	"auth-proxy/web/app/logout"
)

// Handler for our delegation.
func Handler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		operation := ctx.Query("operation")
		returnUrl := ctx.Query("returnUrl")
		salt := ctx.Query("salt")
		sig := ctx.Query("sig")
		delegation_key := os.Getenv("DELEGATION_KEY")

		isRequestFromAPIM := verifyRequestFromAPIM(salt, returnUrl, sig, delegation_key)

		if operation == "SignOut" {
			logout.Handler(ctx)
			return
		}

		if !isRequestFromAPIM {
			// ctx.String(http.StatusInternalServerError, "Request is not from APIM.")
			fmt.Println("Request is not from APIM. But carry on.")
		}

		state, err := generateRandomState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Save the state inside the session.
		session := sessions.Default(ctx)
		session.Set("state", state)
		session.Set("operation", operation)
		session.Set("returnUrl", returnUrl)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
		// if operation == "SignUp" || operation == "SignIn" {
		// 	ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
		// } else {
		// 	ctx.String(http.StatusBadRequest, "Invalid operation.")
		// }

	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func verifyRequestFromAPIM(salt string, returnUrl string, sig string, delegation_key string) bool {
	rurl, _ := url.PathUnescape(returnUrl)
	key, _ := base64.StdEncoding.DecodeString(delegation_key)
	res := hmac.New(sha512.New, key)
	salt, _ = url.PathUnescape(salt)
	sig, _ = url.PathUnescape(sig)
	salt_url := salt + "\n" + rurl
	res.Write([]byte(salt_url))
	sum := res.Sum(nil)
	computed_sig := base64.StdEncoding.EncodeToString(sum)

	// Verify the signature
	if computed_sig == sig {
		return true
	} else {
		return false
	}
}
