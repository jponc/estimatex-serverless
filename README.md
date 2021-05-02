# estimatex-serverless
Serverless Endpoints for [estimatex.io](http://estimatex.io)

# Architecture

- AWS API Gateway for authorising & exposing JSON endpoints
- AWS Lambda for Compute
- AWS DynamoDB for datastore
- AWS SNS for fanout event messaging
- Pusher for pubsub / websockets

![Architecture](https://raw.githubusercontent.com/jponc/estimatex-serverless/master/assets/estimatex-sls.png)

# Built using Serverless Framework

<img src="https://miro.medium.com/max/1400/1*CuALG7dV2rLky1sapJbnUQ.png" width="400" height="140">

# Dashbird for Monitoring & Alerting

<img src="https://mk0dashbirdioprthk8x.kinstacdn.com/wp-content/uploads/2021/03/dashbird-logo@2x.png" width="350" height="140">
