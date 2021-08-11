package graph

import "container/list"

type Graph struct {
  Nodes *list.List;
  Edges map[string][]string;
  BlockedNodes map[string]bool;
  BlockedEdges map[string][]string;
}
