import base64
import json
import logging
import urllib.request
import urllib.error

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import Source
from clu.sources import PARSE_FAIL_MSG, NET_DISABLED_MSG
from clu.auth import get_primary_credentials
from clu.config import get_config

log = logging.getLogger(__name__)
cfg = get_config()

PROD_URL = "https://inventory.nhgri.nih.gov"
TEST_URL = "https://itbinventorytest01.nhgri.nih.gov"
ORIGINAL_DEVICE_NAME = "blueplate"


# https://api.device42.com
class Device42(Source):
    _keys = [
        "ami-id",
        "block-device-mapping",
        "instance-id",
        "instance-type",
        "local-hostname",
        "local-ipv4",
        "mac",
        "profile",
        "reservation-id",
        "security-groups",
    ]

    def provides(self, provides: Provides) -> None:
        # FIX: Provide for the all key as a hack to make the output right
        provides["d42.ALL"] = self
        for key in self._keys:
            provides[f"d42.{key}"] = self

    def requires(self, requires: Requires) -> None:
        return
        # requires.apis.extend(["not sure what goes here"])

    def parse(self, facts: Facts) -> None:
        if not cfg.net:
            facts["d42.ALL"] = NET_DISABLED_MSG
            return

        creds = get_primary_credentials()
        if not creds:
            log.error("Failed to get Device42 credentials, cannot retrieve metadata.")
            facts["d42.ALL"] = PARSE_FAIL_MSG
            return

        # Use the Device42 API to retrieve metadata
        d42_get_device_url = PROD_URL + "/api/1.0/device/name/" + ORIGINAL_DEVICE_NAME
        request = urllib.request.Request(d42_get_device_url)
        credentials = f"{creds[0]}:{creds[1]}"
        encoded_credentials = base64.b64encode(credentials.encode("utf-8")).decode("ascii")
        request.add_header("Authorization", f"Basic {encoded_credentials}")

        try:
            with urllib.request.urlopen(request) as response:
                device_data = json.loads(response.read().decode("utf-8"))
                log.debug(f"Retrieved Device42 metadata: {device_data}")
                print("keys", device_data.keys())
                for key in device_data:
                    print(key)
                    print(f"d42.{key}: {device_data[key]}")
                    facts[f"d42.{key}"] = device_data[key]
        except urllib.error.URLError as e:
            log.error(f"Error retrieving Device42 metadata: {e}")
            facts["d42.ALL"] = PARSE_FAIL_MSG
        print(facts)


