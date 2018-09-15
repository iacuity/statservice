package writer

import (
	"github.com/statservice/data"
)

type IWritter interface {
	Init(*data.Config) error
	Write(map[string]int64) error
}
