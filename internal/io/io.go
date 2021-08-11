package io

import (
  // "fmt"
  "encoding/json"
  "os"
  "io/ioutil"
  "container/list"
  "github.com/PierreLeBlond/secretsanta/pkg/util"
  "github.com/PierreLeBlond/secretsanta/pkg/scc/graph"
  "github.com/PierreLeBlond/secretsanta/internal/user"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

func ImportFromJSON(path string) ([]*user.User) {
  data, err := ioutil.ReadFile(path);
  check(err);

  users := []*user.User{};

  err = json.Unmarshal([]byte(data), &users);
  check(err);

  return users;
}

func UsersToGraph(users []*user.User) (*graph.Graph) {
  graph := &graph.Graph{list.New(), make(map[string][]string), make(map[string]bool), make(map[string][]string)};

  for _, mainUser := range users {
    graph.Nodes.PushBack(mainUser.Name);
    for _, otherUser := range users {
      if mainUser != otherUser &&
      !util.InSlice(mainUser.Always, otherUser.Name) &&
      !util.InSlice(mainUser.Previous, otherUser.Name) {
        graph.Edges[mainUser.Name] = append(graph.Edges[mainUser.Name], otherUser.Name)
      }
    }
  }

  return graph;
}

func GetUpdatedUsersFromPartition(oldUsers []*user.User, partition []*list.List) ([]*user.User) {
  updatedUsers := []*user.User{};
  for _, circuit := range partition {
    for e := circuit.Front(); e != nil; e = e.Next() {
      node := e.Value.(string);
      oldUser := user.GetUserByName(oldUsers, node);
      updatedUser := &user.User{oldUser.Name, oldUser.Always, oldUser.Previous};

      f := e.Next();
      if (f == nil) {
        f = circuit.Front();
      }

      otherNode := f.Value.(string);

      if util.InSlice(updatedUser.Previous, otherNode) {
        updatedUser.Previous = util.RemoveFromOrderedSlice(updatedUser.Previous, otherNode);
      }

      updatedUser.Previous = append(updatedUser.Previous, otherNode);

      updatedUsers = append(updatedUsers, updatedUser);
    }
  }
  return updatedUsers;
}

func GraphToUsers(oldUsers []*user.User, graph *graph.Graph) ([]*user.User) {
  newUsers := []*user.User{};
  for _, oldUser := range oldUsers {
    newUser := &user.User{oldUser.Name, oldUser.Always, oldUser.Previous};
    for _, otherUser := range oldUsers {
      if otherUser.Name != newUser.Name &&
      !util.InSlice(newUser.Always, otherUser.Name) &&
      !util.InSlice(graph.Edges[newUser.Name], otherUser.Name) {
        newUser.Previous = append(newUser.Previous, otherUser.Name);
      }
    }
    newUsers = append(newUsers, newUser);
  }
  return newUsers;
}

func ExportToJSON(path string, users []*user.User) {
  serializedUsers, _ := json.MarshalIndent(users, "", "\t");
  f, err := os.Create(path)
  check(err);
  _, err = f.WriteString(string(serializedUsers))
  check(err);
}
