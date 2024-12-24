---
title: "Introduction to Domain Driven Design"
toc: true
toc_label: "Chapters"
---

Domain Driven Design (DDD) is a software development approach to implement complex software by dividing the solution into domains that serves specific purposes, using tools such as Ubiquitous Language, Bounded Context, Domain Events and Aggregate Patterns.

**Distilling Domain-Driven Design** - Vaughn Vernon, accessed on Kalele 2019
{: .notice--primary}

DDD helps to keep the software maintainable as it scales up in complexity, but also guides the behavior of a cross-functional development team, consisting of software developers and domain experts, towards better collaboration and higher efficiency.

## Chapter 1: Introduction and Overview
There are two parts to DDD, Strategic and Tactical Design.

**Strategic Design**

- Segregate domain models using strategic design pattern called **Bounded Context**
- Develop a **Ubiquitous Language** as your domain model in an explicitly Bounded Context
- Importance of engaging Domain Experts as you develop your model's Ubiquitous Language
- Create one close-knit team of Domain Experts and Software Developers
- Integrate multiple Bounded Contexts using a technique called **Context Mapping**

**Tactical Design**

- Modelling your domain in the most explicit way possible using **Domain Events**. Domain Events share happenings within a model with other systems that needs to be informed, which may exist within the same Bounded Context or a different Bounded Context.
- Aggregate Entities and Values objects together into a right-sized cluster using a tool called **Aggregate Pattern**. (Unfortunately the course video link to this chapter is broken).

## Chapter 2. Strategic Design with Bounded Context:
DDD is essentially modelling an Ubiquitous Language in an explicitly Bounded Context.

Within a contextual boundary, each component of the software model has a specific meaning and does specific action. The software model reflects a language spoken by the team that created the model, called the Ubiquitous Language.

- it consists of a glossary of terms.
- terms are used by the team to communicate with each other about the business logic, and also used in the software model.
- this ensures that when someone uses expressions from the Ubiquitous Language, everyone understand what is meant.

Any components outside the context are deemed out-of-context and are not expected to be consistent with the Ubiquitous Language

**Typical Software Design Challenge**

Typical software development projects, especially those that are highly complex, end up in a growing web of tangled models. Model boundaries should be based on work group definitions in real life, such as the different departments in a business. Attempting to merge concepts across contexts often results in a lose-lose situation for everyone.

Main issues of tangled models:

- growing number of aggregates cross-contamination because of unwarranted connections and dependencies between models.
- changes in one model causes a ripple effect across the entire tangled web of models which makes debugging and implementing new features extremely difficult.
- only tribal knowledge can save such a system from complete collapse.

Therefore, using a strategic design with Bounded Context and Ubiquitous Language helps to unify the mental model for everyone in the team. How to go about designing the software using this approach:

- identify the various software models.
- name the model components within Bounded Context by using terminology that makes sense within the context.
- identify the core domain that the software should focus on.
- put aside model components that do not belong to the context.
- consider the software implementation timeline. Even if some model components are in-context, they can be out-of-scope due to the tight timeline, and we can defer the modelling of these components.
- context may contain components that only implement a part of the feature. The full feature will only be implemented in a separate context.

After taking the above steps, we will be left with a smaller, more focused and relevant, core domain. The other model components that are left out of the context may form their own Bounded Contexts with their own Ubiquitous Language.

A final note, Bounded Context is more than just the domain model. It consists of application and even infrastructure layers.

{% include figure image_path="/assets/images/screenshots/layers-in-bounded-context.png" alt="" caption="Layers inside Bounded Context" %}

## Chapter 3. Strategic design with Context Mapping
Bounded Contexts can be linked through Context Mapping, which can represent different kinds of relationships:

{% include figure image_path="/assets/images/screenshots/context-mapping-partnership.png" alt="" caption="Partnership relation" %}

**Partnership**: two Bounded Context will frequently synchronize what is happening within their context, and will succeed or fail together.

{% include figure image_path="/assets/images/screenshots/context-mapping-shared-kernel.png" alt="" caption="Shared Kernel relation" %}

**Shared Kernel**: two Bounded Context sharing a model. We need to determine which context team is responsible for developing, testing, and maintaining the shared model. 

{% include figure image_path="/assets/images/screenshots/context-mapping-customer-supplier.png" alt="" caption="Customer-Supplier relation" %}

**Customer-Supplier**: two Bounded Context related in an upstream-downstream relationship where the Supplier context will pass information to the Customer context. The downstream context team may inform the upstream context team what they require, but it is still up to the upstream context team to implement those requirements within their own context.

