#!/bin/bash
# Server-side deploy script. Invoked over SSH from the GitHub Actions
# deploy job. Expects these environment variables to be set:
#   DEPLOY_PATH   - base directory on the server (e.g. /opt/smileyan-backend)
#   RELEASE_TAG   - the tag being deployed (e.g. v1.0.0)
#   TARBALL_PATH  - absolute path of the uploaded tarball
#   S_DB_HOST     - SMILEYAN_BACKEND_DB_HOST value
#   S_DB_USER     - SMILEYAN_BACKEND_DB_USER value
#   S_DB_PASSWORD - SMILEYAN_BACKEND_DB_PASSWORD value
#   S_DB_NAME     - SMILEYAN_BACKEND_DB_NAME value
#   S_REDIS_HOST  - SMILEYAN_BACKEND_REDIS_HOST value
#   S_REDIS_USER  - SMILEYAN_BACKEND_REDIS_USERNAME value
#   S_REDIS_PASS  - SMILEYAN_BACKEND_REDIS_PASSWORD value
#   S_EMAIL_PASS  - SMILEYAN_BACKEND_EMAIL_PASSWORD value
#   S_JWT_SECRET  - SMILEYAN_BACKEND_JWT_SECRET value

set -e
: "${DEPLOY_PATH:?DEPLOY_PATH must be set}"
: "${RELEASE_TAG:?RELEASE_TAG must be set}"
: "${TARBALL_PATH:?TARBALL_PATH must be set}"
: "${S_DB_HOST:?S_DB_HOST must be set}"
: "${S_DB_USER:?S_DB_USER must be set}"
: "${S_DB_PASSWORD:?S_DB_PASSWORD must be set}"
: "${S_DB_NAME:?S_DB_NAME must be set}"
: "${S_REDIS_HOST:?S_REDIS_HOST must be set}"
: "${S_REDIS_USER:?S_REDIS_USER must be set}"
: "${S_REDIS_PASS:?S_REDIS_PASS must be set}"
: "${S_EMAIL_PASS:?S_EMAIL_PASS must be set}"
: "${S_JWT_SECRET:?S_JWT_SECRET must be set}"

RELEASE_PATH="$DEPLOY_PATH/releases/$RELEASE_TAG"
CURRENT_LINK="$DEPLOY_PATH/current"

mkdir -p "$DEPLOY_PATH/releases" "$DEPLOY_PATH/shared"

# Write the shared env file that start.sh sources. The file is owned by
# the deploy user and mode 0600. We use an atomic move so a crashed
# half-written file is never picked up.
SHARED_ENV="$DEPLOY_PATH/shared/.env"
SHARED_ENV_TMP="$SHARED_ENV.tmp.$$"
cat > "$SHARED_ENV_TMP" <<EOF
SMILEYAN_BACKEND_DB_HOST=$S_DB_HOST
SMILEYAN_BACKEND_DB_USER=$S_DB_USER
SMILEYAN_BACKEND_DB_PASSWORD=$S_DB_PASSWORD
SMILEYAN_BACKEND_DB_NAME=$S_DB_NAME
SMILEYAN_BACKEND_REDIS_HOST=$S_REDIS_HOST
SMILEYAN_BACKEND_REDIS_USERNAME=$S_REDIS_USER
SMILEYAN_BACKEND_REDIS_PASSWORD=$S_REDIS_PASS
SMILEYAN_BACKEND_EMAIL_PASSWORD=$S_EMAIL_PASS
SMILEYAN_BACKEND_JWT_SECRET=$S_JWT_SECRET
EOF
chmod 600 "$SHARED_ENV_TMP"
mv -f "$SHARED_ENV_TMP" "$SHARED_ENV"
echo "Wrote $SHARED_ENV (mode 0600)"

# Stop the currently running instance (if any).
if [ -x "$CURRENT_LINK/stop.sh" ]; then
  "$CURRENT_LINK/stop.sh" || true
fi

# Extract the new release to its own folder.
rm -rf "$RELEASE_PATH"
mkdir -p "$RELEASE_PATH"
tar -xzf "$TARBALL_PATH" -C "$RELEASE_PATH"

# Generate run scripts (the tarball does not include them).
cat > "$RELEASE_PATH/start.sh" <<'START'
#!/bin/bash
cd "$(dirname "$0")"
# start.sh lives in releases/<tag>/, so two levels up is $DEPLOY_PATH.
# We need $DEPLOY_PATH/shared/.env.
if [ -f ../../shared/.env ]; then
  set -a
  . ../../shared/.env
  set +a
fi
# The release ships a prebuilt binary built by the CI runner; the host
# does not need Go installed.
if [ ! -x ./smileyan-backend ]; then
  echo "FATAL: smileyan-backend binary not found in $(pwd)" >&2
  exit 1
fi
chmod +x ./smileyan-backend
mkdir -p logs
nohup ./smileyan-backend > logs/app.log 2>&1 &
echo $! > logs/app.pid
START

cat > "$RELEASE_PATH/stop.sh" <<'STOP'
#!/bin/bash
cd "$(dirname "$0")"
if [ -f logs/app.pid ]; then
  kill "$(cat logs/app.pid)" 2>/dev/null || true
  rm -f logs/app.pid
fi
pkill -f smileyan-backend || true
STOP

cat > "$RELEASE_PATH/restart.sh" <<'RESTART'
#!/bin/bash
cd "$(dirname "$0")"
./stop.sh
sleep 1
./start.sh
RESTART

chmod +x "$RELEASE_PATH/start.sh" "$RELEASE_PATH/stop.sh" "$RELEASE_PATH/restart.sh"

# Atomically repoint the live symlink.
ln -sfn "$RELEASE_PATH" "$CURRENT_LINK.tmp"
mv -T "$CURRENT_LINK.tmp" "$CURRENT_LINK"

# Start the new release.
"$CURRENT_LINK/start.sh"
