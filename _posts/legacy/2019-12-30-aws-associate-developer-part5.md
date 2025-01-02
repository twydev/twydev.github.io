---
title: "Ultimate AWS Certified Developer Associate 2019 - Part 5"
toc: true
toc_label: "Chapters"
published: false
---

Part 5 of this course covers ECS, Fargate, Encryption Strategies, KMS, Systems Manager SSM, IAM Best Practices, CloudFront, Step Functions, SWF, SES, ACM, and other loose ends that are required for the exam.

## AWS Elastic Container Service (ECS)

### Docker

It is essential to have a basic understanding of Docker before using ECS.

- Docker is a software development platform to deploy apps
- Package apps into containers. Run on any OS.
- Predictable behavior, multi-language, easier to deploy and maintain.
- Docker images are stored in Docker Repositories (Docker Hub, AWS Elastic Container Repository ECR)
- More lightweight than VMs since there are no guest OS involved.
- Dockerfile builds Docker Images, which can be pushed to Repositories, or deployed as Containers.

### ECS Clusters

- Clusters are logical groups of EC2 instances
- EC2 instances runs a special AMI optimized for ECS
- Must configure file **/etc/ecs/ecs.config**
- EC2 instances runs an ECS agent, which registers itself with the ECS Cluster service
- ECS agents will be coordinating and launching the containers on the EC2 instances
- underlying technology is still Docker.
- Note: Security Groups apply to EC2 instances, not tasks / containers running within it

### Task Definitions

- Metadata in JSON format to tell ECS how to run a Docker Container, containing:
  - image name
  - port binding for container and host
  - memory and CPU required
  - environment variables
  - networking information
- For multi container environment, must NOT define host port.

### ECS Service

- Defines how many tasks should run and how to run them
- Ensures the desired number of tasks are running across EC2 instances fleet
- Can be linked to ELB / NLB / ALB if required

### ALB Dynamic Port Forwarding

- Needs to create ALB separately
- Dynamic Host Routing is a feature that comes with ALB
- Adding an ALB to an ECS Service needs to be done at creation, and will need you to recreate a Service for existing ECS Services.
- Needs to configure security group, allowing EC2 fleet and ALB to reach each other.

### Elastic Container Repository (ECR)

- Private Docker Repository by AWS
- Access controlled through IAM
- Run AWS ECR push and pull commands to manage images and interact with ECR

`$(aws ecr get-login --no-include-email --region [region])`

The inner command generates a Docker login command for ECR and the outer parenthesis executes the resultant command immediately. **This needs to be done whenever a new session needs to interact with ECR (push or pull image).** Next we need to build and tag the docker image, before finally pushing it to the ECR.

```bash
docker build -t demo
docker tag <original-name> <aws-account>.dkr.ecr.<region>.amazonaws.com/<original name>
docker push <new-name-from-prev-command>

# to pull image
docker pull <new-name-from-prev-command>
```

- To use ECR images, need to configure Task Definitions
- EC2 instances needs to have permissions to pull from ECR

## AWS Fargate

Without Fargate, to scale ECS clusters require manual launching and adding of EC2 instances to the cluster.

Fargate abstracts away infrastructure management, so that developers only need to manage Task Definitions. Fargate will scale EC2 fleet accordingly. Essentially making ECS serverless!

- Task Definitions need to be created with Fargate compatibility
- Similarly, Clusters and Services needs to be created with Fargate compatibility.
- Note: Fargate tasks / containers can have IAM roles.

### X-Ray Integration

For ECS Clusters (using EC2 fleet):

1. **X-Ray Container as Daemon.** Runs a container of X-Ray Daemon in each EC2 instance.

2. **X-Ray Container as a Side Car.** Runs the X-Ray Daemon in each application container instead.

For Fargate Clusters, we have to use the side car pattern to configure X-Ray Daemon within application container in Fargate Tasks.

> Important Note on Task Definitions:
> X-Ray Daemon port 2000, protocol UDP
> Set Environment Variable called AWS_XRAY_DAEMON_ADDRESS to the daemon defined above
> Link the daemon to the application container

### Elastic Beanstalk Single / Multi Docker Container Mode

- Requires a config file at source code root, **Dockerrun.aws.json**
- This mode helps to create ECS Cluster, EC2 fleet configured with ECS agents, Load Balancer, Task Definitions and Execution.
- Must separately pre-build Docker images in ECR

## AWS Encryption Strategies

