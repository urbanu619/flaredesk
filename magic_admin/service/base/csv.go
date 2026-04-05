package base

type CsvHandler[T any] struct {
}

func NewCsvHandler[T any](t *T) *CsvHandler[T] {
	return &CsvHandler[T]{}
}

// 保存为csv文件

func (*CsvHandler[T]) Assembly(heads []string, data []*T) (string, error) {

	return "", nil
}
