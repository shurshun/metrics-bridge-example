package input

import (
	"github.com/beanstalkd/go-beanstalk"
	"go.uber.org/zap"
	"metrics/bridge/app/types"
	"sync"
	"time"
)

type BeanstalkInput struct {
	log    *zap.SugaredLogger
	mutex  sync.RWMutex
	client *beanstalk.Conn
	source *beanstalk.TubeSet
	tube   *beanstalk.Tube
}

func (b *BeanstalkInput) Connect(addr, tube string) (err error) {
	b.client, err = beanstalk.Dial("tcp", addr)

	if err != nil {
		return err
	}

	b.source = beanstalk.NewTubeSet(b.client, tube)
	b.tube = &beanstalk.Tube{b.client, tube}

	return nil
}

func (b *BeanstalkInput) Disconnect() error {
	return b.client.Close()
}

func (b *BeanstalkInput) Get() (*types.Entity, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	id, body, err := b.source.Reserve(1 * time.Second)
	if err != nil {
		return nil, err
	}

	b.Delete(id)

	return &types.Entity{Id: id, Body: body}, nil
}

func (b *BeanstalkInput) Put(data *types.Entity) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	_, err := b.tube.Put(data.Body, 1, 0, 7*24*time.Hour)

	return err
}

func (b *BeanstalkInput) Delete(id uint64) error {
	return b.client.Delete(id)
}

func (b *BeanstalkInput) DeleteAll(ids []uint64) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.log.Debugf("Q deleting ids: %v", ids)

	for _, id := range ids {

		if err := b.Delete(id); err != nil {
			b.log.Debugf("Q deleting err: %s", err.Error())
			return err
		}
	}

	b.log.Debug("Q deleting ids: ok")

	return nil
}

func (b *BeanstalkInput) Release(id uint64) error {
	return b.client.Release(id, 1, 1*time.Minute)
}

func (b *BeanstalkInput) ReleaseAll(ids []uint64) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.log.Debugf("Q releasing ids: %v", ids)

	for _, id := range ids {
		if err := b.Release(id); err != nil {
			b.log.Debugf("Q releasing err: %s", err.Error())
			return err
		}
	}

	b.log.Debug("Q releasing ids: ok")

	return nil
}
