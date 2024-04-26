package peg

type Tree struct {
	Start int
	End   int
	Child []*Tree
	Tags  map[string]struct{}
	Index int
}

func NewTree(pos int) *Tree {
	t := new(Tree)
	t.Start = pos
	t.End = pos
	t.Tags = make(map[string]struct{})
	return t
}

func (t *Tree) Append(c *Tree) {
	if c != nil && c.End > t.End {
		t.End = c.End
	}
	t.Child = append(t.Child, c)
}

func (t *Tree) SetTag(name string) {
	t.Tags[name] = struct{}{}
}
