
package scc

import (
  "testing"
  "container/list"
  "github.com/PierreLeBlond/secretsanta/pkg/scc/graph"
)

func compareString(t *testing.T, expectedString string, givenString string) {
  if expectedString != givenString {
    t.Log("error should be " + expectedString + ", but got", givenString);
    t.Fail();
  }
}

func compareGraph(t *testing.T, expectedGraph *graph.Graph, givenGraph *graph.Graph) {
  expectedString := GraphToString(expectedGraph);
  givenString := GraphToString(givenGraph);
  compareString(t, expectedString, givenString);
}

func TestComputeStrongConnectedComponents(t *testing.T) {
  mainGraph := &graph.Graph{list.New(), make(map[string][]string), make(map[string]bool), make(map[string][]string)};

  mainGraph.Nodes.PushBack("a");
  mainGraph.Nodes.PushBack("b");
  mainGraph.Nodes.PushBack("c");
  mainGraph.Nodes.PushBack("d");
  mainGraph.Nodes.PushBack("e");
  mainGraph.Nodes.PushBack("f");

  mainGraph.Edges["a"] = append(mainGraph.Edges["a"], "b");
  mainGraph.Edges["a"] = append(mainGraph.Edges["a"], "c");
  mainGraph.Edges["b"] = append(mainGraph.Edges["b"], "a");
  mainGraph.Edges["b"] = append(mainGraph.Edges["b"], "d");
  mainGraph.Edges["c"] = append(mainGraph.Edges["c"], "a");
  mainGraph.Edges["d"] = append(mainGraph.Edges["d"], "e");
  mainGraph.Edges["e"] = append(mainGraph.Edges["e"], "f");
  mainGraph.Edges["f"] = append(mainGraph.Edges["f"], "d");

  rootNode := "a";
  counter := 0;
  visited := make(map[string]bool);
  lowLink := make(map[string]int);
  number := make(map[string]int);
  stack := list.New();

  _, graphs := computeStrongConnectedComponents(mainGraph, rootNode, counter, &lowLink, &number, &visited, stack);
  expected := "d:e/e:f/f:d";
  str := GraphToString(graphs[0]);
  compareString(t, expected, str);
}

func TestGetLowestStrongConnectedComponent(t *testing.T) {
  mainGraph := &graph.Graph{list.New(), make(map[string][]string), make(map[string]bool), make(map[string][]string)};

  mainGraph.Nodes.PushBack("a");
  mainGraph.Nodes.PushBack("b");
  mainGraph.Nodes.PushBack("c");
  mainGraph.Nodes.PushBack("d");
  mainGraph.Nodes.PushBack("e");
  mainGraph.Nodes.PushBack("f");

  mainGraph.Edges["a"] = append(mainGraph.Edges["a"], "b");
  mainGraph.Edges["a"] = append(mainGraph.Edges["a"], "c");
  mainGraph.Edges["b"] = append(mainGraph.Edges["b"], "d");
  mainGraph.Edges["d"] = append(mainGraph.Edges["d"], "e");
  mainGraph.Edges["e"] = append(mainGraph.Edges["e"], "f");
  mainGraph.Edges["f"] = append(mainGraph.Edges["f"], "d");

  component := GetLowestStrongConnectedComponent(mainGraph, mainGraph.Nodes.Front());

  expected := "d:e/e:f/f:d";
  str := GraphToString(component);
  compareString(t, expected, str);
}

func TestGetSubgraph(t *testing.T) {
  mainGraph := &graph.Graph{list.New(), make(map[string][]string), make(map[string]bool), make(map[string][]string)};

  mainGraph.Nodes.PushBack("a");
  mainGraph.Nodes.PushBack("b");
  mainGraph.Nodes.PushBack("c");
  mainGraph.Nodes.PushBack("d");
  mainGraph.Nodes.PushBack("e");
  mainGraph.Nodes.PushBack("f");

  mainGraph.Edges["a"] = append(mainGraph.Edges["a"], "b");
  mainGraph.Edges["b"] = append(mainGraph.Edges["b"], "c");
  mainGraph.Edges["c"] = append(mainGraph.Edges["c"], "a");
  mainGraph.Edges["c"] = append(mainGraph.Edges["c"], "b");
  mainGraph.Edges["e"] = append(mainGraph.Edges["e"], "d");
  mainGraph.Edges["d"] = append(mainGraph.Edges["d"], "f");
  mainGraph.Edges["f"] = append(mainGraph.Edges["f"], "e");


  subGraph := getSubgraph(mainGraph, mainGraph.Nodes.Front());

  compareGraph(t, mainGraph, subGraph);
}
