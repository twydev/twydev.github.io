---
title: "Ultimate AWS Certified Developer Associate 2019 - Part 4"
toc: true
toc_label: "Chapters"
published: false
---

Part 4 covers components offered by AWS that forms the technology stack for serverless applications (Lambda and Step Functions, DynamoDB, Cognito, API Gateway, S3, SNS, SQS, Kinesis, Aurora Serverless).

## AWS Lambda

### Benefits of Lambda

- pay per request and compute time
  - usually very cheap that is why it is popular
- integrated with whole AWS stack
  - easy monitoring on CloudWatch
- compatible with many languages
  - NodeJS, Python, Java, C# (.NET Core), Golang, Powershell
- easy to scale up resources
  - up to 3GB RAM
  - larger RAM sizes also mean larger CPU and network resources

### Lambda Configuration

- Timeout, 3s (default) to 15 mins
- Environment Variables
- Allocated Memory (128M to 3G)
- deploy within VPC, assign security group
- attach IAM execution role to Lambda

### Concurrency and Throttling

- can have up to 1000 concurrent executions per account (request ticket to increase)
- set **reserved concurrency** at function level
- invocation beyond limit will trigger a throttle
  - if synchronous invoke, return ThrottleError 429
  - if asynchronous invoke, lambda will retry automatically or message goes to DLQ
  - DLQ can be just a SNS topic or a SQS queue (ensure IAM permissions)

### Logging, Monitoring, Tracing

CloudWatch

- Lambda execution logs are stored in CloudWatch
- Lambda metrics are displayed in CloudWatch
- ensure Lambda execution role has sufficient permission to write to CloudWatch

X-Ray

- enable in Lambda configuration (Lambda automatically runs X-Ray daemon)
- use AWS SDK in function code
- similarly ensure IAM role has appropriate permissions

### Lambda Limits

1. Memory Allocation

- between 128 MB to 3008 MB (64 MB increments)

2. Maximum Execution Time

- previously 5 mins, now 15 mins

3. Disk Capacity in Function Container

- 512 MB

4. Concurrency Limit

- 1000 executions

5. Function Deployment Size

- 50 MB (compressed .zip)
- 250 MB code and dependencies (uncompressed deployment)
- can use /tmp directory to load other files at startup

6. Environment Variables Size

- 4 KB

### Versioning

- **\$LATEST** version is mutable
- published versions are immutable
  - immutable for both code and configs
  - versions numbers are increasing
  - each version gets their own ARN (Amazon Resource Name)
  - all versions can be accessed using the correct ARN
- **Aliases** are pointers to specific versions of lambda function
  - mutable, can point to a different version
  - can be used to define dev, test, prod etc
  - alias can redirect weighted traffic to different versions
  - enable Blue / Green deployment
  - abstracts away the actual version, enable stable configuration of event triggers and destinations
  - aliases have their own ARN

### External Dependencies

- Dependency packages need to be installed alongside function code and zip together
  - NodeJS, use npm and _node_modules_ directory
  - Python, use _pip --target_ options
  - Java, include relevant _.jar_ files
- Upload zip straight to lambda if less than 50 MB, else upload to S3 first
- Native libraries work, need to be compiled on Amazon Linux first

### Deploy using CloudFormation

- Simply zip and store function code in S3
- Reference S3 location in CF

### Lambda local /tmp directory

- can be used if function needs to download a big file
- or if function needs disk space
- max size of /tmp is 512 MB
- provides transient cache across multiple invocations (if context is frozen, helpful for checkpointing work)
- recommended to use S3 if persistence is required

### Best Practices

- Perform heavy-duty work outside of function handler, such as
  - connect to db
  - initialize AWS SDK
  - pull in dependencies or datasets
- Use environment variables for
  - db connection strings, S3 bucket
  - password, sensitive values (can be encrypted using KMS)
- Minimize deployment package size to its runtime necessities
  - break down the function if necessary
  - don't exceed lambda limits
- avoid using recursive code (have lambda function calls itself)
- recommended not to use VPC unless necessary

### Lambda@Edge

Deploy Lambda functions alongside CloudFront CDN

- more responsive
- deployed globally
- can be used to customize CDN content
- pay for only what you use

Can be used to change CloudFront requests / responses

- after requests received from viewer
- before requests forwarded to origin
- after responses received from origin
- before responses forwarded to viewer
- can even generate responses for viewer without ever forwarding to origin

Use cases

- website security and privacy
- dynamic web app on the edge
- search engine optimization
- intelligently route across origins and data centers
- bot mitigation on the edge
- real-time image transformation
- A/B testing
- user authentication / authorization
- user prioritization
- user tracking and analytics

## AWS DynamoDB

### DynamoDB Overview

Traditional RDB have strong requirements on how data should be modelled. Supports join, aggregations, computations, and have to be scaled vertically with more CPU power and RAM and IO throughput.

