#!/bin/bash

# Install the switchproxy binary
GO111MODULE=on go install github.com/dzianisv/switchproxy/cmd/switchproxy@v0.0.1

# Create a default configuration file if it doesn't exist
CONFIG_DIR="${HOME}/.config/switchproxy"
CONFIG_FILE="${CONFIG_DIR}/switchproxy.yaml"
mkdir -p "$CONFIG_DIR"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Creating default configuration file at $CONFIG_FILE"
    cat > "$CONFIG_FILE" <<EOL
rules:
  - domains:
      - ".*"
    proxy: local
EOL
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    # Create a plist file for the launchd daemon
    PLIST_DIR="${HOME}/Library/LaunchAgents"
    PLIST_FILE="${PLIST_DIR}/com.dzianisv.switchproxy.plist"
    mkdir -p "$PLIST_DIR"

    cat > "$PLIST_FILE" <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.dzianisv.switchproxy</string>
    <key>ProgramArguments</key>
    <array>
        <string>$(go env GOPATH)/bin/switchproxy</string>
        <string>-config</string>
        <string>$CONFIG_FILE</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>${HOME}/Library/Logs/switchproxy/error.log</string>
    <key>StandardOutPath</key>
    <string>${HOME}/Library/Logs/switchproxy/output.log</string>
</dict>
</plist>
EOL

    # Load the daemon
    launchctl load -w "$PLIST_FILE"

    # Print success message
    echo "Switchproxy installed and running as a user daemon on macOS."

elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    # Create a systemd service file
    SERVICE_DIR="${HOME}/.config/systemd/user"
    SERVICE_FILE="${SERVICE_DIR}/switchproxy.service"
    mkdir -p "$SERVICE_DIR"

    cat > "$SERVICE_FILE" <<EOL
[Unit]
Description=Switchproxy Service
After=network.target

[Service]
Type=simple
ExecStart=$(go env GOPATH)/bin/switchproxy -config $CONFIG_FILE
Restart=always
RestartSec=5s

[Install]
WantedBy=default.target
EOL

    # Enable and start the service
    systemctl --user enable --now switchproxy.service

    # Print success message
    echo "Switchproxy installed and running as a user daemon on Linux."

else
    echo "Unsupported OS. The script only supports macOS and Linux."
    exit 1
fi
