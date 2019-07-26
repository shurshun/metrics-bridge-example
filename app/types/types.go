package types

type Entity struct {
	Id   uint64
	Body []byte
}

type QueueEntity struct {
	Ts  int64  `json:"ts"`
	Ip  string `json:"ip"`
	Uri string `json:"uri"`
}

type MetricEntity struct {
	Bucket  string
	Data    map[string]string
	Raw		*Entity
}

type MetricsPipeline interface {
	Commit(data []MetricEntity)
	Run()
}

type MetricsInput interface {
	Connect(addr, tube string) (err error)
	Disconnect() error
	Get() (*Entity, error)
	Delete(id uint64) error
	DeleteAll(ids []uint64) error
	Release(id uint64) error
	ReleaseAll(ids []uint64) error
	Put(data *Entity) error
}

type MetricsStorage interface {
	Connect(dsn, db string) error
	Disconnect() error
	Commit(data []MetricEntity) error
}
