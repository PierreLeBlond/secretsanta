package main

import (
  "fmt"
  "sort"
  "os"
  "container/list"
  "github.com/PierreLeBlond/secretsanta/internal/io"
  "github.com/PierreLeBlond/secretsanta/internal/configuration"
  "github.com/PierreLeBlond/secretsanta/pkg/scc/circuit"
  "github.com/PierreLeBlond/secretsanta/pkg/scc/partition"
  "github.com/PierreLeBlond/secretsanta/pkg/scc/graph"
)

func main() {
  if (len(os.Args) < 3) {
    fmt.Println("Usage: ./secretsanta <path to input file> <path to output file>");
    return;
  }

  fmt.Printf("Computing secret santa's list into %s...\n", os.Args[2]);

  // 1. Import users and sort them by their previous participation
  users := io.ImportFromJSON(os.Args[1]);
  sort.Slice(users, func(i, j int) bool {
    return len(users[i].Previous) > len(users[j].Previous);
  });

  // 2. Build configuration
  initialConfiguration := configuration.GetConfiguration(users);
  configurations := configuration.GetConfigurations(initialConfiguration);

  // 3. Search for partition
  var partitions [][]*list.List;
  var mainGraph *graph.Graph;
  i := 0;
  for (partitions == nil && i < len(configurations)) {
    fmt.Printf("testing configuration %d on %d...\n", i+1, len(configurations));

    // Get a graph from current configuration
    mainGraph = configuration.GetGraphFromConfiguration(users, configurations[i]);

    // Get all circuit of any length within the graph
    circuits := circuit.GetCircuits(mainGraph);

    // Remove unwanted circuit which will results in cycle of length < 4
    var filteredCircuits []*list.List;
    for _, circuit := range circuits {
      if circuit.Len() > 2 && (circuit.Len() == mainGraph.Nodes.Len() || circuit.Len() < mainGraph.Nodes.Len() - 2) {
        filteredCircuits = append(filteredCircuits, circuit);
      }
    }

    // Get all partitions we can build from the filtered circuits
    partitions = partition.FindPartitions(mainGraph, filteredCircuits);
    i++;
  }

  if partitions == nil {
    panic("There is no solution !");
  }

  fmt.Printf("Found %d partition !", len(partitions));

  var partition []*list.List;

  // 4. Keep the partition with least number of circuits
  for _, p := range partitions {
    if (partition == nil || len(partition) > len(p)) {
      partition = p;
    }
  }

  // 5. Update users from new matches
  newUsers := io.GetUpdatedUsersFromPartition(users, partition);

  // 6. Export results
  io.ExportToJSON(os.Args[2], newUsers);
}

