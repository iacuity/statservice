package writer

type IWritter interface {
	Write(map[string]int64) error
}
