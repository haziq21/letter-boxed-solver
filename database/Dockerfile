# Use the official Litestream image as a base
# This image is based on Alpine Linux and includes litestream.
FROM litestream/litestream:latest

# Install dependencies needed by run.sh to download and install sqld
# - curl: to download the sqld release
# - tar: to extract the sqld archive
# ca-certificates: for https downloads
RUN apk add --no-cache curl tar ca-certificates

# Copy the run script into the image
COPY run.sh /app/run.sh

# Make the run script executable (though it should already be)
RUN chmod +x /app/run.sh

# Set the working directory
WORKDIR /app

# Expose the port the app will run on (defined by PORT env var in fly.toml, default 8080)
# This is more for documentation; actual port mapping is done in fly.toml
EXPOSE 8080

# Command to run when the container starts
CMD ["/app/run.sh"]
