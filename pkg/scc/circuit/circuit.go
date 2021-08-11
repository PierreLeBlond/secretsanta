package circuit

import (
  "container/list"
  "github.com/PierreLeBlond/secretsanta/pkg/util"
  "github.com/PierreLeBlond/secretsanta/pkg/scc"
  "github.com/PierreLeBlond/secretsanta/pkg/scc/graph"
)

func findIndex(s []string, v string) int {
  for i, w := range s {
    if v == w {
      return i;
    }
  }
  return -1;
}

func removeFromSlice(s []string, v string) []string {
  i := findIndex(s, v);
  s[i] = s[len(s)-1];
  s[len(s)-1] = "";
  return s[:len(s)-1];
}

func unblock(graph *graph.Graph, node string) {
  graph.BlockedNodes[node] = false;
  var blockedEdges []string;
  copy(blockedEdges, graph.BlockedEdges[node]);
  for _, otherNode := range blockedEdges {
    graph.BlockedEdges[node] = removeFromSlice(graph.BlockedEdges[node], otherNode);
    if graph.BlockedNodes[otherNode] {
      unblock(graph, otherNode);
    }
  }
}

func getNodeElement(graph *graph.Graph, node string) (*list.Element) {
  for e := graph.Nodes.Front(); e != nil; e = e.Next() {
    otherNode := e.Value.(string);
    if (otherNode == node) {
      return e;
    }
  }
  return nil;
}

func hasDoublon(circuit *list.List) bool {
  for e := circuit.Front(); e != nil; e = e.Next() {
    x := e.Value.(string);
    for f := e.Next(); f != nil; f = f.Next() {
      y := f.Value.(string);
      if (x == y) {
        return true;
      }
    }
  }
  return false;
}

func getCircuits(component *graph.Graph, elementaryPath *list.List, startNode string, currentNode string, debug string) (bool, []*list.List) {
  var circuits []*list.List;

  f := false;

  newElementaryPath := list.New();
  newElementaryPath.PushBackList(elementaryPath);
  newElementaryPath.PushBack(currentNode);

  component.BlockedNodes[currentNode] = true;
  debug += " -> block " + currentNode;
  for _, otherNode := range component.Edges[currentNode] {
    if otherNode == startNode {
      circuits = append(circuits, newElementaryPath);
      f = true;
    } else if !component.BlockedNodes[otherNode] {
      success, newCircuits := getCircuits(component, newElementaryPath, startNode, otherNode, debug);
      f = success;
      circuits = append(circuits, newCircuits...);
    }
  }
  if f {
    unblock(component, currentNode);
  } else {
    for _, otherNode := range component.Edges[currentNode] {
      if !util.InSlice(component.BlockedEdges[otherNode], currentNode) {
        component.BlockedEdges[otherNode] = append(component.BlockedEdges[otherNode], currentNode);
      }
    }
  }
  return f, circuits;
}

func GetCircuits(graph *graph.Graph) ([]*list.List) {
  var circuits []*list.List;

  nodeElement := graph.Nodes.Front();

  for nodeElement != nil {
    component := scc.GetLowestStrongConnectedComponent(graph, nodeElement);
    if component != nil {
      nodeElement = getNodeElement(graph, component.Nodes.Front().Value.(string));
      node := nodeElement.Value.(string);

      _, newCircuits := getCircuits(component, list.New(), node, node, "");
      circuits = append(circuits, newCircuits...);

      nodeElement = nodeElement.Next();
    } else {
      nodeElement = nil;
    }
  }

  return circuits;
}

