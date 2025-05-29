# libSQL Server on Fly.io with Native Replication

This directory contains the configuration for deploying a libSQL server (`sqld`) on Fly.io, utilizing `sqld`'s native clustering capabilities for data replication and fault tolerance.

## Overview

- **libSQL Server (`sqld`)**: Provides the SQL database interface, compatible with Turso/libSQL clients. It's configured to run in a clustered mode.
- **Native Clustering**: `sqld` uses a consensus algorithm (e.g., Raft, assumed for this setup) to replicate data across multiple instances. This ensures data persistence, high availability, and allows for scaling read capacity.
- **Fly.io**: Hosts the application, manages scaling, and provides persistent volumes for each node in the cluster.
- **Cluster Dynamics**:
    - The `sqld` cluster elects a leader (primary) node which is responsible for handling write operations.
    - Other nodes in the cluster act as followers (replicas). They replicate data from the leader.
    - Read queries can typically be served by any node in the cluster, distributing the read load. Write operations are automatically forwarded to or handled by the leader.

## Prerequisites

- [flyctl](https://fly.io/docs/hands-on/install-flyctl/) installed and authenticated.
- A Fly.io account.

## Deployment

1.  **Change Primary Region (Optional but Recommended for Initial Launch)**:
    Open `fly.toml` and change `primary_region = "yyz"` to your preferred Fly.io region (e.g., "sjc", "ams", "sin"). While `sqld` manages its own leader election within the cluster, Fly.io uses this `PRIMARY_REGION` for initial instance placement and provides it as an environment variable.

2.  **Launch the App**:
    Navigate to this `database` directory in your terminal.
    Run the following command to create and launch the Fly.io app. Choose a unique name for your app if you don't want to use `letter-unboxed-db` (if so, update `fly.toml`'s `app` name too).
    ```bash
    fly launch --name letter-unboxed-db --region <your_chosen_primary_region> --no-deploy
    ```
    - `--name`: Sets the app name.
    - `--region`: Sets the initial region for the app.
    - `--no-deploy`: We want to set secrets and potentially adjust scaling before the first deployment.

    If you are just updating an existing app, you might skip `fly launch` or use `fly apps create` appropriately.

3.  **Create a Volume for Database Storage**:
    Each instance in the `sqld` cluster requires a persistent volume to store its portion of the SQLite database. You'll need to create a volume for at least your initial primary instance. If you scale out to more instances, each will need its own volume in its respective region.
    Create the first volume in your primary region:
    ```bash
    fly volumes create db_data --region <your_chosen_primary_region> --size 10
    ```
    - Replace `<your_chosen_primary_region>` with the region you set.
    - `--size 10`: Allocates a 10GB volume. Adjust as needed. The name `db_data` matches what's in `fly.toml`.
    - **Note on Scaling**: When scaling to more machines, each new machine will require its own volume in its region. You might need to provision these manually or automate it if Fly.io doesn't automatically create them on scale (behavior can vary).

4.  **Set Authentication Token**:
    Generate a strong, random string to use as your authentication token.
    Set this token as a secret in your Fly.io app:
    ```bash
    fly secrets set LIBSQL_AUTH_TOKEN="YOUR_STRONG_SECRET_TOKEN_HERE"
    ```
    The `run.sh` script will fail to start `sqld` if this secret is not set.

5.  **Deploy the App**:
    ```bash
    fly deploy
    ```
    This will deploy an initial instance. If you plan to have a multi-node cluster immediately, you might want to scale before the first deploy or shortly after.

## Connecting to the Database

-   **Endpoint**: Your database cluster will be accessible via the Fly.io app URL, e.g., `https://letter-unboxed-db.fly.dev`.
    - Client Connections (libSQL over HTTPS/gRPC): `libsql://letter-unboxed-db.fly.dev:5001` (preferred, uses gRPC) or `https://letter-unboxed-db.fly.dev:5000` (HTTPS). Port 5001 is configured for gRPC and port 5000 for HTTPS in `fly.toml`.
    - The `LIBSQL_AUTH_TOKEN` is required for all connections.
    - Example with `turso-cli`:
      ```bash
      turso db shell libsql://letter-unboxed-db.fly.dev:5001 --auth-token "YOUR_STRONG_SECRET_TOKEN_HERE"
      ```
    - Example with libSQL SDK (Node.js):
      ```javascript
      import { createClient } from "@libsql/client";

      const client = createClient({
        url: "libsql://letter-unboxed-db.fly.dev:5001", // Use libsql:// for gRPC
        authToken: process.env.LIBSQL_AUTH_TOKEN,
      });
      ```

-   **Cluster Operation**:
    - Connections to the app URL are routed by Fly.io to an available instance in the cluster.
    - The `sqld` cluster internally handles write operations by ensuring they are processed by the current leader. Read operations can be served by any available node.

## Scaling

-   **To add more nodes to the cluster**:
    ```bash
    fly scale count <N> # e.g., fly scale count 3 to have a 3-node cluster
    ```
    - New `sqld` instances launched by Fly.io should automatically attempt to discover and join the existing cluster. This typically relies on the peer discovery mechanism configured in `run.sh` (e.g., DNS-based using `--raft-join-peers "dns:instances.${FLY_APP_NAME}.internal:${SQLD_RAFT_PORT}"`).
    - Each new node will synchronize its state with the cluster leader through `sqld`'s native replication protocol.
    - **Important**: Ensure each new machine has a persistent volume available in its region. You may need to run `fly volumes create db_data --region <new_instance_region>` for each new instance if not automatically provisioned.

-   **Leader Election**: The `sqld` cluster manages leader election automatically. If a leader node fails, the remaining nodes should elect a new leader. Fly.io will attempt to restart failed instances, which will then rejoin the cluster.

-   **Scaling to Zero**: The `fly.toml` is configured with `min_machines_running = 0` and `auto_stop_machines = true`. This allows the app to scale to zero when not in use. However, for a clustered database, scaling to zero means the entire cluster shuts down. For production workloads, you'll likely want `min_machines_running` to be at least 1 (for a single node setup) or ideally 3 or more for a fault-tolerant cluster, ensuring the cluster remains operational.

## Customization

-   **`sqld` version**: Change `SQLD_VERSION` in `run.sh` to use a different version of `sqld`.
-   **Ports**: Adjust port mappings in `fly.toml` if needed. The key ports are:
    - `8080` (internal HTTP, mapped to 80/443 by Fly proxy)
    - `5000` (external TLS for libSQL)
    - `5001` (external gRPC for libSQL)
    - `5002` (internal TCP for Raft/cluster communication, not exposed externally)
-   **`sqld` Clustering Flags**: The `run.sh` script uses placeholder flags for `sqld` clustering (e.g., `--node-id`, `--raft-advertise-addr`, `--raft-join-peers`, and assumed Raft port `5002`). These are based on common patterns for clustered applications. You should consult the official `sqld` documentation for the most accurate and up-to-date flags and recommended configurations for a robust clustered deployment. You may need to adjust these flags, environment variables, or port configurations in `fly.toml` and `run.sh` accordingly.
