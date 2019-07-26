package storage

import (
	"metrics/bridge/app/types"
)

func New(dsn, db string) (types.MetricsStorage, error) {
	st := &MongoStorage{}

	err := st.Connect(dsn, db)

	if err != nil {
		return nil, err
	}

	return st, nil
}
