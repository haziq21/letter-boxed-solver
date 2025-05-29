# libSQL Server on Fly.io with Litestream Replication

This directory contains the configuration for deploying a libSQL server (`sqld`) on Fly.io, with data replication handled by Litestream. This setup ensures persistence and allows for scaling read replicas.

## Overview

- **libSQL Server (`sqld`)**: Provides the SQL database interface, compatible with Turso/libSQL clients.
- **Litestream**: Continuously replicates the database from the primary instance to a Fly.io volume, enabling disaster recovery and providing a base for read replicas.
- **Fly.io**: Hosts the application, manages scaling, and provides persistent volumes.
- **Primary/Replica**:
    - The instance in the `PRIMARY_REGION` (default: `yyz` - Toronto, Canada) acts as the primary write node.
    - Scaled instances in other regions act as read replicas. They restore their state from Litestream and then attempt to forward write operations to the primary via gRPC.

## Prerequisites

- [flyctl](https://fly.io/docs/hands-on/install-flyctl/) installed and authenticated.
- A Fly.io account.

## Deployment

1.  **Change Primary Region (Optional but Recommended)**:
    Open `fly.toml` and change `primary_region = "yyz"` to your preferred Fly.io region (e.g., "sjc", "ams", "sin"). This is where your primary database instance will run.

2.  **Launch the App**:
    Navigate to this `database` directory in your terminal.
    Run the following command to create and launch the Fly.io app. Choose a unique name for your app when prompted if you don't want to use `letter-unboxed-db` (if so, update `fly.toml`'s `app` name too, though the issue specified `letter-unboxed-db`).
    ```bash
    fly launch --name letter-unboxed-db --region <your_chosen_primary_region> --no-deploy
    ```
    - `--name`: Sets the app name (as specified in the issue).
    - `--region`: Sets the initial region for the app, which should match your `primary_region` in `fly.toml`.
    - `--no-deploy`: We want to set secrets before the first deployment.

    If you are just updating an existing app, you might skip `fly launch` or use `fly apps create` appropriately.

3.  **Create a Volume for Database Storage**:
    The application requires a persistent volume to store the SQLite database and Litestream replicas. Create one with:
    ```bash
    fly volumes create db_data --region <your_chosen_primary_region> --size 10
    ```
    - Replace `<your_chosen_primary_region>` with the region you set.
    - `--size 10`: Allocates a 10GB volume. Adjust as needed. The name `db_data` matches what's in `fly.toml`.

4.  **Set Authentication Token**:
    Generate a strong, random string to use as your authentication token. You can use a password generator for this.
    Set this token as a secret in your Fly.io app:
    ```bash
    fly secrets set LIBSQL_AUTH_TOKEN="YOUR_STRONG_SECRET_TOKEN_HERE"
    ```
    The `run.sh` script will fail to start `sqld` if this secret is not set.

5.  **Deploy the App**:
    ```bash
    fly deploy
    ```

## Connecting to the Database

-   **Endpoint**: Your database will be available at `https://letter-unboxed-db.fly.dev` (or `http://` if you connect to the non-TLS port, though TLS is recommended). `sqld` listens on port 8080 internally, which is mapped to 80/443 by Fly's proxy for HTTP/HTTPS. For direct libSQL client connections, use port 5000 (TLS) or 5001 (for gRPC, also used for replica write-forwarding).
    - Primary connection URL (libSQL over HTTPS): `https://letter-unboxed-db.fly.dev:5000` (or `http://letter-unboxed-db.fly.dev:8080` if not using TLS, but the `fly.toml` configures TLS on 5000)
    - Your `LIBSQL_AUTH_TOKEN` will be required for queries. Example with `turso-cli`:
      ```bash
      turso db shell https://letter-unboxed-db.fly.dev:5000 --auth-token "YOUR_STRONG_SECRET_TOKEN_HERE"
      ```
      Or using a libSQL SDK:
      ```javascript
      // Example for Node.js libSQL SDK
      import { createClient } from "@libsql/client";

      const client = createClient({
        url: "libsql://letter-unboxed-db.fly.dev:5000", // Use libsql:// for TLS
        authToken: process.env.LIBSQL_AUTH_TOKEN,
      });
      ```

-   **Primary vs. Replicas**:
    - Connections to the app URL will be routed by Fly.io, typically to the nearest available instance.
    - If an instance is a replica, it will attempt to forward write operations (INSERT, UPDATE, DELETE) to the primary instance (located in `PRIMARY_REGION`) via an internal gRPC connection. Read queries will be served by the replica itself.
    - This relies on `sqld`'s write-forwarding feature.

## Scaling

-   **To add more read replicas**:
    ```bash
    fly scale count <N> # e.g., fly scale count 3
    ```
    Fly will launch new instances in various regions (if your app is multi-region) or in the same region. Instances not in the `PRIMARY_REGION` will automatically configure themselves as replicas, restore from Litestream, and attempt to forward writes.
-   **Primary Node**: There will only be one primary node, located in the `PRIMARY_REGION` defined in `fly.toml`. If the primary instance fails, Fly.io will attempt to restart it. Litestream ensures that a new primary can recover the database state up to the last successful replication.
-   **Scaling to Zero**: The `fly.toml` is configured with `min_machines_running = 0` and `auto_stop_machines = true`. This means the app can scale to zero when not in use (after a period of inactivity) and will automatically start when a request comes in. For a database, you might want to adjust `min_machines_running` to `1` to keep at least the primary always running, especially for production workloads.

## Litestream Configuration

-   Litestream is configured in `run.sh` to create a `litestream.yml`.
-   It replicates the database at `${DATABASE_PATH}` (e.g., `/data/db.sqlite`) to a file-based replica at `${REPLICA_PATH}` (e.g., `/data/litestream-replica/db.sqlite`) on the mounted volume. This means your primary's data is continuously backed up to its own volume. Replicas restore from this path.

## Customization

-   **`sqld` version**: Change `SQLD_VERSION` in `run.sh` to use a different version of `sqld`.
-   **Ports**: Adjust port mappings in `fly.toml` if needed.
-   **Litestream**: Modify `litestream.yml` generation in `run.sh` for more advanced Litestream configurations (e.g., S3 replication, though this setup uses volume-based replication).
