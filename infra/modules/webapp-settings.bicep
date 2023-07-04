param webApplication string
param auth0ClientId string 
param auth0Domain string 
param auth0CallbackUrl string 
param apimName string 
param rgName string 
param developerPortalUrl string 
param apimResourceUri string
// param azureClientId string
param azureTenantId string = subscription().tenantId
param subscriptionId string = subscription().subscriptionId

// @secure()
// param azureClientSecret string = ''

@description('The Auth0 secret is stored in a key vault and retrieved using a user assigned identity')
@secure()
param auth0Secret string = ''

@description('The APIM delegation key is stored in a key vault and retrieved using a user assigned identity')
@secure()
param delegationKey string = ''

resource appSettings 'Microsoft.Web/sites/config@2022-03-01' = {
  name: '${webApplication}/appsettings'
  properties: {
    AUTH0_CLIENT_ID: auth0ClientId
    AUTH0_DOMAIN: auth0Domain
    AZURE_TENANT_ID: azureTenantId
    AZURE_SUBSCRIPTION_ID: subscriptionId
    AUTH0_CALLBACK_URL: auth0CallbackUrl
    AUTH0_CLIENT_SECRET: auth0Secret
    DELEGATION_KEY: delegationKey
    APIM_NAME: apimName
    APIM_RESOURCE_GROUP: rgName
    APIM_RESOURCE_URI: apimResourceUri
    DEVELOPER_PORTAL_URL: developerPortalUrl
    // AZURE_CLIENT_ID: azureClientId
    // AZURE_CLIENT_SECRET: azureClientSecret
  }
}
