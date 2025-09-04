# AWS IMDS

How to determine IF an ec2 instance at all _without_ virt-what
https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/identify_ec2_instances.html



## auth

TOKEN=`curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600"`
curl -H "X-aws-ec2-metadata-token: $TOKEN" <http://169.254.169.254/latest/meta-data/>

## Python

```python
import urllib.request

# Define the IMDSv2 endpoint and token endpoint
IMDS_ENDPOINT = "http://169.254.169.254"
TOKEN_ENDPOINT = f"{IMDS_ENDPOINT}/latest/api/token"

def get_imds_token(token_ttl_seconds=21600):
    """
    Obtains a session token for IMDSv2.
    The token_ttl_seconds can be between 1 and 21600 (6 hours).
    """
    headers = {"X-aws-ec2-metadata-token-ttl-seconds": str(token_ttl_seconds)}
    req = urllib.request.Request(TOKEN_ENDPOINT, headers=headers, method='PUT')
    try:
        with urllib.request.urlopen(req) as response:
            token = response.read().decode('utf-8')
            return token
    except urllib.error.URLError as e:
        print(f"Error obtaining IMDSv2 token: {e}")
        return None

def get_imds_metadata(path, token):
    """
    Retrieves instance metadata using the provided IMDSv2 token.
    """
    headers = {"X-aws-ec2-metadata-token": token}
    url = f"{IMDS_ENDPOINT}/latest/meta-data/{path}"
    req = urllib.request.Request(url, headers=headers)
    try:
        with urllib.request.urlopen(req) as response:
            metadata = response.read().decode('utf-8')
            return metadata
    except urllib.error.URLError as e:
        print(f"Error retrieving metadata: {e}")
        return None

if __name__ == "__main__":
    token = get_imds_token()
    if token:
        instance_id = get_imds_metadata("instance-id", token)
        if instance_id:
            print(f"Instance ID: {instance_id}")

        instance_type = get_imds_metadata("instance-type", token)
        if instance_type:
            print(f"Instance Type: {instance_type}")
    else:
        print("Failed to get IMDSv2 token, cannot retrieve metadata.")
```

```python
import requests
import json

IMDS_HOST = "http://169.254.169.254"
API_VERSION = "latest"
MAX_TOKEN_TTL = 21600 # 6 hours

def get_imds_token():
    """
    Retrieves a session token from the IMDSv2 service.
    """
    token_url = f"{IMDS_HOST}/{API_VERSION}/api/token"
    headers = {"X-aws-ec2-metadata-token-ttl-seconds": str(MAX_TOKEN_TTL)}
    try:
        response = requests.put(token_url, headers=headers, timeout=5)
        response.raise_for_status()
        return response.text
    except requests.exceptions.RequestException as e:
        print(f"Error retrieving IMDSv2 token: {e}")
        return None

def walk_imds(token, path):
    """
    Recursively walks the IMDS endpoint to retrieve all metadata.
    """
    full_url = f"{IMDS_HOST}/{API_VERSION}/{path}"
    headers = {"X-aws-ec2-metadata-token": token}

    try:
        response = requests.get(full_url, headers=headers, timeout=5)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        print(f"Error fetching metadata for path '{path}': {e}")
        return {path: "ERROR"}

    # Check if the path is a directory
    if response.text.endswith("/"):
        children = response.text.splitlines()
        results = {}
        for child in children:
            if child: # Avoids empty strings
                new_path = f"{path}{child}"
                results.update(walk_imds(token, new_path))
        return results
    else:
        # The path is a leaf node, return its value
        try:
            # Attempt to parse as JSON
            return {path: json.loads(response.text)}
        except json.JSONDecodeError:
            # Not JSON, return as a plain string
            return {path: response.text}

def main():
    token = get_imds_token()
    if not token:
        print("Failed to get IMDSv2 token. Cannot proceed.")
        return

    # Start the walk from the base meta-data endpoint
    full_metadata = walk_imds(token, "meta-data/")

    # Print the results in a readable format
    print(json.dumps(full_metadata, indent=2))

if __name__ == "__main__":
    main()

```

## basic stuff as one query

```text
curl -s http://169.254.169.254/latest/dynamic/instance-identity/document
{
  "accountId": "012345678901",
  "architecture": "x86_64",
  "availabilityZone": "eu-central-1c",
  "billingProducts": null,
  "devpayProductCodes": null,
  "marketplaceProductCodes": null,
  "imageId": "ami-01ff76477b9b30d59",
  "instanceId": "i-0b4ae3f67d725bbe7",
  "instanceType": "t3a.nano",
  "kernelId": null,
  "pendingTime": "2022-06-20T09:51:52Z",
  "privateIp": "172.29.40.136",
  "ramdiskId": null,
  "region": "eu-central-1",
  "version": "2017-09-30"
}

curl -s http://169.254.169.254/latest/meta-data
ami-id
ami-launch-index
ami-manifest-path
block-device-mapping/
events/
hostname
iam/
identity-credentials/
instance-action
instance-id
instance-life-cycle
instance-type
local-hostname
local-ipv4
mac
metrics/
network/
placement/
profile
public-hostname
public-ipv4
reservation-id
security-groups
services/

http://169.254.169.254/latest/meta-data/<metadata-path>
```

## AWS tool from S3

## Download the ec2-metadata script

wget <http://s3.amazonaws.com/ec2metadata/ec2-metadata>

## Modify the permission to execute the bash script

chmod +x ec2-metadata

## Simulator

<https://github.com/aws/amazon-ec2-metadata-mock>
