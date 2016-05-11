package addon

// https://www.hipchat.com/docs/apiv2/dialogs
type Dialog struct {
	Key     string         `json:"key"`
	Options *DialogOptions `json:"options,omitempty"`
	Title   *I18nValue     `json:"title"`
	Url     string         `json:"url"`
}

type DialogOptions struct {
	Filter          *Filter       `json:"filter,omitempty"`
	Hint            *I18nValue    `json:"hint,omitempty"`
	PrimaryAction   *DialogAction `json:"primaryAction,omitempty"`
	SecondaryAction *DialogAction `json:"secondaryAction,omitempty"`
	Size            *Size         `json:"size,omitempty"`
	Style           string        `json:"style,omitempty"`
}

type Filter struct {
	Placeholder *I18nValue `json:"placeholder"`
}

type DialogAction struct {
	Key     string     `json:"key,omitempty"`
	Enabled bool       `json:"enabled,omitempty"`
	Name    *I18nValue `json:"name"`
}

type Size struct {
	Height string `json:"height"`
	Width  string `json:"width"`
}

// TODO: implement dialogs
