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

	"github.com/influx6/npkg/ntrees"
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
<section  id="767h" atid="bj0v2kuaa78t7ijs5070" _tid="bj0v2kuaa78t7ijs5070" _ref="/767h">
	<!-- 
	Commentary
	 -->
	<div  id="0" atid="bj0v2kuaa78t7ijs508g" _tid="bj0v2kuaa78t7ijs508g" _ref="/767h/0" count-target="0" events="click-00 MouseOver-00">
		0
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="1" atid="bj0v2kuaa78t7ijs50a0" _tid="bj0v2kuaa78t7ijs50a0" _ref="/767h/1" count-target="1" events="click-00 MouseOver-00">
		1
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="2" atid="bj0v2kuaa78t7ijs50bg" _tid="bj0v2kuaa78t7ijs50bg" _ref="/767h/2" count-target="2" events="MouseOver-00 click-00">
		2
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="3" atid="bj0v2kuaa78t7ijs50d0" _tid="bj0v2kuaa78t7ijs50d0" _ref="/767h/3" count-target="3" events="click-00 MouseOver-00">
		3
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="4" atid="bj0v2kuaa78t7ijs50eg" _tid="bj0v2kuaa78t7ijs50eg" _ref="/767h/4" count-target="4" events="MouseOver-00 click-00">
		4
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="5" atid="bj0v2kuaa78t7ijs50g0" _tid="bj0v2kuaa78t7ijs50g0" _ref="/767h/5" count-target="5" events="click-00 MouseOver-00">
		5
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="6" atid="bj0v2kuaa78t7ijs50hg" _tid="bj0v2kuaa78t7ijs50hg" _ref="/767h/6" count-target="6" events="click-00 MouseOver-00">
		6
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="7" atid="bj0v2kuaa78t7ijs50j0" _tid="bj0v2kuaa78t7ijs50j0" _ref="/767h/7" count-target="7" events="click-00 MouseOver-00">
		7
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="8" atid="bj0v2kuaa78t7ijs50kg" _tid="bj0v2kuaa78t7ijs50kg" _ref="/767h/8" count-target="8" events="click-00 MouseOver-00">
		8
	</div>
	<!-- 
	Commentary
	 -->
	<div  id="9" atid="bj0v2kuaa78t7ijs50m0" _tid="bj0v2kuaa78t7ijs50m0" _ref="/767h/9" count-target="9" events="click-00 MouseOver-00">
		9
	</div>
