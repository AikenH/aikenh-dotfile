FROM ubuntu:22.04

# Avoid interactive prompts during install
ENV DEBIAN_FRONTEND=noninteractive

# Minimal base: what a fresh server typically has
RUN apt-get update && apt-get install -y \
    git curl sudo bash zsh ca-certificates file \
    && rm -rf /var/lib/apt/lists/*

# Create test user with sudo (simulates real user)
RUN useradd -m -s /bin/zsh tester \
    && echo "tester ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

USER tester
WORKDIR /home/tester/dotfile

# Copy repo contents
COPY --chown=tester:tester . .

# Ensure the pre-built linux binary is executable
RUN chmod +x ./dotsetup 2>/dev/null || true

# Default: interactive shell for manual testing
CMD ["bash"]
