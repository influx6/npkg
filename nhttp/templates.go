package nhttp

import (
	htemplate "html/template"
	"text/template"
)

// TextContextFunctions returns a map of tempalte funcs for usage
// with text/template.Template.
func TextContextFunctions(c *Ctx) template.FuncMap {
	return template.FuncMap{
		"flash":         c.Flash,
		"flashMessages": c.FlashMessages,
		"clearFlashMessages": func() string {
			c.ClearFlashMessages()
			return ""
		},
		"clearFlash": func(name string) string {
			c.ClearFlash(name)
			return ""
		},
		"setFlash": func(name, message string) string {
			c.SetFlash(name, message)
			return ""
		},
	}
}

// HTMLContextFunctions returns a map of tempalte funcs for usage
// with text/template.Template.
func HTMLContextFunctions(c *Ctx) htemplate.FuncMap {
	return htemplate.FuncMap{
		"flash":         c.Flash,
		"flashMessages": c.FlashMessages,
		"clearFlashMessages": func() string {
			c.ClearFlashMessages()
			return ""
		},
		"clearFlash": func(name string) string {
			c.ClearFlash(name)
			return ""
		},
		"setFlash": func(name, message string) string {
			c.SetFlash(name, message)
			return ""
		},
	}
}
