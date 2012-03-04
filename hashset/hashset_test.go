package hashset

import (
  "testing"
)

type Key struct {
  x, y int
}

func (k Key) HashCode() int {
  return k.x * k.y
}

func (k Key) Equals(other interface{}) bool {
  switch other.(type) {
  case Key:
    o := other.(Key)
    return k.x == o.x && k.y == o.y
  }
  return false
}

func TestNew(t *testing.T) {
  s := New()
  if s.Size() != 0 || !s.IsEmpty() {
    t.Errorf("init size not 0 but %d", s.Size())
  }
}

func TestPut(t *testing.T) {
  s := New()
  k := Key{1, 2}
  s.Add(k)
  k2 := &Key{1, 2}
  if !s.Contains(k2) {
    t.Errorf("put failed")
  }
}

func TestRemove(t *testing.T) {
  s := New()
  k := Key{1, 2}
  s.Add(k)
  k2 := &Key{1, 2}
  s.Remove(k2)
  if s.Contains(k) {
    t.Errorf("remove failed")
  }
}