{% include figure image_path="/assets/images/screenshots/context-mapping-conformist.png" alt="" caption="Conformist relation" %}

**Conformist**: another upstream-downstream relationship, where the upstream context has no motivation to provide for the specific requirements of the downstream context, and the downstream context also has no means to translate the upstream information to fit their needs. Therefore, the downstream context simply conforms their model to the upstream model.

{% include figure image_path="/assets/images/screenshots/context-mapping-anti-corruption-layer.png" alt="" caption="Anti-Corruption Layer" %}

**Anti-Corruption Layer**: the most defensive upstream-downstream relationship, where the downstream creates a translation layer to isolate their model from the upstream model. It translates information from the upstream model to what the local model requires specifically.

{% include figure image_path="/assets/images/screenshots/context-mapping-open-host-service.png" alt="" caption="Open Host Service" %}

**Open Host Service**: defines an interface that provides access to a Bounded Context as a set of services. it is open so all other contexts may use it freely.

**Published Language**: is a well documented language that enables translation by any number of consuming Bounded Context to easily translate from their own model into and out of the published language. Often, an Open Host Service will provide a published language.

{% include figure image_path="/assets/images/screenshots/context-mapping-separate-ways.png" alt="" caption="Separate Ways" %}

**Separate Ways**: is a special case where no integration would produce sufficient value to meet your needs. Therefore it is more optimal for each Bounded Context to implement their models in separate ways.

If you ever need to integrate with a legacy system containing tangled models, you may treat it as a single Bounded Context, and try to use an Anti-Corruption Layer to isolate your own Bounded Context from the mess.

**Mechanism to provide integration of Bounded Contexts**

- RPC (using Web Services that communicates through SOAP is the most common approach)
- RESTful API
- Message Queues / Publish-Subscribe
- Database
- File System data exchange

**Robustness of integration mechanism**

1. RPC is the least robust. 
    - If there is a network problem, or system problem in the host running the service, then communication between 2 bounded contexts will fail.
    - However, if all infrastructure problems are prevented, it can work well with Open Host Service + Published Language (upstream) and Anti-Corruption Layer (downstream) integration design.
2. RESTful API is more robust than RPC.
    - One advantage is the API naturally forming an Open Host Service + Published Language.
    - However, RESTful communication will still fail when infrastructure problems occur, just like RPC. RESTful communication has a better performance track record in terms of reliability and scalability since the entire internet is based on this architecture.
    - Be careful not to design API resources that directly reflect aggregates in the domain models, which forces a Conformist relationship on all the API consumers.
3. Messaging mechanism is the most robust.
    - Temporal coupling exists in both RPC and REST. Through a messaging mechanism, we eliminate temporal coupling between a context publishing domain events and a context subscribing to those messages.
    - Downstream context may need to send commands to the upstream context to trigger events or query for additional information, but it should always receive events from the upstream through messages.
    - The quality of the integration depends on the quality of the messaging mechanism. The mechanism should guarantee at-least-once delivery to ensure that all messages will reach the downstream context in any situation.
        - At-least-once delivery typically re-sends a message if there is a message loss, or if the receiver is slow, or if the receiver is down.
        - Downstream receiver needs to be idempotent and needs to handle any repeated messages by de-duplication or by ignoring the message or by a safe re-run of the operation to get the exact same result.

**Example of Context Mapping using Messaging Mechanism**

- Source Bounded Context publishes events and also stores those events.
- Events contain identifier to the corresponding entity object in the source Bounded Context.
- Downstream Bounded Context that is subscribing to the source context acts on the event, and stores the reference identifier so that we have the ability to trace back the what has happened if an issue were to occur.
- Downstream context can query source context using that reference identifier for more information (this query mechanism can be implemented using RESTful API).

## Chapter 4. Tactical Design with Domain Events
Domain Events are used in the communication between Bounded Contexts.

A Domain Event interface should minimally support the following attributes:

- **DateTime**: when event occurred
- **Type**: the type of event. should reflect domain model ubiquitous language. 
- **Name**: this should be a verb in past tense indicating what has happened to what entity.
- **Properties**: are a set of fields to provide details of the event. It contains all command properties of the command that caused the event, such as:
    - entity identifiers
    - command name
    - command description

It is important that when a command triggers a change, one single transaction saves both the state change (called aggregate) and the domain event to ensure consistency that an event has occurred.
