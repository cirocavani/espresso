package graph

import (
	"reflect"
	"testing"
)

func TestEmptyGraph(t *testing.T) {
	g := New()
	if gtype := g.Type(); gtype != DIRECTED {
		t.Errorf("Default graph should be directed graph: %s", gtype)
	}
	if out := g.String(); out != "" {
		t.Errorf("Empty graph shoud output nothing: %s", out)
	}
}

func testGraph(t *testing.T, g *Graph, gt GraphType, ids ...string) {
	if gtype := g.Type(); gtype != gt {
		t.Errorf("Graph should be %s graph: %s", gt, gtype)
	}

	if n := g.VertexCount(); n != 0 {
		t.Errorf("Error graph should start empty: %d", n)
	}

	for _, id := range ids {
		g.Vertex(id)
	}

	if n := g.VertexCount(); n != len(ids) {
		t.Errorf("Error graph setting vertices (%d): %d", len(ids), n)
	}
}

func testSameVertexEdge(t *testing.T, g *Graph, id string, testEdge func(id1, id2 string, e *Edge)) {
	e := g.Edge(id, id)

	if n := g.EdgeCount(); n != 4 {
		t.Errorf("Error graph setting same vertex edge (vertex %s, edges=4): %d", id, n)
	}

	testEdge(id, id, e)

	e.Remove()

	if n := g.VertexCount(); n != 3 {
		t.Errorf("Error graph removing same vertex edge (vertex %s, vertices=3): %d", id, n)
	}
	if n := g.EdgeCount(); n != 3 {
		t.Errorf("Error graph removing same vertex edge (vertex %s, edges=3): %d", id, n)
	}
}

func TestDirectedGraph(t *testing.T) {
	g := NewDirected()

	const id1, id2, id3 = "1", "2", "3"

	testGraph(t, g, DIRECTED, id1, id2, id3)

	testEdge := func(from, to string, edge *Edge) {
		if e := g.Edges(from, to); len(e) != 1 || e[0] != edge {
			t.Errorf("Error graph setting edge (%s)->(%s): %d", from, to, len(e))
		}
	}

	testNoEdge := func(from, to string) {
		if e := g.Edges(from, to); len(e) != 0 {
			t.Errorf("Error graph wrong edge (%s)->(%s): %d", from, to, len(e))
		}
	}

	e12 := g.Edge(id1, id2)
	e23 := g.Edge(id2, id3)
	e31 := g.Edge(id3, id1)

	if n := g.VertexCount(); n != 3 {
		t.Errorf("Error graph setting edges (vertices=3): %d", n)
	}
	if n := g.EdgeCount(); n != 3 {
		t.Errorf("Error graph setting edges (edges=3): %d", n)
	}

	testEdge(id1, id2, e12)
	testNoEdge(id1, id1)
	testNoEdge(id1, id3)
	testEdge(id2, id3, e23)
	testNoEdge(id2, id2)
	testNoEdge(id2, id1)
	testEdge(id3, id1, e31)
	testNoEdge(id3, id3)
	testNoEdge(id3, id2)

	t.Log(g)

	// More tests

	testSameVertexEdge(t, g, id1, testEdge)
	testSameVertexEdge(t, g, id2, testEdge)
	testSameVertexEdge(t, g, id3, testEdge)

	v2 := g.Vertex(id2)
	v2.Remove()

	if n := g.VertexCount(); n != 2 {
		t.Errorf("Error graph removing vertex (vertices=2): %d", n)
	}
	if n := g.EdgeCount(); n != 1 {
		t.Errorf("Error graph removing vertex (edges=1): %d", n)
	}
	if e := g.Edges(id3, id1); len(e) != 1 || e[0] != e31 {
		t.Errorf("Error graph removing vertex, (%s)->(%s): %d", id3, id1, len(e))
	}
}

func TestUndirectedGraph(t *testing.T) {
	g := NewUndirected()

	const id1, id2, id3 = "1", "2", "3"

	testGraph(t, g, UNDIRECTED, id1, id2, id3)

	testEdge := func(v1, v2 string, edge *Edge) {
		if e := g.Edges(v1, v2); len(e) != 1 || e[0] != edge {
			t.Errorf("Error graph setting edge (%s)-(%s): %d", v1, v2, len(e))
		}
	}

	testNoEdge := func(v1, v2 string) {
		if e := g.Edges(v1, v2); len(e) != 0 {
			t.Errorf("Error graph wrong edge (%s)->(%s): %d", v1, v2, len(e))
		}
	}

	e12 := g.Edge(id1, id2)
	e23 := g.Edge(id2, id3)
	e31 := g.Edge(id3, id1)

	if n := g.VertexCount(); n != 3 {
		t.Errorf("Error graph setting edges (vertices=3): %d", n)
	}
	if n := g.EdgeCount(); n != 3 {
		t.Errorf("Error graph setting edges (edges=3): %d", n)
	}

	testEdge(id1, id2, e12)
	testEdge(id1, id3, e31)
	testNoEdge(id1, id1)
	testEdge(id2, id3, e23)
	testEdge(id2, id1, e12)
	testNoEdge(id2, id2)
	testEdge(id3, id1, e31)
	testEdge(id3, id2, e23)
	testNoEdge(id3, id3)

	t.Log(g)

	// More tests

	testSameVertexEdge(t, g, id1, testEdge)
	testSameVertexEdge(t, g, id2, testEdge)
	testSameVertexEdge(t, g, id3, testEdge)

	v2 := g.Vertex(id2)
	v2.Remove()

	if n := g.VertexCount(); n != 2 {
		t.Errorf("Error graph removing vertex (vertices=2): %d", n)
	}
	if n := g.EdgeCount(); n != 1 {
		t.Errorf("Error graph removing vertex (edges=1): %d", n)
	}
	if e := g.Edges(id3, id1); len(e) != 1 || e[0] != e31 {
		t.Errorf("Error graph removing vertex, (%s)-(%s): %d", id3, id1, len(e))
	}
	if e := g.Edges(id1, id3); len(e) != 1 || e[0] != e31 {
		t.Errorf("Error graph removing vertex, (%s)-(%s): %d", id1, id3, len(e))
	}
}

