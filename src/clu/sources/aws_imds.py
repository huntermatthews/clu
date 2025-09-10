import json
import logging
import urllib.request
import urllib.error

from clu.provides import Provides
from clu.requires import Requires
from clu.facts import Facts
from clu.sources import PARSE_FAIL_MSG, NET_DISABLED_MSG, Source
from clu.config import get_config

log = logging.getLogger(__name__)
cfg = get_config()


# https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html#dynamic-data-categories
class AwsImds(Source):
    _IMDS_ENDPOINT = "http://169.254.169.254"
    _TOKEN_ENDPOINT = f"{_IMDS_ENDPOINT}/latest/api/token"
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
        provides["imds.ALL"] = self
        for key in self._keys:
            provides[f"imds.{key}"] = self

    def requires(self, requires: Requires) -> None:
        return
        # requires.apis.extend(["not sure what goes here"])

    def parse(self, facts: Facts) -> None:
        if not cfg.net:
            facts["imds.ALL"] = NET_DISABLED_MSG
            return

        if "aws" not in facts["phy.platform"]:
            log.info("Not an AWS EC2 instance, skipping AWS IMDS parsing")
            return

        token = self._get_imds_token()
        if not token:
            log.error("Failed to get AWS IMDSv2 token, cannot retrieve metadata.")
            facts["imds.ALL"] = PARSE_FAIL_MSG
            return

        for key in self._keys:
            data = self._get_imds_metadata(key, token)
            if data and "\n" in data:
                data = data.replace("\n", ", ")
            facts[f"imds.{key}"] = data if data else PARSE_FAIL_MSG

        identity_doc = self._get_imds_identity_document("document", token)
        if not identity_doc:
            facts["imds.identity"] = PARSE_FAIL_MSG
        else:
            for key in identity_doc:
                if identity_doc[key] is not None:
                    facts[f"imds.identity.{key}"] = identity_doc[key]

    def _get_imds_token(self, token_ttl_seconds=21600):
        """
        Obtains a session token for IMDSv2.
        The token_ttl_seconds can be between 1 and 21600 (6 hours).
        """
        headers = {"X-aws-ec2-metadata-token-ttl-seconds": str(token_ttl_seconds)}
        req = urllib.request.Request(self._TOKEN_ENDPOINT, headers=headers, method="PUT")
        try:
            with urllib.request.urlopen(req) as response:
                token = response.read().decode("utf-8")
                return token
        except urllib.error.URLError as e:
            log.error(f"Error obtaining IMDSv2 token: {e}")
            return None

    def _get_imds_metadata(self, path, token):
        """
        Retrieves instance metadata using the provided IMDSv2 token.
        """
        headers = {"X-aws-ec2-metadata-token": token}
        url = f"{self._IMDS_ENDPOINT}/latest/meta-data/{path}"
        req = urllib.request.Request(url, headers=headers)
        try:
            with urllib.request.urlopen(req) as response:
                metadata = response.read().decode("utf-8")
                return metadata
        except urllib.error.URLError as e:
            log.error(f"Error retrieving metadata (path: {path}): {e}")
            return None

    def _get_imds_identity_document(self, path, token):
        """
        Retrieves instance identity document using the provided IMDSv2 token.
        """
        headers = {"X-aws-ec2-metadata-token": token}
        url = f"{self._IMDS_ENDPOINT}/latest/dynamic/instance-identity/document"
        req = urllib.request.Request(url, headers=headers)
        try:
            with urllib.request.urlopen(req) as response:
                document = response.read().decode("utf-8")
                return json.loads(document)
        except urllib.error.URLError as e:
            log.error(f"Error retrieving identity document (url: {url}): {e}")
            return None
