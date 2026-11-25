package structs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func checkEqual(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func checkOK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func checkError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error")
	}
}

type fileStructure interface {
	SaveText(string) error
	LoadText(string) error
	SaveBin(string) error
	LoadBin(string) error
}

func testFiles(t *testing.T, src, dst fileStructure, verify func()) {
	t.Helper()
	dir := t.TempDir()
	text := filepath.Join(dir, "data.json")
	bin := filepath.Join(dir, "data.bin")
	missing := filepath.Join(dir, "missing", "file")
	checkOK(t, src.SaveText(text))
	checkOK(t, dst.LoadText(text))
	verify()
	checkOK(t, src.SaveBin(bin))
	checkOK(t, dst.LoadBin(bin))
	verify()
	checkError(t, src.SaveText(missing))
	checkError(t, src.SaveBin(missing))
	checkError(t, dst.LoadText(missing))
	checkError(t, dst.LoadBin(missing))
	checkOK(t, os.WriteFile(text, []byte("{"), 0600))
	checkError(t, dst.LoadText(text))
	for _, data := range [][]byte{{}, {255, 255, 255, 255}, {1, 0, 0, 0}} {
		checkOK(t, os.WriteFile(bin, data, 0600))
		checkError(t, dst.LoadBin(bin))
	}
}

func TestExtraArray(t *testing.T) {
	a := NewArray[int32]()
	a.Push(10)
	checkOK(t, a.PushIdx(0, 5))
	checkOK(t, a.Set(1, 20))
	v, err := a.Get(1)
	checkOK(t, err)
	checkEqual(t, v, int32(20))
	checkOK(t, a.Del(0))
	checkEqual(t, a.Len(), 1)
	checkError(t, a.PushIdx(-1, 0))
	checkError(t, a.Set(1, 0))
	checkError(t, a.Del(-1))
	_, err = a.Get(-1)
	checkError(t, err)
	b := NewArray[int32]()
	testFiles(t, a, b, func() { checkEqual(t, b.data, []int32{20}) })
	a.Clear()
	checkEqual(t, a.Len(), 0)
}

func TestExtraForwardList(t *testing.T) {
	f := NewForwardList[int]()
	checkEqual(t, f.FDELBACK(), false)
	f.PushFront(2)
	f.PushBack(4)
	f.FPUSHB(6)
	checkEqual(t, f.FPUSH_BEFORE(1, 3), true)
	checkEqual(t, f.FPUSH_AFTER(2, 5), true)
	checkEqual(t, f.FFIND(5), 3)
	checkEqual(t, f.FDELVAL(2), true)
	checkEqual(t, f.FDELVAL(5), true)
	checkEqual(t, f.FDELBACK(), true)
	checkEqual(t, f.FLen(), 2)
	checkEqual(t, f.FDELBACK(), true)
	checkEqual(t, f.FDELBACK(), true)
	checkEqual(t, f.FLen(), 0)
	checkEqual(t, f.FDEL_AT(-1), false)
	checkEqual(t, f.FDEL_AFTER(-1), false)
	checkEqual(t, f.FDEL_BEFORE(0), false)
	checkEqual(t, f.FPUSH_AFTER(-1, 0), false)
	checkEqual(t, f.FPUSH_BEFORE(-1, 0), false)
	f.PushBack(7)
	f.PushBack(8)
	f.FPRINT()
	f.FPRINTREC(f.Head)
	g := NewForwardList[int]()
	testFiles(t, f, g, func() { checkEqual(t, g.FLen(), 2); checkEqual(t, g.FFIND(7), 0); checkEqual(t, g.FFIND(8), 1) })
}

