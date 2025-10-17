import logging
import pprint
from clu.config import get_config
from clu.device42.api import get_host_by_name


log = logging.getLogger(__name__)
cfg = get_config()

_custom_fields_allow = [
    "ITB Team",
    "Stakeholders",
    "Replaced_Device",
    #     "Provision_Data",    # I don't think we use this anymore.
]

_virtual_fields_allow = [
    "virtual_subtype",
    "virtual_host_name",
]

_physical_fields_allow = [
    "asset_no",
    "physical_subtype",
    "manufacturer",
    "serial_no",
    "start_at",
]

_common_fields_allow = [
    "building",
    "customer",
    "device_id",
    "in_service",
    "last_updated",
    "name",
    "service_level",
    "tags",
    "type",
]


def subcmd_query():
    log.info(f"Querying host information... {cfg.hostname}")
    host_info = get_host_by_name(cfg.hostname)
    if not host_info:
        print(f"No host information found for {cfg.hostname}.")
        return 1
    else:
        print("Host Information:")
        transformed_info = transform_host_info(host_info)
        output_host_info(transformed_info)
        return 0


def transform_custom_fields(custom_fields: list) -> dict:
    # Transform the custom fields into the desired format
    cf_result = {}
    for cf in custom_fields:
        cf_result[cf["key"]] = cf["value"]
    return cf_result


def transform_network_fields(ip_addresses: list, mac_addresses: list) -> dict:
    # Transform the network fields into the desired format
    # TODO: I just have no idea what that format should be
    net_result = {}
    net_result["ip_addresses"] = ip_addresses
    net_result["mac_addresses"] = mac_addresses
    return net_result


def filter_keys(input_dict: dict, allowed_keys: list) -> dict:
    return {k: v for k, v in input_dict.items() if k in allowed_keys}


def transform_host_info(host_info: dict):
    output_host_info = filter_keys(host_info, _common_fields_allow)
    if (type := host_info.get("type")) == "virtual":
        output_host_info.update(filter_keys(host_info, _virtual_fields_allow))
    elif type == "physical":
        output_host_info.update(filter_keys(host_info, _physical_fields_allow))
    # TODO: cluster
    # TODO: unknown
    # TODO: ???

    # Custom field support
    cf = transform_custom_fields(host_info["custom_fields"])
    cf = filter_keys(cf, _custom_fields_allow)
    output_host_info["custom_fields"] = cf

    # network field support
    net = transform_network_fields(host_info["ip_addresses"], host_info["mac_addresses"])
    output_host_info["network"] = net

    return output_host_info


def output_host_info(output_host_info: dict):
    pprint.pprint(output_host_info)


# all_fields = [
#     "asset_no",
#     "building",
#     # "category",
#     # "corethread", # later hw support
#     # "cpucore", # later hw support
#     # "cpucount", # later hw support
#     # "cpuspeed", # later hw support
#     "customer",
#     # "customer_id",
#     # "device_external_links",
#     "device_id",
#     # "device_purchase_line_items",
#     "device_sub_type",
#     # "hdd_details",
#     # "hddcount", # later for storage...
#     # "hddraid",
#     # "hddraid_type",
#     # "hddsize",
#     "hw_depth", # "full"
#     "hw_model",
#     # "hw_model_id",
#     "hw_size", # float 2.0 for a 2U device
#     # "id",
#     "in_service",
#     # "ip_addresses",  # SPECIAL CASE
#     # "is_it_blade_host",
#     # "is_it_switch",
#     # "is_it_virtual_host", # later
#     "last_updated",
#     "manufacturer",
#     # "mac_addresses", # SPECIAL CASE
#     "name",
#     # "nonauthoritativealiases",# later for aliases
#     # "notes", # SPECIAL CASE
#     # "orientation", # later for rack
#     "os",
#     "osarch",
#     "osver",
#     "osverno",
#     "preferred_alias",
#     "rack",  # later for rack
#     # "rack_id",
#     # "ram", # later hw support
#     # "room", # later for rack
#     # "row", # later for rack
#     "serial_no",
#     "service_level",
#     "start_at",
#     "tags",
#     "type",
#     # "ucs_manager",
#     "uuid",
#     "virtual_host_name",
#     "where",
#     # "xpos", # later for rack
# ]
