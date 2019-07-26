package processor

import (
	"encoding/json"
	"errors"
	"metrics/bridge/app/types"
	"net/url"
	"strings"
	"time"
)

func getCollName(ts int64) string {
	loc, _ := time.LoadLocation("Europe/Moscow")
	name := time.Unix(ts, 0).In(loc)
	return name.Format("2006-01-02")
}

func getRecordTime(ts int64) string {
	loc, _ := time.LoadLocation("Europe/Moscow")
	name := time.Unix(ts, 0).In(loc)
	return name.Format("15:04:05")
}

func convertEntity(rawEntity *types.Entity) (*types.MetricEntity, error) {
	if len(string(rawEntity.Body)) == 0 {
		return nil, errors.New("body is null")
	}

	var queueEntity types.QueueEntity

	err := json.Unmarshal(rawEntity.Body, &queueEntity)
	if err != nil {
		return nil, err
	}

	args, err := url.ParseQuery(strings.TrimLeft(queueEntity.Uri, "?"))
	if err != nil {
		return nil, err
	}

	result := &types.MetricEntity{}
	result.Raw = rawEntity
	result.Data = make(map[string]string)

	for k, v := range args {
		if len(v) > 0 {
			result.Data[k] = v[0]
		} else {
			result.Data[k] = ""
		}
	}

	result.Data["ip"] = queueEntity.Ip
	result.Data["_time"] = getRecordTime(queueEntity.Ts)

	result.Bucket = getCollName(queueEntity.Ts)

	return result, nil
}
