package listener

import (
	"testing"
)

func TestTree_Add(t *testing.T) {
	tree := New()
	tree.Add("a.a", "A1")
	tree.Add("a.a", "A2")
	tree.Add("a.b", "B1")
	tree.Add("a.*", "*1")
	tree.Add("a.a.a", "A3")
	tree.Add("a.a.a.b", "A4")
	tree.Add("a.a.b", "A5")
	tree.Add("a.*", "*2")
	tree.Add("a.b.*", "*3")
	tree.Add("a.b.c", "C1")
	tree.Add("*", "overall")
	t.Logf("\n%v", tree)
}

func TestTree_Add2(t *testing.T) {
	tree := New()
	tree.Add("a.a", "A1")
	tree.Add("a.a", "A2")
	tree.Add("a.*", "*1")
	tree.Add("a.a.a", "A3")
	tree.Add("a.a.a.b", "A4")
	tree.Add("a.a.b", "A5")
	tree.Add("a.*", "*2")
	tree.Add("a.b", "B1")
	tree.Add("a.*.b", "OnlyA5")
	t.Logf("\n%v", tree)
}
