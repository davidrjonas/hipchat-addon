package addon

import "sync"

type MemoryInstallationStore struct {
	installations InstallationMap
	mutex         sync.RWMutex
}

func NewMemoryInstallationStore() *MemoryInstallationStore {
	return &MemoryInstallationStore{installations: make(InstallationMap)}
}

func (s *MemoryInstallationStore) Add(id string, installation *Installation) {
	s.mutex.Lock()
	s.installations[id] = installation
	s.mutex.Unlock()
}

func (s *MemoryInstallationStore) Get(id string) (i *Installation) {
	s.mutex.RLock()
	i = s.installations[id]
	s.mutex.RUnlock()
	return
}

func (s *MemoryInstallationStore) GetAll() InstallationMap {
	m := make(InstallationMap, len(s.installations))

	s.mutex.RLock()
	for id, inst := range s.installations {
		m[id] = inst
	}
	s.mutex.RUnlock()

	return m
}

func (s *MemoryInstallationStore) Delete(id string) {
	s.mutex.Lock()
	delete(s.installations, id)
	s.mutex.Unlock()
}
