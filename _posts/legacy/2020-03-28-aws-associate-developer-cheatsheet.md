---
title: "AWS Certified Associate Developer Cheat Sheet"
toc: true
toc_label: "Topics"
published: false
---

Realising that I am likely to fail the certification if I solely rely on a Udemy course, this is a compilation of notes from other reading materials, mostly from TutorialsDojo.

# Topics

## EC2

- basic API calls - RunInstances, DescribeInstances, TerminateInstances
- **EBS vs Local instance store** - AMI of EC2 instance can be loaded to either local instance store (ephemeral) or EBS (persistent), which acts a the root device and boot up from there. Using EBS allows additional API calls - StopInstances, StartInstances (preserves boot partition data)
- We can mount multiple EBS volume to an EC2 instance, but NOT vice versa.
- EBS snapshots can be taken without un-mounting, in realtime.
- **EC2 Auto Scaling** - ensure correct number of EC2 instances running in a group based on pre-defined conditions.
- **Billing** - billed per second, with minimum of 60 seconds.
- **Elastic IPs** - limited to 5 per AWS account.
- **User Data** - Scripts that will be run on instance at every new provision
- **Meta Data** - information about the EC2 instance itself that can be accessed from within the instance (hostname, network info, IAM role and credentials etc).

## ELB

- Application LB supports HTTPS termination (requires SSL certificate on LB)
- SSL certificate can be provisioned from AWS Certificate Manager directly, or privately acquired and uploaded through ACM or **IAM**
- Backend Authentication between LB and servers are not supported. Only encryption.
- Supported redirects: HTTP -> HTTP, HTTP -> HTTPS, HTTPS -> HTTPS
- **Target Group-level Stickiness** - ALB can ensure requests routed to target group remains in the group. **Individual Target Stickiness** - is a configuration of the target group.
- Network LB is used for high performance TCP / UDP load balancing, with support for long-running connections.

## AWS Auto Scaling

