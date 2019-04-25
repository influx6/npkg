package ntrees

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	
	"github.com/stretchr/testify/require"
	
	"github.com/gokit/npkg/natomic"
)

type counter int

func (c *counter) Respond(_ natomic.Signal) {
	*c++
}

func TestTree(t *testing.T) {
	base := NewNode(ElementNode, "red", "red-div")
	require.NotNil(t, base)
	
	textNode := Text("free-bar")
	span := HTMLSpan("content-span", textNode)
	span2 := HTMLSpan("end-span", span)
	
	require.NoError(t, base.AppendChild(span2))
	require.Equal(t, "/red-div/end-span/content-span",span.RefTree())
}

func TestTreeClone(t *testing.T) {
	base := NewNode(ElementNode, "red", "red-div")
	require.NotNil(t, base)
	
	textNode := Text("free-bar")
	span := HTMLSpan("content-span", textNode)
	span2 := HTMLSpan("end-span", span)
	
	require.NoError(t, base.AppendChild(span2))
	require.Equal(t, "/red-div/end-span/content-span",span.RefTree())
	
	var baseClone = base.Clone(true)
	require.Equal(t, base.atid, baseClone.atid)
	require.Equal(t, base.tid, baseClone.tid)
	require.Equal(t, base.crossEvents, baseClone.crossEvents)
	
	var baseContent, baseCloneContent strings.Builder
	require.NoError(t, base.RenderShallowHTML(&baseContent, false))
	require.NoError(t, baseClone.RenderShallowHTML(&baseCloneContent, false))
	require.Equal(t, baseContent.String(), baseCloneContent.String())
	
	baseContent.Reset()
	baseCloneContent.Reset()
	
	require.NoError(t, base.RenderHTMLTo(&baseContent, false))
	require.NoError(t, baseClone.RenderHTMLTo(&baseCloneContent, false))
	require.Equal(t, baseContent.String(), baseCloneContent.String())
	
	baseContent.Reset()
	baseCloneContent.Reset()
	
	require.NoError(t, base.RenderShallowJSON(&baseContent))
	require.NoError(t, baseClone.RenderShallowJSON(&baseCloneContent))
	require.Equal(t, baseContent.String(), baseCloneContent.String())
	
	baseContent.Reset()
	baseCloneContent.Reset()
	
	require.NoError(t, base.RenderJSON(&baseContent))
	require.NoError(t, baseClone.RenderJSON(&baseCloneContent))
	require.Equal(t, baseContent.String(), baseCloneContent.String())
}

func TestDeeperChangedJSON(t *testing.T) {
	base := NewNode(ElementNode, "red", "red-div")
	bag := HTMLDiv("content-span", Text("free-bar"))
	
	require.NotNil(t, base)
	require.NoError(t, base.AppendChild(HTMLSection("end-section", bag)))
	
	baseClone := base.Clone(true)
	var baseContent, baseCloneContent strings.Builder
	require.NoError(t, base.RenderShallowHTML(&baseContent, true))
	require.NoError(t, baseClone.RenderShallowHTML(&baseCloneContent, true))
	require.Equal(t, baseContent.String(), baseCloneContent.String())
	
	section, err := baseClone.Get(0)
	require.NoError(t, err)
	require.NotNil(t, section)
	
	sectionBag, err := section.Get(0)
	require.NoError(t, err)
	require.NotNil(t, sectionBag)
	
	require.NoError(t, sectionBag.AppendChild(HTMLDiv("red-bag")))
	
	baseContent.Reset()
	require.NoError(t, baseClone.RenderChangesJSON(base, &baseContent))
	require.NotEqual(t,"[]", baseContent.String())
	require.Contains(t, baseContent.String(), "/red-div/end-section/content-span/red-bag")
}


func TestChangedJSON(t *testing.T) {
	base := NewNode(ElementNode, "red", "red-div")
	require.NotNil(t, base)
	require.NoError(t, base.AppendChild(HTMLSpan("end-span", HTMLSpan("content-span", Text("free-bar")))))
	
	baseClone := base.Clone(true)
	
	var baseContent, baseCloneContent strings.Builder
	require.NoError(t, base.RenderShallowHTML(&baseContent, true))
	require.NoError(t, baseClone.RenderShallowHTML(&baseCloneContent, true))
	require.Equal(t, baseContent.String(), baseCloneContent.String())
	
	baseContent.Reset()
	
	require.NoError(t, baseClone.RenderChangesJSON(base, &baseContent))
	require.Equal(t,"[]", baseContent.String())
}

