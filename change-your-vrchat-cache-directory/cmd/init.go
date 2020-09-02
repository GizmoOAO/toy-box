package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	appdata string
)

func init() {
	appdata = os.Getenv("AppData")
	if appdata == "" {
		fmt.Println("The AppData environment varible must be set for app to run correctly.")
		os.Exit(2)
	}
	if strings.Contains(appdata, "\\Roaming") {
		appdata = filepath.Dir(appdata)
	}
}
