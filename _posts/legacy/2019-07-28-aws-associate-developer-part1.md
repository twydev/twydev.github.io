---
title: "Ultimate AWS Certified Developer Associate 2019 - Part 1"
toc: true
toc_label: "Chapters"
published: false
---

Part 1 of this Udemy course by Stephane Maarek. Covers the fundamentals of AWS services, including IAM, EC2, ELB, ASG, EBS, Route 53, RDS, ElastiCache, VPC, and S3.

## Exam Coverage

- 22% - Deployment
    - CICD, Beanstalk, Serverless
- 26% - Security
    - each service
- 30% - Development with AWS Services
    - Serverless, API, SDK, CLI
- 10% - Refactoring
    - AWS services best for migration
- 12% - Monitoring & Troubleshooting
    - CloudWatch, CloudTrail, X-Ray

## IAM - Identity & Access Management

### IAM Key Concepts

- **Regions** is a separate geographic area. (us-east-1)
- **Availability Zones** are isolated data centers in the Region. (us-east-1a)
- AWS consoles and resources are region scoped, except IAM and S3. Therefore permissions granted by IAM applies globally.

- **Root Account** is an AWS account used to sign up for the service.
- **Users** are IAM accounts used by physical persons.
- **Groups** are logical grouping of Users, to apply permissions to group for easy management.
- **Roles** are IAM accounts given to applications for internal usage within AWS
- **Never write IAM credentials in code**
- **Never use Root account for anything other than initial AWS set-up**

- **Federation** uses SAML standard to integrate with other enterprise Identity Provider to allow Single Sign-On.

### IAM Best Practice

- Delete Root account access keys
- Activate MFA on Root account
- Create individual IAM user
- Use groups to assign permissions
- Apply an IAM password policy
- Delegate administrative permissions to an "admin" user group, and use those IAM account to manage AWS from now.

## EC2 - Elastic Cloud Compute

The signature service offering of AWS. Understanding EC2 lays the groundwork to understanding cloud computing in general.

EC2 consists of the following capabilities:
- the virtual machines (EC2)
- storing data on virtual drives (EBS)
- distributing load across machines (ELB)
- scaling the services using an auto-scaling group (ASG)

EC2 instances have 5 distinct resources that determines their size:
- RAM (type, amount, generation)
- CPU (type, make, frequency, generation, number of cores)
- I/O (disk performance, EBS, optimization)
- Network (bandwidth, latency)
- GPU

Special mentions: 
- M instances are balanced. 
- T2/T3 instance are burst-able instances, which accumulates CPU credit if it is running at normal load, but will tap into burst credits to handle surge in computation, ensuring good CPU performance until credit runs out.
- T2 Unlimited is a new type that allows unlimited burst credit, which you will be billed after the fact.

### EC2 SSH Private Key

The SSH private key file generated from EC2 console will always encounter the error "Unprotected Private Key: permissions 0644 for file are too open" when you first open it from the local terminal. This means that other users on the machine are also able to access the file, potentially causing security leak. Solution is to run `chmod 0400 the-key-file.pem` to reduce the permission.

Using PuTTy on Windows machine requires the conversion of .pem file to .ppk file.

### Security Groups

Security group controls how traffic is allowed into and out of EC2. (It is like firewall rules).
- *specify access to ports*
- *authorized IP ranges (IPv4 adn IPv6)*
- *Control inbound network (default is all blocked)*
- *Control outbound network (default is all opened)*
- each group can be attached to multiple instances
- tied to a region/VPC combination
- exists outside of the EC2 instance (therefore not possible to troubleshoot from within EC2 if traffic from that instance has already been blocked)

**Specify another Security Group as Source/Destination**

This allows instances associated with the specified security group to access instances associated with current security group. This does not add rules from the source security group to current security group.

Usually, when encountered with network issues and unable to reach instances due to **time out**, it is likely due to misconfiguration of security groups. However, if error is **connection refused** then it is an application level problem.

### Elastic IPs

EC2 instances will be assigned a different IP every time it stops and starts. 

To use a fixed IP, you will need Elastic IP. **Limited to 5 Elastic IPs per AWS account**.