</section>
```

More so, using the `Node.EncodeJSON` we can actually render a JSON format of giving node, attributes
events and children:

```json
{
  "type": 1,
  "ref": "\/767h",
  "typeName": "Element",
  "atid": "bj0v36uaa78tecmm9sn0",
  "name": "section",
  "id": "767h",
  "tid": "bj0v36uaa78tecmm9sn0",
  "attrs": [
    
  ],
  "events": [
    
  ],
  "children": [
    {
      "type": 8,
      "ref": "\/767h\/comment-2042p7w9cx",
      "typeName": "Comment",
      "atid": "bj0v36uaa78tecmm9sng",
      "name": "Comment",
      "content": "Commentary",
      "id": "comment-2042p7w9cx",
      "tid": "bj0v36uaa78tecmm9sng",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": 1,
      "ref": "\/767h\/0",
      "typeName": "Element",
      "atid": "bj0v36uaa78tecmm9sog",
      "name": "div",
      "id": "0",
      "tid": "bj0v36uaa78tecmm9sog",
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
          "type": 3,
          "ref": "\/767h\/0\/text-204250vbsc",
          "typeName": "Text",
          "atid": "bj0v36uaa78tecmm9so0",
          "name": "Text",
          "content": "0",
          "id": "text-204250vbsc",
          "tid": "bj0v36uaa78tecmm9so0",
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
      "type": 8,
      "ref": "\/767h\/comment-2042q98s1h",
      "typeName": "Comment",
      "atid": "bj0v36uaa78tecmm9sp0",
      "name": "Comment",
      "content": "Commentary",
      "id": "comment-2042q98s1h",
      "tid": "bj0v36uaa78tecmm9sp0",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": 1,
      "ref": "\/767h\/1",
      "typeName": "Element",
      "atid": "bj0v36uaa78tecmm9sq0",
      "name": "div",
      "id": "1",
      "tid": "bj0v36uaa78tecmm9sq0",
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
          "type": 3,
          "ref": "\/767h\/1\/text-2042w6hjl8",
          "typeName": "Text",
          "atid": "bj0v36uaa78tecmm9spg",
          "name": "Text",
          "content": "1",
          "id": "text-2042w6hjl8",
          "tid": "bj0v36uaa78tecmm9spg",
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
      "type": 8,
      "ref": "\/767h\/comment-20427kk8bt",
      "typeName": "Comment",
      "atid": "bj0v36uaa78tecmm9sqg",
      "name": "Comment",
      "content": "Commentary",
      "id": "comment-20427kk8bt",
      "tid": "bj0v36uaa78tecmm9sqg",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": 1,
      "ref": "\/767h\/2",
      "typeName": "Element",
      "atid": "bj0v36uaa78tecmm9srg",
      "name": "div",
      "id": "2",
      "tid": "bj0v36uaa78tecmm9srg",
      "attrs": [
        {
          "name": "count-target",
          "value": "2"
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
          "type": 3,
          "ref": "\/767h\/2\/text-2042cxwpm6",
          "typeName": "Text",
          "atid": "bj0v36uaa78tecmm9sr0",
          "name": "Text",
          "content": "2",
          "id": "text-2042cxwpm6",
          "tid": "bj0v36uaa78tecmm9sr0",
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
      "type": 8,
      "ref": "\/767h\/comment-20427ctjrn",
      "typeName": "Comment",
      "atid": "bj0v36uaa78tecmm9ss0",
      "name": "Comment",
      "content": "Commentary",
      "id": "comment-20427ctjrn",
      "tid": "bj0v36uaa78tecmm9ss0",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": 1,
      "ref": "\/767h\/3",
      "typeName": "Element",
      "atid": "bj0v36uaa78tecmm9st0",
      "name": "div",
      "id": "3",
      "tid": "bj0v36uaa78tecmm9st0",
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
          "type": 3,
          "ref": "\/767h\/3\/text-2042gfr78g",
          "typeName": "Text",
          "atid": "bj0v36uaa78tecmm9ssg",
          "name": "Text",
          "content": "3",
          "id": "text-2042gfr78g",
          "tid": "bj0v36uaa78tecmm9ssg",
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
      "type": 8,
      "ref": "\/767h\/comment-2042zf7cmm",
      "typeName": "Comment",
      "atid": "bj0v36uaa78tecmm9stg",
      "name": "Comment",
      "content": "Commentary",
      "id": "comment-2042zf7cmm",
      "tid": "bj0v36uaa78tecmm9stg",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": 1,
      "ref": "\/767h\/4",
      "typeName": "Element",
      "atid": "bj0v36uaa78tecmm9sug",
      "name": "div",
      "id": "4",
      "tid": "bj0v36uaa78tecmm9sug",
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
          "type": 3,
          "ref": "\/767h\/4\/text-2042n5lxrh",
          "typeName": "Text",
          "atid": "bj0v36uaa78tecmm9su0",
          "name": "Text",
          "content": "4",
          "id": "text-2042n5lxrh",
          "tid": "bj0v36uaa78tecmm9su0",
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
      "type": 8,
      "ref": "\/767h\/comment-2042pbtjvl",
      "typeName": "Comment",
      "atid": "bj0v36uaa78tecmm9sv0",
      "name": "Comment",
      "content": "Commentary",
      "id": "comment-2042pbtjvl",
      "tid": "bj0v36uaa78tecmm9sv0",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": 1,
      "ref": "\/767h\/5",
      "typeName": "Element",
      "atid": "bj0v36uaa78tecmm9t00",
      "name": "div",
      "id": "5",
      "tid": "bj0v36uaa78tecmm9t00",
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
          "type": 3,
          "ref": "\/767h\/5\/text-2042pd36rj",
          "typeName": "Text",
          "atid": "bj0v36uaa78tecmm9svg",
          "name": "Text",
          "content": "5",
          "id": "text-2042pd36rj",
          "tid": "bj0v36uaa78tecmm9svg",
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
      "type": 8,
      "ref": "\/767h\/comment-2042dl74kr",
      "typeName": "Comment",
      "atid": "bj0v36uaa78tecmm9t0g",
      "name": "Comment",
      "content": "Commentary",
      "id": "comment-2042dl74kr",
      "tid": "bj0v36uaa78tecmm9t0g",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": 1,
      "ref": "\/767h\/6",
      "typeName": "Element",
      "atid": "bj0v36uaa78tecmm9t1g",
      "name": "div",
      "id": "6",
      "tid": "bj0v36uaa78tecmm9t1g",
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
          "type": 3,
          "ref": "\/767h\/6\/text-2042v0ffkr",
          "typeName": "Text",
          "atid": "bj0v36uaa78tecmm9t10",
          "name": "Text",
          "content": "6",
          "id": "text-2042v0ffkr",
          "tid": "bj0v36uaa78tecmm9t10",
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
      "type": 8,
      "ref": "\/767h\/comment-2042pz33p2",
      "typeName": "Comment",
      "atid": "bj0v36uaa78tecmm9t20",
      "name": "Comment",
      "content": "Commentary",
      "id": "comment-2042pz33p2",
      "tid": "bj0v36uaa78tecmm9t20",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": 1,
      "ref": "\/767h\/7",
      "typeName": "Element",
      "atid": "bj0v36uaa78tecmm9t30",
      "name": "div",
      "id": "7",
      "tid": "bj0v36uaa78tecmm9t30",
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
          "type": 3,
          "ref": "\/767h\/7\/text-2042lvjkbr",
          "typeName": "Text",
          "atid": "bj0v36uaa78tecmm9t2g",
          "name": "Text",
          "content": "7",
          "id": "text-2042lvjkbr",
          "tid": "bj0v36uaa78tecmm9t2g",
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
      "type": 8,
      "ref": "\/767h\/comment-2042qf5ltc",
      "typeName": "Comment",
      "atid": "bj0v36uaa78tecmm9t3g",
      "name": "Comment",
      "content": "Commentary",
      "id": "comment-2042qf5ltc",
      "tid": "bj0v36uaa78tecmm9t3g",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": 1,
      "ref": "\/767h\/8",
      "typeName": "Element",
      "atid": "bj0v36uaa78tecmm9t4g",
      "name": "div",
      "id": "8",
      "tid": "bj0v36uaa78tecmm9t4g",
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
          "type": 3,
          "ref": "\/767h\/8\/text-2042tww0n5",
          "typeName": "Text",
          "atid": "bj0v36uaa78tecmm9t40",
          "name": "Text",
          "content": "8",
          "id": "text-2042tww0n5",
          "tid": "bj0v36uaa78tecmm9t40",
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
      "type": 8,
      "ref": "\/767h\/comment-2042b2xrcd",
      "typeName": "Comment",
      "atid": "bj0v36uaa78tecmm9t50",
      "name": "Comment",
      "content": "Commentary",
      "id": "comment-2042b2xrcd",
      "tid": "bj0v36uaa78tecmm9t50",
      "attrs": [
        
      ],
      "events": [
        
      ],
      "children": [
        
      ]
    },
    {
      "type": 1,
      "ref": "\/767h\/9",
      "typeName": "Element",
      "atid": "bj0v36uaa78tecmm9t60",
      "name": "div",
      "id": "9",
      "tid": "bj0v36uaa78tecmm9t60",
      "attrs": [
        {
          "name": "count-target",
          "value": "9"
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
          "type": 3,
          "ref": "\/767h\/9\/text-20424kwk14",
          "typeName": "Text",
          "atid": "bj0v36uaa78tecmm9t5g",
          "name": "Text",
          "content": "9",
          "id": "text-20424kwk14",
          "tid": "bj0v36uaa78tecmm9t5g",
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
