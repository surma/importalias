ImportAlias is a webservice to manage [vanity remote import paths][1]
as offerd by the `go` tool of the [Go programming language][2]. Its
aim is to keep the necessary effort to maintain your aliases to a
minimum.

## Requirements
The whole apps works on a [MongoDB][3] database. Due to
certain query techniques (like `$elemMatch` projections),
MongoDB 2.2 is needed.

## Quick start
For local development, I usually start the app like this:

    $ go run *.go --hostname importalias.surmair.de -l localhost:80 \
      -m mongodb://localhost/importalias --auth-key github:XXXXX:XXXXX \
      --auth-key google:XXXXX:XXXXX --cookie-key XXXXX \
      --auth-config ./auth.config

* `--hostname` sets the hostname for which the static content will be
  served. Every other hostname will be assumed to be an alias request.
* `--auth-key` defines a `<Provider>:<ClientID>:<Secret>` tuple used
  for [OAuth 2.0][4] login.
* `--cookie-key` is the 32-digit hex key used to sign the cookie based
  sessions.
* `--auth-config` contains the configuration (endpoints etc.) for OAuth providers.

For additional flags, set the `--help` flag.

[1]: http://golang.org/cmd/go/#hdr-Remote_import_path_syntax
[2]: http://golang.org
[3]: http://mongodb.org
[4]: http://oauth.net/

---
Version 1.0.0
