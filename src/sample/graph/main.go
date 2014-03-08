package main

import (
	"espresso/graph"
	"fmt"
)

func main() {
	{
		fmt.Println("Empty Graph")

		g := graph.New()
		g.Set("name", "empty")
		g.Set("type", "nothing")

		fmt.Print(g)

		fmt.Println()
	}
	{
		fmt.Println("Undirected Graph")

		g := graph.NewUndirected()
		g.Set("name", "Pit")
		g.Set("type", "right")
		g.Edge("1", "2").Set("name", "a").Set("size", 3)
		g.Edge("2", "3").Set("name", "b").Set("size", 4)
		g.Edge("3", "1").Set("name", "c").Set("size", 5)

		fmt.Print(g)

		fmt.Println()
	}
	{
		fmt.Println("Directed Graph")

		g := graph.NewDirected()
		g.Set("name", "Me, MySelf And You")
		g.Set("type", "Direct")

		g.Edge("Me", "MySelf").Set("x", 1)
		g.Edge("MySelf", "You").Set("x", 2)
		g.Edge("You", "Me").Set("x", 3)

		g.Vertex("Me").Set("x", 4)
		g.Vertex("MySelf").Set("x", 5)
		g.Vertex("You").Set("x", 6)

		fmt.Print(g)

		fmt.Println()
	}
	{
		fmt.Println("Removing Vertex")

		g := graph.New()
		g.Edge("1", "2")
		g.Edge("2", "3")
		g.Edge("3", "1")

		fmt.Println("all")
		fmt.Print(g)

		g.Vertex("2").Remove()

		fmt.Println("2 removed")
		fmt.Print(g)

		fmt.Println()
	}
	{
		fmt.Println("Removing Vertex 2")

		g := graph.New()
		g.Edge("1", "2")
		g.Edge("2", "3")
		g.Edge("3", "1")
		g.Edge("2", "3")
		g.Edge("2", "2")

		fmt.Println("all")
		fmt.Print(g)

		g.Vertex("2").Remove()

		fmt.Println("2 removed")
		fmt.Print(g)

		fmt.Println()
	}
	{
		fmt.Println("Neo4j tutorial graph")

		g := graph.New()

		v0 := g.Vertex("0").Label("Movie")
		v0.SetMap(map[string]interface{}{
			"id":    "603",
			"title": "The Matrix",
			"year":  "1999-03-31",
		})
		v1 := g.Vertex("1").Label("Movie")
		v1.SetMap(map[string]interface{}{
			"id":    "604",
			"title": "The Matrix Reloaded",
			"year":  "2003-05-07",
		})
		v2 := g.Vertex("2").Label("Movie")
		v2.SetMap(map[string]interface{}{
			"id":    "605",
			"title": "The Matrix Revolutions",
			"year":  "2003-10-27",
		})
		v3 := g.Vertex("3").Label("Actor")
		v3.Set("name", "Keanu Reeves")
		v4 := g.Vertex("4").Label("Actor")
		v4.Set("name", "Laurence Fishburne")
		v5 := g.Vertex("5").Label("Actor")
		v5.Set("name", "Carrie-Anne Moss")

		e30 := g.Edge(v3.Id(), v0.Id()).Label("ACTS_IN")
		e30.Set("role", "Neo")
		e31 := g.Edge(v3.Id(), v1.Id()).Label("ACTS_IN")
		e31.Set("role", "Neo")
		e32 := g.Edge(v3.Id(), v2.Id()).Label("ACTS_IN")
		e32.Set("role", "Neo")

		e40 := g.Edge(v4.Id(), v0.Id()).Label("ACTS_IN")
		e40.Set("role", "Morpheus")
		e41 := g.Edge(v4.Id(), v1.Id()).Label("ACTS_IN")
		e41.Set("role", "Morpheus")
		e42 := g.Edge(v4.Id(), v2.Id()).Label("ACTS_IN")
		e42.Set("role", "Morpheus")

		e50 := g.Edge(v5.Id(), v0.Id()).Label("ACTS_IN")
		e50.Set("role", "Trinity")
		e51 := g.Edge(v5.Id(), v1.Id()).Label("ACTS_IN")
		e51.Set("role", "Trinity")
		e52 := g.Edge(v5.Id(), v2.Id()).Label("ACTS_IN")
		e52.Set("role", "Trinity")

		fmt.Println(g)
	}

}
