#!/usr/bin/env bash

# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0

# docker pull golangci/golangci-lint:v1.27.0
# docker pull hadolint/hadolint:latest-alpine
echo "Setting required environment variables.."
grep -q -F 'export DEVBOX_TRUE="true"' ~/.bash_profile
if [ $? -ne 0 ]; then
tee -a ~/.bash_profile <<EOF
export DEVBOX_TRUE="true"
EOF
fi
source ~/.bash_profile
