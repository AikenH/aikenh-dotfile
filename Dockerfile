FROM ubuntu:22.04

# Avoid interactive prompts during install
ENV DEBIAN_FRONTEND=noninteractive

# Base tools needed for bootstrap and testing
RUN apt-get update && apt-get install -y \
    git curl sudo bash zsh ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create test user with sudo
RUN useradd -m -s /bin/bash tester \
    && echo "tester ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

USER tester
WORKDIR /home/tester/dotfile

# Copy the repo in
COPY --chown=tester:tester . .

# Default: run status check
CMD ["bash", "-c", "./dotsetup status || (go build -o dotsetup ./cmd/dotsetup/ && ./dotsetup status)"]
