// +build js

package nxid

import (
	"crypto/md5"

	"github.com/gopherjs/gopherjs/js"
)

func readPlatformMachineID() (string, error) {
	var platform = js.Global.Get("navigator").Get("platform").String()
	var buildID = js.Global.Get("navigator").Get("buildID").String()
	var summation = md5.Sum([]byte(platform + buildID))
	return string(summation), nil
}