Try to avoid using Elastic IPs as:
- you are probably making poor architectural decisions
- you may use random IPs and register DNS name to it
- or use a Load Balancer without the need for fixed public IP.

### Installing Server

```shell
sudo su #assume superuser
yum update -y #update all packages using YUM
yum install -y httpd.x86_64 #install Apache httpd server
systemctl start httpd.service #starts it as a service using systemd
systemctl enable httpd.service #enable the service across machine reboot
```
At this point remember to allow HTTP access in the Security Group to reach the server from your local browser.

```shell
curl localhost:80 #test to see that server is working. do the same using your local browser.
echo "Hello World from $(hostname -f)" > /var/www/html/index.html #overwrite webpage with custom message and hostname
```

### EC2 User Data Script

- Allows us to bootstrap our EC2 instances.
- Script is only run once at the instance first start.
- Used to automate boot tasks such as installing updates, installing software, downloading common files, anything you want and etc.
- The script will run with machine root user identity.

The script is configured when provisioning a new EC2 instance in the console, under "Configure Instance Details".

Paste the script into the **User Data** text box. The first line is important to indicate the shell to be used.

```shell
#!/bin/bash
# install httpd (Linux 2 version)
yum update -y
yum install -y httpd.x86_64
systemctl start httpd.service
systemctl enable httpd.service
echo "Hello World from $(hostname -f)" > /var/www/html/index.html
```

Navigate to the website served by the Apache service of this new instance to see if User Data Script is working. Likewise, you may SSH to the new instance and inspect the index.html file to verify that the script has been executed.

### EC2 Instances Pricing Options

Instance type refers to sizing, like t2-micro and etc.

- **On Demand Instances** pay as you use, highest cost, but no upfront payment or commitment. For short-term, un-interrupted workloads where you cannot anticipate your application behavior.
- **Reserved Instances** pay upfront, commit between 1-3 years, reserve a specific instance type, discount up to 75% compared to On Demand. Recommended for steady usage applications (database).
    - **Convertible Reserved Instance** - is a Reserved Instance but allows changing of EC2 instance type during commitment. But the added feature reduces the discount to 54%.
    - **Scheduled Reserved Instances** - is a Reserved Instance scheduled to run on recurring time window.
- **Spot Instances** offer up to 90% discount compared to On Demand instances. Bid a price to get an instance, and as long as the current spot price is under the bid, you will get the instance.
    - instances are reclaimed with a 2 minute notification warning when spot price goes above bid price.
    - good for batch jobs, big data analytics, workloads resilient to failures.
    - not good for critical jobs or databases
- **Dedicated Hosts** offers physically dedicated EC2 server. Full control of EC2 instance placement, with visibility of underlying sockets, physical cores of hardware.
    - reserve for 3 years
    - more expensive
    - suitable for Bring Your Own License (BYOL) customers
    - suitable for regulatory compliance requirements
    - billed per host
- **Dedicated Instances** the hardware is still dedicated to you, but unlike Dedicated Host, you have no visibility of hardware. 
    - automatic instance placement
    - may share hardware with other instances in your account
    - on instance stop/start, the hardware may change.
    - billed per instance

### Pricing

EC2 instance prices (per hour) varies based on:
    - AWS region
    - instance type
    - instance pricing option
    - operating system

Billing is by the second, but with a minimum of 60 seconds first.

Final bill may include storage, data transfer, fixed IP public addresses, load balancing, and whatever additional services used.

But if instance is stopped, you will not be billed.

### Amazon Machine Images

Instead of using the default AMIs provided, you may use your own custom AMI for the following:
- pre-install packages needed, which is faster to boot than relying on EC2 User Data Script
- pre-install monitoring/enterprise software, such as Active Directory integration
- security concerns
- Control of maintenance and updates over time
- pre-install application for faster deploy when auto-scaling
- using other AMI that is optimized for specific purposes

**AMI are locked to a specific AWS region!**

### EC2 Exam Checklist
- know how to SSH into instance
- know how to properly use security groups
- know the fundamental differences between private vs public vs elastic IP
- know how to use User Data to customize instance at boot time
- know about custom AMI to enhance the OS
- know that EC2 instances are billed per second, and can be easily spined up/tear down