- AWS Auto Scaling can be used to scale multiple AWS services, but EC2 Auto Scaling is only for EC2
- Can only create target tracking scaling policies (use individual services' scaling configuration instead for scheduled scaling / step scaling)
- Scaling plans sets the target AWS resources, and metrics to track, and target values
- Predictive Scaling - based on historical traffic metric (EC2 only)
- Dynamic Scaling - based on utilization target metric
- **EC2 ASG** can set health checks to use either EC2 service health check on instances, or ELB health checks on target group.

## CloudWatch

- IAM cannot be used to restrict access to data of only specific services on CloudWatch
- **EC2 Basic Monitoring** - EC2 data available in 5 minutes interval
- **EC2 Detailed Monitoring** - Data available in 1 minutes interval
- **EC2 Metrics Defaults** - CPU, Disk Ops, Network and Traffic, CPU credits.
  - monitoring memory consumption requires custom metrics
- **CloudWatch Metrics**
  - Organized into **namespaces**
  - Each metrics can be sliced and diced using different **dimensions**
  - **Custom Metrics** 
    - collecting data from your own custom scripts, apps, or services
    - can be either standard resolution (1-minute, default) or high resolution (1-second)
    - custom metrics stored at 1-second resolution can be read at 1, 5, 10, 30 or 60 seconds aggregation.
    - set **StorageResolution** field in **PutMetricData** API to 1 (high resolution) allows 1-second data point.
    - **only custom metrics can have high resolution**.
- **CloudWatch Logs Agents**
  - by default sends log data every 5 seconds but can be configured
- **CloudWatch Alarm**
  - for custom metrics high resolution, can create **high resolution alarms** that alerts at **10-seconds or 30-seconds intervals**.
  - available alarm actions - SNS, SQS publish, stop/terminate EC2, execute Auto Scaling policies.
  - configuring alarms
    - select CloudWatch metric
    - choose evaluation period and statistical measure (avg, max, etc)
    - set threshold target value
    - set when to trigger alarm (greater than, equal, less than threshold value)
  - alarms are evaluated for every period even after the first alert.
- **CloudWatch Events**
  - supported services that emits events: EC2, Auto Scaling, Cloud Trail.
  - can also configure cron schedules to emit events.
  - configure actions when event matches a rule, such as
    - invoke a Lambda function, relay event to Kinesis
    - notify SNS
    - invoke built-in workflow
  - **Custom Events** can be emitted using PutEvents API

## AWS RDS / Amazon Aurora

- **Enhanced Monitoring** - gives system level metrics such as CPU, memory, file system, diskOps.
- **Backups in S3**
  - Automated Backups are offered by AWS by default, with 7 days retention, point-in-time recovery (will be deleted when original instance is deleted)
  - DB snapshots are user initiated, also offer point-in-time recovery (will be retained even when original instance is deleted)
  - recovery always restores to a **new DB instance**
- **VPC**
  - deployment requires subnet group - which contains subnets from multiple AZ, for fail-over and standby instances.
  - use DNS name has underlying IP may change.
- **master user account** - native database account to control access to the database
- **Multi AZ Deployment** - creates a Primary and a few Standbys in the other AZs for fail-over.
- **Read Replica** - uses DB engine's in-built replication.
- **Encryption in flight** - SSL, to enforce on DB and activate on client.
- **Encryption at rest** - using KMS. Can create from un-encrypted snapshot by first encrypting then restoring.

## Elastic Beanstalk

- Flexibility: select OS, EC2 instance, database & storage, AZ, load-balancer, CloudWatch integration
- **Environment Variables** - Exposes connection to DB to app, environment can be used to switch between test and prod
- **Deployment Process** - Controls Application, Application version, Environment name
  - deploy apps to environment, creates new version
  - apps can be promoted across environments
  - versions can be rolled back (can store up to 1000)
  - **Update Options**
    - all-at-once
    - rolling w/o additional batches (when one of the batch fails and terminates deployment. Requires manual trigger to redeploy)
    - immutable
    - blue/green deployment: deploy to new env, use route 53 to slowly transition traffic to new env.
- **Coding**
  - source must be zipped and uploaded to EB
    - EB will upload to each EC2 machines, then resolve dependencies
    - zipping dependencies together will optimize deployments
  - configurations are in various **.ebextensions/filename.config**, YAML/JSON format
    - modify setting in option_settings
  - commands available `eb create, status, health, events, logs, open, deploy, config, terminate`
- **Deletion** - most resources will be deleted with the env, unless deletion protection is turned on.

## ECS

- allows running of docker containers on a managed cluster of EC2 machines (only pay for EC2 instances)
- schedules placement of containers
- access to ELB, EBS, IAM
- **Compared to Elastic Beanstalk** - EB is a higher level abstraction. EB can deploy ECS without the user provisioning EC2 clusters, managing LB, Auto Scaling, container placements.
- **Design Pattern** encouraged to make each component of app into a container, and use Tasks to launch a set of containers using specific configurations.

- **Components**
  - **ECS Clusters**
    - logical grouping of EC2 instances
    - runs special AMI with ECS agents installed, which de/registers EC2 instance with ECS, and coordinates launching of containers. (agents can be configured using **/etc/ecs/ecs.config** in the instances.)
    - **Security Groups can only apply to EC2 instances**, not tasks / containers running within.
  - **Task Definitions**
    - JSON format metadata to configure the launch of Docker Container
    - includes image name, port binding between container and host, memory, CPU, env vars, network etc.
    - configure to use ECR images
    - **IAM Roles for Tasks**
      - first create a new role, by giving the desired permission to **Amazon EC2 Container Service Task Role** service role (essentially granting visibility of a role to this service)
      - then create new Task Definition, the new role will be available for selection.
  - **ECS Service**
    - defines number of tasks to run, how to run them
    - scales up / down the containers based on capacity requirements
    - react to health checks
    - automatically registering them with ELB
    - **ALB integration**
      - dynamic port forwarding needs separate ALB creation.
      - ALB needs to be added to ECS Service at service creation.
      - Security Group needs to be manually configured to allow ECS Cluster of EC2 fleet and ALB to reach each other.
  - **ECR**
    - Docker repository
    - `$(aws ecr get-login --no-include-email --region [region])`
    - `docker build -t demo`
    - `docker tag <original-name> <aws-account>.dkr.ecr.<region>.amazonaws.com/<original name>`
    - `docker push <new-name-from-prev-command>`
    - `docker pull <new-name-from-prev-command>`
    - EC2 instances needs permissions to pull from ECR (and probably security group access)
  - **ECS Service Roles in IAM**
    - EC2 role for Elastic Container Service: for instances to access ECS
    - Elastic Container Service: for ECS to access other AWS resources
    - Elastic Container Service Autoscale: allows Auto Scaling to access ECS
    - Elastic Container Service Task: allows Tasks to access other AWS resources
  - **Launch Types**
    - EC2: requires management of cluster
      - Has **Task Placement Strategy** (best effort)
        - Binpack - consume all CPU from an instance first
        - Random - select random instance to run Task
        - Spread - evenly distribute Tasks based on certain key-value criteria
      - Has **Task Placement Constraints** (binding)
        - distinctInstance - place each Task on a separate instance
        - memberOf - place Task on instance that satisfies **Cluster Query Language** expression (CQL allows grouping of EC2 instances by attributes like AZ, EC2 metadata)
    - Fargate: abstracts away the provisioning or management of EC2 fleet
      - Task placement spreads across all AZ by default



## Lambda

## API Gateway

## DynamoDB

## AWS X-Ray

- **Elastic Beanstalk**
  - supported platforms comes with X-Ray daemon, just need to enable
  - app code need to use X-Ray SDK to send data to daemon
  - daemon calls X-Ray API to send data to AWS. Requires IAM permission **AWSXrayWriteOnlyAccess** (obtained from EB instance profile).
- **Lambda**
  - already integrated and is running X-Ray daemon
  - when lambda function is invoked for a sampled request, X-Ray daemon will send data.

TODO: to be completed
