package partition

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

func comparePartitions(t *testing.T, expectedPartitions [][]*list.List, givenPartitions [][]*list.List) {
  if (len(expectedPartitions) != len(givenPartitions)) {
    t.Log("expected partitions length doesn't match given one");
    t.Fail();
  }
  for i, expectedPartition := range expectedPartitions {
    comparePaths(t, expectedPartition, givenPartitions[i]);
  }
}

func TestFindPartitions(t *testing.T) {
  mainGraph := &graph.Graph{list.New(), make(map[string][]string), make(map[string]bool), make(map[string][]string)};

  mainGraph.Nodes.PushBack("a");
  mainGraph.Nodes.PushBack("b");
  mainGraph.Nodes.PushBack("c");
  mainGraph.Nodes.PushBack("d");
  mainGraph.Nodes.PushBack("e");
  mainGraph.Nodes.PushBack("f");

  var circuits []*list.List;

  firstExpectedCircuit := list.New();
  firstExpectedCircuit.PushBack("a");
  firstExpectedCircuit.PushBack("b");
  firstExpectedCircuit.PushBack("c");
  circuits = append(circuits, firstExpectedCircuit);

  secondExpectedCircuit := list.New();
  secondExpectedCircuit.PushBack("d");
  secondExpectedCircuit.PushBack("e");
  secondExpectedCircuit.PushBack("f");
  circuits = append(circuits, secondExpectedCircuit);

  thirdCircuit := list.New();
  thirdCircuit.PushBack("d");
  thirdCircuit.PushBack("f");
  circuits = append(circuits, thirdCircuit);

  fourthCircuit := list.New();
  fourthCircuit.PushBack("c");
  fourthCircuit.PushBack("d");
  fourthCircuit.PushBack("e");
  fourthCircuit.PushBack("f");
  circuits = append(circuits, fourthCircuit);

  partitions := FindPartitions(mainGraph, circuits);

  expectedPartition := []*list.List{firstExpectedCircuit, secondExpectedCircuit};

  expectedPartitions := [][]*list.List{expectedPartition};

  comparePartitions(t, expectedPartitions, partitions);
}

