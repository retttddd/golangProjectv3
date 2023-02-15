package storage

type FsStorage struct {
	path string
}

func (st FsStorage) Read(key string) (string, error) {
	return "got it", nil
}

func New() Storage {
	return FsStorage{
		path: "p",
	}
}
