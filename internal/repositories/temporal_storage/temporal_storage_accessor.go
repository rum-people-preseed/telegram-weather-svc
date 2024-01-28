package temporal_storage

type TemporalStorageAccessor struct {
	storage TemporalStorage
	id      int64
}

func CreateTemporalStorageAccessor(storage TemporalStorage, id int64) TemporalStorageAccessor {
	return TemporalStorageAccessor{
		storage: storage, id: id,
	}
}

func (s *TemporalStorageAccessor) Set(key string, value string) {
	err := s.storage.Set(s.id, key, value)
	if err != nil {
		// do something
	}
}
func (s *TemporalStorageAccessor) Get(key string) (string, error) {
	return s.storage.Get(s.id, key)
}

func (s *TemporalStorageAccessor) Del(key string) error {
	return s.storage.DelValue(s.id, key)
}
