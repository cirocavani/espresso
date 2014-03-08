package graph

import (
	"container/list"
	"fmt"
	"strings"
)

type GraphType string

const (
	DIRECTED   GraphType = "DIRECTED"
	UNDIRECTED           = "UNDIRECTED"
)

type data struct {
	values map[string]interface{}
}

func (d *data) string(sep string) string {
	if d.values == nil {
		return ""
	}
	outs := make([]string, len(d.values))
	i := 0
	for k, v := range d.values {
		outs[i] = fmt.Sprintf("%s:%#v", k, v)
		i++
	}
	return strings.Join(outs, sep)
}

func (d *data) String() string {
	return string("\n")
}

func (d *data) DataSize() int {
	return len(d.values)
}

func (d *data) Set(key string, value interface{}) *data {
	if d.values == nil {
		d.values = make(map[string]interface{})
	}
	d.values[key] = value
	return d
}

func (d *data) SetMap(values map[string]interface{}) *data {
	for k, v := range values {
		d.Set(k, v)
	}
	return d
}

func (d *data) Get(key string) (interface{}, bool) {
	if d.values == nil {
		return nil, false
	}
	v, ok := d.values[key]
	return v, ok
}

func (d *data) Unset(key string) {
	if d.values == nil {
		return
	}
	delete(d.values, key)
}

type Vertex struct {
	id, label string
	graph     *Graph
	edges     *list.List
	data
}

type Edge struct {
	label string
	graph *Graph
	link  map[string]*Vertex
	data
}

type Graph struct {
	_type    GraphType
	vertices map[string]*Vertex
	edges    int
	data
}

func New() *Graph {
	return NewDirected()
}

func NewDirected() *Graph {
	return &Graph{_type: DIRECTED, edges: 0}
}

func NewUndirected() *Graph {
	return &Graph{_type: UNDIRECTED, edges: 0}
}

func (g *Graph) Type() GraphType {
	return g._type
}

func (g *Graph) VertexCount() int {
	return len(g.vertices)
}

func (g *Graph) EdgeCount() int {
	return g.edges
}

func (v *Vertex) String() string {
	out := v.id
	if v.label != "" {
		out += ":" + v.label
	}
	if data := v.data.string(","); data != "" {
		out += " {" + data + "}"
	}
	return out
}

func (e *Edge) string() string {
	out := ""
	if e.label != "" {
		out += ":" + e.label + " "
	}
	if data := e.data.string(","); data != "" {
		out += "{" + data + "}"
	}
	return out
}

func (e *Edge) String() string {
	out := e.string()
	if out == "" {
		return out
	}
	return "[" + out + "]"
}

func (g *Graph) string(v1, v2 *Vertex, e *Edge) string {
	switch g.Type() {
	case DIRECTED:
		return fmt.Sprintf("(%s)-%s->(%s)\n", v1, e, v2)
	case UNDIRECTED:
		return fmt.Sprintf("(%s)-%s-(%s)\n", v1, e, v2)
	default:
		return "<unknown>"
	}
}

func (g *Graph) String() string {
	out := ""

	if data := g.data.string("\n"); data != "" {
		out += data + "\n"
	}

	edges := make(map[*Edge]bool)
	for _, v := range g.vertices {
		if v.EdgeCount() == 0 {
			continue
		}
		for i := v.edges.Front(); i != nil; i = i.Next() {
			e := i.Value.(*Edge)

			if edges[e] {
				continue
			}

			adj, ok := e.link[v.id]
			if !ok {
				continue
			}

			out += g.string(v, adj, e)
			edges[e] = true
		}
	}

	return out
}

func (g *Graph) HasVertex(id string) bool {
	if g.vertices == nil {
		return false
	}
	_, ok := g.vertices[id]
	return ok
}

func (g *Graph) getVertex(id string) (*Vertex, bool) {
	if g.vertices == nil {
		return nil, false
	}
	v, ok := g.vertices[id]
	return v, ok
}

func (g *Graph) addVertex(v *Vertex) {
	if g.vertices == nil {
		g.vertices = make(map[string]*Vertex)
	}
	g.vertices[v.id] = v
}

func (g *Graph) Vertex(id string) *Vertex {
	v, ok := g.getVertex(id)
	if ok {
		return v
	}
	v = &Vertex{
		id:    id,
		graph: g,
	}
	g.addVertex(v)
	return v
}

func (v *Vertex) Label(label string) *Vertex {
	v.label = label
	return v
}

func (v *Vertex) Id() string {
	return v.id
}

func (v *Vertex) EdgeCount() int {
	if v.edges == nil {
		return 0
	}
	return v.edges.Len()
}

func (v *Vertex) bind(e *Edge) {
	if e == nil {
		return
	}
	if v.edges == nil {
		v.edges = list.New()
	}
	v.edges.PushBack(e)
}

func (v *Vertex) unbind(e *Edge) {
	if e == nil || v.edges == nil {
		return
	}
	var wrap *list.Element = nil
	for i := v.edges.Front(); i != nil; i = i.Next() {
		if i.Value.(*Edge) == e {
			wrap = i
			break
		}
	}
	if wrap != nil {
		v.edges.Remove(wrap)
	}
}

func (g *Graph) edge(v1, v2 *Vertex) *Edge {
	switch g.Type() {
	case DIRECTED:
		return &Edge{
			graph: g,
			link:  map[string]*Vertex{v1.id: v2},
		}
	case UNDIRECTED:
		return &Edge{
			graph: g,
			link: map[string]*Vertex{
				v1.id: v2,
				v2.id: v1,
			},
		}
	default:
		return nil
	}
}

func (g *Graph) Edge(id1, id2 string) *Edge {
	v1 := g.Vertex(id1)
	v2 := g.Vertex(id2)

	e := g.edge(v1, v2)

	v1.bind(e)
	if v2 != v1 {
		v2.bind(e)
	}

	g.edges++
	return e
}

func (e *Edge) Label(label string) *Edge {
	e.label = label
	return e
}

func (g *Graph) Edges(id1, id2 string) []*Edge {
	v1, ok1 := g.getVertex(id1)
	v2, ok2 := g.getVertex(id2)
	if !ok1 || !ok2 {
		return nil
	}
	min := v1.EdgeCount()
	if m := v2.EdgeCount(); min > m {
		min = m
	}
	if min == 0 {
		return nil
	}
	edges := make([]*Edge, 0, min)

	for i := v1.edges.Front(); i != nil; i = i.Next() {
		e := i.Value.(*Edge)
		if adj, ok := e.link[v1.id]; ok && adj == v2 {
			edges = append(edges, e)
		}
	}

	if len(edges) == 0 {
		return nil
	}
	return edges
}

func (v *Vertex) Remove() {
	if v.edges != nil {
		edges := make([]*Edge, v.edges.Len())
		i := 0
		for e := v.edges.Front(); e != nil; e = e.Next() {
			edges[i] = e.Value.(*Edge)
			i++
		}
		v.edges = nil
		for _, e := range edges {
			e.Remove()
		}
	}

	delete(v.graph.vertices, v.id)
	v.graph = nil
	v.data.values = nil
}

func (e *Edge) Remove() {
	e.graph.edges--
	for k, v := range e.link {
		vk := e.graph.Vertex(k)
		vk.unbind(e)
		if v != vk {
			v.unbind(e)
		}
	}
	e.link = nil
	e.graph = nil
	e.data.values = nil
}
