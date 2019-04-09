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

More so, using the `Node.EncodeJSON` we can actually render a JSON format of giving node, attributes
events and children:

```json


{
  "type": "elem",
  "name": "section",
  "id": "767h",
  "tid": "bimaunv9glholth2tcmg",
  "attrs": [
    
  ],
  "events": [
    
  ],
  "children": [
    {
      "type": "comment",
      "name": "Comment",
      "content": "Commentary",
      "id": "2042p",
      "tid": "bimaunv9glholth2tcn0",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": "elem",
      "name": "div",
      "id": "0",
      "tid": "bimaunv9glholth2tco0",
      "attrs": [
        {
          "name": "count-target",
          "value": "0"
        }
      ],
      "events": [
        {
          "name": "click",
          "preventDefault": false,
          "stopPropagation": false
        },
        {
          "name": "MouseOver",
          "preventDefault": false,
          "stopPropagation": false
        }
      ],
      "children": [
        {
          "type": "text",
          "name": "Text",
          "content": "0",
          "id": "20427",
          "tid": "bimaunv9glholth2tcng",
          "attrs": [
            
          ],
          "events": [
            
          ],
          "children": [
            
          ]
        }
      ]
    },
    {
      "type": "comment",
      "name": "Comment",
      "content": "Commentary",
      "id": "2042w",
      "tid": "bimaunv9glholth2tcog",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": "elem",
      "name": "div",
      "id": "1",
      "tid": "bimaunv9glholth2tcpg",
      "attrs": [
        {
          "name": "count-target",
          "value": "1"
        }
      ],
      "events": [
        {
          "name": "click",
          "preventDefault": false,
          "stopPropagation": false
        },
        {
          "name": "MouseOver",
          "preventDefault": false,
          "stopPropagation": false
        }
      ],
      "children": [
        {
          "type": "text",
          "name": "Text",
          "content": "1",
          "id": "20429",
          "tid": "bimaunv9glholth2tcp0",
          "attrs": [
            
          ],
          "events": [
            
          ],
          "children": [
            
          ]
        }
      ]
    },
    {
      "type": "comment",
      "name": "Comment",
      "content": "Commentary",
      "id": "2042c",
      "tid": "bimaunv9glholth2tcq0",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": "elem",
      "name": "div",
      "id": "2",
      "tid": "bimaunv9glholth2tcr0",
      "attrs": [
        {
          "name": "count-target",
          "value": "2"
        }
      ],
      "events": [
        {
          "name": "click",
          "preventDefault": false,
          "stopPropagation": false
        },
        {
          "name": "MouseOver",
          "preventDefault": false,
          "stopPropagation": false
        }
      ],
      "children": [
        {
          "type": "text",
          "name": "Text",
          "content": "2",
          "id": "2042x",
          "tid": "bimaunv9glholth2tcqg",
          "attrs": [
            
          ],
          "events": [
            
          ],
          "children": [
            
          ]
        }
      ]
    },
    {
      "type": "comment",
      "name": "Comment",
      "content": "Commentary",
      "id": "20425",
      "tid": "bimaunv9glholth2tcrg",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": "elem",
      "name": "div",
      "id": "3",
      "tid": "bimaunv9glholth2tcsg",
      "attrs": [
        {
          "name": "count-target",
          "value": "3"
        }
      ],
      "events": [
        {
          "name": "click",
          "preventDefault": false,
          "stopPropagation": false
        },
        {
          "name": "MouseOver",
          "preventDefault": false,
          "stopPropagation": false
        }
      ],
      "children": [
        {
          "type": "text",
          "name": "Text",
          "content": "3",
          "id": "20420",
          "tid": "bimaunv9glholth2tcs0",
          "attrs": [
            
          ],
          "events": [
            
          ],
          "children": [
            
          ]
        }
      ]
    },
    {
      "type": "comment",
      "name": "Comment",
      "content": "Commentary",
      "id": "2042v",
      "tid": "bimaunv9glholth2tct0",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": "elem",
      "name": "div",
      "id": "4",
      "tid": "bimaunv9glholth2tcu0",
      "attrs": [
        {
          "name": "count-target",
          "value": "4"
        }
      ],
      "events": [
        {
          "name": "click",
          "preventDefault": false,
          "stopPropagation": false
        },
        {
          "name": "MouseOver",
          "preventDefault": false,
          "stopPropagation": false
        }
      ],
      "children": [
        {
          "type": "text",
          "name": "Text",
          "content": "4",
          "id": "2042b",
          "tid": "bimaunv9glholth2tctg",
          "attrs": [
            
          ],
          "events": [
            
          ],
          "children": [
            
          ]
        }
      ]
    },
    {
      "type": "comment",
      "name": "Comment",
      "content": "Commentary",
      "id": "2042s",
      "tid": "bimaunv9glholth2tcug",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": "elem",
      "name": "div",
      "id": "5",
      "tid": "bimaunv9glholth2tcvg",
      "attrs": [
        {
          "name": "count-target",
          "value": "5"
        }
      ],
      "events": [
        {
          "name": "click",
          "preventDefault": false,
          "stopPropagation": false
        },
        {
          "name": "MouseOver",
          "preventDefault": false,
          "stopPropagation": false
        }
      ],
      "children": [
        {
          "type": "text",
          "name": "Text",
          "content": "5",
          "id": "2042c",
          "tid": "bimaunv9glholth2tcv0",
          "attrs": [
            
          ],
          "events": [
            
          ],
          "children": [
            
          ]
        }
      ]
    },
    {
      "type": "comment",
      "name": "Comment",
      "content": "Commentary",
      "id": "2042q",
      "tid": "bimaunv9glholth2td00",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": "elem",
      "name": "div",
      "id": "6",
      "tid": "bimaunv9glholth2td10",
      "attrs": [
        {
          "name": "count-target",
          "value": "6"
        }
      ],
      "events": [
        {
          "name": "click",
          "preventDefault": false,
          "stopPropagation": false
        },
        {
          "name": "MouseOver",
          "preventDefault": false,
          "stopPropagation": false
        }
      ],
      "children": [
        {
          "type": "text",
          "name": "Text",
          "content": "6",
          "id": "20429",
          "tid": "bimaunv9glholth2td0g",
          "attrs": [
            
          ],
          "events": [
            
          ],
          "children": [
            
          ]
        }
      ]
    },
    {
      "type": "comment",
      "name": "Comment",
      "content": "Commentary",
      "id": "20428",
      "tid": "bimaunv9glholth2td1g",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": "elem",
      "name": "div",
      "id": "7",
      "tid": "bimaunv9glholth2td2g",
      "attrs": [
        {
          "name": "count-target",
          "value": "7"
        }
      ],
      "events": [
        {
          "name": "click",
          "preventDefault": false,
          "stopPropagation": false
        },
        {
          "name": "MouseOver",
          "preventDefault": false,
          "stopPropagation": false
        }
      ],
      "children": [
        {
          "type": "text",
          "name": "Text",
          "content": "7",
          "id": "2042s",
          "tid": "bimaunv9glholth2td20",
          "attrs": [
            
          ],
          "events": [
            
          ],
          "children": [
            
          ]
        }
      ]
    },
    {
      "type": "comment",
      "name": "Comment",
      "content": "Commentary",
      "id": "20421",
      "tid": "bimaunv9glholth2td30",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": "elem",
      "name": "div",
      "id": "8",
      "tid": "bimaunv9glholth2td40",
      "attrs": [
        {
          "name": "count-target",
          "value": "8"
        }
      ],
      "events": [
        {
          "name": "click",
          "preventDefault": false,
          "stopPropagation": false
        },
        {
          "name": "MouseOver",
          "preventDefault": false,
          "stopPropagation": false
        }
      ],
      "children": [
        {
          "type": "text",
          "name": "Text",
          "content": "8",
          "id": "2042h",
          "tid": "bimaunv9glholth2td3g",
          "attrs": [
            
          ],
          "events": [
            
          ],
          "children": [
            
          ]
        }
      ]
    },
    {
      "type": "comment",
      "name": "Comment",
      "content": "Commentary",
      "id": "2042w",
      "tid": "bimaunv9glholth2td4g",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": "elem",
      "name": "div",
      "id": "9",
      "tid": "bimaunv9glholth2td5g",
      "attrs": [
        {
          "name": "count-target",
          "value": "9"
        }
      ],
      "events": [
        {
          "name": "MouseOver",
          "preventDefault": false,
          "stopPropagation": false
        },
        {
          "name": "click",
          "preventDefault": false,
          "stopPropagation": false
        }
      ],
      "children": [
        {
          "type": "text",
          "name": "Text",
          "content": "9",
          "id": "20426",
          "tid": "bimaunv9glholth2td50",
          "attrs": [
            
          ],
          "events": [
            
          ],
          "children": [
            
          ]
        }
      ]
    }
  ]
}
```