func TestExtraList(t *testing.T) {
	l := NewList[int]()
	l.LPUSH(2)
	l.LPUSH(1)
	l.LPUSHB(4)
	checkEqual(t, l.LPUSH_AFTER(1, 3), true)
	checkEqual(t, l.LPUSH_AFTER(3, 5), true)
	checkEqual(t, l.LPUSH_BEFORE(0, 0), true)
	checkEqual(t, l.LPUSH_BEFORE(2, 9), true)
	checkEqual(t, l.LPUSH_BEFORE(l.LLen(), 6), true)
	checkEqual(t, l.LFIND(9), 2)
	checkEqual(t, l.LDELNODE(2), true)
	checkEqual(t, l.LDEL_AFTER(1), true)
	checkEqual(t, l.LDEL_BEFORE(1), true)
	checkEqual(t, l.LDELVAL(6), true)
	checkEqual(t, l.LDEL_AFTER(2), true)
	checkEqual(t, l.LLen(), 3)
	checkEqual(t, l.LDELVAL(99), false)
	checkEqual(t, l.LPUSH_AFTER(-1, 0), false)
	checkEqual(t, l.LPUSH_AFTER(99, 0), false)
	checkEqual(t, l.LPUSH_BEFORE(99, 0), false)
	checkEqual(t, l.LDEL_AT(99), false)
	checkEqual(t, l.LDEL_AFTER(-1), false)
	checkEqual(t, l.LDEL_AFTER(99), false)
	checkEqual(t, l.LDEL_BEFORE(0), false)
	l.LPRINT()
	l.LPRINTR()
	l.LPRINTREC(l.Head)
	g := NewList[int]()
	testFiles(t, l, g, func() {
		checkEqual(t, g.LLen(), 3)
		checkEqual(t, g.Head.Value, 1)
		checkEqual(t, g.Tail.Value, 4)
		checkEqual(t, g.Tail.Prev.Value, 3)
	})
	l.Clear()
	checkEqual(t, l.LLen(), 0)
	l.LPUSHB(1)
	l.LPUSHB(2)
	checkEqual(t, l.LDELBACK(), true)
	checkEqual(t, l.Tail.Next == nil, true)
	checkEqual(t, l.LDEL(), true)
	checkEqual(t, l.Tail == nil, true)
}

func TestExtraStackQueue(t *testing.T) {
	s := NewStack[int32]()
	s.Push(1)
	s.Push(2)
	s2 := NewStack[int32]()
	testFiles(t, s, s2, func() { checkEqual(t, s2.data, []int32{1, 2}) })
	q := NewQueue[int32]()
	q.Push(1)
	q.Push(2)
	q2 := NewQueue[int32]()
	testFiles(t, q, q2, func() { checkEqual(t, q2.data, []int32{1, 2}) })
}

func TestExtraHashTree(t *testing.T) {
	h := NewHashTable()
	h.Insert("a", "one")
	h.Insert("k", "two")
	checkEqual(t, h.Delete("a"), true)
	_, ok := h.Find("a")
	checkEqual(t, ok, false)
	h.Insert("a", "three")
	h2 := NewHashTable()
	testFiles(t, h, h2, func() {
		v, ok := h2.Find("a")
		checkEqual(t, ok, true)
		checkEqual(t, v, "three")
		v, ok = h2.Find("k")
		checkEqual(t, ok, true)
		checkEqual(t, v, "two")
	})
	tr := NewTree()
	for i := 1; i <= 7; i++ {
		tr.Insert(i)
	}
	tr2 := NewTree()
	testFiles(t, tr, tr2, func() {
		checkEqual(t, tr2.Size, 7)
		for i := 1; i <= 7; i++ {
			checkEqual(t, tr2.getNode(i).Value, i)
		}
	})
	checkEqual(t, tr.getNode(0) == nil, true)
}

func TestExtraBinaryTypes(t *testing.T) {
	a := NewArray[int]()
	a.Push(12)
	b := NewArray[int]()
	testFiles(t, a, b, func() { checkEqual(t, b.data, []int{12}) })
	f := NewForwardList[int32]()
	f.PushBack(12)
	f2 := NewForwardList[int32]()
	testFiles(t, f, f2, func() { checkEqual(t, f2.Head.Value, int32(12)) })
	l := NewList[int32]()
	l.PushBack(12)
	l2 := NewList[int32]()
	testFiles(t, l, l2, func() { checkEqual(t, l2.Head.Value, int32(12)) })
	s := NewStack[int]()
	s.Push(1)
	s.Push(2)
	s2 := NewStack[int]()
	testFiles(t, s, s2, func() { checkEqual(t, s2.data, []int{1, 2}) })
	v, ok := s.Top()
	checkEqual(t, ok, true)
	checkEqual(t, v, 2)
	checkEqual(t, s.Pop(), true)
	checkEqual(t, s.Pop(), true)
	checkEqual(t, s.Pop(), false)
	checkEqual(t, s.Empty(), true)
	_, ok = s.Top()
	checkEqual(t, ok, false)
	q := NewQueue[int]()
	q.Push(1)
	q.Push(2)
	q2 := NewQueue[int]()
	testFiles(t, q, q2, func() { checkEqual(t, q2.data, []int{1, 2}) })
	v, ok = q.Front()
	checkEqual(t, ok, true)
	checkEqual(t, v, 1)
	checkEqual(t, q.Pop(), true)
	checkEqual(t, q.Pop(), true)
	checkEqual(t, q.Pop(), false)
	checkEqual(t, q.Empty(), true)
	_, ok = q.Front()
	checkEqual(t, ok, false)
}
