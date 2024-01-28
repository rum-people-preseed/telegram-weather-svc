package temporal_storage

type TemporalStorage interface {
	Set(id int64, key string, value string) error
	Del(id int64) error
	Get(id int64, key string) (string, error)
	DelValue(id int64, key string) error
}
