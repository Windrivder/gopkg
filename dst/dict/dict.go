package dict

import (
	"encoding/json"

	"github.com/windrivder/gopkg/cast"
	"github.com/windrivder/gopkg/syncx"
	"github.com/windrivder/gopkg/typex"
)

type Dict struct {
	mu   *syncx.RWMutex
	data typex.DictType
}

// NewDict returns an empty Dict object.
// The parameter <safe> is used to specify whether using map in concurrent-safety,
// which is false in default.
func NewDict(safe ...bool) *Dict {
	return &Dict{
		mu:   syncx.CreateRWMutex(safe...),
		data: make(typex.DictType),
	}
}

// NewDictFrom creates and returns a hash map from given map <data>.
// Note that, the param <data> map will be set as the underlying data map(no deep copy),
// there might be some concurrent-safe issues when changing the map outside.
func NewDictFrom(data typex.DictType, safe ...bool) *Dict {
	return &Dict{
		mu:   syncx.CreateRWMutex(safe...),
		data: data,
	}
}

// Iterator iterates the hash map readonly with custom callback function <f>.
// If <f> returns true, then it continues iterating; or false to stop.
func (m *Dict) Iterator(f func(k string, v typex.GenericType) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.data {
		if !f(k, v) {
			break
		}
	}
}

// Clone returns a new hash map with copy of current map data.
func (m *Dict) Clone() *Dict {
	return NewDictFrom(m.MapCopy(), m.mu.IsSafe())
}

// Map returns the underlying data map.
// Note that, if it's in concurrent-safe usage, it returns a copy of underlying data,
// or else a pointer to the underlying data.
func (m *Dict) Map() typex.DictType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.mu.IsSafe() {
		return m.data
	}
	data := make(typex.DictType, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// MapStrAny returns a copy of the underlying data of the map as map[string]interface{}.
func (m *Dict) MapStrAny() typex.DictType {
	return m.Map()
}

// MapCopy returns a copy of the underlying data of the hash map.
func (m *Dict) MapCopy() typex.DictType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data := make(typex.DictType, len(m.data))
	for k, v := range m.data {
		data[k] = v
	}
	return data
}

// FilterEmpty deletes all key-value pair of which the value is empty.
// Values like: 0, nil, false, "", len(slice/map/chan) == 0 are considered empty.
func (m *Dict) FilterEmpty() {
	m.mu.Lock()
	for k, v := range m.data {
		if cast.IsEmpty(v) {
			delete(m.data, k)
		}
	}
	m.mu.Unlock()
}

// FilterNil deletes all key-value pair of which the value is nil.
func (m *Dict) FilterNil() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.data {
		if cast.IsNil(v) {
			delete(m.data, k)
		}
	}
}

// Set sets key-value to the hash map.
func (m *Dict) Set(key string, val typex.GenericType) {
	m.mu.Lock()
	if m.data == nil {
		m.data = make(typex.DictType)
	}
	m.data[key] = val
	m.mu.Unlock()
}

// Sets batch sets key-values to the hash map.
func (m *Dict) Sets(data typex.DictType) {
	m.mu.Lock()
	if m.data == nil {
		m.data = data
	} else {
		for k, v := range data {
			m.data[k] = v
		}
	}
	m.mu.Unlock()
}

// Search searches the map with given <key>.
// Second return parameter <found> is true if key was found, otherwise false.
func (m *Dict) Search(key string) (value typex.GenericType, found bool) {
	m.mu.RLock()
	if m.data != nil {
		value, found = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Get returns the value by given <key>.
func (m *Dict) Get(key string) (value typex.GenericType) {
	m.mu.RLock()
	if m.data != nil {
		value, _ = m.data[key]
	}
	m.mu.RUnlock()
	return
}

// Pop retrieves and deletes an item from the map.
func (m *Dict) Pop() (key string, value typex.GenericType) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, value = range m.data {
		delete(m.data, key)
		return
	}
	return
}

// Pops retrieves and deletes <size> items from the map.
// It returns all items if size == -1.
func (m *Dict) Pops(size int) typex.DictType {
	m.mu.Lock()
	defer m.mu.Unlock()
	if size > len(m.data) || size == -1 {
		size = len(m.data)
	}
	if size == 0 {
		return nil
	}
	var (
		index  = 0
		newMap = make(typex.DictType, size)
	)
	for k, v := range m.data {
		delete(m.data, k)
		newMap[k] = v
		index++
		if index == size {
			break
		}
	}
	return newMap
}

// doSetWithLockCheck checks whether value of the key exists with mutex.Lock,
// if not exists, set value to the map with given <key>,
// or else just return the existing value.
//
// When setting value, if <value> is type of <func() interface {}>,
// it will be executed with mutex.Lock of the hash map,
// and its return value will be set to the map with <key>.
//
// It returns value with given <key>.
func (m *Dict) doSetWithLockCheck(key string, value typex.GenericType) typex.GenericType {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(typex.DictType)
	}
	if v, ok := m.data[key]; ok {
		return v
	}
	if f, ok := value.(func() typex.GenericType); ok {
		value = f()
	}
	if value != nil {
		m.data[key] = value
	}
	return value
}

