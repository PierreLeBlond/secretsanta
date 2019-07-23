package main

var counter = 0;

var visited []bool;
var stack []int;
var lowLink []int;
var number []int;

var sccs [][]int;

func min(x, y int) int {
    if x < y {
        return x
    }
    return y
}

func max(x, y int) int {
    if x > y {
        return x
    }
    return y
}

func contains(a []int, x int) bool {
    for _, n := range a {
        if x == n {
            return true
        }
    }
    return false
}

func getSubgraph(graph Graph, node int) Graph {
    subGraph := make(Graph, len(graph));

    for i := node; i < len(graph); i++ {
        var childs []int;
        for _, child := range graph[i].childs {
            if child >= node {
                childs = append(childs, child);
            }
        }
        subGraph[i] = &Node{graph[i].name, childs, false, nil}
    }

    return subGraph;
}

func computeStrongConnectedComponent(graph Graph, root int) {
    counter++;
    lowLink[root] = counter;
    number[root] = counter;
    visited[root] = true;

    stack = append(stack, root);

    for _, child := range graph[root].childs {
        if !visited[child] {
            computeStrongConnectedComponent(graph, child);
            lowLink[root] = min(lowLink[root], lowLink[child]);
        } else if number[child] < number[root] {
            if contains(stack, child) {
                lowLink[root] = min(lowLink[root], number[child]);
            }
        }
    }

    if lowLink[root] == number[root] && len(stack) > 0 {
        next := -1;
        var component []int;

        next = stack[len(stack) - 1];
        stack = stack[0:len(stack) - 1];
        component = append(component, next);
        for number[next] > number[root] {
            next = stack[len(stack) - 1];
            stack = stack[0:len(stack) - 1];
            component = append(component, next);
        }

        if (len(component) > 1) {
            sccs = append(sccs, component);
        }
    }
}

func getLowestIdComponent(graph Graph) []int {
    min := len(graph);
    var component []int;

    for _, scc := range sccs {
        for _, node := range scc {
            if (node < min) {
                component = scc;
                min = node;
            }
        }
    }

    return component;
}

func getAdjacencyListCore(graph Graph, nodes []int) [][]int {
    var lowestIdAdjacencyList [][]int;

    if nodes != nil {
        lowestIdAdjacencyList = make([][]int, len(graph));
        for _, node := range nodes {
            for _, child := range graph[node].childs {
                if contains(nodes, child) {
                    lowestIdAdjacencyList[node] = append(lowestIdAdjacencyList[node], child);
                }
            }
        }

    }

    return lowestIdAdjacencyList;
}

func getAdjacencyList(graph Graph, node int) ([][]int, int) {
    visited = make([]bool, len(graph));
    lowLink = make([]int, len(graph));
    number = make([]int, len(graph));

    stack = nil;
    sccs = nil;

    subGraph := getSubgraph(graph, node);

    for i := node; i < len(graph); i++ {
        if !visited[i] {
            computeStrongConnectedComponent(subGraph, i);
            nodes := getLowestIdComponent(graph);
            if nodes != nil && !contains(nodes, node) && !contains(nodes, node + 1) {
                return getAdjacencyList(graph, node + 1);
            } else {
                adjacencyList := getAdjacencyListCore(subGraph, nodes);
                if adjacencyList != nil {
                    for j, _ := range graph {
                        if len(adjacencyList[j]) > 0 {
                            return adjacencyList, j;
                        }
                    }
                }
            }
        }
    }
    return nil, -1
}