## ELB - Elastic Load Balancer

### Purpose of ELB
Similar to classic load balancers:
- spread load across multiple downstream instances
- expose single point of access via DNS to your application
- resilient to downstream failures
- perform regular health checks on downstream (can specify route for ALB)
- SSL termination (HTTPS) for website, and forward traffic as normal HTTP to downstream
- enforce stickiness with cookies so that client communicate to the same instance
- high availability across zones
- separate public from private traffic
- *AWS managed (guaranteed uptime, upgrades and maintenance, HA, easily configured)*
- *integrated with many AWS services*

### Application Load Balancer v2
- LB to multiple HTTP apps across machines (target groups)
- LB to multiple apps on same machines (containers)
- LB based on route in URL
- LB based on hostname in URL
- port mapping feature to redirect to dynamic port
- stickiness provided by ALB directly, and can be enabled at target group level, to ensure same requests go to same instance.
- supports HTTP/HTTPS/Web sockets protocol
- applications don't see client IP directly. ALB terminates client connection and start a new one with downstream.
    - client true IP is in the header **X-Forwarded-For**
    - **X-Forwarded-Port** and **X-Forwarded-Proto** contains port and proto information.

### Network Load Balancer v2
- forward TCP traffic
- millions of request per seconds
- support static or elastic IP
- less latency (100ms vs 400ms for ALB)
- used for extreme performance and hardly the default choice
- directly sees client IP

### Classic Load Balancer
- is deprecated
- all 3 LBs have health check capabilities
- only CLB and ALB supports SSL certificates and SSL termination
- all LB uses a static host name. Do not resolve and use the underlying IP.
- LBs can scale, but not instantaneously
- HTTP 4xx errors are client induced
- HTTP 5xx errors are application induced
    - 503 means LB error, either capacity issue or no registered target
- if LB cannot connect to application, check security groups

## ASG - Auto Scaling Group

### Purpose of ASG
- scale out more EC2 instances on load increase/scale in on load decrease
- ensure minimum/maximum number of instances running
- automatically register new instances to LB

### ASG attributes
It has a launch configuration which states the details of EC2 instances it will be provisioning:
- AMI, Instance Types, User Data
- EBS Volumes
- Security Groups
- SSH Key Pair

You can also configure min size/max size/initial capacity, network information, LB information and scaling policies.

### Auto Scaling Alarms
- CloudWatch monitors a metric and raise alarm on conditions
    - (if an ASG metric is selected, it will be computed for the overall ASG instances)
- based on the alarm, we can create scale out/scale in policies
- can also create alarm using custom metric sent from applications to CloudWatch

### ASG other details
- Scaling policies can be based on CloudWatch metrics, custom metrics, even schedules
- IAM roles attached to ASG will get assigned to EC2 instances
- ASG are free. You only pay for underlying instances
- Having instances under ASG provides extra guarantee
- ASG can terminate instances marked by LB as unhealthy

## EBS - Elastic Block Store
EBS Volume is a network drive you can attach to your instance while they run, ensures the drive persists even if instance terminates unexpectedly.

### EBS details
- network drive, so might have a bit of latency
- but is able to detach and attach to other instances quickly
- is locked to AZ (rather than region. To move volume across AZ requires you to create snapshot).
- have a provisioned capacity (in GB space and IOPS) which will be billed regardless of usage. You may increase capacity over time.

### EBS Volume Types
(not testing specifics)
- GP2 is general purpose SSD balancing price and performance
- IO1 is highest-performance SSD for low latency high throughput
- ST1 is low cost HDD volume for throughput intensive usage
- SC1 is lowest cost HDD for less frequent access

### Volume Resize
- can only increase EBS volume
- increase by size
- increase by IOPS is only for IO1 type
- after resizing, need to repartition the drive

### Snapshot
- snapshot only take the actual space of the data, not the entire EBS value
- used for backups
- used for migration (resizing down, changing volume type, encrypting, changing AZ)

