---
title: "Ultimate AWS Certified Developer Associate 2019 - Part 2"
toc: true
toc_label: "Chapters"
published: false
---

Part 2 covers software development using AWS. Topics include the use of AWS CLI and AWS SDK, tools for troubleshooting IAM roles and policies, using IAM permissions with EC2, Elastic BeanStalk, and AWS CICD suite (CodeCommit, CodePipeline, CodeBuild, CodeDeploy, CodeStar).

## AWS CLI - Command Line Interface

### Set Up AWS CLI

1. Install AWS CLI from official AWS website to always get the most updated installer. Install instructions differ for each operating system.

2. Check AWS CLI version
```
aws --version
```

3. Configure credentials that will be used to access the AWS. After configuring, the configuration files can be found in `~/.aws` or the Windows equivalent of user's home directory.
```
aws configure
# enter access key
# enter secret key
# enter default region
```

### Credentials Profile

We can manage our credentials for multiple AWS accounts by storing them as different profiles

```bash
aws configure --profile myProfile1

# we can now run AWS CLI command using specific profiles
aws s3 ls --profile myProfile1
```

### AWS S3 CLI

1. Use AWS S3 CLI document to reference the available commands.
2. Tip: use `--help` whenever we are unsure of the correct command usage.
3. Commands covered in tutorial:
```bash
# list all buckets
aws s3 ls 
# list all objects in a bucket
aws s3 ls s3://<bucket> 
# copy files to local machine
aws s3 cp s3://<bucket>/<object> localFileName
# make bucket
aws s3 mb s3://<newBucket>
# remove bucket
aws s3 rb s3://<bucket>
```

### AWS CLI on EC2

As a good security practice, **NEVER STORE PERSONAL CREDENTIALS ON EC2 INSTANCES**. Follow the best practice recommended by AWS:

- To use AWS CLI on EC2, we can attach IAM Roles to the EC2 instances.
- Policy can be used to authorize precisely what the EC2 IAM Role can/cannot do.
- EC2 attached with IAM Roles can use these access without additional configurations.

1. Create IAM Role. Select **AWS Service**, **EC2**, attached the desired policies and give this role a unique name.
2. In the EC2 instance settings, attach IAM Role and select the newly created Role.
3. SSH into EC2 instance (AWS CLI is already installed) and test out the AWS CLI commands with the newly granted permissions.
4. NOTE: Attaching IAM Role to EC2 instance may take a while to take effect as the changes are replicated through AWS infrastructure.
5. NOTE: Each EC2 instance can only be attached with one role at a time.

## IAM Roles and Policies

- After creating an IAM Role, Policies can be attached to it.
- We may also create Policies directly on top of Roles (called inline policy) 
- However this is not recommended as it is hard to manage. (the policies created cannot be accessed by other roles)

### AWS Managed Policies

- These are policies pre-built by AWS and available for use by default.
- Policies statements has 3 components:
  - Effect: whether actions are allowed or denied
  - Action: the actions that can be carried out (e.g. Get, List)
  - Resource: the resource that will be affected by the action
- Actions can use wildcard character to give access to all actions (e.g. Get*, List*)
- Specific resources can be specified using ARN.
- Policies can be created using Visual Editor or JSON Editor (under IAM console)
- We may also use another tool outside of IAM console, called AWS Policy Generator.

Creating policies on IAM console has an added advantage of viewing policy usage, and policy version history.

### AWS Policy Simulator

Can be used to test actions on selected roles and policies, to see if those actions are currently allowed.

This tool will help to save many hours of debugging application permission issues.

### CLI Dry-Run

- Some AWS CLI commands contains `--dry-run` options
- This allow us to test if the current permissions to perform the command is available, instead of actually running the commands (as some commands are expensive to test, such as creating a new EC2 instance
)

```bash
# dry run creation of new  EC2 instance
aws ec2 run-instances --dry-run --image-id <ami-id> --instance-type t2.micro
```

