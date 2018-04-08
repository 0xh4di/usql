package clickhouse

import (
	// DRIVER: clickhouse
	_ "github.com/kshvakov/clickhouse"

	"github.com/xo/usql/drivers"
)

func init() {
	drivers.Register("clickhouse", drivers.Driver{
		AllowMultilineComments: true,
	})
}