### Encryption
Creating an encrypted EBS volume lets you benefit from the following:
- **data at rest encryption**
- **data in flight encryption** between the instance and the volume
- all snapshots are encrypted
- all volumes created from snapshot are encrypted
- is transparent to the user
- AWS guarantees minimal latency
- encryption leverages keys from KMS
- copying an un-encrypted snapshot allows encryption to create a new volume

### EBS other notes
- some EC2 instances use instance store physically attached to the machine
    - faster IO performance
    - on termination, instance store will be lost
    - cannot be resized
    - backups must be operated by user
- EBS should meet most use cases compared to instance store unless you need extreme performance
- EBS only attach to one instance at a time
- EBS backup uses IO so should be run during your application off-peak window
- by default Root EBS Volumes on EC2 instances will be deleted on termination, which you can disable this option.

## Route 53 - Managed Domain Name System

Route 53 is a global service not bounded by AWS regions.

### Types of DNS records

- A record: maps URL to IPv4
- AAAA record: maps URL to IPv6
- CNAME record: maps URL to URL
- Alias record: maps URL to AWS resource

### Route 53 usage

- public DNS and private DNS (within AWS VPC)
- DNS client load balancing
- limited health checks
- Routing policies:
    - simple
    - failover
    - geolocation
    - geoproximity
    - latency
    - weighted
- Alias record provides better performance (than CNAME) when routing AWS resources

## RDS - Relational Database Service

### RDS vs Hosting DB on EC2

- Managed Service with hassle-free benefits
    - OS patching
    - Continuous Backup with Point in Time Restore
        - Automated daily full snapshots.
        - capture transaction logs in real time to enable restore to any point.
        - 7 days retention (up to 35 days).
        - Manually triggered snapshot has no expiry in retention.
    - Monitoring Dashboards
    - Read replicas for improved read performance
        - replication is Asynchronous and data in replica will be eventually consistent.
        - replicas may be promoted to become separate DB instances.
        - application needs to explicitly connect to the read replicas.
    - Multi-Availability Zones for Disaster Recovery
        - Synchronous replication to a replica in a different AZ
        - both master and replica only expose **one single DNS name**. 
        - This allows automatic failover for the master DB.
    - Maintenance windows for upgrades
    - Horizontal and Vertical Scaling

But since it is a managed DB, we are not able to SSH into the instance to administrate it.

### Data Encryption

Encryption at rest can be enabled with AWS KMS using AES-256 encryption.

Encryption in flight can be enabled using SSL.

- first, DB needs to **enforce** SSL, through the RDS console (PostgreSQL) or through SQL execution (MySQL).
- next, clients needs to **connect** to DB using SSL, by obtaining the SSL Trust certificate from AWS, and specifying the option during connection.

### Security

- deploy RDS in private subnet
- control network access using Security Groups
- control permission to manage RDS using IAM
- control user access through traditional DB username / password login
    - grant IAM users access to DB is only supported for MySQL and Aurora.

### Aurora

- Aurora is a proprietary technology from AWS, compatible with PostgreSQL and MySQL drivers.
- Cloud optimized with better performance on AWS than other RDS DB.
- Supports incremental storage.
- More replicas, and faster replication than other RDS.
- Instantaneous failover and native High Availability
- but costs 20% more than RDS.

## ElastiCache - Managed Cache

Comes in **Redis** or **Memcached** variants.

- in-memory, high performance, low latency
- reduce load off database for read intensive workloads
- makes application stateless (by storing states in cache)
- Write Scaling using sharding
- Read Scaling using read replicas
- Multi AZ failover capabiity
- Managed OS maintenance, patching, optimizations, setup, configuration, monitoring, failure recovery and backup

### Use Cases

1. DB Cache
    - application first queries ElastiCache for data. 
    - Only queries DB if data is not found. 
    - Cache response from DB so that subsequent call is faster.
    - Requires a cache invalidation strategy.

2. User Session Store
    - user login to any application in cluster.
    - app stores user session in ElastiCache.
    - subsequent user access through a separate app instance retrieves the same user session.
    - allows app instances to remain stateless.

### Redis vs Memcached

Redis is more popular than Memcached due to better features.

- In-memory key-value store with super low latency and persistence across reboots.
- Multi AZ automatic failover for DR
- Support Read Replicas

