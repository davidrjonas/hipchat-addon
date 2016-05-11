package addon

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jmoiron/jsonq"
)

type HipchatAddon struct {
	descriptor          *CapabilitiesDescriptor
	installations       InstallationStore
	installCallback     InstallationsPrechangeCallback
	installedCallback   InstallationsChangedCallback
	uninstalledCallback InstallationsChangedCallback
	logger              AddonLogger
	http                HttpDoer
}

// If an error is returned the handler will return a 500 error to the HipChat
// server and change to the store will be aborted.
type InstallationsPrechangeCallback func(a *HipchatAddon, i *Installation) error

// Called *after* the change to the store.
type InstallationsChangedCallback func(a *HipchatAddon, i *Installation)

type HipchatAddonOption func(*HipchatAddon) error

func Logger(log AddonLogger) HipchatAddonOption {
	return func(a *HipchatAddon) error {
		a.logger = log
		return nil
	}
}

type HttpDoer interface {
	Do(h *http.Request) (resp *http.Response, err error)
}

func HttpClient(client HttpDoer) HipchatAddonOption {
	return func(a *HipchatAddon) error {
		a.http = client
		return nil
	}
}

func InstallCallback(fn InstallationsPrechangeCallback) HipchatAddonOption {
	return func(a *HipchatAddon) error {
		a.installCallback = fn
		return nil
	}
}

func InstalledCallback(fn InstallationsChangedCallback) HipchatAddonOption {
	return func(a *HipchatAddon) error {
		a.installedCallback = fn
		return nil
	}
}

func UninstalledCallback(fn InstallationsChangedCallback) HipchatAddonOption {
	return func(a *HipchatAddon) error {
		a.uninstalledCallback = fn
		return nil
	}
}

func NewWithStateFile(descriptor *CapabilitiesDescriptor, stateFilename string, options ...HipchatAddonOption) *HipchatAddon {
	return New(descriptor, NewFileInstallationStore(stateFilename), options...)
}

func New(descriptor *CapabilitiesDescriptor, store InstallationStore, options ...HipchatAddonOption) *HipchatAddon {
	addon := &HipchatAddon{
		descriptor:          descriptor,
		installations:       store,
		installCallback:     func(a *HipchatAddon, i *Installation) error { return nil },
		installedCallback:   func(a *HipchatAddon, i *Installation) {},
		uninstalledCallback: func(a *HipchatAddon, i *Installation) {},
	}

	addon.setOptions(options...)
	addon.initializeDefaults()

	return addon
}

func (a *HipchatAddon) setOptions(options ...HipchatAddonOption) {
	for _, op := range options {
		err := op(a)
		if err != nil {
			panic(err)
		}
	}
}

func (a *HipchatAddon) initializeDefaults() {
	if a.logger == nil {
		a.logger = NewStdLogger()
	}

	if a.http == nil {
		a.http = &http.Client{}
	}
}

func (a *HipchatAddon) install(installation *Installation) error {

	if err := a.installCallback(a, installation); err != nil {
		return err
	}

	req, err := http.NewRequest("GET", installation.CapabilitiesUrl, nil)

	if err != nil {
		return err
	}

	r, err := a.http.Do(req)

	if err != nil {
		return err
	}

	var data map[string]interface{}

	if err = json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	jq := jsonq.NewQuery(data)

	if installation.TokenUrl, err = jq.String("capabilities", "oauth2Provider", "tokenUrl"); err != nil {
		return err
	}
	if installation.ApiUrl, err = jq.String("capabilities", "hipchatApiProvider", "url"); err != nil {
		return err
	}

	a.installations.Add(installation.OauthId, installation)

	a.installedCallback(a, installation)

	return nil
}

func (a *HipchatAddon) uninstallFromUrl(url string) error {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}

	r, err := a.http.Do(req)

	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	oauthId := data["oauthId"].(string)

	installation := a.installations.Get(oauthId)
	if installation == nil {
		return nil
	}

	a.installations.Delete(oauthId)

	a.uninstalledCallback(a, installation)

	return nil
}

func (a *HipchatAddon) getAccessToken(installation *Installation) (*AccessToken, error) {
	return installation.GetAccessToken(a.http)
}

func (a *HipchatAddon) postWithToken(installation *Installation, url string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, err
	}

	token, err := a.getAccessToken(installation)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.String())

	return a.http.Do(req)
}
