param location string
@description('Name of the APIM service')
param name string

@description('The name of the owner of the service')
@minLength(1)
param publisherName string

@description('The email address of the owner of the service')
@minLength(1)
param publisherEmail string

param identityWebAppUrl string
param delegationKey string

resource apiManagementInstance 'Microsoft.ApiManagement/service@2022-08-01' = {
  name: name
  location: location
  sku:{
    capacity: 1
    name: 'Developer'
  }
  properties:{
    virtualNetworkType: 'None'
    publisherEmail: publisherEmail
    publisherName: publisherName
  }
}

resource apiManagementIdentityDelegation 'Microsoft.ApiManagement/service/portalsettings@2022-08-01' = {
  parent: apiManagementInstance
  name: 'delegation'
  properties: {
    subscriptions: {
      enabled: false
    }
    url: '${identityWebAppUrl}/delegation'
    userRegistration: {
      enabled: true
    }
    validationKey: delegationKey
  }
}

output name string = apiManagementInstance.name
output id string = apiManagementInstance.id
