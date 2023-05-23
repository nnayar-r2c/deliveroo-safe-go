#!/bin/bash

GOBIN=$1
VERSION=$2

LINTER_EXISTS=1
if [ -f ${GOBIN}/golangci-lint ]; then
  ${GOBIN}/golangci-lint --version | grep -q "version ${VERSION}"
  LINTER_EXISTS=$(echo $?)
fi
if [ "${LINTER_EXISTS}" -gt 0 ]; then
  echo "Installing golangci-lint v${VERSION}"
  curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOBIN} v${VERSION}
fi
