import base64
import json
import logging
import urllib.request
import urllib.error

from clu.config import get_config
from clu.auth import get_primary_credentials

log = logging.getLogger(__name__)
cfg = get_config()


def _api_call(endpoint: str) -> dict:
    # Create a request object
    request = urllib.request.Request(endpoint)

    # Add authentication headers
    creds = get_primary_credentials()
    creds = f"{creds[0]}:{creds[1]}"
    encoded_creds = base64.b64encode(creds.encode("utf-8")).decode("ascii")
    request.add_header("Authorization", f"Basic {encoded_creds}")

    # Make the API call
    try:
        with urllib.request.urlopen(request) as response:
            data = json.loads(response.read().decode("utf-8"))
            log.debug(f"Retrieved Device42 host info: {data}")
    except urllib.error.URLError as e:
        log.error(f"Error calling endpoint {endpoint}: {e}")
        data = {}
    return data


def get_host_by_name(hostname):
    d42_url = f"https://{cfg.d42_server}/api/1.0/device/name/{hostname}"
    return _api_call(d42_url)


def get_host_by_serial(serial_no):
    d42_url = f"https://{cfg.d42_server}/api/1.0/device/serial/{serial_no}"
    return _api_call(d42_url)


def get_host_by_asset(asset_no):
    d42_url = f"https://{cfg.d42_server}/api/1.0/device/asset/{asset_no}"
    return _api_call(d42_url)