### Cache Patterns

ElastiCache is used for read-heavy application workload (caching) or for compute-intensive workloads (memoization), and there are generally 2 patterns for implementing such architecture, and both can be implemented together.

**Lazy Loading**

1. application query cache first.
2. if cache miss, query from DB, then write DB response to cache.

This ensures that only relevant data is stored in cache, and failures in the cache will not be fatal.

However, every cache miss will incur 3 round trips, causing delay in the app. Data will also go stale and needs invalidation strategy.

**Write Through**

1. on update, app write to DB
2. app also writes to cache

This ensures that cache data is never stale. This incurs a write penalty (longer writing time for app) instead of read penalty.

However, data will be missing from cache until it is written / updated to DB (can be mitigated by implementing Lazy Loading in conjunction). Cache churn (data in cache may never be requested).

## VPC - Virtual Private Cloud

- **VPC** - are logically isolated network on AWS within a region. (Per Account Per Region).
- **Subnets** - VPC contains subnets. Each subnet must be mapped to an AZ in the region. Subnets in the same VPC can communicate with each other. (Per VPC Per AZ).
- **AZ** - we can have many subnets per AZ.
- **Public & Private Subnet** - this is a common set up to separate components in a software architecture.
    - Public subnets usually contain load balancers, static websites, files, public authentication layer.
    - Private subnets usually contain web app servers and DB.

- All new accounts come with default VPC.
- **VPN** - can be set up to connect to a VPC.
- **Flow Logs** - VPC comes with Flow Logs to monitor traffic in and out of VPC.
- Some AWS resouces cannot be deployed into a VPC.
- **Peering VPC** - we can set up VPC peering to make separate VPC looks like part of a single network (even across AWS account).

## S3 - Simple Storage Service

### Buckets

- Defined at **Region** level
- has **Globally** unique name. Not per Account, but globally across entire AWS.
- act as "directories" to store objects (files).
- Has naming convention restrictions:
    - no uppercase
    - no underscore
    - 3-63 characters long
    - not an IP
    - must start with lowercase letter or number

### Objects

- identified by their keys.
- keys represent the full path (/folder1/folder2/my_file.txt).
- the concept of folders is just a UI trick in the console. All folders are part of the key name.
- Max file size 5GB
    - more than 5GB requires "multi-part upload" to S3
- has Metadata (list of key-value pairs in text), can be system or user metadata.
- can be Tagged (unicode key-value pairs, up to 10 tags)
- has Version ID (if object versioning enabled)

For Objects with no public access,

1. the "Open" action from S3 console allows opening of file in browser, as it evaluates currently logged in AWS user permission.
2. however, accessing the same file using that Object's URL will result in an error since we did not allow public access.

### Versioning

- enabled at Bucket level
- same key overwrite of the same Object will increment version ID
- allows easy roll back, protecting against unintended deletes / modifications.
    - "Delete" works as a soft delete, which merely place a new delete marker version for the key.
    - deleting the delete marker version restores the Object.
- files not versioned prior to versioning enabled will have version "null"

### Encryption

1. **SSE-S3** stands for Server-Side Encryption by S3
    - Client sends Object to S3
    - Client request set header `"x-amz-server-side-encryption":"AES256"`
    - S3 receives Object and check request header
    - S3 creates an S3-managed encryption key
    - S3 encrypts Object using the managed key before storing it in bucket
2. **SSE-KMS** stands for Server-Side Encryption by Key Management System
    - works like SSE-S3.
    - encryption is done using key from KMS instead.
    - with more visibility of audit trail and user control over key rotations.
    - request must set header `"x-amz-server-side-encryption":"aws:kms"`.
3. **SSE-C** stands for Server-Side Encryption by Customer key
    - works like SSE-KMS.
    - encryption is done using key supplied by customer, generated from outside of AWS.
    - HTTPS **must** be used to make the request.
    - encryption key is provided in request header.
    - S3 will receive the key and Object, encrypts object, throws away the key, and store Object.
    - S3 do not store the key.
    - only the customer will have a decryption key to decrypt object.
