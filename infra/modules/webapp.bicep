param suffix string
param location string
param acrName string
param acrRepoName string
param imageTag string
param sku string = 'P1v3'

resource appServicePlan 'Microsoft.Web/serverfarms@2022-03-01' = {
  name: 'asp-${suffix}'
  location: location
  kind: 'app,linux,container'

  sku: {
    name: sku
    capacity: 1
  }

  properties: {
    reserved: true
  }
}

resource webApplication 'Microsoft.Web/sites@2022-03-01' = {
  name: 'webapp-${suffix}'
  location: location
  tags: {
    'hidden-related:${resourceGroup().id}/providers/Microsoft.Web/serverfarms/appServicePlan': 'Resource'
  }
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    serverFarmId: appServicePlan.id
    siteConfig:  {
      linuxFxVersion: 'DOCKER|${acrName}.azurecr.io/${acrRepoName}:${imageTag}'
      alwaysOn: true
    } 
    
  }
}
output hostname string = webApplication.properties.defaultHostName
output name string = webApplication.name
output appServiceId string = webApplication.id
output webappSystemAssignedIdentityId string = webApplication.identity.principalId
