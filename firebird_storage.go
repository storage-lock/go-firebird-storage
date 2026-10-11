package firebird_storage

import (
	"context"
	"strings"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"time"
)

// FirebirdStorage 基于 Firebird 的存储引擎实现
type FirebirdStorage struct {
	*sql_based_storage.SqlBasedStorage
}

func NewFirebirdStorage(ctx context.Context, options *sql_based_storage.SqlBasedStorageOptions) (*FirebirdStorage, error) {
	// 如果调用方没有显式指定方言（还是默认的 Sql92Provider），换成 Firebird 自己的方言
	// 提供者：Sql92Provider.NowTimestampSql 是 MySQL 专有语法，在 Firebird 上会直接报语法错误
	if _, ok := options.SqlProvider.(*sql_based_storage.Sql92Provider); ok || options.SqlProvider == nil {
		options.SetSqlProvider(NewFirebirdSqlProvider())
	}
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

// CreateWithVersion Firebird 唯一约束冲突时返回 ErrLockAlreadyExists
// 注意：未接入真实 Firebird 驱动验证，错误消息按 Firebird 官方文档做的最佳努力匹配。
func (x *FirebirdStorage) CreateWithVersion(ctx context.Context, lockId string, version storage.Version, lockInformation *storage.LockInformation) error {
	err := x.SqlBasedStorage.CreateWithVersion(ctx, lockId, version, lockInformation)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "violation of PRIMARY or UNIQUE KEY constraint") || strings.Contains(msg, "attempt to store duplicate value") {
			return storage_lock.ErrLockAlreadyExists
		}
	}
	return err
}

const StorageName = "firebird-storage"

func (x *FirebirdStorage) GetName() string {
	return StorageName
}

func (x *FirebirdStorage) Capabilities() []storage.StorageCapability {
	return []storage.StorageCapability{
		storage.CapabilityCAS,
		storage.CapabilityReliableTime,
		storage.CapabilityWatch,
	}
}

func (x *FirebirdStorage) Watch(ctx context.Context, lockId string) (<-chan *storage.LockInformation, error) {
	return storage.NewPollingLockWatcher(x, 500*time.Millisecond).Watch(ctx, lockId)
}
