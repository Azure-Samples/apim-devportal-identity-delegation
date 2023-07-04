# Azure API Management (APIM) Developer portal identity delegation with Auth0
<p align="center">
  <img src="./src/images/APIMIdentityDelegationSample.drawio.png">
</p>
This project is created as an example for using identity delegation with Auth0 and Azure API Management (APIM) Developer portal.

## 📝 Demo
### SignUp
![SignUp](./src/images/sign-up.gif)
### SignIn
![SignIn](./src/images/sign-in.gif)

## 🛠️ Prerequisites
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)
- [Docker](https://docs.docker.com/get-docker/)
- [Go](https://golang.org/doc/install)
- [Bicep VSCode extension](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.vscode-bicep)

## 💻 Setup instructions
### 🧩 Clone the repo 🧩
    
    git clone https://github.com/zoeyzuo-se/azure-apim-identity-delegation-sample.git
    

### 🔒 Configure Auth0 🔒

![ConfigureAuth0](./src/images/configure-auth0.gif)

Detail steps:
1. Create a new Auth0 account at [Auth0](https://auth0.com/).
    > If prompted, select `Database` for authentication.
2. Select Applications in the left menu and click the Create Application button.
3. Name your new app and select `Regular Web App Applications`.
4. Click the Create button.
5. Select the Settings tab.
6. Add the following URL to the `Allowed Callback URLs` list. Separated by a comma:
    ```
    http://localhost:3000/callback,
    https://\<your-web-app-name\>.azurewebsites.net/callback,
    https://\<your-apim-name\>.developer.azure-api.net/callback
    ```
7. Add the following URL to the `Allowed Logout URLs` list:
    ```
    http://localhost:3000,
    https://\<your-web-app-name\>.azurewebsites.net
    https://\<your-apim-name\>.developer.azure-api.net
    ```
8. Save the changes.
9. Copy and paste the Domain, Client ID, and Client Secret values into the `.env` file in the root of this project. They should match to the following keys respectively: `AUTH0_DOMAIN`, `AUTH0_CLIENT_ID`, `AUTH0_CLIENT_SECRET`.

### 🌳 Environment Variables 🌳
To run the environment successfully, rename the `.env.example` file to `.env` and provide the following values:

Here's a brief description of each environment variable:

- `PERSONAL_SUB_ID`: Your personal subscription ID.
- `SUFFIX`: The suffix for your environment.
- `PUB_EMAIL`: Your public email address.
- `PUB_NAME`: Your public name.
- `DELEGATION_KEY`: Your delegation key.
- `ACR_NAME`: The name of your Azure Container Registry (ACR). This should be globally unique.
- `ACR_REPO_NAME`: The name of your ACR repository.
- `IMAGE_TAG`: The tag for the Docker image.

Additionally, you will need to obtain the following values from Auth0:

- `AUTH0_CLIENT_ID`: Your Auth0 client ID.
- `AUTH0_DOMAIN`: Your Auth0 domain.
- `AUTH0_CLIENT_SECRET`: Your Auth0 client secret.

Make sure to provide the correct values for these variables to ensure proper authentication and authorization within the environment.

Note: It's important to keep sensitive information, such as personal IDs and secrets, private and secure. Be cautious when sharing your `.env` file or these values with others.

### 🚘 Deploy Azure resources 🚘
- Run the following command to deploy the Azure resources:
    ```bash
    make deploy
    ```
- This will deploy the following resources:
    - Azure Container Registry
    - Azure APIM instance
    - Azure App Service Plan
    - Azure Web App
- This will also set up the following permissions:
    - Give Azure Web App `ACRPull` role access to the Azure Container Registry.
    - Give Azure Web App `Contributor` role access to the Azure APIM instance. This is for creating users in the APIM instance.
- Run `make buildimage` and `make pushimage` to build and push the docker image to the Azure Container Registry.

## 🚀 Usage
### 🌐 Publish developer portal on APIM
Publish your API Management developer portal following the tutorial [here](https://learn.microsoft.com/en-us/azure/api-management/api-management-howto-developer-portal-customize#publish-from-the-azure-portal).
### 🔐 Signup and Login to the developer portal!
SignUp will create a user in the Auth0 Database and also create a user in the APIM instance.

Login will authenticate the user with Auth0 and then delegate the user to the APIM instance.

## 🏃‍♂️ Run the app locally
- Create an .env file under `src/identityApp` and provide the following values:
    ```
    AUTH0_DOMAIN=<your-auth0-domain>
    AUTH0_CLIENT_ID=<your-auth0-client-id>
    AUTH0_CLIENT_SECRET=<your-auth0-client-secret>
    ```
- Once you've set your Auth0 credentials in the `.env` file, run `go mod vendor` to download the Go dependencies.
- Run `go run main.go` to start the app and navigate to [http://localhost:3000/](http://localhost:3000/).
- If everything is working correctly, you should be able to see the following page:
![homepage](./src/images/local-success.png)

## 📝 Notes
### AzureDefaultCredentials
The AzureDefaultCredentials is used to authenticate with Azure resources.
```go
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
```
The NewDefaultAzureCredential will attempt to authenticate with each of these credential types, in the following order, stopping when one provides a token:
- EnvironmentCredential
- WorkloadIdentityCredential

    If environment variable configuration is set by the Azure workload identity webhook. Use WorkloadIdentityCredential directly when not using the webhook or needing more control over its configuration.
- ManagedIdentityCredential
- AzureCLICredential

Details can be found [here](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#DefaultAzureCredential)

In the deployed web app we are using the ManagedIdentityCredential.

### 🐛 Debug tips
If you are seeing unexpected errors, please go to Azure web app service and enable `Application logging` in the `App Service Log` page under `Monitor` and using [golang print statements](https://pkg.go.dev/fmt#Println). You should be able to see the logs in the Log Stream.

Common issues:
- 403 Forbidden: Make sure the web app has `Contributor` role access to the APIM instance.
- 500 Internal Server Error: make sure the web app is up and running. You can check the logs in the Log Stream.

## 🚩 Contact
If you have any questions, please feel free to reach out to me at zoeyzuouk@gmail.com or create an issue in this repo.