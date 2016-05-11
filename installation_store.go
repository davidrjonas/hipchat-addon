package addon

type InstallationMap map[string]*Installation

type InstallationStore interface {
	Add(id string, installation *Installation)
	Get(id string) *Installation
	GetAll() InstallationMap
	Delete(id string)
}
