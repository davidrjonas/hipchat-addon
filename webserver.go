package addon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Route struct {
	Path    string
	Handler http.HandlerFunc
}

func (a *HipchatAddon) capabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	response, err := json.MarshalIndent(a.descriptor, "", "  ")

	if err != nil {
		a.logger.Errorf("Failed to encode capabilities json: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(response))
}

func (a *HipchatAddon) installHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	installation := new(Installation)

	if err := json.NewDecoder(r.Body).Decode(installation); err != nil {
		a.logger.Errorf("unable to decode json: %v", err)
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
		return
	}

	if err := a.install(installation); err != nil {
		a.logger.Errorf("failed to install: %v", err)
		http.Error(w, "failed to install", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *HipchatAddon) installedHandler(w http.ResponseWriter, r *http.Request) {
	// The URL to redirect the browser to after the integration has been
	// installed. The redirected URL will also contain the query parameters
	// 'redirecturl' and 'installableurl' for the integration configuration
	// page and REST resource for installable information, respectively.
}

func (a *HipchatAddon) uninstalledHandler(w http.ResponseWriter, r *http.Request) {
	installableUrl := r.FormValue("installable_url")
	redirectUrl := r.FormValue("redirect_url")

	if redirectUrl == "" || installableUrl == "" {
		a.logger.Errorf("missing query param redirect_url or installable_url")
		http.Error(w, "missing query params", http.StatusBadRequest)
		return
	}

	if err := a.uninstallFromUrl(installableUrl); err != nil {
		http.Error(w, "Failed to uninstall from url", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectUrl, http.StatusFound)
}

func (a *HipchatAddon) adminPageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func getUrlPath(s string) string {
	parsed, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return parsed.Path
}

func maybeAddRoute(r *[]Route, s string, handler http.HandlerFunc) {
	if s == "" {
		return
	}

	*r = append(*r, Route{getUrlPath(s), handler})
}

func (a *HipchatAddon) Routes() (routes []Route) {
	routes = make([]Route, 0, 50) // FIXME: what capacity do we need?

	routes = append(routes, Route{getUrlPath(a.descriptor.Links.Self), http.HandlerFunc(a.capabilitiesHandler)})

	maybeAddRoute(&routes, a.descriptor.Capabilities.Installable.CallbackUrl, http.HandlerFunc(a.installHandler))
	maybeAddRoute(&routes, a.descriptor.Capabilities.Installable.InstalledUrl, http.HandlerFunc(a.installedHandler))
	maybeAddRoute(&routes, a.descriptor.Capabilities.Installable.UninstalledUrl, http.HandlerFunc(a.uninstalledHandler))
	//maybeAddRoute(&routes, a.descriptor.Capabilities.AdminPage.Url, http.HandlerFunc(a.adminPageHandler))

	for _, glance := range a.descriptor.Capabilities.Glance {
		maybeAddRoute(&routes, glance.QueryUrl, a.JwtAuthHandlerFunc(a.newGlanceHandler(glance)))
	}

	for _, webhook := range a.descriptor.Capabilities.WebHook {
		maybeAddRoute(&routes, webhook.Url, a.JwtAuthHandlerFunc(a.newWebHookHandler(webhook)))
	}
	return
}

type JsonLogRequest struct {
	RemoteAddr string      `json:"client"`
	Method     string      `json:"method"`
	RequestURI string      `json:"uri"`
	Headers    http.Header `json:"headers"`
}

func JsonLogger(handler http.Handler, logger AddonLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(JsonLogRequest{
			RemoteAddr: r.RemoteAddr,
			Method:     r.Method,
			RequestURI: r.RequestURI,
			Headers:    r.Header,
		})

		logger.Info(string(b[:]))

		handler.ServeHTTP(w, r)
	})
}

func (a *HipchatAddon) Serve(listenOn string) {

	mux := http.NewServeMux()

	for _, route := range a.Routes() {
		mux.HandleFunc(route.Path, route.Handler)
	}

	http.ListenAndServe(listenOn, JsonLogger(mux, a.logger))
}
