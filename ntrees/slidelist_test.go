package ntrees

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArrayList(t *testing.T) {
	list := newNodeArrayList()

	var err error

	// adding new items into list.
	var f1 = NewNode(ElementNode, "bam", "097")
	_, err = list.Add(f1)
	require.NoError(t, err)

	var f2 = NewNode(ElementNode, "cam", "y676")
	_, err = list.Add(f2)
	require.NoError(t, err)

	var f3 = NewNode(ElementNode, "yam", "5445")
	_, err = list.Add(f3)
	require.NoError(t, err)

	require.Equal(t, 2, list.LastIndex())
	require.Equal(t, 3, list.Length())

	require.Equal(t, f1, list.First())
	require.Equal(t, f3, list.Last())

	var rmNode, rmErr = list.RemoveAndSwap(0)
	require.NoError(t, rmErr)
	require.NotNil(t, rmNode)
	require.Equal(t, f1, rmNode)

	require.Equal(t, f2, list.First())
	require.Equal(t, f3, list.Last())

	var opErr = list.SwapIndex(0)
	require.Error(t, opErr)
	require.Equal(t, ErrIndexNotEmpty, opErr)

	_, err = list.Add(f1)
	require.NoError(t, err)

	require.Equal(t, f2, list.First())
	require.Equal(t, f1, list.Last())

	var rmNode2, rmErr2 = list.RemoveAndSwap(2)
	require.NoError(t, rmErr2)
	require.NotNil(t, rmNode2)
	require.Equal(t, f1, rmNode2)

	_, err = list.Add(f1)
	require.NoError(t, err)

	var f4 = NewNode(ElementNode, "vam", "67")
	_, err = list.Add(f4)
	require.NoError(t, err)

	var f5 = NewNode(ElementNode, "gam", "98")
	_, err = list.Add(f5)
	require.NoError(t, err)

	var f6 = NewNode(ElementNode, "iam", "7676")
	_, err = list.Add(f6)
	require.NoError(t, err)

	require.Equal(t, 6, list.Length())

	var rmNode3, rmErr3 = list.RemoveAndSwap(3)
	require.NoError(t, rmErr3)
	require.NotNil(t, rmNode3)
	require.Equal(t, f4, rmNode3)

	require.Equal(t, f6, list.Last())

	_, err = list.Add(f4)
	require.NoError(t, err)

	var order = []*Node{f2, f3, f1, f5, f6, f4}
	var listItems = list.ToList()
	require.Equal(t, order, listItems)

	var madeList []*Node
	list.Each(func(node *Node, i int) bool {
		madeList = append(madeList, node)
		return true
	})
	require.Equal(t, order, madeList)
}

type m struct {
	m []*Node
}

func (im *m) add(i *Node) {
	im.m = append(im.m, i)
}

func BenchmarkSlice_Add(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	node := NewNode(ElementNode, "bam", "676")
	var bm []*Node
	for i := 0; i < b.N; i++ {
		bm = append(bm, node)
	}
}

// NOTE: Using interface{} as a slice type causes one allocation on itself,
// due to the fact that an interface is a 2 word data type and will get used
// by the system to store both type and value.

func BenchmarkSlice_Add_Through_Struct(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	var bm m
	node := NewNode(ElementNode, "bam", "8987")
	for i := 0; i < b.N; i++ {
		bm.add(node)
	}
}

func BenchmarkSlice_Get_Sequential(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	node := NewNode(ElementNode, "bam", "656")
	var bm []*Node
	for i := 0; i < 1001; i++ {
		bm = append(bm, node)
	}

	for i := 0; i < 1000; i++ {
		_ = bm[i]
	}
}

func BenchmarkSlice_Get_Random(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	node := NewNode(ElementNode, "bam", "tee")
	var bm []interface{}
	for i := 0; i < 1001; i++ {
		bm = append(bm, node)
	}

	var r = 1000
	for i := 0; i < 1000; i++ {
		_ = bm[rand.Intn(r)]
	}
}

func BenchmarkArrayList_Add(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	node := NewNode(ElementNode, "bam", "iu")
	list := newNodeArrayList()
	for i := 0; i < b.N; i++ {
		list.Add(node)
	}
}

func BenchmarkArrayList_Add_Random_RemoveIndex_Swap_ByHistory(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	node := NewNode(ElementNode, "bam", "gyt")
	list := newNodeArrayList()
	for i := 0; i < b.N+1; i++ {
		list.Add(node)
	}

	var history []int
	var sample = b.N / 2
	for i := 0; i < sample; i++ {
		nxt := rand.Intn(sample)

		if _, err := list.RemoveIndex(nxt); err == nil {
			history = append(history, nxt)
		}
	}

	for _, tx := range history {
		list.SwapIndex(tx)
	}

	list.SortList()
}

func BenchmarkArrayList_Add_Remove_Sort(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	node := NewNode(ElementNode, "bam", "gyt")
	list := newNodeArrayList()
	for i := 0; i < b.N+1; i++ {
		list.Add(node)
	}

	for i := 0; i < b.N; i++ {
		list.RemoveAndSwap(i)
		list.SortList()
	}
}

func BenchmarkArrayList_Add_Remove(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	node := NewNode(ElementNode, "bam", "gyt")
	list := newNodeArrayList()
	for i := 0; i < b.N+1; i++ {
		list.Add(node)
	}

	for i := 0; i < b.N; i++ {
		list.RemoveAndSwap(i)
	}
}

func BenchmarkArrayList_Sequential_Add_Get(b *testing.B) {
	b.ReportAllocs()

	node := NewNode(ElementNode, "bam", "ds")
	list := newNodeArrayList()
	for i := 0; i < 1001; i++ {
		list.Add(node)
	}

	for i := 0; i < 1000; i++ {
		list.Get(i)
	}
}

func BenchmarkArrayList_Random_Add_Get(b *testing.B) {
	b.ReportAllocs()

	node := NewNode(ElementNode, "bam", "dsa")
	list := newNodeArrayList()
	for i := 0; i < 1001; i++ {
		list.Add(node)
	}

	r := 1000
	for i := 0; i < 1000; i++ {
		list.Get(rand.Intn(r))
	}
}