In comparison, NoSQL databases (not-only-SQL) are non-relational, do not support joins, and are distributed (scale horizontally). In essence, all the data needed for a query is present within the row.

DynamoDB is

- fully managed, highly available (replication across 3 AZs)
- millions of requests per seconds, hundreds of TB of storage.
- Fast and consistent (low retrieval latency)
- Integrated with IAM
- Event driven programming with DynamoDB Streams
- low cost and auto-scaling

### Tables

- has primary key defined at creation
  - can also use partition key + sort key (aka range key) for unique primary key (both not nullable)
  - ideally partition key should have high cardinality to produce sparse hashes
- each items (rows) has attributes (nullable, extendable with new fields)
- max item size **400 KB**
- supported data types:
  - Scalar Types: String, Number, Binary, Boolean, Null
  - Document Type: List, Map
  - Set Types: String Set, Number Set, Binary Set

### Throughputs

Read Capacity Units (RCU) and Write Capacity Units (WCU):

- if throughput exceeds limit, will use burst credit for auto-scaling if configured.
- else, will encounter "ProvisionedThroughputException" which requires client to exponential back-off retry

Key throughput rates (needed for calculations)

- **1 WCU = 1 write per second for an item up to 1 KB**
- Eventually Consistent Read **1 RCU = 2 read per second for items up to 4 KB**
- Strong Consistent Read **1 RCU = 1 read per second for items up to 4 KB**
- There is no pooling of throughput rate. All item sizes will be rounded up to the nearest rate unit and charged per read / write
- **RCU and WCU are evenly distributed across all partitions of the dynamoDB table** therefore hot partitions will suffer amplified performance penalty

Solutions to throttling:

- exponential back-off retry
- distribute partition keys better to avoid hot partitions
- for higher Read performance, use DynamoDB Accelerator (DAX)

### APIs

- PutItem - full replace / write new item
- UpdateItem - partial update of some attributes
- ConditionalWrites - accept write only if conditions are met (checks in a single transaction)
- DeleteItem - on a row, can be conditional
- DeleteTable - faster than DeleteItem
- BatchWriteItem - up to 25 PutItem / DeleteItem in one call (up to 16MB data).
  - Improves latency by reducing API calls round trip
  - Server performs writes in parallel
  - partial failure is possible. Need to retry failed items
- GetItem - using primary keys
  - option for strong consistent reads.
  - can have ProjectionExpression to only return requested attributes over network
- BatchGetItem - up to 100 items (16 MB data), in parallel
- Query - return items based on PartitionKey value (= operator) and optional SortKey value (comparison operators)
  - FilterExpression can be used (performed on client side however)
  - returns up to 1 MB data
  - able to do pagination on results
  - can query the table, LSI or GSI
- Scan - filter the entire table, consumes a lot of RCU
  - not encouraged, to be avoided with better table design
  - use Limit or reduce size of the result
  - scans can be done in parallel (but still consume same RCU)
  - can use ProjectionExpression + FilterExpression (but still consume same RCU)
  - **--page-size** option will make more api calls in the background, retrieving a smaller set of results in each background call, to prevent network timeout
  - **--max-item** option specifies the number of items to be returned in total by the api call.
  - **--starting-token** option, for the next api call to continue scanning the table from a location

### Local Secondary Index (LSI)

- up to 5 LSI per table
- it is an alternate range key, local to each hash key
- consists of exactly 1 scalar attribute
- must be defined at table creation

### Global Secondary Index (GSI)

- purpose is to speed up query on a table for non-key attributes
- defines a new partition key + optional sort key
- index is a "new table" and can project attributes on it
  - (KEYS_ONLY) project original partition key and sort key of table
  - (INCLUDE) specify some attributes to project
  - (ALL) use all attributes from main table
- must define RCU / WCU for the GSI
  - throttled writes on GSI throttles the main table as well
  - even if main table has sufficient WCU
  - LSI uses RCU / WCU of main table, no special considerations
- possible to add GSI after table creation

### Concurrency

DynamoDB provides conditional updates / deletes to achieve **optimistic locking**

### DynamoDB Accelerator (DAX)

- Seamless cache with no application re-writes, click on button
- micro-second latency
- solves hot key problem
- 5 minutes TTL for cache by default
- up to 10 nodes, multi-AZ, IAM integrated

### DynamoDB Streams

- Changes in DB (create / update / delete) ends up in Stream
  - stream of events can trigger Lambda for real-time processing or analytics
  - cross region replication
  - 24 hours data retention in Stream
  - streams keys only / new change / old image / new and old

### DynamoDB TTL

Automatically delete an item after expiry time (TTL, in epoch time)

- background delete task with no RCU / WCU costs
- may take 48 hours to fully delete items
- items that expire are also deleted from GSI and LSI
- activating DynamoDB Streams can help recover expired items

### Transactions

New feature to perform have a single transaction for multiple create / update / delete on

- multiple items
- across multiple tables
- Consumes 2x of WCU / RCU of standard write mode / read mode

