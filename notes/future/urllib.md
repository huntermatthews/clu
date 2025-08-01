# URLIB

```python

import urllib2
import json

def basic_authorization(user, password):
    s = user + ":" + password
    return "Basic " + s.encode("base64").rstrip()

def submit_pull_request(user, repo):
    auth = (settings.username, settings.password)
    url = '<https://api.github.com/repos/>' + user + '/' + repo + '/pulls'
    params = {'title': 'My Title', 'body': 'My Boday'}
    req = urllib2.Request(url,
        headers = {
            "Authorization": basic_authorization(settings.username, settings.password),
            "Content-Type": "application/json",
            "Accept": "*/*",
            "User-Agent": "Myapp/Gunio",
        }, data = json.dumps(params))
    f = urllib2.urlopen(req)
```
