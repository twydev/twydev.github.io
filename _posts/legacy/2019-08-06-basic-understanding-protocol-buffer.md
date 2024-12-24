---
title: "Basic understanding of Protocol Buffer"
---

Understanding what is Protocol Buffer (Protobuf) and how we can use it.

**Is Protocol Buffer a data format, like JSON?**

It is not a data format. It is a method of serializing and deserializing structured data. Along with this method is an *interface definition language*, which can be used to define the structure of the data. 

A code generator tool can then parse the definition file (a `.proto` file) and generate code in various supported programming languages that serialize or deserialize the data.

The serialized data message is dense and is in binary format.

**Where did this method come from?**

According to Wikipedia, Google developed this method for internal use, but has opened source it, along with the code generator for multiple programming languages.

**When do we use ProtoBuf?**

When we want to define the communication between 2 systems. Instead of defining the message that gets passed, and implement the serializer and deserializer separately in both systems, we can now make use of ProtoBuf and automatically generate the necessary code.

**Why use ProtoBuf if we can use JSON messages?**

JSON messages are human readable, and can be parsed without having any knowledge of the schema (since this information is implicit in the message).

However, ProtoBuf offers huge performance improvement. Serializing and deserializing a ProtoBuf is much faster as compared to JSON messages. Also, the ProtoBuf messages are smaller since it is in binary.

Also, as mentioned earlier, the code to handle communication messages can be auto generated. This means that the ProtoBuf definition file can be the basis for maintaining cross-platform, multi-language library code.

**Any good library/framework to get started with?**

TBC
