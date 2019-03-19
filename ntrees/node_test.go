package ntrees

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type counter int

func (c counter) Respond(s Signal) {
	c++
}

func TestNodeEvent(t *testing.T) {
	var c counter
	var parent = HTMLDiv("page-block", NewEventDescriptor("click", c))

	var clicked = NewClickEvent("1")
	for i := 0; i < 5; i++ {
		parent.Respond(clicked)
	}

	require.Equal(t, 5, int(c))
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
			continue
		}

		child.SwapNode(newChild)
		child = newChild
	}
}
