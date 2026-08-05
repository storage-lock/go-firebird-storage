package firebird_storage

import (
	"context"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
)

// FirebirdStorage 基于 Firebird 的存储引擎实现
type FirebirdStorage struct {
	*sql_based_storage.SqlBasedStorage
}

func NewFirebirdStorage(ctx context.Context, options *sql_based_storage.SqlBasedStorageOptions) (*FirebirdStorage, error) {
	baseStorage, err := sql_based_storage.NewSqlBasedStorage(options)
	if err != nil {
		return nil, err
	}
	s := &FirebirdStorage{SqlBasedStorage: baseStorage}
	if err := s.Init(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

const StorageName = "firebird-storage"

func (x *FirebirdStorage) GetName() string {
	return StorageName
}

func (x *FirebirdStorage) Capabilities() []storage.StorageCapability {
	return []storage.StorageCapability{
		storage.CapabilityCAS,
		storage.CapabilityReliableTime,
	}
}
