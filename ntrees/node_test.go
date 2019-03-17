package ntrees_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/gokit/trees"
	"github.com/stretchr/testify/require"
)

func TestNode(t *testing.T) {
	base := trees.NewNode(trees.ElementNode, "red", "767h")
	require.NotNil(t, base)

	for i := 0; i < 1000; i++ {
		require.NoError(t, base.AppendChild(trees.NewNode(trees.ElementNode, fmt.Sprintf("red%d", i), "65jnj")))
	}

	require.Equal(t, 1000, base.ChildCount())
}

func TestNode_Remove(t *testing.T) {
	base := trees.NewNode(trees.ElementNode, "red", "767h")
	require.NotNil(t, base)

	for i := 0; i < 1000; i++ {
		require.NoError(t, base.AppendChild(trees.NewNode(trees.ElementNode, fmt.Sprintf("red%d", i), "65jnj")))
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
	base := trees.NewNode(trees.ElementNode, "red", "767h")
	require.NotNil(t, base)

	for i := 0; i < 1000; i++ {
		require.NoError(t, base.AppendChild(trees.NewNode(trees.ElementNode, fmt.Sprintf("red%d", i), "65jnj")))
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

	base := trees.NewNode(trees.ElementNode, "red", "232")
	for i := 0; i < 1001; i++ {
		base.AppendChild(trees.NewNode(trees.ElementNode, strconv.Itoa(i), "45g"))
	}
}

func BenchmarkNode_Append_Remove_Balance(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	base := trees.NewNode(trees.ElementNode, "red", "232")

	var count = 4001
	for i := 0; i < count; i++ {
		base.AppendChild(trees.NewNode(trees.ElementNode, "sf", strconv.Itoa(i)))
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

	base := trees.NewNode(trees.ElementNode, "red", "232")

	var count = 4001
	for i := 0; i < count; i++ {
		base.AppendChild(trees.NewNode(trees.ElementNode, strconv.Itoa(i), "343f"))
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

	base := trees.NewNode(trees.ElementNode, "red", "322")

	var child *trees.Node
	for i := 0; i < 1001; i++ {
		newChild := trees.NewNode(trees.ElementNode, strconv.Itoa(i), "331")
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

	base := trees.NewNode(trees.ElementNode, "red", "112")

	var child *trees.Node
	for i := 0; i < 1001; i++ {
		newChild := trees.NewNode(trees.ElementNode, strconv.Itoa(i), "132")
		if child == nil {
			base.AppendChild(newChild)
			continue
		}

		child.SwapNode(newChild)
		child = newChild
	}
}
