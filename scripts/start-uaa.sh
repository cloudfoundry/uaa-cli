#!/usr/bin/env bash
# Builds and launches a real UAA server in the background using the default
# (hsqldb, no external DB/Docker) profile documented in the UAA README, then
# waits until it's ready to serve requests.
#
# Usage: scripts/start-uaa.sh [path-to-uaa-checkout]
#   path-to-uaa-checkout defaults to ./uaa (matches the checkout path used by
#   .github/workflows/integration-test.yml)
set -eu -o pipefail

UAA_DIR="${1:-./uaa}"
READY_URL="http://localhost:8080/uaa/login"
TIMEOUT_SECONDS=300

if [ ! -d "${UAA_DIR}" ]; then
  echo "UAA checkout not found at ${UAA_DIR}" >&2
  exit 1
fi

echo "Building UAA bootWar from ${UAA_DIR}..."
(cd "${UAA_DIR}" && ./gradlew bootWar)

WAR_FILE="$(find "${UAA_DIR}" -name 'cloudfoundry-identity-uaa-*.war' -path '*/build/libs/*' | head -1)"
if [ -z "${WAR_FILE}" ]; then
  echo "Could not find a built UAA war file under ${UAA_DIR}" >&2
  exit 1
fi

echo "Starting UAA from ${WAR_FILE}..."
nohup java \
  -DCLOUDFOUNDRY_CONFIG_PATH="${UAA_DIR}/scripts/boot" \
  -DSECRETS_DIR="${UAA_DIR}/scripts/boot" \
  -Dserver.servlet.context-path=/uaa \
  -Dsmtp.host=localhost \
  -Dsmtp.port=2525 \
  -Dspring.profiles.active=hsqldb \
  -Djava.security.egd=file:/dev/./urandom \
  -jar "${WAR_FILE}" > uaa-boot.log 2>&1 &

echo "Waiting up to ${TIMEOUT_SECONDS}s for UAA to respond at ${READY_URL}..."
elapsed=0
while [ "${elapsed}" -lt "${TIMEOUT_SECONDS}" ]; do
  if curl --silent --fail --max-time 5 --output /dev/null "${READY_URL}"; then
    echo "UAA is up."
    exit 0
  fi
  sleep 5
  elapsed=$((elapsed + 5))
done

echo "UAA failed to start within ${TIMEOUT_SECONDS}s. Boot log:" >&2
cat uaa-boot.log >&2
exit 1
