#!/bin/bash
set -euo pipefail

: "${MAKEFLAGS?'🔥 Whoa there! Script should be run from makefile, not directly silly!'}"

: "${ACR_NAME?'💥 Check ACR NAME is defined in ACR_NAME in .env'}"
: "${ACR_REPO_NAME?'💥 Check repo name is defined in ACR_REPO_NAME in .env'}"
: "${IMAGE_TAG?'💥 Check image tag is defined in IMAGE_TAG in .env'}"

HOST_HOME=${HOST_HOME:-$HOME}

# Build container image from Dockerfile
docker build . --file src/identityApp/Dockerfile \
  --tag "$ACR_NAME/$ACR_REPO_NAME:$IMAGE_TAG"