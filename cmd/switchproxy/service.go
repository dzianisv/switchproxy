package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func getServicefile() string {
	return fmt.Sprintf("%s/Library/LaunchAgents/com.dzianisv.switchproxy.plist", os.Getenv("HOME"))
}

func serviceReload() {
	serviceStart()
	serviceStop()
}

func serviceStart() {
	if runtime.GOOS == "darwin" {
		exec.Command("launchctl", "load", "-w", getServicefile()).Run()
	} else {
		panic("not implemented")
	}
	return
}

func serviceStop() {
	if runtime.GOOS == "darwin" {
		exec.Command("launchctl", "unload", getServicefile()).Run()
	} else {
		panic("not implemented")
	}
	return
}

func serviceInstall() error {
	// Create the plist file.
	plistFile, err := os.Create(getServicefile())
	if err != nil {
		return err
	}
	defer plistFile.Close()
	homeDir := os.Getenv("HOME")

	// Write the plist file contents.
	plistData := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
	<dict>
		<key>Label</key>
		<string>com.dzianisv.switchproxy</string>
		<key>ProgramArguments</key>
		<array>
			<string>/usr/local/bin/switchproxy</string>
			<string>-config</string>
			<string>%s/.config/switchproxy/switchproxy.yaml</string>
		</array>
		<key>RunAtLoad</key>
		<true/>
		<key>KeepAlive</key>
		<true/>
		<key>StandardErrorPath</key>
		<string>%s/Library/Logs/switchproxy/error.log</string>
		<key>StandardOutPath</key>
		<string>%s/Library/Logs/switchproxy/output.log</string>
	</dict>
	</plist>`, homeDir, homeDir, homeDir)

	_, err = plistFile.Write([]byte(plistData))
	if err != nil {
		return err
	}

	return nil
}
