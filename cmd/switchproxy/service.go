package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func serviceReload() {
	serviceStart()
	serviceStop()
}

func serviceStart() {
	if runtime.GOOS == "darwin" {
		serviceFile := fmt.Sprintf("%s/Library/LaunchAgents/com.dzianisv.switchproxy.plist", os.Getenv("HOME"))
		exec.Command("launchctl", "load", "-w", serviceFile).Run()
	} else {
		panic("not implemented")
	}
	return
}

func serviceStop() {
	if runtime.GOOS == "darwin" {
		serviceFile := fmt.Sprintf("%s/Library/LaunchAgents/com.dzianisv.switchproxy.plist", os.Getenv("HOME"))
		exec.Command("launchctl", "unload", serviceFile).Run()
	} else {
		panic("not implemented")
	}
	return
}
