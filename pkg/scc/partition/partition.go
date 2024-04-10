package partition

import (
  "container/list"
  "math/rand"
  "time"
  "github.com/PierreLeBlond/secretsanta/pkg/scc/graph"
)

// Tells whether a & b does not have any common element
func areDisjoinct(a *list.List, b *list.List) bool {
  for e := a.Front(); e != nil; e = e.Next() {
    for f := b.Front(); f != nil; f = f.Next() {
      if e.Value == f.Value {
        return false
      }
    }
  }
  return true
}

func isDisjoinctFromPartition(circuit *list.List, partition []*list.List) bool {
  for _, otherCircuit := range partition {
    if (!areDisjoinct(circuit, otherCircuit)) {
      return false;
    }
  }
  return true;
}

func findPartitions(partition []*list.List, circuits []*list.List, length int, max int) ([][]*list.List) {
  if (len(circuits) == 0) {
    return nil;
  }
  var partitions [][]*list.List;
  newCircuit := circuits[0];
  newLength := length + newCircuit.Len();
  if newLength <= max && isDisjoinctFromPartition(newCircuit, partition) {
    if newLength == max {
      newPartition := append(partition, newCircuit);
      // This is needed to prevent bugs I do not understand yet
      newPartition = append([]*list.List(nil), newPartition...);
      partitions = append(partitions, newPartition);
    } else {
      newPartition := append(partition, newCircuit);
      newPartitions := findPartitions(newPartition, circuits[1:], newLength, max);
      partitions = append(partitions, newPartitions...);
    }
  }
  newPartitions := findPartitions(partition[:], circuits[1:], length, max);
  partitions = append(partitions, newPartitions...);
  return partitions;
}

func FindPartitions(graph *graph.Graph, circuits []*list.List) ([][]*list.List) {
  return findPartitions(nil, circuits, 0, graph.Nodes.Len());
}

func FindFinalPartition(graph *graph.Graph, circuits []*list.List) ([]*list.List) {
  var partitions [][]*list.List;
  var partition []*list.List;

  partitions = FindPartitions(graph, circuits);

  if partitions != nil {
    rand.Seed(time.Now().UnixNano());
    partition = partitions[rand.Intn(len(partitions))];
  }

  return partition;
}

