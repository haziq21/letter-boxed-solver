#!/bin/sh
set -e # Exit immediately if a command exits with a non-zero status.

# Check for required environment variables
if [ -z "${LIBSQL_AUTH_TOKEN}" ]; then
  echo "Error: LIBSQL_AUTH_TOKEN is not set or is empty." >&2
  exit 1
fi

# Environment Variables (some will be provided by Fly.io, others by fly.toml or secrets)
# FLY_PRIMARY_INSTANCE (true/false) - Fly.io specific, not standard. We'll use PRIMARY_REGION and FLY_REGION.
# PRIMARY_REGION (e.g., "yyz") - From fly.toml, defines where the primary should run (though this script won't use it directly anymore).
# FLY_REGION (e.g., "ams") - Provided by Fly.io, the current region.
# DATABASE_PATH (e.g., "/data/db.sqlite") - From fly.toml
# LIBSQL_AUTH_TOKEN - From Fly.io secrets.
# PORT (e.g., "8080") - The port for libsql-server to listen on.

echo "Starting run.sh script..."
echo "Current Region: ${FLY_REGION}" # Primary Region and Replica Path removed
echo "Database Path: ${DATABASE_PATH}"
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

# Ensure database directory exists
DB_DIR=$(dirname "${DATABASE_PATH}")
mkdir -p "${DB_DIR}"

# 3. Start libsql-server
echo "Starting libsql-server..."
# Placeholder for starting libsql-server with authentication and other options.
# Example: sqld --http-listen-addr 0.0.0.0:${PORT} --auth-token "${LIBSQL_AUTH_TOKEN}" "${DATABASE_PATH}"

# Basic command for now, will be expanded in later steps
# SQLD_NODE_TYPE="standalone" # default
# TODO: Add --primary-grpc-url for replicas if we implement write forwarding

# Authentication will be added in a later step.
# For now, start without auth to ensure basic setup works.
# When auth is added: --auth-token "${LIBSQL_AUTH_TOKEN}"
# When TLS is added: --https-listen-addr "0.0.0.0:5000" (or whatever TLS port is configured in fly.toml)
# For now, simple HTTP listener on the internal port.

# Construct the primary's gRPC URL. FLY_APP_NAME is provided by Fly.io.
# Note: Using http scheme for internal gRPC for now. If TLS is enforced internally, this would be https.
# SQLD_PRIMARY_GRPC_URL="http://${PRIMARY_REGION}.${FLY_APP_NAME}.internal:5001" # Removed for now
# echo "Primary gRPC URL for replicas will be: ${SQLD_PRIMARY_GRPC_URL}" # Removed for now

# The exec sqld command will be updated in the next step to use Raft.
# For now, a simplified placeholder:
echo "Starting libsql-server (sqld) with native clustering..."

# FLY_PRIVATE_IP is automatically available in the Fly.io environment.
# FLY_APP_NAME is also automatically available.
# FLY_ALLOC_ID is also automatically available.

# Placeholder for the port sqld uses for internal Raft/cluster communication.
# This should match a port defined in fly.toml for TCP traffic.
SQLD_RAFT_PORT="5002" # Example, ensure this matches fly.toml in a later step.

# Placeholder for the gRPC port for client connections and possibly some cluster operations.
SQLD_GRPC_PORT="5001" # Example, ensure this matches fly.toml.

exec sqld \
  --node-id "${FLY_ALLOC_ID}" \
  --db-path "${DATABASE_PATH}" \
  --auth-token "${LIBSQL_AUTH_TOKEN}" \
  --http-listen-addr "0.0.0.0:${PORT}" \            # For client HTTP connections
  --grpc-listen-addr "0.0.0.0:${SQLD_GRPC_PORT}" \  # For client gRPC connections
  # Conceptual Raft / Clustering flags:
  # These flags are placeholders and need to be verified against actual sqld documentation.
  # Assumes sqld uses Raft and internal DNS for peer discovery.
  --raft-advertise-addr "${FLY_PRIVATE_IP}:${SQLD_RAFT_PORT}" \
  # The --join flag usually takes a list of known peers to bootstrap the cluster.
  # Using internal DNS for instances.<app>.internal can provide these.
  # The exact format or mechanism might vary for sqld.
  --raft-join-peers "dns:instances.${FLY_APP_NAME}.internal:${SQLD_RAFT_PORT}"
  # Alternative or additional flags might be needed, e.g., for initial cluster formation,
  # or if sqld uses a different discovery mechanism.
  # Example: --cluster-advertise-url "http://${FLY_PRIVATE_IP}:${SQLD_GRPC_PORT}"
  # Example: --cluster-peers "dns:instances.${FLY_APP_NAME}.internal:${SQLD_GRPC_PORT}"
