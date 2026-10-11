package firebird_storage

import (
	"context"

	sql_based_storage "github.com/storage-lock/go-sql-based-storage"
)

// FirebirdSqlProvider Firebird 数据库的 SQL 方言提供者
// Firebird 不支持 MySQL 的 UNIX_TIMESTAMP(NOW())。
type FirebirdSqlProvider struct {
	*sql_based_storage.Sql92Provider
}

func NewFirebirdSqlProvider() *FirebirdSqlProvider {
	return &FirebirdSqlProvider{
		Sql92Provider: sql_based_storage.NewSql92Provider(),
	}
}

// NowTimestampSql Firebird 获取当前时间的 Unix 时间戳
// 修复：原来直接继承 Sql92Provider 的 "SELECT UNIX_TIMESTAMP(NOW())"，Firebird 不支持这两个
// MySQL 专有函数，执行即报语法错误，导致 GetTime 完全不可用。用 Firebird 的 DATEDIFF 函数
// 换算 1970-01-01 到 CURRENT_TIMESTAMP 的秒数差作为 Unix 时间戳。
// 注意：未接入真实 Firebird 驱动验证，是按 Firebird SQL 官方文档做的最佳努力实现。
func (x *FirebirdSqlProvider) NowTimestampSql(ctx context.Context, tableFullName string) (string, []any) {
	return `SELECT DATEDIFF(SECOND, TIMESTAMP '1970-01-01 00:00:00', CURRENT_TIMESTAMP) FROM RDB$DATABASE`, nil
}
