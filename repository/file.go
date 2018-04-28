package repository

import (
	"fmt"
	"github.com/wschenk/archiver"
	"io/ioutil"
	"os"
	"path/filepath"
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
	path := filepath.Join(os.TempDir(), hash)

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
	filename := filepath.Join(repo.path, path)

	if _, err := os.Stat(filename); err == nil {
		return ioutil.ReadFile(filename)
	}
	return nil, nil
}

func (repo *FileRepository) Put(key string, data []byte) error {
	repo.dirty = true

	outfile := filepath.Join(repo.path, key)
	baseDir := filepath.Dir(outfile)

	err := os.MkdirAll(baseDir, 0755)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(outfile, data, 0644)
}

func (repo *FileRepository) Dirty() bool {
	return repo.dirty
}

func (repo *FileRepository) Clean() error {
	return os.RemoveAll(repo.path)
}

func (repo *FileRepository) Dir() string {
	return repo.path
}
