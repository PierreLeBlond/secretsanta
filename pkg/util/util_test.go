package util

import (
  "testing"
)

func TestRemoveFromOrderedSlice(t *testing.T) {
  input := []string{"a", "b", "c", "d", "e"};
  expectedOutput := []string{"a", "b", "d", "e"};

  output := RemoveFromOrderedSlice(input, "c");

  if len(expectedOutput) != len(output) {
    t.Log("expected length doesn't match given one");
    t.Fail();
  }

  for i, x := range expectedOutput {
    if x != output[i] {
      t.Log("error at position " + string(i) + " ,should be " + x + ", but got", output[i]);
      t.Fail();
    }
  }
}

