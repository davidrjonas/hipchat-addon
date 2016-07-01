HipChat Connect Addon
=====================

Usage
-----

Fill out a `CapabilitiesDescriptor` struct and then call `New()`. Set an install and an uninstall callback if you like. Run the server and you're golden!

See [examples/](examples/) for a complete addon.

```go
import "github.com/davidrjonas/hipchat-addon"

func main() {
    a := addon.NewWithStateFile(
        &addon.CapabilitiesDescriptor{
            Name:        "My Awesome AddOn",
            Description: "Adds on my awesome.",
            // ...
        },
        "/var/tmp/awesome-addon.db",
        addon.InstalledCallback(func(a *addon.HipChatAddon, installation *addon.Installation) {
            if err := a.SendNotification(installation, &addon.Notification{Message: "Installed!"}); err != nil {
                log.Error(err)
            }
        }),
    )

    a.Serve("127.0.0.1:3000")
}
```

If you want more control over the routes you can provide the mux and server

```go
mux := http.NewServeMux()

for _, route := range a.Routes() {
    mux.HandleFunc(route.Path, route.Handler)
}

http.ListenAndServe("127.0.0.1:3000", mux)
```

### Options

Options are passed to the `New*()` constructors as `HipchatAddonOption` functions. The functional options produce the right functions. See https://commandcenter.blogspot.com.au/2014/01/self-referential-functions-and-design.html

| Option Function                                         | Default |
| ---------------                                         | ------- |
| `Logger(logger AddonLogger)`                            | Standard library adapter (See [logger.go](logger.go)) |
| `HttpClient(client HttpDoer)`                           | [net/http][] (See [core.go](core.go)) |
| `InstallCallback(fn func(i *addon.Installation) error)` | none (always success) |
| `InstalledCallback(fn func(i *addon.Installation))`     | none |
| `UninstalledCallback(fn func(i *addon.Installation))`   | none |


[net/http]: https://golang.org/pkg/net/http/

Persistence
-----------

Storage backends are pluggable. An in-memory store is available for easy testing and a file backed store that is good enough for most addons. The file store uses [encoding/gob][] for serialization.

[encoding/gob]: https://golang.org/pkg/encoding/gob/

```go
// Memory
a := addon.New(descriptor, addon.NewMemoryInstallationStore())

// File
a := addon.NewWithStateFile(descriptor, "/var/tmp/state.db")
```

Logging
-------

By default HipchatAddon will use the standard library logger. But personally I like [logrus][] and did development using it. So, without any extra work, you can pass any logger that implements the [AddonLogger](logger.go) interface.

[logrus]: http://github.com/Sirupsen/logrus

```go
import "github.com/Sirupsen/logrus"

a := addon.New(descriptor, store, addon.Logger(logrus.StandardLogger()))
```

Testing
-------

The http client can be set when initializing the addon. This makes it easy to mock the responses even if you aren't using net/http, as long as the client can handle the [HttpDoer](core.go) interface.

```go
a := addon.New(descriptor, store, addon.HttpClient(fauxhttp))
```

As for this library, I recognize there are as yet no tests.

Development
-----------

I use [ngrok][] to test directly with Atlassian during development. In one window I start up the ngrok proxy to listen locally on port 3000.

```bash
ngrok http 3000
```

It will show you a URL to use to access whatever service is running locally on the port you specify (in this case 3000).

In another window I run my extension on port 3000 with the ngrok URL.

```bash
./echo-addon -port 3000 -url https://deadbeef.ngrok.io
```

I can test it by hitting the capabilities URL locally and over ngrok.

```bash
curl -i http://localhost:3000/capabilities.json

curl -i http://deadbeef.ngrok.io/capabilities.json
```

To install the addon with HipChat click the "Configure Integrations" link in the lower right corner of the HipChat application. Login and it will take you to the HipChat Integrations page for the room you were in. Scroll to the very bottom and click the link "Install an integration from a descriptor URL". Enter your capabilities URL and follow the wizard prompts. Remember to have your service running on the same ngrok URL when you remove it from HipChat!

[ngrok]: https://ngrok.io

Todos and Known Issues
----------------------

- [ ] Safe Installation Access. Right now it is actually unsafe to access Installation that have been looked up from the InstallationStore. I haven't decided whether to return a copy or enforce locking semantics. Advice appreciated.
- [ ] Test all the units
- [ ] Implement Admin Page Handler
- [ ] Implement Actions
- [ ] Implement Dialogs
- [ ] Implement External Pages
- [ ] Implement WebPanel
- [ ] Implement Card-style Notifications
- [ ] Go 1.7 http context for parsed jwt
- [ ] Research and correct room-specific UpdateGlanceData()
- [ ] Add Glance example / docs
- [ ] Implement uninstallCallback to stop uninstallations

Thanks
------

Thanks to the Go authors and the awesome package maintainers that made this easier.

- https://github.com/dgrijalva/jwt-go
- https://github.com/jmoiron/jsonq