### Encryption In Flight (SSL)

- SSL Certificates required (HTTPS)
- Ensures no Man In The Middle attack

### Server Side Encryption at Rest

- Server encrypts data received, and decrypts data before sending out
- Service will manage the encryption / descryption, using a data key it has access to

### Client Side Encryption

- Data is encrypted before sending to Server
- Server is not able to decrypt data it is storing
- Data retrieved from Server will be decrypted on Client-Side

## AWS Key Management Service (KMS)

- Control access to data by having AWS manage the keys
- Fully integrated with IAM for Authorization
- Integrated with many AWS services, with CLI / SDK support

### Overview of KMS

- Used for sharing sensitive information
- Data is encrypted by KMS using Customer Master Key (CMK), that will never be revealed to users
- CMK can be rotated for extra security
- Secrets can be encrypted by KMS, and be stored in code / environment variables
- **KMS can only encrypt up to 4KB data per call**, if larger than 4KB use envelope encryption
- to grant access to KMS:
  - Key Policy allows the user
  - IAM Policy allows the API calls
- Able to manage keys & policies to:
  - Create Keys
  - Rotation Policies
  - Disable / Enable Keys
  - see key usage (via CloudTrail)
- CMK options:
  - default CMK - free
  - user created keys in KMS / externally imported keys - \$1 per month
- every API call to KMS to encrypt / decrypt has a cost.

### Envelope Encryption

Encryption

- Client calls GenerateDataKey API to KMS
- KMS creates new data key, and also encrypts data key using CMK and sends over to Client
- Client use data key to encrypt data, and destroys key
- Encrypted data key and encrypted data are bundled as a final output

Decryption

- Client calls Decrypt API to KMS
- KMS receives encrypted data key, and decrypt using CMK
- Client receives plaintext data key to decrypt data

## AWS Systems Manager (SSM) Parameter Store

- secure storage of configurations and secrets
- serverless, scalable, durable, easy SDK, free
- version tracking of configurations / secrets
- configuration management using path and IAM
  - organized in a tree-like structure
- nofications with CloudWatch Events
- integrations with CloudFormation

## IAM Best Practices

- never use root credentials, enable MFA on root account
- grant least priviledge
  - never grant "\*" access to a service
  - monitor API calls denied explicitly by policies on CloudTrail
- never store IAM key credentials on other machines
- EC2 machines, Lambda functions, ECS tasks, should have their own role
- CodeBuild should have its own service role
- Cross Account Access should make use of STS AssumeRole API
  - define IAM role in target account
  - define which source account can access this role
  - source account user calls AWS Security Token Service (STS) to retrieve credentials
  - source account user impersonate role and access target account resources

### IAM Policies Evaluation Sequence

1. Authorization starts at DENY by default
2. Evaluate all policies related
3. If explicit DENY exists, evaluate to DENY
4. If explicit ALLOW exists, evaluate to ALLOW

### IAM Policies with S3 Bucket Policies Evaluation

AWS will evaluate the UNION of both IAM and S3 policies and provide an authorization.

### IAM Dynamic Policies

IAM policies allow the use of AWS variables to authorize access to resources dynamically. Example to allow respective users to only access folders in a bucket that are named after their user name, so instead of having one policy per user, having one dynamic policy will grant the same least priviledge access.

### IAM Policy Variants

**AWS Managed Policy**

- maintained by AWS
- updated for new services / APIs
- good for assigning standard power users / administrators

**Customer Managed Policy**

- best practice, re-usable
- version controlled with rollback, central change management

**Inline Policy**

- strict one-to-one relationship between IAM principal and policy
- deletion of IAM principal will also delete the policy

## AWS CloudFront

It is a Content Delivery Network (CDN) with 136 Points of Presence around the global, caching contents on the edge.

- Popular with S3 integration, but works with EC2 and LoadBalancer
- Can help protect against DDOS
- Provide SSL encryption (HTTPS) for incoming connections, and communications to applications
- Support RTMP protocol for videos and media

## AWS Step Functions

- JSON state machine to orchestrate serverless visual workflow
- for lambda functions but also works with EC2, ECS, on-premise servers, API gateway
- features sequential, parallel, conditions, timeout, error handling and etc.
- maximum execution time of 1 year
- can implement human approval workflows.
- Use cases: order fulfillment, data processing, web applications, any workflows

## AWS Simple Workflow Service (SWF)

