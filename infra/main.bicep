targetScope = 'subscription'
param location string = deployment().location
param suffix string
param acrName string
param acrRepoName string = 'identity'
param imageTag string = 'latest'
// param apprClientId string
// Params passed for APIM
@minLength(1)
param publisherEmail string

@minLength(1)
param publisherName string

@secure()
param delegationKey string = ''
// @secure()
// param apprClientSecret string = ''
@secure()
param auth0ClientSecret string = ''

// param azureClientId string
param auth0ClientId string = ''
param auth0Domain string

// param apprSpOID string // App registration "apim-sample-zzuo"'s service principal object id
param acrPullRoleId string = '7f951dda-4ed3-4680-a7ca-43fe172d538d' // AcrPull
param apimContributorRoleId string = 'b24988ac-6180-42a0-ab88-20f7382dd24c' // API Management Service Contributor
// Create a resource group
resource resourceGroup 'Microsoft.Resources/resourceGroups@2022-09-01' = {
  name: 'rg-apim-sample-${suffix}'
  location: location
}

output resourceGroupName string = resourceGroup.name

module apim 'modules/apim.bicep' = {
  scope: resourceGroup
  name: 'apim'
  params: {
    identityWebAppUrl: 'https://webapp-${suffix}.azurewebsites.net'
    location: location
    name: 'apim-sample-${suffix}'
    publisherEmail: publisherEmail
    publisherName: publisherName
    delegationKey: delegationKey
  }
}

module webapp 'modules/webapp.bicep' = {
  scope: resourceGroup
  name: 'app'
  params: {
    location: location
    suffix: suffix
    acrName: acr.outputs.name
    acrRepoName: acrRepoName
    imageTag: imageTag
  }
  dependsOn: [
    acr
  ]
}

module webappsettings 'modules/webapp-settings.bicep' = {
  scope: resourceGroup
  name: 'webapp-settings'
  params: {
    webApplication: webapp.outputs.name
    auth0ClientId: auth0ClientId
    auth0Domain: auth0Domain
    auth0CallbackUrl: 'https://webapp-${suffix}.azurewebsites.net/callback'
    auth0Secret: auth0ClientSecret
    delegationKey: delegationKey
    apimName: apim.outputs.name
    rgName: resourceGroup.name
    developerPortalUrl: 'https://${apim.outputs.name}.developer.azure-api.net/'
    apimResourceUri: apim.outputs.id
    // azureClientId: apprClientId
    // azureClientSecret: apprClientSecret
  }
}

module acr 'modules/acr.bicep' = {
  scope: resourceGroup
  name: 'acr'
  params: {
    location: location
    acrName: acrName
  }
}

module acrRoleAssignment './modules/acr-role-assignment.bicep' = {
  name: 'container-registry-role-assignment'
  scope: resourceGroup
  params: {
    roleId: acrPullRoleId
    principalId: webapp.outputs.webappSystemAssignedIdentityId
    registryName: acr.outputs.name
  }
  dependsOn: [
    acr
    webapp
  ]
}


module apimRoleAssignmentSystemAssign './modules/apim-role-assignment.bicep' = {
  name: 'apim-role-assignment'
  scope: resourceGroup
  params: {
    principalId: webapp.outputs.webappSystemAssignedIdentityId
    apimName: apim.outputs.name
    roleId: apimContributorRoleId
  }
  dependsOn: [
    apim
    webapp
  ]
}

// module apimRoleAssignmentAppreg './modules/apim-role-assignment.bicep' = {
//   name: 'apim-role-assignment-2'
//   scope: resourceGroup
//   params: {
//     principalId: apprSpOID
//     apimName: apim.outputs.name
//     roleId: apimContributorRoleId
//   }
//   dependsOn: [
//     apim
//     webapp
//   ]
// }
