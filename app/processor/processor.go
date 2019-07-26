package processor

import (
	"go.uber.org/zap"
	"metrics/bridge/app/mbuffer"
	"metrics/bridge/app/types"
	"time"
)

type MetricsProcessor struct {
	log          *zap.SugaredLogger
	input        types.MetricsInput
	output       types.MetricsStorage
	buffer       mbuffer.MetricsBuffer
	batchSize    int
	batchTimeout time.Duration
}

func New(input types.MetricsInput, output types.MetricsStorage, batchSize, batchTimeout int, log *zap.SugaredLogger) types.MetricsPipeline {
	return &MetricsProcessor{
		input:        input,
		output:       output,
		batchSize:    batchSize,
		batchTimeout: time.Duration(batchTimeout),
		log:          log}
}

func (p *MetricsProcessor) InitBuffer() {
	p.buffer = mbuffer.New(p.batchSize, p.batchTimeout, p)
	go p.buffer.Run()
}

func (p *MetricsProcessor) Commit(data []types.MetricEntity) {
	var ids []uint64

	start := time.Now()

	if err := p.output.Commit(data); err != nil {
		p.log.Error(err.Error())

		for _, e := range data {
			p.input.Put(e.Raw)
		}
	} else {
		p.log.Infof("%d metrics have been inserted [%s]", len(ids), time.Since(start))
		//if err := p.input.DeleteAll(ids); err != nil {
		//	p.log.Error(err.Error())
		//}
	}
}

func (p *MetricsProcessor) InitQueue() {
	for {
		rec, err := p.input.Get()
		if err != nil {
			//p.log.Error(err.Error())
			//time.Sleep(time.Millisecond * 500)
			continue
		}

		p.log.Debugf("id: %d %s", rec.Id, rec.Body)

		kv, err := convertEntity(rec)

		if err != nil {
			p.log.Info(err.Error())
			continue
			//_ = p.input.Delete(rec.Id)
		}

		for !p.buffer.Add(kv) {
			p.log.Debugf("Storage buffer is full!")
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (p *MetricsProcessor) Run() {
	p.InitBuffer()
	p.InitQueue()
}
