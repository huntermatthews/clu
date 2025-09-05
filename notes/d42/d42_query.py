import urllib.request
import base64
import json

from creds import D42_URL, D42_USERNAME, D42_PASSWORD

ORIGINAL_DEVICE_NAME = "blueplate"  # original device name


d42_get_devices_url = D42_URL + "/api/1.0/devices/all/"
d42_get_device_url = D42_URL + "/api/1.0/device/name/" + ORIGINAL_DEVICE_NAME  # + '/'
request = urllib.request.Request(d42_get_device_url)
credentials = f"{D42_USERNAME}:{D42_PASSWORD}"
encoded_credentials = base64.b64encode(credentials.encode("utf-8")).decode("ascii")
request.add_header("Authorization", f"Basic {encoded_credentials}")

r = urllib.request.urlopen(request)
if r.getcode() == 200:
    obj = r.read()
    data = json.loads(obj)

    # print(data)
    # print('----')
    print(json.dumps(data, indent=4))
else:
    print(f"Error: {r.getcode()}")

#

# data = urllib.urlencode(params)
# headers = {
#     'Authorization': 'Basic ' + base64.b64encode(D42_USERNAME + ':' + D42_PASSWORD),
#     'Content-Type': 'application/x-www-form-urlencoded'}
# req = urllib2.Request(D42_API_URL, data, headers)

# def uploader(self, data, url):
#         payload = data
#         headers = {
#             'Authorization': 'Basic ' + base64.b64encode(self.username + ':' + self.password),
#             'Content-Type': 'application/json'
#         }

#         r = requests.post(url, data=json.dumps(payload), headers=headers, verify=False)
#         msg =  unicode(payload)
#         if self.debug:
#             print msg
#         msg = 'Status code: %s' % str(r.status_code)
#         print msg
#         msg = str(r.text)
#         if self.debug:
#             print msg
