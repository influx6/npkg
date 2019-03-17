// +build ignore

// The generation of this package was inspired by Neelance work on DOM (https://github.com/neelance/dom)

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type event struct {
	Name string
	Link string
	Desc string
}

func main() {
	nameMap := map[string]string{
		"CssRuleViewCSSLinkClicked": "CSSRuleViewCSSLinkClicked",
		"CssRuleViewRefreshed":      "CSSRuleViewRefreshed",
		"CssRuleViewChanged":        "CSSRuleViewChange",
		"cssRuleViewCSSLinkClicked": "CSSRuleViewCSSLinkClicked",
		"cssRuleViewRefreshed":      "CSSRuleViewRefreshed",
		"cssRuleViewChanged":        "CSSRuleViewChange",
		"afterprint":                "AfterPrint",
		"animationend":              "AnimationEnd",
		"animationiteration":        "AnimationIteration",
		"animationstart":            "AnimationStart",
		"appinstalled":              "ApplicationInstalled",
		"audioprocess":              "AudioProcess",
		"audioend":                  "AudioEnd",
		"audiostart":                "AudioStart",
		"beforeprint":               "BeforePrint",
		"beforeunload":              "BeforeUnload",
		"canplay":                   "CanPlay",
		"canplaythrough":            "CanPlayThrough",
		"chargingchange":            "ChargingChange",
		"chargingtimechange":        "ChargingTimeChange",
		"compassneedscalibration":   "CompassNeedsCalibration",
		"compositionend":            "CompositionEnd",
		"compositionstart":          "CompositionStart",
		"compositionupdate":         "CompositionUpdate",
		"contextmenu":               "ContextMenu",
		"dblclick":                  "DoubleClick",
		"devicechange":              "DeviceChange",
		"devicelight":               "DeviceLight",
		"devicemotion":              "DeviceMotion",
		"deviceorientation":         "DeviceOrientation",
		"deviceproximity":           "DeviceProximity",
		"dischargingtimechange":     "DischargingTimeChange",
		"dragend":                   "DragEnd",
		"dragenter":                 "DragEnter",
		"dragleave":                 "DragLeave",
		"dragover":                  "DragOver",
		"dragstart":                 "DragStart",
		"durationchange":            "DurationChange",
		"focusin":                   "FocusIn",
		"focusout":                  "FocusOut",
		"fullscreenchange":          "FullScreenChange",
		"fullscreenerror":           "FullScreenError",
		"gamepadconnected":          "GamepadConnected",
		"gamepaddisconnected":       "GamepadDisconnected",
		"gotpointercapture":         "GotPointerCapture",
		"hashchange":                "HashChange",
		"keydown":                   "KeyDown",
		"keypress":                  "KeyPress",
		"keyup":                     "KeyUp",
		"languagechange":            "LanguageChange",
		"levelchange":               "LevelChange",
		"loadeddata":                "LoadedData",
		"loadedmetadata":            "LoadedMetadata",
		"loadend":                   "LoadEnd",
		"loadstart":                 "LoadStart",
		"lostpointercapture":        "LostPointerCapture",
		"messageerror":              "MessageError",
		"mousedown":                 "MouseDown",
		"mouseenter":                "MouseEnter",
		"mouseleave":                "MouseLeave",
		"mousemove":                 "MouseMove",
		"mouseout":                  "MouseOut",
		"mouseover":                 "MouseOver",
		"mouseup":                   "MouseUp",
		"noupdate":                  "NoUpdate",
		"nomatch":                   "NoMatch",
		"notificationclick":         "NotificationClick",
		"orientationchange":         "OrientationChange",
		"pagehide":                  "PageHide",
		"pageshow":                  "PageShow",
		"pointercancel":             "PointerCancel",
		"pointerdown":               "PointerDown",
		"pointerenter":              "PointerEnter",
		"pointerleave":              "PointerLeave",
		"pointerlockchange":         "PointerLockChange",
		"pointerlockerror":          "PointerLockError",
		"pointermove":               "PointerMove",
		"pointerout":                "PointerOut",
		"pointerover":               "PointerOver",
		"pointerup":                 "PointerUp",
		"popstate":                  "PopState",
		"pushsubscriptionchange":    "PushSubscriptionChange",
		"ratechange":                "RateChange",
		"readystatechange":          "ReadyStateChange",
		"resourcetimingbufferfull":  "ResourceTimingBufferFull",
		"selectstart":               "SelectStart",
		"selectionchange":           "SelectionChange",
		"slotchange":                "SlotChange",
		"soundend":                  "SoundEnd",
		"soundstart":                "SoundStart",
		"speechend":                 "SpeechEnd",
		"speechstart":               "SpeechStart",
		"timeupdate":                "TimeUpdate",
		"touchcancel":               "TouchCancel",
		"touchend":                  "TouchEnd",
		"touchenter":                "TouchEnter",
		"touchleave":                "TouchLeave",
		"touchmove":                 "TouchMove",
		"touchstart":                "TouchStart",
		"transitionend":             "TransitionEnd",
		"updateready":               "UpdateReady",
		"upgradeneeded":             "UpgradeNeeded",
		"userproximity":             "UserProximity",
		"versionchange":             "VersionChange",
		"visibilitychange":          "VisibilityChange",
		"voiceschanged":             "VoicesChanged",
		"volumechange":              "VolumeChange",
		"vrdisplayconnected":        "VRDisplayConnected",
		"vrdisplaydisconnected":     "VRDisplayDisconnected",
		"vrdisplaypresentchange":    "VRDisplayPresentChange",
	}

	ignore := map[string]bool{
		"error": true,
	}

	doc, err := goquery.NewDocument("https://developer.mozilla.org/en-US/docs/Web/Events")
	if err != nil {
		panic(err)
	}

	events := make(map[string]*event)

	tables := doc.Find(".standard-table")

	tables.Each(func(ind int, item *goquery.Selection) {
		item.Eq(0).Find("tr").Each(func(i int, s *goquery.Selection) {
			cols := s.Find("td")
			if cols.Length() == 0 || cols.Find(".icon-thumbs-down-alt").Length() != 0 {
				return
			}
			link := cols.Eq(0).Find("a").Eq(0)
			var e event

			// fmt.Printf("Name: %q -> %q\n", nameMap[link.Text()], link.Text())

			if newName, ok := nameMap[link.Text()]; ok {
				e.Name = newName
			} else {
				e.Name = link.Text()
			}

			e.Link, _ = link.Attr("href")
			e.Desc = strings.TrimSpace(cols.Eq(3).Text())
			if e.Desc == "" {
				e.Desc = "(no documentation)"
			}

			funName := e.Name
			if strings.Contains(funName, "-") {
				parts := strings.Split(funName, "-")
				for ind, sm := range parts {
					parts[ind] = capitalize(sm)
				}

				funName = strings.Join(parts, "")
			}

			funName = capitalize(funName)

			if e.Name == "" || ignore[e.Name] {
				return
			}

			events[funName] = &e
		})
	})

	var names []string
	for name := range events {
		names = append(names, name)
	}

	sort.Strings(names)

	file, err := os.Create("event.gen.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprint(file, `// The generation of this package was inspired by Neelance work on DOM (https://github.com/neelance/dom)

//go:generate go run generate.go

// Documentation source: "Event reference" by Mozilla Contributors, https://developer.mozilla.org/en-US/docs/Web/Events, licensed under CC-BY-SA 2.5.

//Package events defines the event binding system that combines different libraries to create a interesting event system.
package events

import (
	"github.com/gokit/trees"
	"github.com/gokit/npkg/natomic"
)

`)

	for _, name := range names {
		e := events[name]
		fmt.Fprintf(file, `
// %sEvent Documentation is as below: %q
// https://developer.mozilla.org%s
func %sEvent(responder natomic.SignalResponder) EventDescription {
	var handler EventHandler
	ops := append([]trees.EventOptions{trees.EventType(%q)}, options...)

	ev := trees.NewEvent(ops...)

	eventHandler := common.NewEventBroadcastHandler(func (evm common.EventBroadcast){
		if ev.ID() != evm.EventID{
			return
		}

		handler(evm.Event, ev.Tree)
	})

	return ev
}
`, name, descToComments(e.Desc), e.Link[6:], name, e.Name)
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

func descToComments(desc string) string {
	c := ""
	length := 80
	for _, word := range strings.Fields(desc) {
		if length+len(word)+1 > 80 {
			length = 3
			c += "\n//"
		}
		c += " " + word
		length += len(word) + 1
	}
	return c
}