- similar to Step Functions, but AWS seems to be stopping support for this service
- code runs on EC2, with 1 year max runtime
- concept of "activity step" and "decision step" and built-in "human intervention" step
- only use SWF when
  - we have external signals to intervene in the processes
  - we require child processes that return values to parent processes

## AWS Simple Email Service (SES)

Service allows you to send email using SMTP interface or AWS SDK.

Can also receive emails. Integrates with S3, SNS and Lambda. Uses IAM to control email sending permissions.

## Summary of Databases Available

- Relational Databases (RDS), largely used for OLTP
  - PostgreSQL, MySQL, Oracle ...
  - Aurora + Aurora Serverless
  - Provisioned database
- DynamoDB, NoSQL DB
  - managed, key-value and document store.
  - serverless
- ElastiCache, in-memory DB
  - Redis / Memcached variants
  - Cache Capabilities
- Redshift, largely used for OLAP
  - data warehouse, data lake
  - analytics queries
- Neptune, graph database
- Data Migration Service, DMS
  - to move existing DB to any of the AWS services

## AWS Certificate Manager (ACM)

- host public SSL certificates in AWS
  - either upload your own certificates to ACM
  - or let ACM provision and renew public SSL certificates (free)
- ACM is integrated with
  - LoadBalancers (including those provisioned by Elastic Beanstalk)
  - CloudFront Distributions
  - API Gateway
- ACM makes it easy to manage and replace expiring SSL certificates without disrupting operations

## AWS Exam CheatSheet

### ECS

- Task Placement Strategy
  - Binpack = consume all CPU from an instance first
  - Random = select random instance to run Task
  - Spread = evenly distribute Tasks based on certain key-value criteria
- Service-Linked Role: IAM role for ECS service itself to communicate with other AWS services.
- Fargate Launch Type
  - tasks placement default spreads across availability zones
- EC2 Launch Type
  - tasks placement customizable
  - cluster query language allows grouping tasks by attributes
- Used with Elastic Beanstalk
  - if need to provisioning of the resources, load balancing, auto-scaling, monitoring, and placing the containers across the cluster.

### CloudWatch

- Metrics granularity
  - Standard Resolution = 1 minute
  - High Resolution = 1 second
  - turn on **detailed monitoring** is only 1 minute
- Custom Metrics
  - use **PutMetricData** API
  - **--storage-resolution** options determines resolution. Default is 1 minute (standard)
  - required to find out about EC2 memory and swap usage
- Alarm
  - Period = interval of one data point
  - Evaluation Period = number of data point to evaluate
  - Datapoints to alarm = threshold in evaluation period
- Cloudwatch Metrics for RDS = hypervisor of servers running DB sending the metrics
- RDS Enhanced Monitoring = agent running on DB instance sends metrics (more precise)

### SAM

- Resources
  - AWS::Serverless::Function = lambda function
  - AWS::Serverless::LayerVersion = lambda function layer version
  - AWS::Serverless::Api = API Gateway
  - AWS::Serverless::Application = nested Application
- Transform
  - this section is used to specify resources using SAM syntax. Will be converted to CF syntax before deployment.

### CodeDeploy

- Deploys to:
  - EC2, on-premise
  - lambda
  - ECS
- Deployment behavior:
  - Blue/Green EC2: will use new instances. Therefore not available for on-premise
  - Blue/Green Lambda: uses traffic shifting. This is the default.
  - Blue/Green ECS: traffic shifted to new tasks set using load balancer. 
  - uses CodeDeploy agent (HTTPS) if using compute instance
- Deployment Traffic Shift
  - Canary = only two increments, specify first increment percentage and interval between second increment
  - Linear = fixed percentage and interval between increments
  - All-at-once
- AppSpec hooks
  - beforeAllowTraffic, afterAllowTraffic

### CloudFront

- Viewer >requests> CloudFront >requests> Origin
  - each way of communication maintains its own SSL encryption
- SSL Cert for CloudFront needs to first be managed by **AWS Certificate Manager (ACM)** or **IAM Certificate Store**

### DynamoDB

- Throughput
  - 1 write = 1 KBps
  - 1 read = 2x 4 KBps (eventual consistent)
  - 1 read = 1x 4 KBps (strong consistent)
  - 2 read = 1x 4 KBps (transactional)

- **projection-expression** = read returns some attributes in a record
- **condition-expression** = write on conditions
- **filter-expression** = filter which records to read
- Scan
  - to be avoided and use Query if possible
  - if not, reduce page size to slow down the read throughput using `limit` parameter
