package addon

// https://www.hipchat.com/docs/apiv2/externalPages
type ExternalPage struct {
	Key  string     `json:"key"`
	Name *I18nValue `json:"name,omitempty"`
	Url  string     `json:"url"`
}

// TODO: implement externalPages
