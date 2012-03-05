// port java hashmap to Go
// Author: dccmx<dccmx@dccmx.com>

// Package hashmap provides a general hash based map for any type that implements
// Hashable, it is ported from jdk
package hashmap

const (
  // The default initial capacity - MUST be a power of two.
  defaultInitialCapacity = 16

  // The maximum capacity, used if a higher value is implicitly specified
  // by either of the constructors with arguments.
  // MUST be a power of two <= 1<<30.
  maximumCapacity = 1 << 30

  // The load factor used when none specified in constructor.
  defaultLoadFactor = 0.75
)

type Hashable interface {
  // HashCode is the hash function that compute the hash code of your type
  HashCode() int

  // Equals defines the == operator of key
  Equals(other interface{}) bool
}

type entry struct {
  // hashCode cached the hash code of key
  hashCode int

  // key
  key Hashable

  // the mapped value of key
  value interface{}

  // next entry in entry list
  next *entry
}

type HashMap struct {
  // The table, resized as necessary. Length MUST Always be a power of two.
  table []*entry

  // The number of key-value mappings contained in this map.
  size int

  // The next size value at which to resize (capacity * load factor).
  threshold int

  // The load factor for the hash table.
  loadFactor float64
}

// New creates hashmap with default settings.
func New() *HashMap {
  h := new(HashMap)
  h.loadFactor = defaultLoadFactor
  h.threshold = int(defaultInitialCapacity * defaultLoadFactor)
  h.table = make([]*entry, defaultInitialCapacity)
  return h
}

// Applies a supplemental hash function to a given hashCode, which
// defends against poor quality hash functions.  This is critical
// because HashMap uses power-of-two length hash tables, that
// otherwise encounter collisions for hashCodes that do not differ
// in lower bits. Note: Null keys always map to hash 0, thus index 0.
//
func hash(h int) int {
  // This function ensures that hashCodes that differ only by
  // constant multiples at each bit position have a bounded
  // number of collisions (approximately 8 at default load factor).
  h ^= (h >> 20) ^ (h >> 12);
  return h ^ (h >> 7) ^ (h >> 4);
}

// Returns index for hash code h.
func indexFor(h, length int) int {
  return h & (length - 1)
}

// Returns the number of key-value mappings in this map.
func (h HashMap) Size() int {
  return h.size
}

// Returns true if this map contains no key-value mappings.
func (h HashMap) IsEmpty() bool {
  return h.size == 0
}

// Returns the find result & value to which the specified key is mapped,
// or nil if this map contains no mapping for the key.
//
// More formally, if this map contains a mapping from a key
// to a value such that (key==null ? k==null :
// key.equals(k)), then this method returns v; otherwise
// it returns nil.  (There can be at most one such mapping.)
func (h HashMap) Get(key Hashable) (interface{}, bool) {
  if key == nil {
    return h.getForNullKey()
  }
  hash := hash(key.HashCode())
  for e := h.table[indexFor(hash, len(h.table))]; e != nil; e = e.next {
    if e.hashCode == hash && (e.key == key || key.Equals(e.key)) {
      return e.value, true
    }
  }
  return nil, false
}

// Offloaded version of get() to look up nil keys.  Null keys map
// to index 0.  This nil case is split out into separate methods
// for the sake of performance in the two most commonly used
// operations (get and put), but incorporated with conditionals in
// otherst.
func (h HashMap) getForNullKey() (interface{}, bool) {
  for e := h.table[0]; e != nil; e = e.next {
    if e.key == nil {
      return e.value, true
    }
  }
  return nil, false
}

// Returns true if this map contains a mapping for the
// specified key.
func (h HashMap) ContainsKey(key Hashable) bool {
  return h.getEntry(key) != nil
}

// Returns the entry associated with the specified key in the
// HashMap.  Returns nil if the HashMap contains no mapping
// for the key.
func (h HashMap) getEntry(key Hashable) *entry {
  hashCode := 0
  if key != nil {
   hashCode = hash(key.HashCode());
  }
  for e := h.table[indexFor(hashCode, len(h.table))]; e != nil; e = e.next {
    if e.hashCode == hashCode && (e.key == key || (key != nil && key.Equals(e.key))) {
      return e
    }
  }
  return nil
}

