package ntrees_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/gokit/trees"
	"github.com/stretchr/testify/require"
)

func TestCascadiaQuery(t *testing.T) {
	firstChild := trees.Element("div", "div1", trees.Element("span", "span-2"), trees.Element("span", "span-3"))
	require.NotNil(t, firstChild)

	query, err := trees.Query("div")
	require.NoError(t, err)
	require.NotNil(t, query)

	require.True(t, query.Match(firstChild))

	query2, err := trees.Query("div span#span-2")
	require.NoError(t, err)
	require.NotNil(t, query2)

	dn := query2.MatchAll(firstChild)
	require.Len(t, dn, 1)
}

var elems = []string{"section", "div", "span", "em", "p", "a", "img"}

func generateRandomNodeTree(n int, base *trees.Node) {
	node := trees.NewNode(trees.ElementNode, "section", "section-43")
	for i := 0; i < n; i++ {
		elem := elems[rand.Intn(len(elems))]
		node.AppendChild(trees.NewNode(trees.ElementNode, elem, strconv.Itoa(i)))
	}
	base.AppendChild(node)
}
