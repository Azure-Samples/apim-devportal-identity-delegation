package azureClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func fetchAzureEnv() (string, string, string, error) {
	subscriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	rg := os.Getenv("APIM_RESOURCE_GROUP")
	apimName := os.Getenv("APIM_NAME")

	return subscriptionId, rg, apimName, nil
}

func getTimeAfterHours(hours int) string {
	// Get the current time
	currentTime := time.Now()

	// Add the specified number of hours to the current time
	futureTime := currentTime.Add(time.Duration(hours) * time.Hour)

	// Format the resulting time in the desired format
	formattedTime := futureTime.UTC().Format("2006-01-02T15:04:05.9999999Z")

	return formattedTime
}

// Default Azure credential will use the following credential chain to authenticate:
// 1. Environment Credential
// 2. Workload Identity Credential
// 3. Managed Identity Credential
// 4. Azure CLI Credential
// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#DefaultAzureCredential
func GetTokenViaGoSDK(ctx *gin.Context) (string, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return "", err
	}

	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return "", err
	}

	return token.Token, nil
}

func GetSharedAccessTokenFromAPIM(accessToken, uid string) (string, error) {
	// Construct the Azure API Management API request URL
	subscriptionId, rg, apimName, err := fetchAzureEnv()
	userId := uid
	expriyTime := getTimeAfterHours(12)
	requestBody := fmt.Sprintf(`{"properties": {"keyType": "primary", "expiry": "%s"}}`, expriyTime)
	apiEndpoint := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/users/%s/token?api-version=2022-08-01", subscriptionId, rg, apimName, userId)
	req, err := http.NewRequest("POST", apiEndpoint, strings.NewReader(requestBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Execute the Azure API Management API request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Parse the API response and print the result
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %v", err)
	}
	sharedAccessToken := gjson.GetBytes(body, "value").String()
	return sharedAccessToken, nil
}

type createUserRequest struct {
	Properties createUserProperties `json:"properties"`
}

type createUserProperties struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	State     string `json:"state"`
	UserName  string `json:"userName"`
}

func CreateUserInAPIM(username, uid, accessToken string) error {
	// construct request payload
	payload := createUserRequest{
		Properties: createUserProperties{
			Email:     username,
			FirstName: strings.Split(username, "@")[0],
			LastName:  strings.Split(username, "@")[0],
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	subscriptionId, rg, apimName, err := fetchAzureEnv()
	userId := uid

	// make HTTP POST request to create user in APIM
	client := http.Client{}
	apiEndpoint := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/users/%s?api-version=2022-08-01", subscriptionId, rg, apimName, userId)
	req, err := http.NewRequest("PUT", apiEndpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	// check response status code for success
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Println("Failed to create or update user in APIM")
		return fmt.Errorf("Failed to create or update user in APIM: %s", resp.Status)
	}

	return nil
}

func UserExistInApim(uid, accessToken string) (bool, error) {
	// Construct the request URL with the given uid
	subscriptionId, rg, apimName, err := fetchAzureEnv()
	userId := uid
	url := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/users/%s?api-version=2021-04-01", subscriptionId, rg, apimName, userId)

	// Create a new HTTP request with the constructed URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	// Add authentication and other required headers to the request
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("Content-Type", "application/json")

	// Send the request to Azure API Management
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// If the response status code is not 200 OK, return false
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	// Parse the response JSON to check if the user exists
	var user map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return false, err
	}

	if user["name"] != nil && user["name"] == uid {
		// The user with the given uid exists in Azure API Management
		return true, nil
	}

	// The user with the given uid does not exist in Azure API Management
	return false, nil
}
