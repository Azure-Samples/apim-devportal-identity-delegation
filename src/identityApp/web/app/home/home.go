package home

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gin-gonic/gin"
)

// Handler for our home page.
func Handler(ctx *gin.Context) {
	getTokenViaGoSDK(ctx)
	ctx.HTML(http.StatusOK, "home.html", nil)
}
func getTokenViaGoSDK(ctx *gin.Context) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	// Default Azure credential will use the following credential chain to authenticate:
	// 1. Environment Credential
	// 2. Workload Identity Credential
	// 3. Managed Identity Credential
	// 4. Azure CLI Credential

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Acquire an access token using the managed identity
	// opts := policy.TokenRequestOptions{Scopes: []string{scope}}
	// token, err := cred.GetToken()
	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("Access Token from managed identity:", token)
}
