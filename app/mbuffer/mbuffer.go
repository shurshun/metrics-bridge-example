package mbuffer

import (
	"metrics/bridge/app/types"
	"sync"
	"time"
)

const (
	TIMER_LIMIT = 0
	BATCH_LIMIT = 1
)

type Batch struct {
	Size    int
	Timeout time.Duration
}

type MetricsBuffer struct {
	mutex     sync.RWMutex
	data      []types.MetricEntity
	MaxSize   int
	batch     *Batch
	lastParse time.Time
	writer    types.MetricsPipeline
}

func New(size int, timeout time.Duration, writer types.MetricsPipeline) MetricsBuffer {
	return MetricsBuffer{
		MaxSize: size * 10,
		batch:   &Batch{Size: size, Timeout: timeout},
		writer:  writer}
}

func (m *MetricsBuffer) GetSize() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.data)
}

func (m *MetricsBuffer) Add(entity *types.MetricEntity) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.data) > m.MaxSize {
		return false
	}

	m.data = append(m.data, *entity)

	return true
}

func (m *MetricsBuffer) parse(count int, source int) {
	m.mutex.Lock()
	batch := m.data[0:count]
	m.data = m.data[count:]
	m.mutex.Unlock()

	m.writer.Commit(batch)
	m.lastParse = time.Now()
}

func (m *MetricsBuffer) Run() {
	ticker := time.Tick(m.batch.Timeout * time.Second)

	for {
		size := m.GetSize()
		select {
		case <-ticker:
			if (time.Since(m.lastParse).Seconds() >= float64(m.batch.Timeout)) && (size > 0) {
				m.parse(size, TIMER_LIMIT)
			}
		default:
			if size >= m.batch.Size {
				m.parse(m.batch.Size, BATCH_LIMIT)
			}
		}

		time.Sleep(time.Millisecond * 100)
	}
}