func TestNodeEvent(t *testing.T) {
	var c = new(counter)
	var parent = HTMLDiv("page-block", NewEventDescriptor("click", c))

	var clicked = NewClickEvent("1")
	for i := 0; i < 5; i++ {
		parent.Respond(clicked)
	}

	require.Equal(t, 5, int(*c))
}

func TestNode(t *testing.T) {
	base := NewNode(ElementNode, "red", "767h")
	require.NotNil(t, base)

	for i := 0; i < 1000; i++ {
		require.NoError(t, base.AppendChild(NewNode(ElementNode, fmt.Sprintf("red%d", i), "65jnj")))
	}

	require.Equal(t, 1000, base.ChildCount())
}

func TestNode_Remove(t *testing.T) {
	base := NewNode(ElementNode, "red", "767h")
	require.NotNil(t, base)

	for i := 0; i < 1000; i++ {
		require.NoError(t, base.AppendChild(NewNode(ElementNode, fmt.Sprintf("red%d", i), "65jnj")))
	}

	require.Equal(t, 1000, base.ChildCount())

	deletes := 500
	for i := 0; i < deletes; i++ {
		var target = rand.Intn(deletes)
		if node, err := base.Get(target); err == nil {
			node.Remove()
		}
	}

	require.Equal(t, 500, base.ChildCount())
}

func TestNode_Remove_Balance(t *testing.T) {
	base := NewNode(ElementNode, "red", "767h")
	require.NotNil(t, base)

	for i := 0; i < 1000; i++ {
		require.NoError(t, base.AppendChild(NewNode(ElementNode, fmt.Sprintf("red%d", i), "65jnj")))
	}

	require.Equal(t, 1000, base.ChildCount())

	deletes := 500
	for i := 0; i < deletes; i++ {
		var target = rand.Intn(deletes)
		if node, err := base.Get(target); err == nil {
			node.Remove()
			node.Balance()
		}
	}

	require.Equal(t, 500, base.ChildCount())
}

func BenchmarkNode_Append(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	base := NewNode(ElementNode, "red", "232")
	for i := 0; i < 1001; i++ {
		base.AppendChild(NewNode(ElementNode, strconv.Itoa(i), "45g"))
	}
}

func BenchmarkNode_Append_Remove_Balance(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	base := NewNode(ElementNode, "red", "232")

	var count = 4001
	for i := 0; i < count; i++ {
		base.AppendChild(NewNode(ElementNode, "sf", strconv.Itoa(i)))
	}

	var deletes = count / 2
	for i := 0; i < deletes; i++ {
		var target = rand.Intn(deletes)
		if node, err := base.Get(target); err == nil {
			node.Remove()
		}
	}
	base.Balance()
}

func BenchmarkNode_Append_Remove_With_Balance(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	base := NewNode(ElementNode, "red", "232")

	var count = 4001
	for i := 0; i < count; i++ {
		base.AppendChild(NewNode(ElementNode, strconv.Itoa(i), "343f"))
	}

	var deletes = count / 2
	for i := 0; i < deletes; i++ {
		var target = rand.Intn(deletes)
		if node, err := base.Get(target); err == nil {
			node.Remove()
			base.Balance()
		}
	}
}

func BenchmarkNode_Append_SwapAll(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	base := NewNode(ElementNode, "red", "322")

	var child *Node
	for i := 0; i < 1001; i++ {
		newChild := NewNode(ElementNode, strconv.Itoa(i), "331")
		if child == nil {
			base.AppendChild(newChild)
			child = newChild
			continue
		}

		child.SwapAll(newChild)
		child = newChild
	}
}

func BenchmarkNode_Append_SwapNode(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	base := NewNode(ElementNode, "red", "112")

	var child *Node
	for i := 0; i < 1001; i++ {
		newChild := NewNode(ElementNode, strconv.Itoa(i), "132")
		if child == nil {
			base.AppendChild(newChild)
			child = newChild
			continue
		}

		child.SwapNode(newChild)
		child = newChild
	}
}
