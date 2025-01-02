---
title: "SQS Triggered Lambda Functions"
published: false
---

Important configuration and concepts to keep in mind when working with SQS and Lambda integration, which affects retries, processing time and message visibility, de-duplication, and idempotency design.

## SQS and Lambda Integration

Unlike a service that is polling from SQS, which needs to explicitly delete the message from the queue once it is done with processing, a Lambda function triggered by SQS does not need to do that. Here are some great stackoverflow Q&A, and some additional points to note:

- [Prevent SQS from deleting a message when Lambda function fails](https://stackoverflow.com/questions/56212311/how-to-prevent-aws-sqs-from-deleting-a-message-when-lambda-function-triggered-fa)
- [How Lambda polls from SQS](https://stackoverflow.com/questions/52904643/aws-lambda-triggered-by-sqs-increases-sqs-request-count)

## Only works for Standard Queue

Unfortunately this feature does not work for FIFO queue.

## Message VisibilityTimeout vs Lambda function processing time

It is important to ensure that the visibilityTimeout of the SQS message is configured to be larger than the expected processing time of the Lambda function.

If the visibilityTimeout is smaller than Lambda function processing time, the message will reappear in the queue and potentially be picked up by another handler, resulting in your system processing the same message twice. And eventually after multiple processing, the retry count will exceed the redrive policy threshold, and the message will be moved to the dead letter queue.

Lambda function configurations has a property `Timeout` which determines the maximum duration the function is allowed to run. AWS do not allow you to configure an SQS-triggered Lambda function with a larger Timeout than the message visibilityTimeout.

## Fail fast to guarantee re-processing of message

Since Lambda assumes that a successful function execution (no timeout and no error) means the message has been processed successfully, it will delete the message from the queue. 

If we want to ensure that the message reappears in the queue, our function must exit with an error. Therefore, we have to be careful with our error handling within our Lambda function and ensure that our function only terminates successfully when the message is completely process. We may have to resort to throwing errors explicitly in our function.

## How Lambda is triggered, under the hood

Lambda will make 5 concurrently long polling request to the queue. The maximum long polling window is 20 seconds. The pricing for SQS per 1 million requests is $0.40 Singapore dollars. For Lambda, it costs $0.20 per 1 million requests. To implement any simple event driven systems without having to worry about infrastructure, this SQS+Lambda combo is probably the most cost-effective option in the market.


*Note: At the time of writing, this is the interaction between AWS SQS and AQS Lambda.*
