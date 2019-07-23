package main

import (
    "fmt"
    "os"
    "bufio"
    "io/ioutil"
    "strings"
    "math/rand"
    "time"
)

var users map[string]*User
var circuits [][]int
var circuit []int

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// Contains tells whether a contains x.
func containsUser(a []string, x string) bool {
    for _, n := range a {
        if x == n {
            return true
        }
    }
    return false
}

func areDisjoinct(a []int, b[]int) bool {
    for _, i := range a {
        for _, j := range b {
            if i == j {
               return false
           }
        }
    }
    return true
}

func unblock(graph Graph, u int) {
    graph[u].blocked = false;
    for _, w := range graph[u].blockedUsers {
        if graph[w].blocked {
            unblock(graph, w);
        }
    }
    graph[u].blockedUsers = nil;
}

func findCircuit(graph Graph, s int, v int) bool {
    f := false;
    circuit = append(circuit, v);
    graph[v].blocked = true;
    for _, w := range graph[v].childs {
        if w == s {
            c := make([]int, len(circuit));
            copy(c, circuit);
            circuits = append(circuits, c);
            f = true;
        } else if !graph[w].blocked {
            f = findCircuit(graph, s, w);
        }
    }
    if f {
        unblock(graph, v);
    } else {
        for _, w := range graph[v].childs {
            if !contains(graph[w].blockedUsers, w) {
                graph[w].blockedUsers = append(graph[w].blockedUsers, w);
            }
        }
    }
    circuit = circuit[0:len(circuit)-1];
    return f;
}

func findPartition(ids []int, partition [][]int, index int, length int, max int) ([]int, [][]int) {
    for i := index+1; i < len(circuits); i++ {
        if len(circuits[i]) > 2 && (length + len(circuits[i]) == max || length + len(circuits[i]) < max - 1) {
            f := true;
            for _, circuit := range partition {
                f = f && areDisjoinct(circuit, circuits[i]);
            }
            if f {
                ids = append(ids, i);
                partition = append(partition, circuits[i]);
                if length + len(circuits[i]) == max {
                    return ids, partition;
                } else {
                    nextIds, nextPartition := findPartition(ids, partition, i, length + len(circuits[i]), max);
                    if nextPartition != nil {
                        return nextIds, nextPartition;
                    }
                }
            }
        }
    }
    return nil, nil;
}

func findPartitions(graph Graph) ([][]int, [][][]int) {
    s := 0;

    for s < len(graph) {
        adjacencyList, lowestId := getAdjacencyList(graph, s);
        if adjacencyList != nil {
            s = lowestId;
            for i, component := range adjacencyList {
                if component != nil {
                    graph[i].blocked = false;
                    graph[i].blockedUsers = nil;
                }
            }
            findCircuit(graph, s, s);
            s = s + 1;
        } else {
            s = len(graph);
        }
    }

    var partitions [][][]int;
    var ids [][]int;

    for i, _ := range circuits {
        id, partition := findPartition(nil, nil, i, 0, len(graph));
        if partition != nil {
            ids = append(ids, id);
            partitions = append(partitions, partition);
        }
    }

    return ids, partitions;
}

func main() {
    users := make(map[string]*User);
    var graph []*Node;

    if (len(os.Args) < 2) {
        fmt.Println("Error: Please provide ouput file name.");
        return;
    }

    fmt.Printf("Computing secret santa's list into %s...\n", os.Args[1]);

    // Retrieve users
    var usersNames []string;
    f, err := os.Open("../data/user.txt");
    check(err);
    defer f.Close();

    id := 0;
    usersScanner := bufio.NewScanner(f);
    for usersScanner.Scan() {
        var s = usersScanner.Text();
        users[s] = &User{id, nil, nil};
        usersNames = append(usersNames, s);

        id++;
    }

    for i, u := range usersNames {
        g, err := ioutil.ReadFile("../data/" + u + ".txt");
        check(err);

        t := string(g);

        parts := strings.SplitN(t, "\n\n", 2);

        always := strings.Split(parts[0], "\n");
        for _, name := range always {
            if name != "" {
                users[u].always = append(users[u].always, users[name].id);
            }
        }

        previous := strings.Split(parts[1], "\n");
        for _, name := range previous {
            if name != "" {
                users[u].previous = append(users[u].previous, users[name].id);
            }
        }

        graph = append(graph, &Node{u, nil, false, nil});
        for j, v := range usersNames {
            if !contains(users[u].always, j) && !contains(users[u].previous, j) && u != v {
                graph[i].childs = append(graph[i].childs, j);
            }
        }
    }

    var partitions [][][]int;
    var ids [][]int;

    for partitions == nil {
        ids, partitions = findPartitions(graph);
        if partitions == nil {
            previous := 0;
            previousCount := 0;
            for i, node := range graph {
                node.blocked = false;
                node.blockedUsers = nil;
                if len(users[node.name].previous) > previousCount {
                    previous = i;
                    previousCount = len(users[node.name].previous);
                }
            }
            if users[graph[previous].name].previous == nil {
                panic("There is no solution !");
            }
            fmt.Printf("Removing %s from %s\n...", graph[users[graph[previous].name].previous[0]].name, graph[previous].name);
            graph[previous].childs = append(graph[previous].childs, users[graph[previous].name].previous[0]);
            users[graph[previous].name].previous = users[graph[previous].name].previous[1:len(users[graph[previous].name].previous)];
        }
    }

    for i, partition := range partitions {
        fmt.Print("\n");
        for j, circuit := range partition {
            fmt.Printf("%d ", ids[i][j]);
            fmt.Print(circuit);
            fmt.Print("\n");
        }
        fmt.Print("\n");
    }

    rand.Seed(time.Now().UnixNano());

    partition := partitions[rand.Intn(len(partitions))];

    h, err := os.Create("../output/" + os.Args[1] + ".txt")
    check(err)

    defer h.Close()

    for _, circuit := range partition {
        for j, node := range circuit {
            k := j + 1;
            if j == len(circuit) - 1 {
                k = 0;
            }
            next := circuit[k];
            _, err := h.WriteString(graph[node].name + " to " + graph[next].name + "\n");
            check(err);
            users[graph[node].name].previous = append(users[graph[node].name].previous, next);
        }
    }

    for _, node := range graph {
        g, err := os.Create("../data/" + node.name + ".txt")
        check(err)

        for _, always := range users[node.name].always {
            _, err := g.WriteString(graph[always].name + "\n");
            check(err);
        }
        if users[node.name].always == nil {
            _, err = g.WriteString("\n");
            check(err);
        }
        _, err = g.WriteString("\n");
        check(err);
        for _, previous := range users[node.name].previous {
            _, err = g.WriteString(graph[previous].name + "\n");
            check(err);
        }

        defer g.Close()
    }
}