// Associates the specified value with the specified key in this map.
// If the map previously contained a mapping for the key, the old
// value is replaced.
// return a bool indicate if there already has a mapping for the key
// and the previous value associated with key, or nil if there was no mapping for key.
func (h HashMap) Put(key Hashable, value interface{}) (interface{}, bool) {
  if key == nil {
    return h.putForNullKey(value)
  }
  hashCode := hash(key.HashCode())
  i := indexFor(hashCode, len(h.table))
  for e := h.table[i]; e != nil; e = e.next {
    if e.hashCode == hashCode && (e.key == key || key.Equals(e.key)) {
      oldValue := e.value
      e.value = value
      return oldValue, true
    }
  }

  h.addentry(hashCode, key, value, i)
  return nil, false
}

// Offloaded version of put for nil keys
func (h HashMap) putForNullKey(value interface{}) (interface{}, bool) {
  for e := h.table[0]; e != nil; e = e.next {
    if e.key == nil {
      oldValue := e.value
      e.value = value
      return oldValue, true
    }
  }
  h.addentry(0, nil, value, 0)
  return nil, false
}

// Adds a new entry with the specified key, value and hash code to
// the specified bucket.  It is the responsibility of this
// method to resize the table if appropriate.
func (h HashMap) addentry(hashCode int, key Hashable, value interface{}, bucketIndex int) {
  e := h.table[bucketIndex]
  h.table[bucketIndex] = &entry{hashCode, key, value, e}
  h.size++
  if h.size >= h.threshold {
    h.resize(2 * len(h.table))
  }
}

// Rehashes the contents of this map into a new array with a
// larger capacity.  This method is called automatically when the
// number of keys in this map reaches its threshold.
//
// If current capacity is MAXIMUM_CAPACITY, this method does not
// resize the map, but sets threshold to Integer.MAX_VALUE.
// This has the effect of preventing future calls.
//
// newCapacity is the new capacity, MUST be a power of two;
// must be greater than current capacity unless current
// capacity is MAXIMUM_CAPACITY (in which case value
// is irrelevant).
func (h HashMap) resize(newCapacity int) {
  oldTable := h.table
  oldCapacity := len(oldTable)
  if oldCapacity == maximumCapacity {
    h.threshold = maximumCapacity
    return
  }

  newTable := make([]*entry, newCapacity)
  h.transfer(newTable)
  h.table = newTable
  h.threshold = int(float64(newCapacity) * h.loadFactor)
}

// Transfers all entries from current table to newTable.
func (h HashMap) transfer(newTable []*entry) {
  src := h.table
  newCapacity := len(newTable)
  for j := 0; j < len(src); j++ {
    e := src[j]
    if e != nil {
      src[j] = nil
      for e != nil{
        next := e.next
        i := indexFor(e.hashCode, newCapacity)
        e.next = newTable[i]
        e, newTable[i] = e, next
      }
    }
  }
}

// Removes the mapping for the specified key from this map if present.
//
// return if there was a mapping for key and 
// the previous value associated with key, or
func (h HashMap) Remove(key Hashable) (interface{}, bool) {
  e, found := h.removeEntryForKey(key)
  if found {
    return e.value, true
  }
  return nil, false
}

// Removes and returns the entry associated with the specified key
// in the HashMap.
//
// Returns if there was a mapping for key and the associated entry
func (h HashMap) removeEntryForKey(key Hashable) (*entry, bool) {
  hashCode :=  0
  if key != nil {
    hashCode = hash(key.HashCode())
  }
  i := indexFor(hashCode, len(h.table))
  prev := h.table[i]
  e := prev

  for e != nil {
    next := e.next
    if e.hashCode == hashCode && (e.key == key || (key != nil && key.Equals(e.key))) {
      h.size--
      if prev == e {
        h.table[i] = next
      } else {
        prev.next = next
      }
      return e, true
    }
    prev = e
    e = next
  }

  return nil, false
}

// Removes all of the mappings from this map.
// The map will be empty after this call returns.
func (h HashMap) Clear() {
  tab := h.table
  for i := 0; i < len(tab); i++ {
    tab[i] = nil
  }
  h.size = 0
}
