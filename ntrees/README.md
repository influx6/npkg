Trees
-------
Trees provide a virtual tree like DOM structure, for the structure of connecting components
able to react to signals (aka events), fast in iteration and navigation.


## Examples


```go
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gokit/npkg/ntrees"
)

func main() {
	var base = ntrees.Element("section", "767h")

	for i := 0; i < 10; i++ {
		var digit = fmt.Sprintf("%d", i)
		if err := base.AppendChild(ntrees.Comment(ntrees.TextContent("Commentary"))); err != nil {
			log.Fatalf("bad things occured: %+s\n", err)
		}
		if err := base.AppendChild(
			ntrees.Element(
				"div",
				digit,
				ntrees.ClickEvent(nil),
				ntrees.MouseOverEvent(nil),
				ntrees.NewStringAttr("count-target", digit),
				ntrees.Text(
					ntrees.TextContent(digit),
				),
			),
		); err != nil {
			log.Fatalf("bad things occured: %+s\n", err)
		}
	}

	// Render html into giving builder.
	var content strings.Builder
	if err := base.RenderNodeTo(&content, true); err != nil {
		log.Fatalf("failed to render: %+s\n", err)
	}

	fmt.Println(content.String())
}


```

The code above produces the underline html:

```html
<section  id="767h" _tid="bima7q79glhupsnqjo90">
	<!-- 
	Commentary
	 -->
	<div  id="0" _tid="bima7q79glhupsnqjoag" count-target="0" events="click-00 MouseOver-00">
		0
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="1" _tid="bima7q79glhupsnqjoc0" count-target="1" events="MouseOver-00 click-00">
		1
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="2" _tid="bima7q79glhupsnqjodg" count-target="2" events="click-00 MouseOver-00">
		2
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="3" _tid="bima7q79glhupsnqjof0" count-target="3" events="click-00 MouseOver-00">
		3
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="4" _tid="bima7q79glhupsnqjogg" count-target="4" events="click-00 MouseOver-00">
		4
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="5" _tid="bima7q79glhupsnqjoi0" count-target="5" events="click-00 MouseOver-00">
		5
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="6" _tid="bima7q79glhupsnqjojg" count-target="6" events="click-00 MouseOver-00">
		6
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="7" _tid="bima7q79glhupsnqjol0" count-target="7" events="click-00 MouseOver-00">
		7
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="8" _tid="bima7q79glhupsnqjomg" count-target="8" events="click-00 MouseOver-00">
		8
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="9" _tid="bima7q79glhupsnqjoo0" count-target="9" events="click-00 MouseOver-00">
		9
	</div>
</section>
```