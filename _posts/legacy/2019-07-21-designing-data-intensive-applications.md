---
title: "Designing Data-Intensive Applications"
toc: true
toc_label: "Chapters"
---

*Work in progress.*

Explores the tradeoffs of data storage and processing technologies available today, and the hard problems facing distributed computing. Key topics include data structures of databases (B-Trees and LSM-Trees), Replication, Partitioning, Transactions, Consistency, Consensus, MapReduce, Stream Processing, and Event Driven.

___

# Designing Data-Intensive Applications

**Designing Data-Intensive Applications** *The Big Ideas Behind Reliable, Scalable, and Maintainable Systems* - Martin Kleppmann, 2017
{: .notice--primary}

I like the opening quote from Alan Kay. While we dive into a mad rush to implement new technologies, we forgot about the fundamental principles that make good software and applications. This book aims to help everyone understand which kind of data storage and processing technology is appropriate for which purpose. The focus is on software architecture. Very important in our current industry filled with buzzwords and hype.

## Foundations of Data Systems

### Chapter 1: Reliable, Scalable, and Maintainable Applications
This book focus on three primary concerns:

#### 1. Reliability
The system should continue to work correctly (performing the correct function at desired level of performance) even in the face of adversity (hardware or software faults, and even human error). The trend for modern software is to be resilient to hardware faults at the application level, such as ability to quickly scale up and down, or migrate instances without data loss.

#### 2. Scalability
The system should have options to keep performance good even when load increases. The key metrics here are load and performance, which are measured using different parameters depending on the software. Could be requests traffic vs response time, or data volume vs throughput, and etc.

#### 3. Maintainability
Making life better for engineering and operations team. General rules are to use good abstractions, reduce complexity of software (does not mean weak features), make it easy to extend and modify for new use cases. While the software is operating, everyone should have easy and good visibility of system health to address issues before failures.

> Good operations can often work around the limitations of bad software, but good software cannot run reliably with bad operations

### Chapter 2: Data Models and Query Languages

#### Relational Model

Classic RDB. Data are stored in rows of records. Sequence of record usually do not matter (due to external structure like indexes). Tables can have relations to other tables.

**PROS**

- Great for representing data with fixed relation / schema.
- Normalized data provides consistent values across all references, avoiding ambiguity, ease of updating only one record to affect all references
- Better at searching since indexes are easy to set up and maintain
- Strong at high performance joins, good for modelling many-to-one relationships (many users live in one country), and many-to-many relationships (many users, each live in a different state)

**CONS**

- Objects in applications require awkward translation layer to convert it to relational data models (despite help of ORM)
- Schema-on-write, therefore provides less flexibility to store evolving data, or heterogeneous data. (it may be impractical to store every object as a table)
- Modern RDB supports high performance DB migration to step up or step down the schema, but may still take a long time for large DBs.

#### NoSQL Model / Document Model

Think JSON. Storing the entire JSON into a data storage.

**PROS**

- Data is represented with closer resemblance to Objects in applications
- Greater scalability than RDB, with higher write throughput (due to sequential writes instead of splitting into tables)
- In some cases, better locality provide higher read performance as well.
- More dynamic and expressive data model. Easy to change the data model as software features and requirements evolve. (some implementations offer schema-on-read).
- Great at modelling one-to-many relationships (e.g. one user has many assets, but different users have different number and types of assets, which may have sub-assets)

**CONS**

- Not so great at joins. Requires executing scans and loops to match data (although some implementations offer join features).
- Typically requires data to be denormalized.
- Not good at searching across all entities.

#### What to choose?

As a general rule, select the data storage that helps you build simpler and cleaner applications. And this really depends on your features and use-cases, there is no one-size-fit-all, and some have adopted a polyglot persistence strategy.

#### Convergence

It seems like RDB and NoSQL technologies are converging, as the platforms release increasingly similar features.

- RDB like PostGreSQL and MySQL allows storing and querying of fields nested in JSON.
- RethinkDB supports joins. Some MongoDB emulates join by performing the loop on the client side.

#### Query Languages

- Declarative languages (e.g. SQL) are more concise, easier to work with, and abstracts implementation details (which makes it possible for DBMS to optimize the query), and hence supports parallel execution better.
- Imperative languages (e.g. most programming languages) explicitly describes the query execution, which may not offer higher performance, but allows you to describe complex query more easily.

MapReduce is an example that combines both declarative and imperative query.

#### Graph-Like Data Model

Essentially consists of Vertices (or Objects) and Edges (relationships between Objects, which can be directional). Since this data model serves very specific use-cases, they cannot really be compared side-by-side with the above two storage models when you are designing your application.

- Great for social network modelling, genealogical databases.
- Serves specific use-cases where everything and anything in your database can be related to everything and anything.
- there are specialized query languages to concisely declare the search query instead of having to write a recursive function.
- schema-on-read

### Chapter 3. Storage and Retrieval



# Additional Notes for Further Reading

- Why is Heap File used to store values for indexes in database systems, instead of having the index value hold a reference to the actual data row in the underlying data structure such as a B-Tree or a LSM-Tree? (Page 86)
- What is R-tree and why is it better at multi-dimensional data storage and searching? Is R-tree the solution to storing locations coordinates and distance matrices for travelling salesman problems? (Page 87)
- What are the trie-like index and Levenshtein automaton used in Lucene to support fuzzy full text searches? (Page 88)
