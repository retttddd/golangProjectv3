package service

type storage interface {
	Read(key string) (*StorageModel, error)
	Write(key string, model *StorageModel) error
}