- If we do not have permission to perform this action, we will get a long cryptic error message that needs to be decoded to make sense (via AWS Security Token Service)

```bash
aws sts decode-authorization-message --encoded-message <message>
```

- STS decode authorization message permissions need to first be granted to the Role that is attempting to run this command.

## EC2 Metadata

- All EC2 instances are able to access their own metadata from within the instances, without using any IAM Role.
- EC2 metadata are accessible from **http://169.254.169.254/latest/meta-data**
- This allows us to retrieve the attached IAM Role name (part of the instance metadata), but not the associated IAM Policy (as such retrieval would require an authorized IAM Role).
- Metadata = information about EC2 instance
- Userdata = launch script of EC2 instance
- Being able to access an instance's own metadata will be helpful for automation in future.

```bash
# query from within EC2 instance

# shows all available API versions
curl http://169.254.169.254/
# use latest version and list all available metadata
curl http://169.254.169.254/latest/meta-data/

# show security credentials information
curl http://169.254.169.254/latest/meta-data/iam/security-credentials/<iam-role-name>
```

- the last API query will return security credentials information of the currently attached IAM Role.
- It will reveal the access key and secret key that the EC2 instance is using.
- These are temporary values that will expire and be refreshed, to ensure security. But these are all handled by EC2 and IAM service under the hood, and the user do not need to do anything other than attaching the Roles and Policies.

## AWS SDK