4. **Client Side Encryption**
    - client uses library such as Amazon S3 Encryption Client.
    - client must encrypt object before sending request to S3.
    - client must decrypt object after retrieving from S3.
    - customer fully manage the keys and encryption cycle.
5. **Encryption In Flight**
    - basically uses HTTPS endpoint exposed by S3, instead of HTTP.
    - it is mandatory when using SSE-C pattern.
    - also called SSL / TLS.

Encryption can be configured on Bucket as a default, or on each Object uploads.

### Security

- **User Based** - through IAM policies
- **Resource Based**
    - through Bucket Policies
    - through Object ACL
    - through Bucket ACL (less common)

**Bucket Policy JSON**

- Resources: buckets and objects
- Actions: Set of API to Allow or Deny
- Effect: Allow / Deny
- Principal: The account or user to apply the policy to

This policy can be used to grant public access to the bucket, force objects to be encrypted at upload, or grant access to another AWS account.

The required JSON information can be generated in a web UI called **AWS Policy Generator**.

**Other Security Features**

1. Networking
    - S3 provides VPC endpoints, to allow other resources to connect to it within the VPC (without internet).
2. Logging and Audit
    - S3 access logs can be stored within S3 buckets (recommended to store in separate bucket from the data)
    - API calls can be logged in AWS CloudTrail
3. User Security
    - force MFA when user wants to delete Objects in versioned Buckets.
    - generate Signed S3 URLs that are only valid for limited time period.

### Website

- static website can be hosted on S3, accessible on www.
- accessed using URL: `<bucket-name>.s3-website-<AWS-region>.amazonaws.com` or `<bucket-name>.s3-website.<AWS-region>.amazonaws.com`
- getting a 403 (Forbidden) error likely means bucket policy did not allow public read access.

To give public access, in the Bucket Permission config we first needs to disable restriction for public access (set by default).

Then we need to create bucket policy that allows public access.

### CORS - Cross Origin Resource Sharing

- For web apps using S3 to fetch resources, CORS needs to be enabled since the web app is likely from a different origin.
- In S3 Bucket Permissions, CORS Configuration will allow a list of permitted origin to request for objects in this bucket.

(Reminder, Same-Origin Policy is enforced by the browser. Client will check with S3 bucket if the requesting web app is a permitted for CORS)

### Consistency Model

- Read after Write is consistent for NEW object PUTs
    - condition on no prior GET request on the same object. (the GET response before PUT request may still be cached and yet to be invalidated).
    - if there was a prior GET request, Read after Write is eventually consistent.
- eventual consistency for DELETES and PUT of EXISTING objects.
    - Read after a series of Writes on the same object may fetch an older version of the object.

### Performance

Traditional S3 usage recommends havng **Randomized Object Key Prefix** to force S3 to distribute your objects in separate partition (to improve transaction throughputs). Using dates as prefix is discouraged as it usually leads to sequential access in the same partition.

(probably not in exam) Throughput performance has significantly improved since July 2018, and randomized prefix is no longer necessary.

- For faster upload (file size > 100 MB), use multi-part upload. (for file size > 5 GB, this must be used)
    - parts will be parallelized in multiple PUTs
    - maximize bandwidth and efficiency
    - decrease time to retry (retry by parts).
- For improved read performance, accessed globally, use AWS CloudFront to cache S3 objects.
- S3 Transfer Acceleration uses edge location to upload a file destined for an S3 bucket in a separate region, for faster upload. Only endpoint needs to be changed.
- SSE-KMS may throttle S3 throughput as KMS is unable to keep up with the necessary encryption / decryption.
    - KMS usage limit needs to be adjusted.

### S3 Glacier - Long Term Archival product tier

S3 and S3 Glacier provides a SELECT feature, to retrieve data straight from objects.

- up to 80% cost saving and 400% performance improvements.
- runs SELECT SQL queries on the data in S3 / Glacier. Allows simple filtering on columns only (no sub-queries or joins).
- only needs to return a small subset of data (less rows and less columns).
- minimize network consumption.
- works on CSV, JSON, Parquet files.
- files can be compressed in GZIP or BZIP2.
- can only be used in code (no UI in S3 console).
