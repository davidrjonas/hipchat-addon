package addon

// https://www.hipchat.com/docs/apiv2/capabilities
type CapabilitiesDescriptor struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Key          string        `json:"key"`
	Vendor       *Vendor       `json:"vendor,omitempty"`
	Links        *Links        `json:"links"`
	Capabilities *Capabilities `json:"capabilities"`
	ApiVersion   string        `json:"apiVersion,omitempty"`
}

type Vendor struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Links struct {
	Homepage string `json:"homepage,omitempty"`
	Self     string `json:"self"`
}

type Capabilities struct {
	HipchatApiConsumer *HipchatApiConsumer `json:"hipchatApiConsumer"`

	Installable *Installable `json:"installable"`

	AdminPage *Page `json:"adminPage,omitempty"`

	Oauth2Consumer *Oauth2Consumer `json:"oauth2Consumer,omitempty"`
	Oauth2Provider *Oauth2Provider `json:"oauth2Provider,omitempty"`

	Configurable *Configurable `json:"configurable,omitempty"`

	Action       []*Action       `json:"action,omitempty"`
	Dialog       []*Dialog       `json:"dialog,omitempty"`
	ExternalPage []*ExternalPage `json:"externalPage,omitempty"`
	Glance       []*Glance       `json:"glance,omitempty"`
	WebPanel     []*WebPanel     `json:"webPanel,omitempty"`
	WebHook      []*WebHook      `json:"webhook,omitempty"`
}

type HipchatApiConsumer struct {
	Scopes   []string `json:"scopes"`
	FromName string   `json:"fromName,omitempty"`
	Avatar   *Image   `json:"avatar,omitempty"`
}

type Installable struct {
	AllowGlobal       bool   `json:"allowGlobal"`
	AllowRoom         bool   `json:"allowRoom"`
	CallbackUrl       string `json:"callbackUrl,omitempty"`
	InstalledUrl      string `json:"installedUrl,omitempty"`
	UninstalledUrl    string `json:"uninstalledUrl,omitempty"`
	UpdateCallbackUrl string `json:"updateCallbackUrl,omitempty"`
	UpdatedUrl        string `json:"updatedUrl,omitempty"`
}

type Page struct {
	Url string `json:"url"`
}

type Oauth2Consumer struct {
	RedirectionUrls []string `json:"redirectionUrls"`
}

type Oauth2Provider struct {
	AuthorizationUrl string `json:"authorizationUrl"`
	TokenUrl         string `json:"tokenUrl"`
}

// https://www.hipchat.com/docs/apiv2/configurable
type Configurable struct {
	Url                     string `json:"url"`
	AllowAccessToRoomAdmins bool   `json:"url,omitempty"`
}

type Image struct {
	Url   string `json:"url"`
	Url2x string `json:"url@2x,omitempty"` // required
}

type I18nValue struct {
	I18n  string `json:"i18n,omitempty"`
	Value string `json:"value"`
}
