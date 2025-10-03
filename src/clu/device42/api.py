import base64
import json
import logging
import urllib.request
import urllib.error

from clu.auth import get_primary_credentials
from clu.config import get_config


log = logging.getLogger(__name__)
cfg = get_config()


class Device42API:
    _credentials = None

    @property
    def credentials(self):
        if not Device42API._credentials:
            Device42API._credentials = get_primary_credentials()
            if not Device42API._credentials:
                log.error("Failed to get Device42 credentials, cannot make API calls.")
        return Device42API._credentials

    def __init__(self):
        self.d42_server = cfg.d42_server

    def _api_call(self, endpoint: str) -> dict:
        request = urllib.request.Request(endpoint)
        creds = f"{self.credentials[0]}:{self.credentials[1]}"
        encoded_creds = base64.b64encode(creds.encode("utf-8")).decode("ascii")
        request.add_header("Authorization", f"Basic {encoded_creds}")

        try:
            with urllib.request.urlopen(request) as response:
                data = json.loads(response.read().decode("utf-8"))
                log.debug(f"Retrieved Device42 host info: {data}")
        except urllib.error.URLError as e:
            log.error(f"Error calling endpoint {endpoint}: {e}")
            data = {}
        return data

    def get_host_by_name(self, hostname):
        d42_url = f"{self.d42_server}/api/1.0/device/name/{hostname}"
        return self._api_call(d42_url)

    def get_host_by_serial(self, serial_no):
        d42_url = f"{self.d42_server}/api/1.0/device/serial/{serial_no}"
        return self._api_call(d42_url)

    def get_host_by_asset(self, asset_no):
        d42_url = f"{self.d42_server}/api/1.0/device/asset/{asset_no}"
        return self._api_call(d42_url)


# def get_host_info(hostname: str) -> dict:
#     creds = get_primary_credentials()
#     if not creds:
#         log.error("Failed to get Device42 credentials, cannot retrieve host info.")
#         return {}

#     # Use the Device42 API to retrieve basic host info
#     d42_get_device_url = (
#         "https://itbinventorytest01.nhgri.nih.gov" + "/api/1.0/device/name/" + hostname
#     )
#     request = urllib.request.Request(d42_get_device_url)
#     credentials = f"{creds[0]}:{creds[1]}"
#     encoded_credentials = base64.b64encode(credentials.encode("utf-8")).decode("ascii")
#     request.add_header("Authorization", f"Basic {encoded_credentials}")

#     try:
#         with urllib.request.urlopen(request) as response:
#             device_data = json.loads(response.read().decode("utf-8"))
#             log.debug(f"Retrieved Device42 host info: {device_data}")
#     except urllib.error.URLError as e:
#         log.error(f"Error retrieving Device42 host info: {e}")
#         device_data = {}
#     return device_data
