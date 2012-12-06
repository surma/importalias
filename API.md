# Documentation
## API
### v1
Root: `/api/v1`

#### Alias object
An alias object holds all the relevant metadata for an entry.

The Go definition is:

```Go
type Alias struct {
	// URL of the repository to link to
	RepoURL    string `json:"repo_url"`
	// VCS of the repository ("git", "hg" or "bzr")
	RepoType   string `json:"repo_type"`
	// URL to forward to if the URL is being accessed
	// by something else than the go tool.
	ForwardURL string `josn:"forward_url"`
	Alias      string `json:"alias"`
}
```

Accordingly, an JSON example would be:

```JSON
{
	"repo_url": "https://github.com/surma/stacksignal",
	"repo_type": "git",
	"forward_url": "http://surmas-nonexistend-homepage.de/about_stacksignal",
	"alias": "go.surmair.de/stacksignal",
}
```

#### Authentication
Authentication is done via basic auth which has to be your [API-Key](http://notthereyet.com).

#### GET /
Lists all current aliases in use. Returns an array of alias objects.

#### GET /<alias>
Returns a single alias object.

#### PUT /<alias>
Add or update a new alias to your account. The data must be a single alias object.

#### DELETE /<alias>
Delete an alias from your account.
