// port java hashmap to Go
// Author: dccmx<dccmx@dccmx.com>

// Package hashset provides a general hash based map for any type that implements
// Hashable, it is ported from jdk
package hashset

import "github.com/dccmx/hashcontainer/hashmap"

type Hashable hashmap.Hashable

type HashSet struct {
  hashmap *hashmap.HashMap
}

// New creates hashset with default settings.
func New() *HashSet {
  s := new(HashSet)
  s.hashmap = hashmap.New()
  return s
}

// Returns the number of elements in this set (its cardinality).
func (s HashSet) Size() int {
  return s.hashmap.Size()
}

// Returns true if this set contains no elements.
func (s HashSet) IsEmpty() bool {
  return s.hashmap.IsEmpty()
}

// Returns true if this set contains the specified element.
func (s HashSet) Contains(e Hashable) bool {
  return s.hashmap.ContainsKey(e)
}

// Adds the specified element to this set if it is not already present.
// If this set already contains the element, the call leaves the set
// unchanged and returns false.
func (s HashSet) Add(e Hashable) bool {
  already, _ := s.hashmap.Put(e, true)
  return !already
}

// Removes the specified element from this set if it is present.
// if this set contains such an element.  Returns true if
// this set contained the element (or equivalently, if this set
// changed as a result of the call).
// This set will not contain the element once the call returns.
func (s HashSet) Remove(e Hashable) bool {
  already, _ := s.hashmap.Remove(e)
  return already
}

// Removes all of the elements from this set.
// The set will be empty after this call returns.
func (s HashSet) Clear() {
  s.hashmap.Clear()
}
