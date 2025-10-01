#! /usr/bin/env python3
import urllib.request
import base64


D42_URL = "https://swaggerdemo.device42.com"
D42_USERNAME = "api_user"
D42_PASSWORD = "ap!_user_pr0d"


d42_get_devices_url = D42_URL + "/api/1.0/devices/all/"
request = urllib.request.Request(d42_get_devices_url)
credentials = f"{D42_USERNAME}:{D42_PASSWORD}"
encoded_credentials = base64.b64encode(credentials.encode("utf-8")).decode("ascii")
request.add_header("Authorization", f"Basic {encoded_credentials}")

r = urllib.request.urlopen(request)
if r.getcode() == 200:
    obj = r.read().decode(encoding="utf-8")
    print(obj)

else:
    print(f"Error: {r.getcode()}")
