package azureClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

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

func GetBearerTokenFromAzAD() (string, error) {
	// Azure AD Application ID and Secret
	clientId := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	tenantId := os.Getenv("AZURE_TENANT_ID")
	// Construct the OAuth2.0 token request URL
	tokenUrl := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantId)
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Set("scope", "https://management.azure.com/.default")

	// Execute the token request
	req, err := http.NewRequest("POST", tokenUrl, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the HTTP request to obtain an Azure AD access token
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}
	defer resp.Body.Close()

	// Read the HTTP response body and extract the Azure AD access token
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %v", err)
	}

	accessToken := gjson.GetBytes(body, "access_token").String()
	fmt.Println("Azure AD Access Token: " + accessToken)

	return accessToken, nil
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
