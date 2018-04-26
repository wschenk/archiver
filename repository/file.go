package repository

import (
	"fmt"
	"github.com/wschenk/archiver"
	"io/ioutil"
	"os"
)

type FileRepository struct {
	path  string
	dirty bool
}

func CreateFileRepository(path string) (*FileRepository, error) {
	err := os.MkdirAll(path, 0755)

	if err != nil {
		return nil, err
	}
	return &FileRepository{path: path}, nil
}

func CreateFileRepoFromStore(store archiver.ArchiveStore, hash string) (*FileRepository, error) {
	path := os.TempDir() + "/" + hash

	fmt.Printf("Checking out the repo to %s\n", path)

	err := os.MkdirAll(path, 0755)

	if err != nil {
		return nil, err
	}

	err = store.Get(hash, path)

	if err != nil {
		return nil, err
	}

	return CreateFileRepository(path)
}

func (repo *FileRepository) Get(path string) ([]byte, error) {
	return ioutil.ReadFile(repo.path + "/" + path)
}

func (repo *FileRepository) Put(key string, data []byte) error {
	repo.dirty = true
	return ioutil.WriteFile(repo.path+"/"+key, data, 0644)
}

func (repo *FileRepository) Dirty() bool {
	return repo.dirty
}

func (repo *FileRepository) Dir() string {
	return repo.path
}
