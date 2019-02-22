// +build js

package nxid

import "github.com/gopherjs/gopherjs/js"

func readPlatformMachineID() (string, error) {
	return js.Global.Get("navigator").Get("platform").String(), nil
}
