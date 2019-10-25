package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/influx6/npkg/ntrees"
)

func main() {
	var base = ntrees.Element("section", "767h")

	for i := 0; i < 10; i++ {
		var digit = fmt.Sprintf("%d", i)
		if err := base.AppendChild(ntrees.Comment("Commentary")); err != nil {
			log.Fatalf("bad things occured: %s\n", err)
		}
		if err := base.AppendChild(
			ntrees.Element(
				"div",
				digit,
				ntrees.ClickEvent(nil),
				ntrees.MouseOverEvent(nil),
				ntrees.NewStringAttr("count-target", digit),
				ntrees.Text(digit),
			),
		); err != nil {
			log.Fatalf("bad things occured: %+s\n", err)
		}
	}

	// Render html into giving builder.
	var content strings.Builder
	if err := base.RenderHTMLTo(&content, true); err != nil {
		log.Fatalf("failed to render: %s\n", err)
	}

	fmt.Println(content.String())
}
