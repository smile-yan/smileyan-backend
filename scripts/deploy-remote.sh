#!/bin/bash
# Server-side deploy script. Invoked over SSH from the GitHub Actions
# deploy job. Expects three environment variables to be set:
#   DEPLOY_PATH  - base directory on the server (e.g. /opt/smileyan-backend)
#   RELEASE_TAG  - the tag being deployed (e.g. v1.0.0)
#   TARBALL_PATH - absolute path of the uploaded tarball

set -e
: "${DEPLOY_PATH:?DEPLOY_PATH must be set}"
: "${RELEASE_TAG:?RELEASE_TAG must be set}"
: "${TARBALL_PATH:?TARBALL_PATH must be set}"

RELEASE_PATH="$DEPLOY_PATH/releases/$RELEASE_TAG"
CURRENT_LINK="$DEPLOY_PATH/current"

mkdir -p "$DEPLOY_PATH/releases" "$DEPLOY_PATH/shared"

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
if [ -f ../shared/.env ]; then
  set -a
  . ../shared/.env
  set +a
fi
# Build the binary on first run; the host has Go installed.
if [ ! -x ./smileyan-backend ]; then
  go build -ldflags="-s -w" -o smileyan-backend .
fi
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
