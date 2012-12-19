# Documentation
## API
### v1
Root: `/api/v1`

#### Authentication
Authentication is done via basic auth which has to be your [API-Key](http://notthereyet.com).

#### /domains
To be able to create an alias under a certain domain, domains have to be claimed.

##### POST /domains/<domain>
Claim a domain
Returns: 204 on success

##### GET /domains/
Returns: 200 and a list of all claimed domains


##### DELETE /domains/<domain>
Delete a domain (and all its aliases) from the list, effictively makeing it reclaimable
Returns: 204 on success

##### GET /domains/<domain>
Returns: 200 and all aliases defined for a domain on success

##### PUT /domains/<domain>
Add an alias under the given domain.
An alias object holds all the relevant metadata for an entry.
Payload:

```JSON
POST /domains/go.surmair.de

{
	"repo_url": "https://github.com/surma/stacksignal",
	"repo_type": "git",
	"forward_url": "http://surmas-nonexistend-homepage.de/about_stacksignal",
	"alias": "/stacksignal",
}
```

Returns: 201 and the domain object with an additional `id` field

##### Delete /domains/<domain>/<id>
Delete a single alias
Return: 204 on success

