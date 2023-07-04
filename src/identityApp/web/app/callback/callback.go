package callback

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"auth-proxy/platform/authenticator"
	"auth-proxy/platform/azureClient"
)

// Handler for our callback.
func Handler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}
		operation := session.Get("operation").(string)
		tokenFromManagedIdentity, err := azureClient.GetTokenViaGoSDK(ctx)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		handleSignUpAndSignIn(ctx, auth, tokenFromManagedIdentity, operation)
	}
}

func handleSignUpAndSignIn(ctx *gin.Context, auth *authenticator.Authenticator, accessToken, operation string) {
	token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
	if err != nil {
		ctx.String(http.StatusUnauthorized, "Failed to convert an authorization code into a token.")
		return
	}

	idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	session := sessions.Default(ctx)
	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)
	fmt.Println("profile: ", profile)
	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	parts := strings.Split(profile["sub"].(string), "|")
	uid := parts[1]
	username := profile["name"].(string)

	userExistInApim, err := azureClient.UserExistInApim(uid, accessToken)
	if userExistInApim == false {
		err = azureClient.CreateUserInAPIM(username, uid, accessToken)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Error creating user: %s", username, "error message:", err.Error())
			return
		}
	}

	sharedAccessToken, err := azureClient.GetSharedAccessTokenFromAPIM(accessToken, uid)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error fetching shared access token for uid %s", uid)
		return

	} else if sharedAccessToken == "" {
		ctx.String(http.StatusInternalServerError, "SharedAccessToken for uid %s is empty", uid)
		return
	}
	// Redirect to logged in page.
	encodedToken := url.QueryEscape(sharedAccessToken)
	returnURL := "/"
	encodedReturnURL := url.QueryEscape(returnURL)
	developerPortalUrl := os.Getenv("DEVELOPER_PORTAL_URL")
	fmt.Println("developerPortalUrl: ", developerPortalUrl)
	url := fmt.Sprintf("%ssignin-sso?token=%s&returnUrl=%s", developerPortalUrl, encodedToken, encodedReturnURL)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
