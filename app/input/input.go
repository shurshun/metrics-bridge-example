package input

import (
	"go.uber.org/zap"
	"metrics/bridge/app/types"
)

func New(addr, tube string, log *zap.SugaredLogger) (result types.MetricsInput, err error) {
	result = &BeanstalkInput{log: log}

	err = result.Connect(addr, tube)

	if err != nil {
		return nil, err
	}

	return result, nil
}
