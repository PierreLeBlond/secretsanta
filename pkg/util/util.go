package util

import (
  "container/list"
)

func Min(x, y int) int {
  if x < y {
    return x
  }
  return y
}

func Max(x, y int) int {
  if x > y {
    return x
  }
  return y
}

func InList(a *list.List, x string) bool {
  for e := a.Front(); e != nil; e = e.Next() {
    if x == e.Value.(string) {
      return true;
    }
  }
  return false;
}

func InSlice(a []string, x string) bool {
  for _, y := range a {
    if x == y {
      return true;
    }
  }
  return false;
}

func RemoveFromOrderedSlice(a []string, x string) []string {
  i := -1;
  for j, y := range a {
    if x == y {
      i = j;
    }
  }
  return append(a[:i], a[i+1:]...);
}

