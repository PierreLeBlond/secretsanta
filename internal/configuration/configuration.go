package configuration

import (
  "github.com/PierreLeBlond/secretsanta/pkg/scc/graph"
  "github.com/PierreLeBlond/secretsanta/internal/io"
  "github.com/PierreLeBlond/secretsanta/internal/user"
)

func min(x, y int) int {
  if (x < y) {
    return x;
  }
  return y;
}

func GetGraphFromConfiguration(users []*user.User, configuration []int) (*graph.Graph) {
  newUsers := []*user.User{};
  for i, u := range users {
    newUser := &user.User{u.Name, nil, nil};
    newUser.Always = u.Always[:];
    newUser.Previous = u.Previous[len(u.Previous) - configuration[i]:];
    newUsers = append(newUsers, newUser);
  }
  return io.UsersToGraph(newUsers);
}

func GetConfiguration(users []*user.User) ([]int) {
  configuration := []int{};
  for _, currenUser := range users {
    configuration = append(configuration, len(currenUser.Previous));
  }
  return configuration;
}

func getSubConfigurations(configuration []int, max int) ([][]int) {
  configurations := make([][]int, 0);

  if (len(configuration) == 0 || configuration[0] < max) {
    return append(configurations, configuration);
  }

  subConfigurations := getSubConfigurations(configuration[1:len(configuration)], max);

  maxRootConfiguration := []int{min(max, configuration[0])};
  minRootConfiguration := []int{min(max-1, configuration[0])};
  for _, subConfiguration := range subConfigurations {
    configurations = append(configurations, append(maxRootConfiguration, subConfiguration...));
    if (max > 0 && maxRootConfiguration[0] != minRootConfiguration[0]) {
      configurations = append(configurations, append(minRootConfiguration, subConfiguration...));
    }
  }

  return configurations;
}

func GetConfigurations(configuration []int) ([][]int) {
  configurations := make([][]int, 0);

  if (len(configuration) == 0) {
    return configurations;
  }

  for i := configuration[0]; i > -1; i-- {
    configurations = append(configurations, getSubConfigurations(configuration, i)...);
  }

  return configurations;
}

