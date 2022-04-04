#!/usr/bin/env bash

ROOT_DIR=$(dirname "$0")/..
readonly ROOT_DIR

# shellcheck disable=SC1091,SC1090
source "${ROOT_DIR}/vendor/knative.dev/hack/library.sh"

set -o errexit
set -o nounset
set -o pipefail

cd "${ROOT_DIR}"

# This controls the knative release version we track.
KN_VERSION="release-1.2"
EVENTING_VERSION="release-v1.0"
EVENTING_KAFKA_VERSION="release-v1.0"
EVENTING_KAFKA_BROKER_VERSION="release-v1.0"
SERVING_VERSION="release-v1.2.0"

# The list of dependencies that we track at HEAD and periodically
# float forward in this repository.
FLOATING_DEPS=(
  "knative.dev/hack@v0.0.0-20220224013837-e1785985d364"
  "knative.dev/networking@${KN_VERSION}"
  "knative.dev/operator@${KN_VERSION}"
  "knative.dev/pkg@v0.0.0-20220223180940-4fcbc1bc12e8"
)

FLOATING_FORK_DEPS=(
  "knative.dev/eventing-kafka-broker=github.com/openshift-knative/eventing-kafka-broker@${EVENTING_KAFKA_BROKER_VERSION}"
  "knative.dev/eventing=github.com/openshift/knative-eventing@${EVENTING_VERSION}"
  "knative.dev/eventing-kafka=github.com/openshift-knative/eventing-kafka@${EVENTING_KAFKA_VERSION}"
  "knative.dev/serving=github.com/openshift/knative-serving@${SERVING_VERSION}"
)

# Parse flags to determine if we need to update our floating deps.
GO_GET=0
while [[ $# -ne 0 ]]; do
  parameter=$1
  case ${parameter} in
    --upgrade) GO_GET=1 ;;
    *) abort "unknown option ${parameter}" ;;
  esac
  shift
done
readonly GO_GET

if (( GO_GET )); then
  export GOPROXY="https://proxy.golang.org,direct"
  # Treat forks specifically due to https://github.com/golang/go/issues/32721
  for dep in "${FLOATING_FORK_DEPS[@]}"; do
    go mod edit -replace "${dep}"
    # Let the dependency update the magic SHA otherwise the
    # following "go mod edit" will fail.
    go mod vendor
  done
  go get -d "${FLOATING_DEPS[@]}"
fi

# Prune modules.
go mod tidy
go mod vendor

# Remove unnecessary files.
find vendor/ \( -name "OWNERS" \
  -o -name "OWNERS_ALIASES" \
  -o -name "BUILD" \
  -o -name "BUILD.bazel" \
  -o -name "*_test.go" \) -exec rm -fv {} +

find vendor -type f -name '*.sh' -exec chmod +x {} +

# Apply patches
git apply "${ROOT_DIR}"/hack/patches/*