func testData(t *testing.T, d *data) {
	const k1, k2, k3, k4 = "string", "int", "float", "map"
	v1, v2, v3, v4 := "text", 100, 1.618, map[string]string{"a": "b"}

	test := func(k string, equal func(v interface{}) bool) {
		if v, ok := d.Get(k); !ok || !equal(v) {
			t.Errorf("Error setting value for key '%s': %#v, (%t)", k, v, ok)
		}
		d.Unset(k)
		if v, ok := d.Get(k); ok {
			t.Errorf("Error unsetting value for key '%s': %#v, (%t)", k, v, ok)
		}
	}

	if d.DataSize() != 0 {
		t.Errorf("Error data should start empty: %v", d)
	}

	d.SetMap(map[string]interface{}{
		k1: v1,
		k2: v2,
		k3: v3,
		k4: v4,
	})

	if d.DataSize() != 4 {
		t.Errorf("Error setting data: %v", d)
	}

	test(k1, func(v interface{}) bool { return v == v1 })
	test(k2, func(v interface{}) bool { return v == v2 })
	test(k3, func(v interface{}) bool { return v == v3 })
	test(k4, func(v interface{}) bool { return reflect.DeepEqual(v, v4) })

	if d.DataSize() != 0 {
		t.Errorf("Error data should end empty: %v", d)
	}
}

func TestEdges(t *testing.T) {
	test := func(g *Graph) {
		if n := g.Edges("1", "3"); n != nil {
			t.Errorf("(%s) Error no edge from 1 to 3: %+v", g.Type(), n)
		}
		if n := g.Edges("3", "1"); n != nil {
			t.Errorf("(%s) Error no edge from 3 to 1: %+v", g.Type(), n)
		}
		if n := g.Edges("1", "4"); n != nil {
			t.Errorf("(%s) Error no edge from 1 to 4 (alone): %+v", g.Type(), n)
		}
		if n := g.Edges("4", "1"); n != nil {
			t.Errorf("(%s) Error no edge from 4 (alone) to 1: %+v", g.Type(), n)
		}
		if n := g.Edges("1", "5"); n != nil {
			t.Errorf("(%s) Error no edge from 1 to 5 (missing): %+v", g.Type(), n)
		}
		if g.HasVertex("5") {
			t.Errorf("(%s) Error vertex 5 created (edge in)!", g.Type())
		}
		if n := g.Edges("5", "1"); n != nil {
			t.Errorf("(%s) Error no edge from 5 (missing) to 1: %+v", g.Type(), n)
		}
		if g.HasVertex("5") {
			t.Errorf("(%s) Error vertex 5 created (edge out)!", g.Type())
		}
		if n := g.Edges("5", "6"); n != nil {
			t.Errorf("(%s) Error no edge from 5 to 6 (both missing): %+v", g.Type(), n)
		}
		if g.HasVertex("5") || g.HasVertex("6") {
			t.Errorf("(%s) Error vertex 5 or 6 created!", g.Type())
		}
	}

	{
		g := NewDirected()
		e := g.Edge("1", "2")
		g.Edge("3", "2")
		g.Vertex("4")

		if n := g.Edges("1", "2"); len(n) != 1 || n[0] != e {
			t.Errorf("Error edges size (1): %+v", n)
		}
		if n := g.Edges("2", "1"); n != nil {
			t.Errorf("Error no edge from 2 to 1: %+v", n)
		}

		test(g)
	}
	{
		g := NewUndirected()
		e := g.Edge("1", "2")
		g.Edge("3", "2")
		g.Vertex("4")

		if n := g.Edges("1", "2"); len(n) != 1 || n[0] != e {
			t.Errorf("Error edges size (1): %+v", n)
		}
		if n := g.Edges("2", "1"); len(n) != 1 || n[0] != e {
			t.Errorf("Error edges size (1): %+v", n)
		}

		test(g)
	}
}

func TestData(t *testing.T) {
	g := New()
	testData(t, &g.data)
	v := g.Vertex("1")
	testData(t, &v.data)
	e := g.Edge("1", "2")
	testData(t, &e.data)
}
