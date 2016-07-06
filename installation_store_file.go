package addon

import (
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type FileInstallationStore struct {
	stateFilename string
	installations InstallationMap
	mutex         sync.RWMutex
}

func NewFileInstallationStore(stateFilename string) *FileInstallationStore {

	return &FileInstallationStore{
		stateFilename: stateFilename,
		installations: mustReadFile(stateFilename),
	}
}

func mustReadFile(filename string) InstallationMap {

	im := InstallationMap{}

	file, err := os.Open(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return im
		}

		panic(err)
	}

	if err = gob.NewDecoder(file).Decode(&im); err != nil {
		if err == io.EOF {
			return im
		}

		panic(err)
	}

	return im
}

func (s *FileInstallationStore) write() {

	dir, _ := path.Split(s.stateFilename)

	tmpfile, err := ioutil.TempFile(dir, "state")

	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpfile.Name())

	if err := gob.NewEncoder(tmpfile).Encode(s.installations); err != nil {
		panic(err)
	}

	if err := tmpfile.Close(); err != nil {
		panic(err)
	}

	if err := os.Rename(tmpfile.Name(), s.stateFilename); err != nil {
		panic(err)
	}
}

func (s *FileInstallationStore) Add(id string, installation *Installation) {
	s.mutex.Lock()
	s.installations[id] = installation
	s.write()
	s.mutex.Unlock()
}

func (s *FileInstallationStore) Get(id string) (i *Installation) {
	s.mutex.RLock()
	i = s.installations[id]
	s.mutex.RUnlock()
	return
}

func (s *FileInstallationStore) GetAll() InstallationMap {
	m := make(InstallationMap, len(s.installations))

	s.mutex.RLock()
	for id, inst := range s.installations {
		m[id] = inst
	}
	s.mutex.RUnlock()

	return m
}

func (s *FileInstallationStore) Delete(id string) {
	s.mutex.Lock()
	delete(s.installations, id)
	s.write()
	s.mutex.Unlock()
}
