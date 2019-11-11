# Meraki event handler (for Openhab integration)

## Configure Meraki

Network-wide -> Alerts


## Overview

* Use AWS SAM
* Local testing

`sam local start-api --host 0.0.0.0`

Expose publicly using ngrok.
Free: 40 connections / minute -- is that enough?

Implement SSM for config and X-Ray for impact:
https://aws.amazon.com/blogs/compute/sharing-secrets-with-aws-lambda-using-aws-systems-manager-parameter-store/

Install SAM:
https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install-linux.html


