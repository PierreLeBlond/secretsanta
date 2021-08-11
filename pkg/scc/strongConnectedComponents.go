package scc

import (
  "strings"
  "container/list"
  "github.com/PierreLeBlond/secretsanta/pkg/scc/graph"
  "github.com/PierreLeBlond/secretsanta/pkg/util"
)

func GraphToString(graph *graph.Graph) (string) {
  s := "";
  for e := graph.Nodes.Front(); e != nil; e = e.Next() {
    node := e.Value.(string);
    s += node;
    s += ":";
    for _, otherNode := range graph.Edges[node] {
      s += otherNode;
    }
    if (e.Next() != nil) {
      s += "/";
    }
  }
  return s;
}

func StringToGraph(s string) (*graph.Graph) {
  mainGraph := &graph.Graph{list.New(), make(map[string][]string), make(map[string]bool), make(map[string][]string)};

  nodeStrings := strings.Split(s, "/");
  for _, nodeString := range nodeStrings {
    nodeStringSplit := strings.Split(nodeString, ":");
    node := nodeStringSplit[0];
    otherNodeStrings := strings.Split(nodeStringSplit[1], "");
    mainGraph.Nodes.PushBack(node);
    for _, otherNode := range otherNodeStrings {
      mainGraph.Edges[node] = append(mainGraph.Edges[node], otherNode);
    }
  }

  return mainGraph;
}

func getSubgraph(mainGraph *graph.Graph, nodeElement *list.Element) *graph.Graph {
  subGraph := &graph.Graph{list.New(), make(map[string][]string), make(map[string]bool), make(map[string][]string)};

  for e := nodeElement; e != nil; e = e.Next() {
    node := e.Value.(string)
    subGraph.Nodes.PushBack(node);
  }
  for e := subGraph.Nodes.Front(); e != nil; e = e.Next() {
    node := e.Value.(string);
    for f := subGraph.Nodes.Front(); f != nil; f = f.Next() {
      otherNode := f.Value.(string);
      if (node != otherNode && util.InSlice(mainGraph.Edges[node], otherNode)) {
        subGraph.Edges[node] = append(subGraph.Edges[node], otherNode)
      }
    }
  }

  return subGraph;
}

func computeStrongConnectedComponents(
  mainGraph *graph.Graph, rootNode string, counter int,
  lowLink *map[string]int, number *map[string]int, visited *map[string]bool, stack *list.List,
) (int, []*graph.Graph) {
  var components []*graph.Graph;
  var newComponents []*graph.Graph;

  (*lowLink)[rootNode] = counter;
  (*number)[rootNode] = counter;
  (*visited)[rootNode] = true;

  stack.PushBack(rootNode);

  for _, node := range mainGraph.Edges[rootNode] {
    if !(*visited)[node] {
      counter, newComponents = computeStrongConnectedComponents(
        mainGraph, node, counter + 1, lowLink, number, visited, stack,
      );
      components = append(
        components,
        newComponents...
      );
      (*lowLink)[rootNode] = util.Min((*lowLink)[rootNode], (*lowLink)[node]);
    } else if (*number)[node] < (*number)[rootNode] {
      if util.InList(stack, node) {
        (*lowLink)[rootNode] = util.Min((*lowLink)[rootNode], (*number)[node]);
      }
    }
  }

  if (*lowLink)[rootNode] == (*number)[rootNode] && stack.Len() > 0 {
    component := &graph.Graph{
      list.New(), make(map[string][]string), make(map[string]bool), make(map[string][]string),
    };

    nextElement := stack.Back();
    stack.Remove(nextElement);
    nextNode := nextElement.Value.(string);
    component.Nodes.PushFront(nextNode);
    for (*number)[nextNode] > (*number)[rootNode] {
      nextElement = stack.Back();
      stack.Remove(nextElement);
      nextNode = nextElement.Value.(string);
      component.Nodes.PushFront(nextNode);
    }

    for e := component.Nodes.Front(); e != nil; e = e.Next() {
      node := e.Value.(string);
      for f := component.Nodes.Front(); f != nil; f = f.Next() {
        otherNode := f.Value.(string);
        if (node != otherNode && util.InSlice(mainGraph.Edges[node], otherNode)) {
          component.Edges[node] = append(component.Edges[node], otherNode)
        }
      }
    }

    if (component.Nodes.Len() > 1) {
      components = append(components, component);
    }
  }
  return counter, components;
}

func getLowestComponent(mainGraph *graph.Graph, components []*graph.Graph) (*graph.Graph) {
  for e := mainGraph.Nodes.Front(); e != nil; e = e.Next() {
    node := e.Value.(string);
    for _, component := range components {
      if (util.InList(component.Nodes, node)) {
        return component;
      }
    }
  }
  return nil;
}

func GetLowestStrongConnectedComponent(mainGraph *graph.Graph, nodeElement *list.Element) (*graph.Graph) {
  visited := make(map[string]bool);
  lowLink := make(map[string]int);
  number := make(map[string]int);

  stack := list.New();

  subGraph := getSubgraph(mainGraph, nodeElement);

  var components []*graph.Graph;
  var newComponents []*graph.Graph;

  node := nodeElement.Value.(string);
  for e := nodeElement; e != nil; e = e.Next() {
    otherNode := e.Value.(string);
    if !visited[otherNode] {
      _, newComponents = computeStrongConnectedComponents(
        subGraph, otherNode, 0, &lowLink, &number, &visited, stack,
      );
      components = append(components, newComponents...);
      component := getLowestComponent(mainGraph, components);
      if (component != nil && !util.InList(component.Nodes, node) && !util.InList(component.Nodes, nodeElement.Next().Value.(string))) {
        return GetLowestStrongConnectedComponent(mainGraph, nodeElement.Next());
      } else {
        if component != nil {
          return component;
        }
      }
    }
  }
  return nil;
}
