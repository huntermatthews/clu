# URLIB

<https://dev.to/bowmanjd/http-calls-in-python-without-requests-or-other-external-dependencies-5aj1>

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

```python
import json

newConditions = {"con1":40, "con2":20, "con3":99, "con4":40, "password":"1234"}
params = json.dumps(newConditions).encode('utf8')
req = urllib.request.Request(conditionsSetURL, data=params,
                             headers={'content-type': 'application/json'})
response = urllib.request.urlopen(req)

```
