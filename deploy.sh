#!/bin/bash
set -euo pipefail

# Check to see if a suffix has been provided
if [[ -z "${SUFFIX}" ]]; then
    echo "[ERROR] SUFFIX is not set.";
    exit 1;
fi

# Get the current subscription ID
SUBID=$(az account show --query "id" -o tsv)
SUBNAME=$(az account show --query "name" -o tsv)

# Show the user what they are about to deploy
echo -e "You are deploying into your ${SUBNAME} sub with \n suffix=${SUFFIX} \n publisherEmail=${PUB_EMAIL} \n publisherName=${PUB_NAME}"

no_prompt=${NO_PROMPT:-false}

# Check the user wants to continue
if [[ "${no_prompt}" == "false" ]]; then
  read -p "❓ Are you sure you want to continue (y/n)? " -n 1 -r
  echo    # move to a new line
  if [[ ! $REPLY =~ ^[Yy]$ ]]
  then
      exit 1
  fi
fi

echo -e "\n ⏳ Deploying to personal subscription ${SUBID}...\n"

# Run the Bicep deployment
az deployment sub create \
            --template-file ./infra/main.bicep \
            -l westeurope \
            --name "deploy-${SUFFIX}" \
            --parameters \
                suffix="${SUFFIX}" \
                publisherEmail="${PUB_EMAIL}" \
                publisherName="${PUB_NAME}" \
                auth0ClientId="${AUTH0_CLIENT_ID}" \
                auth0Domain="${AUTH0_DOMAIN}" \
                auth0ClientSecret="${AUTH0_CLIENT_SECRET}" \
                delegationKey="${DELEGATION_KEY}" \
                acrName="${ACR_NAME}" \
                acrRepoName="${ACR_REPO_NAME}" \
                imageTag="${IMAGE_TAG}" \
                apprClientId="${APPR_CLIENT_ID}" \
                apprClientSecret="${APPR_CLIENT_SECRET}" \
                apprSpOID="${APPR_SP_OID}" \

echo -e "\n 🎉 Deployment complete\n"