### Security and Extra Features

- VPC endpoint for private subnet access
- IAM access control
- encryption at rest using KMS
- encryption in transit using SSL
- Point-in-time restore with no performance impact
- Global Tables = multi-region, fully replicated, high performance
- Amazon DMS helps to migrate from other DBs straight to DynamoDB
- Local DynamoDB for development purposes

## AWS API Gateway

- managed service
- handles API versioning, different environment (dev, prod, ...)
- handles security, API keys
- request throttling
- Swagger / Open API imports
- transform and validate requests and responses
- generate SDK and API specifications
- Cache API responses

### Deployment Stages

- API changes needs to be deployed, and deployed to Stages
- Each stages have its own config params, and can be rolled back
  - Stage variables are like environment variables
  - can be used to map to integration like Lambda function ARN, HTTP endpoints etc
  - therefore API of different stages can call different endpoints
  - can also be passed to lambda "context" object
  - Common use case is to use Stage Variable to point to Lambda Aliases, to invoke the right function for the right Stage
- Enable canary deployments (to split traffic to certain canary channel)

### Mapping Templates

Used to modify integration requests / responses

- rename parameters
- modify body content
- add headers
- Map JSON to XML
- uses Velocity Template Language (VTL)
- filter output results

### API Specifications

- Import / Export to Swagger
- Import / Export to OpenAPI spec
- generate SDK using Swagger file

### Caching

- Default TTL 300 seconds
- Cache per Stage, encryption option, capacity between 0.5 to 237 GB
- Override cache settings for specific methods
- Ability to invalidate cache immediately
- **Authorized clients** can invalidate / by-pass cache with request header "Cache-Control: max-age=0"

### Logging, Monitoring, Tracing

- CloudWatch Logs
  - log at Stage level
  - can override settings per API basis (log ERROR, DEBUG, INFO)
  - Log can contain request / response body
- CloudWatch Metrics
  - by Stage
  - detailed metrics
- X-Ray
  - tracing for end-to-end full picture

### Enabling Cross-Origin Resource Sharing (CORS)

- must be enabled to receive API calls from another domain
- OPTIONS pre-flight request must contain these headers
  - Access-Control-Allow-Methods
  - Access-Control-Allow-Headers
  - Access-Control-Allow-Origin

### Usage Plans & API Keys

Usage plans can be used to define:

- Throttling - overall capacity and burst capacity
- Quotas - requests per day / week / month
- Associated to API Stages

API Keys can be generated for each user client (your external customer), and can be associated to usage plans

API Keys usage can also be tracked

### IAM Permissions

For internal applications interacting with API Gateway within AWS cloud

- IAM policy attached to User / Role of client app
- client passes credentials in requests header using "Sig v4"
- API Gateway checks access policy with IAM

### Lambda Authorizers / Custom Authorizers

For external client usage, especially with OAuth / SAML / third party authentication

- API Gateway called by client with authorization token
- Gateway passes token to Lambda Authorizer to evaluate, and Lambda returns an IAM policy (will include authorization)
- Gateway can cache result of this authentication

### Cognito

- Client must first authenticate with Cognito
- CLient will pass authentication token to API Gateway in requests
- Gateway will verify the token with Cognito directly
- does not handle authorization, therefore authorization will need to be build into backend

## AWS Cognito

### Cognito User Pools

- User sign in functionality
  - creates a serverless database of users
  - MFA enabled
  - can enable Federated Identities (Facebook, SAML, ...)
- Integrates with API Gateway
  - returns JWT token after login for clients to forward to Gateway

### Cognito Identity Pools (Federated Identity)

- Provide AWS credentials to users to access AWS resources directly
  - login to federated identity provider (may / may not be Cognito User Pools)
  - forwards token to Cognito Federated Identity Pool (Cognito will verify with identity providers)
  - get temporary AWS credentials from Cognito
  - credentials come with pre-defined IAM policy stating their permissions, from STS
  - use-case: allow temporary access to write to S3 using Facebook login
- Integrate with Cognito User Pools as an identity provider

### Cognito Sync

- Synchronize data from device to Cognito
  - stores user preferences, configurations
  - requires Cognito Federated Identity Pool to work
  - supports offline mode and cross device sync
- Deprecated by AppSync

## Serverless Application Model (SAM)

- Framework for developing and deploying serverless applications
- Uses simple YAML template to generate complex CloudFormation
- Only 2 commands, package (which also zips and uploads to S3, and generate an output template) and deploy (which creates a stack based on output template)
- can use CodeDeploy to deploy lambda functions
- can run Lambda, API Gateway, and DynamoDB locally
- _Transform_ header indicates a SAM template
- Three helpers for resources, _Function_ (lambda), _Api_ (Gateway), _SimpleTable_ (DynamoDB)

{% include figure image_path="/assets/images/screenshots/aws-sam-flow.png" alt="" caption="Serverless Application Model Flow" %}