// GetOrSet returns the value by key,
// or sets value with given <value> if it does not exist and then returns this value.
func (m *Dict) GetOrSet(key string, value typex.GenericType) typex.GenericType {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, value)
	} else {
		return v
	}
}

// GetOrSetFunc returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
func (m *Dict) GetOrSetFunc(key string, f func() typex.GenericType) typex.GenericType {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f())
	} else {
		return v
	}
}

// GetOrSetFuncLock returns the value by key,
// or sets value with returned value of callback function <f> if it does not exist
// and then returns this value.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function <f>
// with mutex.Lock of the hash map.
func (m *Dict) GetOrSetFuncLock(key string, f func() typex.GenericType) typex.GenericType {
	if v, ok := m.Search(key); !ok {
		return m.doSetWithLockCheck(key, f)
	} else {
		return v
	}
}

// SetIfNotExist sets <value> to the map if the <key> does not exist, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *Dict) SetIfNotExist(key string, value typex.GenericType) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, value)
		return true
	}
	return false
}

// SetIfNotExistFunc sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
func (m *Dict) SetIfNotExistFunc(key string, f func() typex.GenericType) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f())
		return true
	}
	return false
}

// SetIfNotExistFuncLock sets value with return value of callback function <f>, and then returns true.
// It returns false if <key> exists, and <value> would be ignored.
//
// SetIfNotExistFuncLock differs with SetIfNotExistFunc function is that
// it executes function <f> with mutex.Lock of the hash map.
func (m *Dict) SetIfNotExistFuncLock(key string, f func() typex.GenericType) bool {
	if !m.Contains(key) {
		m.doSetWithLockCheck(key, f)
		return true
	}
	return false
}

// Removes batch deletes values of the map by keys.
func (m *Dict) Removes(keys ...string) {
	m.mu.Lock()
	if m.data != nil {
		for _, key := range keys {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
}

// Remove deletes value from map by given <key>, and return this deleted value.
func (m *Dict) Remove(key string) (value typex.GenericType) {
	m.mu.Lock()
	if m.data != nil {
		var ok bool
		if value, ok = m.data[key]; ok {
			delete(m.data, key)
		}
	}
	m.mu.Unlock()
	return
}

// Keys returns all keys of the map as a slice.
func (m *Dict) Keys() []string {
	m.mu.RLock()
	var (
		keys  = make([]string, len(m.data))
		index = 0
	)
	for key := range m.data {
		keys[index] = key
		index++
	}
	m.mu.RUnlock()
	return keys
}

// Values returns all values of the map as a slice.
func (m *Dict) Values() []typex.GenericType {
	m.mu.RLock()
	var (
		values = make([]typex.GenericType, len(m.data))
		index  = 0
	)
	for _, value := range m.data {
		values[index] = value
		index++
	}
	m.mu.RUnlock()
	return values
}

// Contains checks whether a key exists.
// It returns true if the <key> exists, or else false.
func (m *Dict) Contains(key string) bool {
	var ok bool
	m.mu.RLock()
	if m.data != nil {
		_, ok = m.data[key]
	}
	m.mu.RUnlock()
	return ok
}

// Size returns the size of the map.
func (m *Dict) Size() int {
	m.mu.RLock()
	length := len(m.data)
	m.mu.RUnlock()
	return length
}

// IsEmpty checks whether the map is empty.
// It returns true if map is empty, or else false.
func (m *Dict) IsEmpty() bool {
	return m.Size() == 0
}

// Clear deletes all data of the map, it will remake a new underlying data map.
func (m *Dict) Clear() {
	m.mu.Lock()
	m.data = make(typex.DictType)
	m.mu.Unlock()
}

// Replace the data of the map with given <data>.
func (m *Dict) Replace(data typex.DictType) {
	m.mu.Lock()
	m.data = data
	m.mu.Unlock()
}

// LockFunc locks writing with given callback function <f> within RWMutex.Lock.
func (m *Dict) LockFunc(f func(m typex.DictType)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(m.data)
}

// RLockFunc locks reading with given callback function <f> within RWMutex.RLock.
func (m *Dict) RLockFunc(f func(m typex.DictType)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f(m.data)
}

// Flip exchanges key-value of the map to value-key.
func (m *Dict) Flip() {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := make(typex.DictType, len(m.data))
	for k, v := range m.data {
		n[cast.ToString(v)] = k
	}
	m.data = n
}

// Merge merges two hash maps.
// The <other> map will be merged into the map <m>.
func (m *Dict) Merge(other *Dict) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = other.MapCopy()
		return
	}
	if other != m {
		other.mu.RLock()
		defer other.mu.RUnlock()
	}
	for k, v := range other.data {
		m.data[k] = v
	}
}

// String returns the map as a string.
func (m *Dict) String() string {
	b, _ := m.MarshalJSON()
	return cast.ByteSliceToString(b)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (m *Dict) MarshalJSON() ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return json.Marshal(m.data)
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (m *Dict) UnmarshalJSON(b []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data == nil {
		m.data = make(typex.DictType)
	}
	if err := json.Unmarshal(b, &m.data); err != nil {
		return err
	}
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for map.
func (m *Dict) UnmarshalValue(value typex.GenericType) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = cast.ToStringMap(value)
	return
}
