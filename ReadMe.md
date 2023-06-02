## Configure Auth0
1. Create a new Auth0 account at [Auth0](https://auth0.com/).
2. Select Applications in the left menu and click the Create Application button.
3. Name your new app and select Regular Web App Applications.
4. Click the Create button.
5. Select the Settings tab.
6. Add the following URL to the Allowed Callback URLs list:
    Localhost:
    - http://localhost:3000/callback

    Web app URL:
    - https://<your-web-app-name>.azurewebsites.net/callback
    - Or fetch it from the Overview tab of your web app in the Azure portal. Copy the Default Domain value and append /callback to it.

    APIM Developer Portal URL:
    - https://<your-apim-name>.developer.azure-api.net/callback
    - Or fetch it from the Overview tab of your APIM instance in the Azure portal. Copy the Developer Portal URL value and append /callback to it.
7. Add the following URL to the Allowed Logout URLs list:
    Localhost:
    - http://localhost:3000

    Web app URL:
    - https://<your-web-app-name>.azurewebsites.net

    APIM Developer Portal URL:
    - https://<your-apim-name>.developer.azure-api.net
8. Scroll down and click the Save Changes button.
9. Copy and paste the Domain, Client ID, and Client Secret values into the `.env` file in the root of this project. They should match to the following keys respectively: `AUTH0_DOMAIN`, `AUTH0_CLIENT_ID`, `AUTH0_CLIENT_SECRET`.

## Configure Role Based Access Control (RBAC)

### Create Azure App registration 
We need to create an Azure App registration to allow our web app to access the APIM instance. This is done by following these steps:
1. Go to the Azure portal and select App registrations.
2. Click the New registration button.
3. Name your app and select Accounts in this organizational directory only (Default Directory only - Single tenant) for the Supported account types.
4. Click the Register button.
5. In Authentication tab, click Add a platform and select Web. 
6. Put "http://localhost:8000" as redirect URI. (Not used but required)
7. Click the Configure button.
8. Check the Access tokens and ID tokens boxes under Implicit grant.
9. Click the Save button.

### Keep note of the following values:
1. Application (client) ID. This matches to the `APPR_CLIENT_ID` key in the `.env` file.
2. Go to the Certificates & secrets tab and click the New client secret button. Create a new secret and keep note of the value. This matches to the `APPR_CLIENT_SECRET` key in the `.env` file.
3. Go to Azure portal home page and search for Azure AD. Select Azure Active Directory.
4. Click on the Enterprise applications tab and search for your app registration. Select it.
5. In the Enterprise Application's overview page, you will find the "Object ID" field. This is the object ID of the associated service principal. This matches to the `APPR_SP_OID` key in the `.env` file.