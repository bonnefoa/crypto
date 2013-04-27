package list

import "testing"
import "testing/assert"

func TestSublist(t *testing.T) {
       lst := []byte{1,2,3,4}
       expected := [][]byte{ {1, 2}, {3, 4} }
       sublists := list.Sublist(lst, 2)
       for i := range expected {
               assert.AssertBytesEquals(t, expected[i], sublists[i])
       }
}

func TestSublistInequal(t *testing.T) {
       lst := []byte{1,2,3,4,5}
       expected := [][]byte{ {1, 2}, {3, 4}, {5} }
       sublists := list.Sublist(lst, 2)
       for i := range expected {
               assert.AssertBytesEquals(t, expected[i], sublists[i])
       }
}

