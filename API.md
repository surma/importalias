# Documentation
## API
### v1
Root: `/api/v1`

#### Authentication
Authentication is done via basic auth which has to be your
[API-Key](http://importalias.surmair.de/#/wtf).

#### /me
##### GET /me
Show your API-Key and your user ID

#### /domains
##### POST /domains/[domain]
Claim a domain. If this call succeedes, you can manage this domain’s
aliases. The server will look for a TXT-Record on this domain’s DNS
which has to be names `_importalias.[domain]` and has to contain
your user ID.
Returns: 204 on success

##### GET /domains/
Returns: 200 and a list of all claimed domains


##### DELETE /domains/[domain]
Delete a domain (and all its aliases) from the list, effictively
makeing it reclaimable.
Returns: 204 on success

##### GET /domains/[domain]
Returns: 200 and all aliases defined for a domain on success

##### PUT /domains/[domain]
Add an alias under the given domain. An alias object holds all the
relevant metadata for an entry.
Payload:

```JSON
PUT /domains/go.surmair.de

{
	"repo_url": "https://github.com/surma/stacksignal",
	"repo_type": "git",
	"forward_url": "http://surmas-nonexistent-homepage.de/about_stacksignal",
	"alias": "/stacksignal",
}
```

Returns: 201 and the domain object with an additional `id` field

##### DELETE /domains/[domain]/[id]
Delete a single alias
Return: 204 on success

