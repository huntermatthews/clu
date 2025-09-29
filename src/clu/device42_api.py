import base64
import json
import logging
import urllib.request
import urllib.error

from clu.auth import get_primary_credentials


log = logging.getLogger(__name__)

TEST_URL = "http://swaggerdemo.device42.com"
ORIGINAL_DEVICE_NAME = "abqhxgpc01"

# THIS IS NOT A CREDENTIAL LEAK!
# THESE ARE PUBLICLY AVAILABLE ON THE DEVICE42 SWAGGER DEMO SITE
USERNAME = "guest"
PASSWORD = "device42_rocks!"


def get_host_info() -> dict:
    creds = get_primary_credentials()
    if not creds:
        log.error("Failed to get Device42 credentials, cannot retrieve host info.")
        return {}

    # Use the Device42 API to retrieve basic host info
    d42_get_device_url = TEST_URL + "/api/1.0/device/name/" + ORIGINAL_DEVICE_NAME
    request = urllib.request.Request(d42_get_device_url)
    credentials = f"{creds[0]}:{creds[1]}"
    encoded_credentials = base64.b64encode(credentials.encode("utf-8")).decode("ascii")
    request.add_header("Authorization", f"Basic {encoded_credentials}")

    try:
        with urllib.request.urlopen(request) as response:
            device_data = json.loads(response.read().decode("utf-8"))
            log.debug(f"Retrieved Device42 host info: {device_data}")
    except urllib.error.URLError as e:
        log.error(f"Error retrieving Device42 host info: {e}")
        device_data = {}
    return device_data
