#!/bin/sh
set -e # Exit immediately if a command exits with a non-zero status.

# Check for required environment variables
if [ -z "${LIBSQL_AUTH_TOKEN}" ]; then
  echo "Error: LIBSQL_AUTH_TOKEN is not set or is empty." >&2
  exit 1
fi

# Environment Variables (some will be provided by Fly.io, others by fly.toml or secrets)
# FLY_PRIMARY_INSTANCE (true/false) - Fly.io specific, not standard. We'll use PRIMARY_REGION and FLY_REGION.
# PRIMARY_REGION (e.g., "yyz") - From fly.toml, defines where the primary should run.
# FLY_REGION (e.g., "ams") - Provided by Fly.io, the current region.
# DATABASE_PATH (e.g., "/data/db.sqlite") - From fly.toml
# REPLICA_PATH (e.g., "/data/litestream-replica/db.sqlite") - From fly.toml, for Litestream.
# LIBSQL_AUTH_TOKEN - From Fly.io secrets.
# PORT (e.g., "8080") - The port for libsql-server to listen on.

echo "Starting run.sh script..."
echo "Primary Region: ${PRIMARY_REGION}"
echo "Current Region: ${FLY_REGION}"
echo "Database Path: ${DATABASE_PATH}"
echo "Litestream Replica Path: ${REPLICA_PATH}"
echo "Port: ${PORT}"

# 0. Ensure necessary tools are available
# litestream is in the base image. We need to install libsql-server (sqld).
SQLD_VERSION="v0.20.1" # Specify a version for sqld
SQLD_URL="https://github.com/tursodatabase/libsql/releases/download/${SQLD_VERSION}/sqld-${SQLD_VERSION}-x86_64-unknown-linux-gnu.tar.gz"

if ! command -v sqld > /dev/null; then
    echo "sqld not found. Installing..."
    apk add --no-cache curl tar
    curl -L -o sqld.tar.gz ${SQLD_URL}
    tar -xzf sqld.tar.gz -C /usr/local/bin sqld
    rm sqld.tar.gz
    echo "sqld installed successfully."
else
    echo "sqld already installed."
fi
sqld --version

# 1. Determine role: Primary or Replica
IS_PRIMARY=false
if [ "${FLY_REGION}" = "${PRIMARY_REGION}" ]; then
  IS_PRIMARY=true
  echo "This instance is the PRIMARY in region ${FLY_REGION}."
else
  echo "This instance is a REPLICA in region ${FLY_REGION}."
fi

# 2. Configure and Start Litestream
# Placeholder for Litestream configuration (litestream.yml or env vars)
# Example: litestream replicate -config /path/to/litestream.yml or using env vars
# For now, we'll just use a simple command structure.

DB_DIR=$(dirname "${DATABASE_PATH}")
mkdir -p "${DB_DIR}" # Ensure database directory exists

LITESTREAM_CONFIG_DIR="/etc/litestream.d"
LITESTREAM_CONFIG_PATH="${LITESTREAM_CONFIG_DIR}/config.yml"
mkdir -p "${LITESTREAM_CONFIG_DIR}"

# Create a basic litestream.yml
# Litestream will use S3-compatible storage. For Fly.io, this means using a mounted volume
# as the "S3 bucket" by pointing to a directory on it.
# The REPLICA_PATH from fly.toml will be used by litestream as its replication destination.
cat <<EOF > "${LITESTREAM_CONFIG_PATH}"
dbs:
  - path: ${DATABASE_PATH}
    replicas:
      - type: file # Using "file" type to replicate to another path on the mounted volume
        path: ${REPLICA_PATH}
        # sync-interval: 1s # Optional: How often to sync, default is 1s
EOF

echo "Litestream configuration generated at ${LITESTREAM_CONFIG_PATH}:"
cat "${LITESTREAM_CONFIG_PATH}"

if [ "$IS_PRIMARY" = true ]; then
  echo "Primary: Starting Litestream for replication..."
  # For primary, just ensure the database exists, then replicate
  # If DB exists, Litestream replicates it. If not, sqld will create it, then Litestream replicates.
  litestream replicate -config "${LITESTREAM_CONFIG_PATH}" &
else
  echo "Replica: Attempting to restore database with Litestream..."
  # For replicas, try to restore. If no backup, it will wait.
  # -if-db-not-exists: only restore if DATABASE_PATH doesn't exist
  # -if-replica-exists: only restore if a replica is found at REPLICA_PATH
  litestream restore -if-db-not-exists -if-replica-exists -config "${LITESTREAM_CONFIG_PATH}" "${DATABASE_PATH}"
  echo "Replica: Starting Litestream for continuous replication..."
  litestream replicate -config "${LITESTREAM_CONFIG_PATH}" &
fi

# 3. Start libsql-server
echo "Starting libsql-server..."
# Placeholder for starting libsql-server with authentication and other options.
# Example: sqld --http-listen-addr 0.0.0.0:${PORT} --auth-token "${LIBSQL_AUTH_TOKEN}" "${DATABASE_PATH}"

# Basic command for now, will be expanded in later steps
SQLD_NODE_TYPE="standalone" # default
# TODO: Add --primary-grpc-url for replicas if we implement write forwarding
# if [ "$IS_PRIMARY" = false ]; then
# SQLD_NODE_TYPE="replica" # This isn't a direct flag for sqld yet, more conceptual for our setup
# fi

# Authentication will be added in a later step.
# For now, start without auth to ensure basic setup works.
# When auth is added: --auth-token "${LIBSQL_AUTH_TOKEN}"
# When TLS is added: --https-listen-addr "0.0.0.0:5000" (or whatever TLS port is configured in fly.toml)
# For now, simple HTTP listener on the internal port.

# Construct the primary's gRPC URL. FLY_APP_NAME is provided by Fly.io.
# Note: Using http scheme for internal gRPC for now. If TLS is enforced internally, this would be https.
SQLD_PRIMARY_GRPC_URL="http://${PRIMARY_REGION}.${FLY_APP_NAME}.internal:5001"
echo "Primary gRPC URL for replicas will be: ${SQLD_PRIMARY_GRPC_URL}"

if [ "$IS_PRIMARY" = true ]; then
  echo "Primary: Starting sqld with HTTP and gRPC listeners..."
  exec sqld \
    --http-listen-addr "0.0.0.0:${PORT}" \
    --grpc-listen-addr "0.0.0.0:5001" \
    --auth-token "${LIBSQL_AUTH_TOKEN}" \
    "${DATABASE_PATH}"
else
  # Check if FLY_APP_NAME is set, which is needed for SQLD_PRIMARY_GRPC_URL
  if [ -z "${FLY_APP_NAME}" ]; then
    echo "Error: FLY_APP_NAME is not set. Cannot construct primary gRPC URL." >&2
    exit 1
  fi
  echo "Replica: Starting sqld with HTTP and gRPC listeners, pointing to primary ${SQLD_PRIMARY_GRPC_URL}..."
  exec sqld \
    --http-listen-addr "0.0.0.0:${PORT}" \
    --grpc-listen-addr "0.0.0.0:5001" \
    --primary-grpc-url "${SQLD_PRIMARY_GRPC_URL}" \
    --auth-token "${LIBSQL_AUTH_TOKEN}" \
    "${DATABASE_PATH}"
fi
