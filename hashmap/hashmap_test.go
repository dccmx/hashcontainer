// Author: dccmx<dccmx@dccmx.com>

package hashmap

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
  h := New()
  if h.Size() != 0 || !h.IsEmpty() {
    t.Errorf("init size not 0 but %d", h.Size())
  }
}

func TestPutGet0(t *testing.T) {
  h := New()
  k := Key{1, 2}
  h.Put(k, 1)
  if _, v := h.Get(k); v != 1 {
    t.Errorf("put 1 get %d", v)
  }

  k2 := &Key{1, 2}
  if _, v := h.Get(k2); v != 1 {
    t.Errorf("put 1 get %d", v)
  }
}

func TestRemove(t *testing.T) {
  h := New()
  k := Key{1, 2}
  h.Put(k, 1)
  k2 := &Key{1, 2}
  h.Remove(k2)
  if _, v := h.Get(k); v != nil {
    t.Errorf("remove failed")
  }
}