- **Sparse Index**
  - refers to index only used for a small subset of data
  - only gets updated if incoming record contains index sort key value
- **Local Secondary Index**
  - uses same partition key as main table, only different sort key
  - cannot be added after table creation
  - good for alternative sorting use-cases on same partition key
  - attributes fetched from base table
  - uses provisioned throughput of base table
- **Global Secondary Index**
  - can be entirely different attribute as partition as sort key
  - can be added after table creation
  - good for entirely different use-case using the same data in table
  - **Projected Attributes**
    - stored on the index itself, GSI always use this
  - uses its own throughput apart from the main table
- **Global Table**
  - reconciliation method is last-write-wins

- **ReturnConsumedCapacity** when write request
  - TOTAL = total consumed WCU across all indexes and tables
  - INDEXES = with subtotal breakdowns
  - NONE

### Elastic Beanstalk

- Deployment
  - all at once
  - rolling = hot swap existing services in batches
  - rolling with additional batches = rolling, with one extra new batch
    - only cause inconsistency, if first batch succeeds, but subsequent batches fails
  - immutable = deploy brand new services
  - blue/green = deploy brand new environment, then redirect traffic using CNAME
    - only swap traffic when succeed
- Configs
  - most configs are stored in `.ebextensions` directory

### AWS CLI

- pagination
  - CLI still retrieves all items, but makes more API calls in background to retrieve pages

### X-Ray

- Segment Document
  - default contains host, request, response, work done and issues occurred.
  - **subsegments** provide more details on work down by downstream services
    - **inferred segments** are generated to fill gaps for services that don't send segments
    - **X-Forwarded-For** header used to show source client IP
  - config the data collected
    - **annotations** = key-value pair. can filter traces by annotations in future. indexed for search.
    - **metadata** = key-value pair. any type of value. not indexed for search.
  - errors
    - **Error** = 400 client side
    - **Fault** = 500 server side
    - **Throttle** = 429 too many requests
- Grouping
  - traces can be grouped using filter expressions
- Sampling
  - used to select a representative subset of data
  - **reservoir size** = target sample per second before applying fixed rate.
  - **fixed rate** = applied to sample outstanding requests beyond reservoir size.
  - total sampling size per second = reservoir size + (incoming requests - reservoir size) * fixed rate
- Environment Variables
  - _X_AMZN_TRACE_ID = tracing header
  - AWS_XRAY_CONTEXT_MISSING = behaviour in event of missing tracing header
  - AWS_XRAY_DAEMON_ADDRESS = IP_ADDRESS:PORT 

### Lambda

- **Concurrent Execution** = (max number of invocation per second) x (duration of each execution)
- AWS default unreserved concurrency = 100 (minimum)
- Integration with DynamoDB / Kinesis
  - number of shards determines concurrency
- Invocation Type
  - RequestResponse = default. Synchronous.
  - Event = asynchronous.
  - DryRun = validate permissions to run.
- Alias = pointer to a specific lambda function
  - can use `routing-config` to route traffic between alias

### ElastiCache

- Memcached
  - simplest model
  - run one large node with multi-threading
  - scale out and in, remove and add nodes
  - cache objects such as a database
- Redis
  - richer set of features
  - more durable across reboots
  - single threaded. scales by having multiple instances.

### API Gateway

- Integration Types
  - HTTP = calls your own HTTP endpoint. Transforms input to custom HTTP request.
  - HTTP Proxy = pass through.
  - AWS = for lambda integration. Transforms input to custom lambda input.
  - AWS_PROXY = pass through.
- **504 error Gateway Timeout**
  - **INTEGRATION_FAILURE** server flaw in backend integrated app.
  - **INTEGRATION_TIMEOUT** occasional occurrence.

### S3

- **SSE-KMS**
  - regardless of KMS key id in request header, S3 will use KMS key id stated in policy.
- **Cross Region Replication (CRR)**
  - requires versioning

### CloudFormation

- AWS CloudFormation StackSets
  - Cross-AWS Account CF Stacks CRUD
- inline lambda coding
  - can be done by adding code string to property of lambda resource > Code > ZipFile

### Admin

- **AWS Organizations**
  - manage across AWS accounts
  - **Service Control Policy**
    - introduce cap on usage of AWS resources across accounts
- **AWS Secrets Manager**
  - like Systems Manager Parameter Store but with rotation feature