- The SDK is used for our application code to access AWS.
- Fun Fact: the AWS CLI is actually a wrapper around the boto3 SDK (a.k.a the Python SDK).
- NOTE: if region is not configured when using SDK, the default will be us-east-1.
- For security credentials when using SDK, it is recommended to use the **default credential provider chain**:
  - the credentials file (~/.aws/credentials) for local machine
  - instance profile credentials using IAM Role (for EC2 instances or AWS services)
  - Environment Variables (export AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
- **NEVER STORE AWS CREDENTIALS IN SOURCE CODE**
- Exponential Back-off retry is included in most SDK API as a retry mechanism.

## EB - ElasticBeanStalk 

### EB Overview

EB abstracts away a lot of the non-essential concerns from the developer.
- Managing Infrastructure
- Configuring multiple databases, load balancers on each deployment
- Scaling concerns
- Configuration in multiple environments and versions

EB is a managed service:
- make use of underlying AWS resources (EC2, RDS, ELB, ASG etc.)
- instance configuration/OS is handled by EB
- deployment strategy configurable but performed by EB
- developers only need to focus on the application code
- comes in 3 architecture model: 
  1. Single instance deployment (good for dev)
  2. LB + ASG (good for web app)
  3. ASG only (good for workers)

### EB Deployment Lifecycle

EB has 3 components:
- Application
- Application version
- Environment name

General Deployment process:
- We can deploy applications to environment
- Each deployment will be assigned a version
- We can promote a deployed app to the next environment
- We can rollback to previous app version
- We have full control over life cycle of our environments

### EB Environment

When creating an environment, we have the freedom to set the configuration for each of the following areas:
- Software (platform version, AWS X-Ray, server)
- Instances (EC2 config)
- Capacity (toggle ASG)
- Load Balancer
- Rolling Updates and Deployment
- Security (IAM roles, VM key pair)
- Monitoring
- Managed Updates
- Notifications
- Network (for VPC)
- Database (RDS created here is tied to this env, and will be lost if env gets deleted)
- Tags

### EB - Single Instance Deployment Mode

- Great for development
- one EC2 instance
- one Auto-Scaling Group
- one DNS name mapped to one Elastic IP
- one database if provisions
- one Security Group for EC2 (and one for DB)
- all the above within one Availability Zone

### EB - High Availability with Load Balance Deployment Mode

- Great for Production
- one Auto-Scaling Group
- spans multiple Availability Zones
- one EC2 instance(s) per AZ
- each EC2 instance has one Security Group
- database also spans multiple AZ (Master-Read Replica set up)
- one Load Balancer communicates with ASG
- Load Balancer exposes one DNS name, which will be the DNS name for the EB app.

## EB - Deployment Update Options

### All At Once

EB will stop all of the old versions of the application at once, and deploy the newer version.

**Pros and Cons:**

- Fastest deployment
- App has downtime
- Great for quick iterations in dev env
- No additional cost

### Rolling

EB will stop a part of the application (depending on configured bucket size) and upgrade those instances into the newer version, before moving to the next bucket to update.

**Pros and Cons:**

- Long deployment (especially for large EB with many instances)
- No downtime
- App running at lower capacity during deployment
- App will be running both versions simultaneously
- No additional cost

### Rolling with additional batches

Same as Rolling option, but EB will first create one bucket size (depends of configuration) of additional new app version instances, before starting the Rolling update. After the update is complete, the additional bucket of instances will be terminated.

**Pros and Cons:**

- Longer deployment
- App is always running at capacity
- App is running both versions simultaneously
- Small additional cost (due to the extra bucket)
- Good for prod env

### Immutable

EB will deployment an entire set of new version instances on a temporary ASG (EB will make sure health check on the first instance passes before proceeding), then migrate all the newer version instances into the current ASG, then terminate the older instances and delete the temporary ASG.

**Pros and Cons:**

- Longest deployment
- Zero downtime
- High cost due to double capacity running
- Quick rollback in case of failures (just terminate new ASG)
- Great for prod

### Blue / Green deployment

This is not a feature of EB, but an approach which is possible to implement on EB.

**Steps to Blue / Green deployment (often manual)**

1. Have an existing prod environment (Blue)
2. Create a new environment (Green)
3. Deploy newer version of app in Green env
4. Validate and test the app, rollback by simply terminating the Green env without affecting the Blue env
5. Use Route 53 to set up weighted traffic redirect policies, to push a little bit of traffic to Green env, and test it to ensure it is working.
6. Gradually transition traffic from Blue to Green env. Or use EB to swap the URLs to instantly cut over.
7. Finally, all traffic is on Green env, and terminate the Blue env.

## EB - Advanced Concepts

- Code uploaded to EB must be zipped.
- All params set in the UI of EB console can also be configured in code.
- Config files must be located in **.ebextensions/** folder, in the root directory of the source code.
- files must have **filename.config** extensions
- files must be in YAML / JSON format
- modify default settings using option_settings
- able to add resources to EB (e.g. RDS, ElastiCache, DynamoDB)(these added resources will get deleted if env gets deleted)

To simplify the management of EB, we can use CLI tool called EB CLI. Available commands are 
`eb create, status, health, events, logs, open, deploy, config, terminate`

EB CLI is useful for automated deployment pipeline

Under the hood, ElasticBeanstalk uses CloudFormation. Creation of an EB env creates a CloudFormation stack. CloudFormation will perform all the heavy lifting of deploying resources.

When an app is uploaded as zip file, it is uploaded to each EC2 machine in the EB app. Each EC2 machine will resolve dependencies for the source code. 

Packaging the dependencies together with source code in the zip file will help optimize long deployments.

## EB - Exam Tips

1. Configure BeanStalk with HTTPS
  - load SSL cert on to the Load Balancer
  - can be done from EB console > Load Balancer config
  - can be done from source code > .ebextensions/securelistener-alb.config
  - SSL cert can be provisioned using AWS Certificate Manager or CLI
  - Must configure Security Group to allow incoming port 443

2. BeanStalk redirect HTTP to HTTPS
  - can configure instance in the code
  - can configure App Load Balancer with a rule (but make sure health checks are not redirected)

3. BeanStalk Lifecycle Policy
  - EB can store at most 1000 versions of an app
  - Lifecycle Policy can be configured to purge older versions
    - can based on time (oldest versions are removed)
    - can based on space (when you have too many versions)
  - Versions in use will not be deleted
  - Ability to retain source bundle in S3 even if version gets deleted.

4. Web Server vs Worker Environment
  - worker env is for more intensive workload that takes longer time to complete
  - offload tasks from web server to worker env is a common two tier architecture
  - periodic tasks for workers can be defined in cron.yaml

5. RDS with BeanStalk
  - Provisioning RDS with EB ties the RDS instance to the env (convenient for dev)
  - But for production, ideal way is to separately provision RDS, then provide DB connection string to the EB app
  - Steps to migrate from RDS coupled in EB to standalone RDS:
    1. take RDS DB snapshot (just in case)
    2. enable deletion protection in RDS
    3. create new EB env without RDS, point to existing old RDS
    4. perform blue/green deployment to migrate to new EB env
    5. terminate old env (old RDS will remain thanks to protection)
    6. delete CloudFormation stack manually (stuck in DELETE_FAILED state due to protected RDS)

## CICD 

### CICD on AWS Overview

Components provided by AWS:

- CodeCommit = storing of code
- CodePipeline = automating pipeline from code to EB
- CodeBuild = building and testing code
- CodeDeploy = deploying code to EC2 fleet

5 steps of CICD:

1. Code - AWS CodeCommit
2. Build - AWS CodeBuild
3. Test - AWS CodeBuild
4. Deploy - AWS EB or AWS CodeDeploy
5. Provision - AWS EB or EC2 Fleet / CloudFormation

Every steps can be orchestrated with AWS CodePipeline

{% include figure image_path="/assets/images/screenshots/aws-cicd-overview.jpg" alt="" caption="AWS CICD Overview" %}

### AWS CodeCommit

- private Git repository
- no size limit (scales seamlessly)
- fully managed, highly available
- code stored in AWS cloud account = increased security and compliance
- security features (encrypted, access control)
- integrated with Jenkins, CodeBuild, other CI tools

**Security**

- Authentication
  - SSH keys from IAM user console
  - HTTPS through CLI Authentication / generated HTTPS credentials
  - MFA
- Authorization
  - IAM policies
- Encryption
  - at rest = repositories encrypted using KMS
  - in transit = SSH / HTTPS
- Cross AWS account access
  - IAM role and AWS STS (with AssumeRole API)

**Notifications**

- Triggers notifications on SNS / Lambda:
  - deletion of branches
  - push to master branch
  - notify external build system
  - trigger lambda to perform codebase analysis
- Triggers CloudWatch Event Rules:
  - pull request updates (create / update / delete / comment)
  - commit comment events
  - CloudWatch Event Rules subsequently trigger SNS topic

### AWS CodePipeline

- Continuous Delivery Visual Workflow
- Sources can be from GitHub, CodeCommit, S3 etc.
- Build can be run by CodeBuild, Jenkins etc.
- Load Testing can be run by 3rd party tools
- Deploy can be through CodeDeploy, Beanstalk, CloudFormation, ECS etc.
- It is recommended to trigger CodePipeline through CloudWatch event instead (if using CodeCommit) for efficient performance

**Stages in Pipeline**

- Each stage can have sequential / parallel actions
- Each stage can have multiple action groups within
- Configure manual approval (one of the action option) at any stage

**Artifacts passing**

- Each stage can pass artifacts as output to be processed by the next stage
- Artifacts are passed through S3

**Troubleshooting**

- State changes in the pipeline generate AWS CloudWatch Event
- Event can trigger SNS notification
- Detailed information on pipeline failure can be obtained from Console
- AWS CloudTrail can be used to audit AWS API calls made by pipeline

Tips: If pipeline cannot perform certain actions, it is usually a problem with insufficient permission of IAM Service Role attached to pipeline

### AWS CodeBuild

- Fully managed build service (continuous scaling, no build queue!)
- Alternative to other build tools like Jenkins
- Pay for usage (build time)
- Uses Docker under the hood (AWS managed imsages available), so we are able to build with our own custom base Docker images
- fully integrated with AWS security
  - KMS to encrypt build artifacts
  - IAM to manage build permissions
  - VPC to secure the network
  - CloudTrail to audit API call logs

**Usage**

- build instructions can be provided in source code (*buildspec.yml* at root of code)
- build output logs to S3 and CloudWatch Logs
- able to use S3 to cache files to be shared by muliple builds (e.g. dependencies)
- metrics available to monitor CloudBuild statistics
- CloudWatch Alarms can be used to detect build failure and trigger notifications
- CloudWatch Events / AWS Lambda as a Glue
- trigger SNS notifications
- able to reproduce CodeBuild locally to troubleshoot in case of errors
- able to define CodeBuild in either CodePipeline or CodeBuild (which may cause conflict)

**Supported Environments**

- Java, Ruby, Python, Go, NodeJS, Android, .NET Core, PHP
- or use Docker to extend for any environment

**BuildSpec**

- must be at the root of project source code
- Define environment variables
  - plain text
  - secured secrets using SSM parameter store
- Phases of build (to execute commands)
  - install: install dependencies
  - pre build: execute final commands before build
  - build: actual build commands
  - post build: finishing touches (e.g. zipping of files)
- Define Artifacts, what to upload to S3 as output (encrypted using KMS)
- Define Cache, for files to store in S3 to speed up future builds

**Local Troubleshooting**

- Requires Docker to run CodeBuild locally using CodeBuild Agent

### AWS CodeDeploy

- We want to deploy application to many EC2 instances
- but not using Elastic Beanstalk
- CodeDeploy is an alternative to Ansible, Terraform, Chef, Puppet etc.
- It is a managed service

**Usage**

- Each EC2 machine must be running CodeDeploy Agent
- Agent continuously poll for work from CodeDeploy
- Developer commit new source code to repo, with *appspec.yml* file at the root of project
- Developer trigger CodeDeploy
- Agent polls and receive repo address, and pulls code
- EC2 will run deployment instructions based on spec
- Agent will report success / failure of deployment

**Features**

- EC2 instances can be grouped by deployment groups (dev / test / prod)
- Can be chained into CodePipeline and use artifact from previous stage
- Can work with any application and auto scaling integration (more flexible than Elastic Beanstalk, but therefore more complex)
- Able to have Blue / Green deployment (only for EC2 instances, not on-premise)
- Supports Lambda deployment

Note: CodeDeploy does not provision resources

**AppSpec**

- File Section: how to source and copy files from repo
- Hooks: sets of instructions to perform (in sequence)
  1. ApplicationStop
  2. DownloadBundle
  3. BeforeInstall
  4. AfterInstall
  5. ApplicationStart
  6. ValidateService (crucial)

**Configs**

- deploy one instance at a time. stops on failure
- deploy half of all instances at a time
- deploy all. quick but will require application downtime
- deploy custom. e.g. number of healthy host > 75%
- Failures:
  - instances stays in *failed state*
  - subsequent new deployment will deploy to failed state intances first
  - rollback by deploying an older deployment
  - or rollback by enabling automated rollbacks
- Deployment Targets
  - can be EC2 instances with Tags
  - can be an ASG
  - can be a mixed of ASG and Tagged instances
  - can have customization in scripts
- Deployment strategy is the same as Elastic Beanstalk (EB probably used CodeDeploy under the hood)
  - In Place Deployment
  - Blue Green Deployment

### AWS CodeStar

- integrated solution that provides a wrapper around:
  - GitHub, CodeCommit, 
  - CodeBuild, CodeDeploy
  - CloudFormation, CodePipeline
  - CloudWatch
- quickly creates CI/CD-ready projects
- issue tracking integrations e.g. GitHub issues / JIRA
- integrate with Cloud9 to obtain a web IDE (not all regions)
- single dashboard to view all components
- only pay for underlying resources
- limited customization of underlying components

