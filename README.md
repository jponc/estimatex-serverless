# estimatex-serverless
Serverless Endpoints for Estimatex.io

# Architecture

- AWS Lambda for Compute
- AWS DynamoDB for datastore
- AWS SNS for fanout event messaging
- Pusher for pubsub / websockets

![Architecture](https://raw.githubusercontent.com/jponc/estimatex-serverless/master/assets/estimatex-sls.png)
