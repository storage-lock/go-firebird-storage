package firebird_storage

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/storage-lock/go-sql-based-storage"
	"github.com/storage-lock/go-storage"
	"github.com/stretchr/testify/assert"
)

type testSqlite3ConnectionManager struct {
	db *sql.DB
}

func (t *testSqlite3ConnectionManager) Name() string {
	return "test-sqlite3"
}

func (t *testSqlite3ConnectionManager) Take(ctx context.Context) (*sql.DB, error) {
	return t.db, nil
}

func (t *testSqlite3ConnectionManager) Return(ctx context.Context, db *sql.DB) error {
	return nil
}

func (t *testSqlite3ConnectionManager) Shutdown(ctx context.Context) error {
	return nil
}

func TestNewFirebirdStorage_Valid(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?mode=memory")
	assert.NoError(t, err)
	defer db.Close()

	opts := sql_based_storage.NewSqlBasedStorageOptions()
	opts.SetConnectionManager(&testSqlite3ConnectionManager{db: db})
	s, err := NewFirebirdStorage(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, StorageName, s.GetName())
}

func TestFirebirdStorage_GetName(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?mode=memory")
	assert.NoError(t, err)
	defer db.Close()

	opts := sql_based_storage.NewSqlBasedStorageOptions()
	opts.SetConnectionManager(&testSqlite3ConnectionManager{db: db})
	s, err := NewFirebirdStorage(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, StorageName, s.GetName())
}

func TestFirebirdStorage_Capabilities(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?mode=memory")
	assert.NoError(t, err)
	defer db.Close()

	opts := sql_based_storage.NewSqlBasedStorageOptions()
	opts.SetConnectionManager(&testSqlite3ConnectionManager{db: db})
	s, err := NewFirebirdStorage(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	caps := s.Capabilities()
	assert.Contains(t, caps, storage.CapabilityCAS)
	assert.Contains(t, caps, storage.CapabilityReliableTime)
}

func TestFirebirdStorage_Init(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?mode=memory")
	assert.NoError(t, err)
	defer db.Close()

	opts := sql_based_storage.NewSqlBasedStorageOptions()
	opts.SetConnectionManager(&testSqlite3ConnectionManager{db: db})
	s, err := NewFirebirdStorage(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	err = s.Init(context.Background())
	assert.NoError(t, err)
}

func TestFirebirdStorage_Close(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test.db?mode=memory")
	assert.NoError(t, err)
	defer db.Close()

	opts := sql_based_storage.NewSqlBasedStorageOptions()
	opts.SetConnectionManager(&testSqlite3ConnectionManager{db: db})
	s, err := NewFirebirdStorage(context.Background(), opts)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	err = s.Close(context.Background())
	assert.NoError(t, err)
}
