---
title: "AWS S3 Select - powerful, but tricky"
classes: wide
toc: true
---

S3 Select is a powerful feature that allows us to query a subset of data from within an S3 object using SQL expressions. However, I did not find this feature well documented. Here are my insights.

## Use Cases

S3 Select feature works nicely for our use case that fits the following descriptions:

- We need to update data files in S3.
- We are unable to identify the delta difference between updates easily due to the limitations of the data source.
- Updates occur frequently (15 mins interval).
- Data files contain large number of records, using databases with strong update consistency will be costly/reduce performance of overall system.
- Reading from the data file does not occur frequently (up to 6 times a day).
- For each read request, we only need to retrieve a small subset of records (150 records out of 50,000).

It sounds like S3 Select works perfectly for us. The write operation conveniently overwrites the data file in a single request, and read operations that occurs at low frequency allow us to query specific record in the file. We are paying for the Select query on demand instead of processing the data into a database on each update.

Here are some of the tricky issues I have faced working with the S3 Select feature.

## Promise-Based Approach (JavaScript SDK)

According to the [SDK documentation](https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/S3.html#selectObjectContent-property), only the _Async Iterator EventStream_ approach allows us to make use of the promise-ify feature of the SDK when calling S3 Select. If we would like to use Promises in other EventStream approaches, we would need to DIY by wrapping our client code in a Promise.

Unfortunately, the example provided in the documentation do not inform us how the data content may be retrieved from the EventStream. The use of Async Iterator is too sophisticated for me, and I have yet to wrap my head around how we may subscribe to the EventStream events using this approach. Eventually I settled with using Node.js EventStream approach:

```javascript
/* a TypeScript example */

/* Unfortunately this interface was not defined in the SDK */
interface SelectObjectContentEvent {
  Records?: AWS.S3.RecordsEvent;
  Stats?: AWS.S3.StatsEvent;
  Progress?: AWS.S3.ProgressEvent;
  Cont?: AWS.S3.ContinuationEvent;
  End?: AWS.S3.EndEvent;
}

/* Wrap EventStream subscription in a Promise */
new Promise((resolve, reject) => {
    this.s3client.selectObjectContent(params, (err, data) => {
        if (err) {
            console.error('S3 Client Error', err);
            reject(err);
        }
        if (!data) {
            console.error('S3 Data Object is Empty');
            reject(new Error('Empty data object'));
        }

        const records = [];
        let eventStream = data.Payload as StreamingEventStream<SelectObjectContentEvent>;

        eventStream.on('data', event => {
            if (event.Records) {
            records.push(event.Records.Payload);
            }
        });

        eventStream.on('error', (err) => {
            console.error('S3 Server Error', err);
            reject(err);
        });

        eventStream.on('end', () => {
            let bufferString = Buffer.concat(records).toString('utf8');
            /* string manipulation depends on s3 client query settings.
            My settings simply return individual JSON objects separated by commas,
            therefore I faced the issue of trailing comma, and the objects are not wrapped in a list */
            bufferString = bufferString.replace(/\,$/, '');
            bufferString = `[${bufferString}]`;
            try {
                const outputJson = JSON.parse(bufferString);
                resolve(outputJson);
            } catch (err) {
                reject(new Error('Unable to convert S3 data to JSON object'));
            }
        });
    });
});
```

## OverMaxRecordSize Error

When querying the S3 object, I faced the following error:

```shell
OverMaxRecordSize: The character number in one record is more than our max threshold, maxCharsPerRecord: 1,048,576
```

This was caused by my query requesting for too much data that exceeds the threshold. My data looks like this:

```json
{ "project_name": "project1", "completed": false }
{ "project_name": "project2", "completed": true }
```

And my query was simply:

```sql
select * from s3object s where s."project_name" = 'project1';
```

This issue puzzled me for a long time as I was unable to structure a SQL expression that overcomes the limitation. A guide to creating SQL expressions for S3 can be found in [AWS documentation](https://docs.aws.amazon.com/AmazonS3/latest/dev/s3-glacier-select-sql-reference.html).

I eventually found a [forum response](https://forums.aws.amazon.com/thread.jspa?messageID=933637&tstart=0) that provided the crucial information missing from AWS documentation.

> Definition of what's considered as a record depends on the file type you are querying. For a CSV file, it will be defined by the presence of a record delimiter (the default delimiter is a line break) that you can specify in your query parameters. For JSON file though, it will be defined by the path of your FROM clause. If your query looks like "FROM S3Object[*].array" then each element of your array will be considered as a record. But if it looks like "FROM S3Object" then the root element of your document will be considered as a single record.

I am not sure if the information in the response is fully accurate, but it did give me some ideas:

- CSV file implicitly contains additional information, which is the schema of the data.
- JSON file do not contain schema information, and the S3 Select query feature do not make any assumptions on the schema of the data. It probably takes a Schema-on-Read approach to handling query.
- Therefore, in order to execute my query, which is to filter for records that meet the condition `s."project_name" = 'project1'`, S3 Select feature must parse the entire JSON file as a record in order to apply the schema, therefore exceeding the query threshold.
- I suspect that this was not the case for querying from CSV file. With a well defined schema, the query request can be satisfied using some kind of index, or optimization that do not require us to parse the entire file.

Therefore, when I change my data file format from JSON to CSV, the exact same query worked perfectly.

```csv
"project_name","completed"
"project1","false"
"project2","true"
```
