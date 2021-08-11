package circuit

import (
  "testing"
  "container/list"
  "github.com/PierreLeBlond/secretsanta/pkg/scc/graph"
)

func pathToString(path *list.List) (string) {
  s := "";
  for e := path.Front(); e != nil; e = e.Next() {
    s += e.Value.(string);
  }
  return s;
}

func compareString(t *testing.T, expectedString string, givenString string) {
  if expectedString != givenString {
    t.Log("error should be " + expectedString + ", but got", givenString);
    t.Fail();
  }
}

func comparePaths(t *testing.T, expectedPaths []*list.List, givenPaths []*list.List) {
  if (len(expectedPaths) != len(givenPaths)) {
    t.Log("expected paths length doesn't match given one");
    t.Fail();
  }
  for i, expectedPath := range expectedPaths {
    compareString(t, pathToString(expectedPath), pathToString(givenPaths[i]));
  }
}

func TestGetCircuits(t *testing.T) {
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

  mainGraph.Edges["a"] = append(mainGraph.Edges["a"], "d");
  mainGraph.Edges["d"] = append(mainGraph.Edges["d"], "e");
  mainGraph.Edges["e"] = append(mainGraph.Edges["e"], "f");
  mainGraph.Edges["f"] = append(mainGraph.Edges["f"], "a");

  node := mainGraph.Nodes.Front().Value.(string);
  _, circuits := getCircuits(mainGraph, list.New(), node, node);

  firstExpectedCircuit := list.New();
  firstExpectedCircuit.PushBack("a");
  firstExpectedCircuit.PushBack("b");
  firstExpectedCircuit.PushBack("c");

  secondExpectedCircuit := list.New();
  secondExpectedCircuit.PushBack("a");
  secondExpectedCircuit.PushBack("d");
  secondExpectedCircuit.PushBack("e");
  secondExpectedCircuit.PushBack("f");

  expectedCircuits := []*list.List{firstExpectedCircuit, secondExpectedCircuit};

  comparePaths(t, expectedCircuits, circuits);
}

