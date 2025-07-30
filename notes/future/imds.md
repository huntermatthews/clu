# AWS MSDS

## basic stuff as one query

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


## AWS tool from S3

## Download the ec2-metadata script

wget http://s3.amazonaws.com/ec2metadata/ec2-metadata

## Modify the permission to execute the bash script

chmod +x ec2-metadata

## Simulator

https://github.com/aws/amazon-ec2-metadata-mock
