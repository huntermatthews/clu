import logging
import pprint
from clu.config import get_config

log = logging.getLogger(__name__)
cfg = get_config()

_custom_fields_allow = [
    "ITB Team",
    "Stakeholders",
    "Replaced_Device",
    #     "Provision_Data",    # I don't think we use this anymore.
]

_primary_fields_allow = [
    "asset_no",
    "building",
    # "category",
    "corethread",
    "cpucore",
    "cpucount",
    "cpuspeed",
    "customer",
    # "customer_id",
    # "device_external_links",
    "device_id",
    # "device_purchase_line_items",
    "device_sub_type",
    "hdd_details",
    "hddcount",
    "hddraid",
    "hddraid_type",
    "hddsize",
    "hw_depth",
    "hw_model",
    "hw_model_id",
    "hw_size",
    "id",
    "in_service",
    "is_it_blade_host",
    "is_it_switch",
    "is_it_virtual_host",
    "last_updated",
    "manufacturer",
    "name",
    "nonauthoritativealiases",
    # "notes",
    # "orientation",
    "os",
    "osarch",
    "osver",
    "osverno",
    "preferred_alias",
    "rack",
    # "rack_id",
    "ram",
    "room",
    "row",
    "serial_no",
    "service_level",
    "start_at",
    "tags",
    "type",
    # "ucs_manager",
    "uuid",
    "virtual_host_name",
    "where",
    "xpos",
]


def transform_custom_fields(custom_fields: list) -> dict:
    # Transform the custom fields into the desired format
    cf_result = {}
    for cf in custom_fields:
        cf_result[cf["key"]] = cf["value"]
    return cf_result


def filter_keys(input_dict: dict, allowed_keys: list) -> dict:
    return {k: v for k, v in input_dict.items() if k in allowed_keys}


def output_host_info(host_info: dict):
    # pprint.pprint(host_info)
    print("---")
    output_host_info = filter_keys(host_info, _primary_fields_allow)
    cf = transform_custom_fields(host_info["custom_fields"])
    cf = filter_keys(cf, _custom_fields_allow)
    output_host_info.update(cf)

    pprint.pprint(output_host_info)
