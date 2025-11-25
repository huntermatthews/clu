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


x = {
    "osverno": "15",
    "preferred_alias": None,
    "is_it_switch": "no",
    "type": "physical",
    "manufacturer": "DellEMC",
    "device_external_links": [],
    "hddraid_type": "Raid 5",
    "building": "Building 12",
    "osver": "9",
    "is_it_blade_host": "no",
    "hw_model_id": 307,
    "row": "F",
    "tags": [],
    "device_purchase_line_items": [],
    "hddcount": 6,
    "in_service": True,
    "asset_no": "02313337",
    "ram": 512.0,
    "ucs_manager": None,
    "hw_model": "PowerEdge R660xs",
    "id": 1926,
    "osarch": 64,
    "hw_depth": 1,
    "cpucount": 2,
    "uuid": "",
    "hw_size": 1.0,
    "where": 5,
    "customer": "ITB",
    "hddsize": 2000.0,
    "virtual_host_name": None,
    "customer_id": 40,
    "name": "blueplate",
    "notes": "iDRAC MAC: d0:46:0c:70:51:b8",
    "service_level": "Production",
    "device_id": 1926,
    "rack_id": 24,
    "room": "CSA1",
    "device_sub_type": "Rackable",
    "is_it_virtual_host": "no",
    "rack": "F26",
    "xpos": 0,
    "corethread": 48,
    "ip_addresses": [
        {
            "label": None,
            "ip": "165.112.174.7",
            "macaddress": "d4:04:e6:ef:9d:90",
            "subnet": "DC-B12-NHGRI-CSA2-165.112.174.0/24",
            "type": 1,
            "subnet_id": 1,
        },
        {
            "label": None,
            "ip": "2607:f220:404:2101::14",
            "macaddress": "d4:04:e6:ef:9d:90",
            "subnet": "DC-B12-NHGRI-CSA2-2607:f220:404:2101::/64",
            "type": 1,
            "subnet_id": 9,
        },
        {
            "label": "mgmt",
            "ip": "10.136.28.96",
            "macaddress": "d0:46:0c:70:51:b8",
            "subnet": "DC-B12-NHGRI-CSA2-mgmt-10.136.28.0/24",
            "type": 1,
            "subnet_id": 2,
        },
    ],
    "last_updated": "2025-08-02T23:06:39.642217Z",
    "hddraid": "Hardware",
    "cpucore": 24,
    "nonauthoritativealiases": [],
    "category": "",
    "cpuspeed": 2000.0,
    "custom_fields": [
        {"value": None, "notes": None, "key": "aws:autoscaling:groupName"},
        {"value": None, "notes": None, "key": "aws:cloudformation:logical-id"},
        {"value": None, "notes": None, "key": "aws:cloudformation:stack-id"},
        {"value": None, "notes": None, "key": "aws:cloudformation:stack-name"},
        {"value": None, "notes": None, "key": "aws:ec2:fleet-id"},
        {"value": None, "notes": None, "key": "aws:ec2launchtemplate:id"},
        {"value": None, "notes": None, "key": "aws:ec2launchtemplate:version"},
        {"value": None, "notes": None, "key": "aws:eks:cluster-name"},
        {"value": None, "notes": None, "key": "Created"},
        {"value": None, "notes": None, "key": "Creator"},
        {"value": None, "notes": None, "key": "eks:cluster-name"},
        {"value": None, "notes": None, "key": "eks:nodegroup-name"},
        {"value": None, "notes": None, "key": "GuardDutyMalwareExclusion"},
        {"value": None, "notes": None, "key": "InspectorEc2Exclusion"},
        {"value": "Unix Team", "notes": None, "key": "ITB Team"},
        {"value": None, "notes": None, "key": "nhgri:itb:automation"},
        {"value": None, "notes": None, "key": "nhgri:itb:customer"},
        {"value": None, "notes": None, "key": "nhgri:itb:ic"},
        {"value": None, "notes": None, "key": "nhgri:itb:itb_team"},
        {"value": None, "notes": None, "key": "nhgri:itb:ssa:automation"},
        {"value": None, "notes": None, "key": "nhgri:itb:ssa:base_ami_name"},
        {"value": None, "notes": None, "key": "nhgri:itb:ssa:build_release"},
        {"value": None, "notes": None, "key": "nhgri:itb:ssa:os_name"},
        {"value": None, "notes": None, "key": "nhgri:itb:ssa:os_version"},
        {"value": None, "notes": None, "key": "Provision_Data"},
        {"value": "blueplate-old", "notes": None, "key": "Replaced_Device"},
        {"value": "Storage Team", "notes": None, "key": "Stakeholders"},
    ],
    "serial_no": "95KQ144",
    "hdd_details": None,
    "os": "Rocky Linux",
    "orientation": 1,
    "start_at": 28.0,
    "aliases": [],
    "mac_addresses": [
        {"mac": "d4:04:e6:ef:9d:90", "port": None, "vlan": "v174-B12", "port_name": ""},
        {"mac": "d0:46:0c:70:51:b8", "port": None, "vlan": "v330-B12", "port_name": "mgmt"},
    ],
}
