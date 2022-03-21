package cache

type Status struct {
	KeySize   int
	ValueSize int
	Count     int
}

func (s *Status) add(key string, value []byte) {
	s.KeySize += len(key)
	s.ValueSize += len(value)
	s.Count++
}

func (s *Status) del(key string, value []byte) {
	s.KeySize -= len(key)
	s.ValueSize -= len(value)
	s.Count--
}
