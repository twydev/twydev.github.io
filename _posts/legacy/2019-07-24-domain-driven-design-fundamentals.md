---
title: "Domain-Driven Design Fundamentals"
toc: true
toc_label: "Chapters"
published: false
---

Comprehensive course covering Domain Modelling, Entities and Value Objects, Aggregate Pattern, Domain Services, Bounded Context, Repository Pattern, and Domain Events.

**Domain-Driven Design Fundamentals** - Steve Smith and Julie Lerman, accessed on PluralSight 2019
{: .notice--primary}

## Introduction

I think so far this is the only course I have given a 5-star rating on PluralSight. Extremely clear explanations with walk through of source code applying concepts of DDD. This is an introductory course that I would highly recommend. The dog pictures used in the demo is a plus! 

### Focus of DDD

1. Interaction with Domain Experts
2. Focus on a single sub-domain at a time
3. Implementation of sub-domains

### Resources

- DDDCommunity.org
- Domain-Driven Design: Tackling Complexity in the Heart of Software, by Eric Evans
- Applying Domain-Driven Design and Patterns: with examples in C# and .NET, by Jimmy Nilsson
- Implementing Domain-Driven Design, by Vaughn Vernon

## Modeling Problems in Software

1. First step is to understand as much of the business domain as possible from domain experts.
2. Ideally all important business domain terms are captured (nouns and verbs especially)
3. With sufficient terms captured, we can start to see exact same terms in the business domain used 
	- by different people to mean different things
	- by different people to achieve different goals
4. This is when we can start setting up Bounded Context to group and divide up the business domain according to specific purposes.
	- now it is fine for the same term may exist in separate bounded contexts and carry different meaning
5. Consider how the separate bounded context would interact with each other via Context Mapping
6. Create the Ubiquitous Language for each bounded context, used throughout conversations and code.

**Key Terms**
- **Problem Domain** is the specific problem the software is trying to solve.
- **Core Domain** is the key differentiator for the business. Something the software MUST do well.
- **Sub-Domains** are the separate applications or features that the software must support or interact with.
- **Bounded Context** is a specific responsibility, with explicit boundaries that separate it from other parts of the system.

## Elements of a Domain Model
The focus of this chapter is on capturing business behavior of the real world.

### Entities
Entities are objects in the system that are defined by a thread of continuity and identity rather than their attributes.

**If an entity captures all the behavior and becomes inflated, does this violate the Single Responsibility Principle?**

Eric Evans suggests that he would focus the design of such an entity on the responsibilities of **identity** and **life cycle**. If there are more complex logic that derives from the identity or the life cycle (such as to match if two objects have the same identity, or what processing should be done at each stage of the life cycle), which would cause the entity to be too large, then we may consider placing those derived logic in separate services.

**Should entities implement Equality comparer?**

Eric Evans responded that it makes sense for Value Objects to have equality, but applying that concept to Entities is a non-trivial question, because what kind of system or domain processes and use-cases would result in two identical entities, and what does it even mean to compare them for equality? In most cases, entity comparison is never necessary.

**Bi-directional Relationships between Entities in Bounded Context?**

Bi-directional relationships often leads to complications, and is best to avoid it unless that truly describes the business domain. On close inspection, most bi-directional relationships may be reduced to uni-directional relationships, by focusing on whether the purpose and responsibility of the Bounded Context can still be achieved without one of the entities. The direction of the relationship therefore should traverse from the entity that helps meet the Bounded Context purpose to the entity that is weakly associated with the purpose of the context.

### Value Objects
Value Objects are defined by composition of all their attributes. They are immutable once created, does not have an identity, and do not produce side effects to the state of the system.

Value Objects can have specific logic defined to process it. Classic value objects are date time ranges, money value with currency, and etc.

We should strive towards using Value Objects in our model as much as possible. Entities can serve as a simple value container, with properties defined as value objects, and any processing logic are encapsulated within value objects. This approach to modeling the domain allows us to take full advantage of the immutability to of value objects, and keeps our system clean.

**Should Value Objects have methods and logic?**

Eric Evans reasoned that it is better to have logic reside in Value Objects, than in Entities, to capitalize on the advantage of immutability and no side effects.

Testing a value object is tremendously easier than an entity, and you will feel much more encouraged to write and perform tests.

Loading most of the logic in value objects make entities perform a pure orchestration role, and the entire logic of the system will become very concise. The most ideal state is to end up with source code of entities and methods that read like use-cases written in the ubiquitous language, and is easily understood even by non-technical stakeholders.

### Domain Services
Domain Services are pieces of logic not naturally part of any entities or value objects, have an interface defined in terms of other domain model elements (in other words it interacts with other entities and value objects), and they are stateless processes by themselves, even though they may produce side effects on the system.

NOTE: Do not abuse domain services or you will end up with a procedural software or an anemic model. Try as much as possible to find a natural home for every piece of logic in a value object or an entity.

**Key Terms**
- **Anemic Domain Model** uses classes focused on state management. Good for CRUD.
- **Rich Domain Model** creates logic focused on behavior, not just state.
- **Entity** is a mutable class with an identity used for tracking and persistence.
- **Value Object** is an immutable class whose identity depends on the combination of all its attributes.
- **Domain Services** hold behavior that does not belong anywhere else in the domain.
- **Side Effects** are changes in the state of the application or interactions with the outside world (e.g. UI or Infrastructure Layer)

NOTE: So far, modeling the domain did the consider persistence of data, or how the database table would look like, how fields in tables are named or relate to each other, at all. These concerns will be addressed later. Following the approach of DDD, you can already implement the application logic and domain logic without even dealing with the persistent data. This is truly a code-first but behavior-centric approach.

## Aggregate Pattern
This pattern is the primary way of logically grouping entities and value objects, so that our classes and methods don't end up in a ball of mud.

### Aggregate
An aggregate is a cluster of associated objects that we treat as a unit for the purpose of data changes.

**Aggregate Root**

Every aggregate must have an Aggregate Root, an entity which acts as the single entry point to access any aggregate elements.
- access to elements within the aggregate must obey the ACID principles (Atomic, Consistent, Isolated, Durable). 
- the Aggregate Root is responsible for maintaining the aggregate's Invariant.

**Aggregate Invariant**

An Invariant is a condition that should always be true for a system to be in consistent state.

For example, some elements must be initialized in order for an aggregate to be in a consistent state. The Aggregate Root is responsible for ensuring those conditions are always met whenever changes takes place in the aggregate.

**Identifying Aggregate Roots**

An easy way to identify Aggregate Roots is to see whether deletion of an entity leads to cascading deletion of other elements in the aggregate for the system to maintain a consistent state. If yes, that is a likely candidate to be an Aggregate Root. If not, you may also be modeling the domain inaccurately, and some of these elements that do not get deleted in a cascade should just be a reference to a separate aggregate (as it may possibly even be an Aggregate Root itself), instead of existing in its entirety as a member of the current aggregate.

Another way is to identify Aggregate Roots is to see if an element within the aggregate will still make sense if other entities are removed. An Aggregate Root should still make meaningful sense even if it is the only member in an aggregate. It is the relationship with the Aggregate Root that gives other element its meaning (e.g. modeling next-of-kin information for employees. You will likely never need to access those information separately without relating to any employees)

Checklist for Aggregate Root:
- it enforces invariant
- changes to aggregate root leads to changes in other elements and aggregate state changes
- deletion of aggregate root leads to cascading deletion of other elements in the aggregate

**Relationships between Aggregates**

Aggregates should only access each other's member objects by referencing the Aggregate Root.

### Breakthrough
Occasionally, we may encounter a breakthrough as we continually improve our domain models to more accurately capture the business domain. We may suddenly realize that an entity we initially thought was an Aggregate Root should be modelled differently, and we would need to make changes to our model and even introduce new entities and concepts that never once occurred in our domain conversations.

This is totally fine.

### Tips for using Aggregates
Use of aggregate is not strictly necessary, and you don't want to use it just for the sake of using it. Also, use of foreign key to reference non-root entities in other aggregates is perfectly acceptable (I don't really get what this means. Does it refer to persistence layer.)

**Key Terms**

- **Aggregate** is a transactional graph of objects
- **Aggregate Root** is the entry point of an aggregate which ensures integrity of the entire graph
- **Invariant** is a condition that should always be true for a system to be in consistent state
- **Persistence Ignorant Classes** are classes with no knowledge how they are being persisted

## Repositories
Repositories can be used to manage the persistence of objects throughout their life cycle.

{% include figure image_path="/assets/images/screenshots/persistence-object-life-cycle.png" alt="" caption="Life cycle of Objects with Persistence" %}

From Eric Evans, a repository represents all objects of a certain type as a conceptual set, like a collection with more elaborate querying capabilities.

### Tips for using Repositories
- think of it as an in-memory collection
- possibly implement a common access interface that serves all repositories
- interface contains methods to add and remove objects
- methods that predefined criteria for object selection
- repositories should only be for aggregate roots
- repositories focuses on persistence

**Common Repository Blunders**

If implementation of repositories are not fine-tuned, it will commonly result in:
- **N+1 Query Problems**: in order to retrieve a list of items, one query will first generate a list from the database, followed by one query for each item, to retrieve the full item from the database.
- **Inappropriate use of eager loading or lazy loading**
- **Fetching more data than required**

Therefore it is important to know the data model of underlying persistence when implementing repositories.

**Repositories are not Factories**

Even though both return objects when called, Factories generate new objects and does not deal with persistence. Repositories are responsible for finding and updating existing objects. Repositories may or may not use Factories when handling objects, but they are definitely not the same.

**Key Terms**
- **Repository** is a class that encapsulates data persistence for an aggregate root.
- **ACID** refers to Atomic, Consistent, Isolated, and Durable.

## Domain Events

### Tips for using Domain Events

- Each event is its own class
- Include time when the event took place
- Capture event-specific details
- Event fields are initialized in constructor
- No behavior or side effects

### Example of how Events logic works

I have chosen to use JavaScript. If other OO language are chosen, the Services, Entities, Events, Repositories, Handlers should be established as Interfaces first, before implementing as Classes. Also this is just an example that I quickly came up with. Tons of smell.

```javascript
function mainApp() {
	ScheduleApptService.scheduleAppt(email, apptTime);
}

class ScheduleApptService { // Service
	static scheduleAppt(email, apptTime) {
		const appt = Appt.create(email);
		ApptRepository.save(appt); 
	}
}

class Appt { // Entity
	constructor(id = new Guid()) {
		this.id = id;
		this.email = null;
	}
	create(email) {
		const appt = new Appt();
		appt.email = email;
		DomainEvents.raise(new ApptCreated(appt));
		return appt;
	}
}

class ApptCreated { // Event
	constructor(appt) {
		this.appt = appt;
		this.datetimeOccurred = new Date();
		this.type = 'ApptCreated';
	}
}

class NotifyUIApptCreated { // Handler
	static handle(apptCreated) {
		/* trigger UI that appt has been created */
	}
	static subscribes(event) {
		return event.type === 'ApptCreated';
	}
}


// Event Manager
const handlerContainer = [ 
	NotifyUIApptCreated,
];
let eventActions = [];
class DomainEvents { 
	static register(callback) {
		actions.push(callback);
	}
	static clearCallbacks() {
		eventActions = [];
	}
	static raise(event) {
		handlerContainer.map(handler => {
			if (handler.subscribes(event)) {
				handler.handle(event);
			}
		});
		eventActions.map(action => {
			if (action.subscribes(event)) {
				action.handle(event);
			}
		});
	}
}

class Schedule { // Entity
	constructor(id, dateRange, apptList) {
		/* create list of appts and set other attributes */
		DomainEvents.register(handleConflict); // register the handler following Hollywood Principle
	}
	const handleConflict = { 
		const subscribes = event => {
			return event.type === 'ApptUpdated';
		}
		const handle = event => {
			/* de-conflict the appt */
		}
	}
}
```

### Event Boundaries

To prevent events from becoming bloated, consider using separate events for different clients. Also, sometimes we may need to receive an event, translate and publish as a separate new event.
- separate event objects for specific clients
- separate events for external clients
- specific application events for presentation layer

### Anti-corruption Layer
Additional adapters and services to interface whatever is happening inside the bounded context with external clients.

**Key Terms**

- **Domain Event** is a class that captures the occurrence of an event in a domain object.
- **Hollywood Principle** says that "Don't call us, we'll call you". Domain logic should not call events, events should trigger domain logic. (essentially the callback concept).
- **Inversion of Control** is a pattern for loosely coupling a dependent object with an object it will need at runtime.


## Final Notes

### Consider the UI

Sometimes, the problems we are trying to tackle in the domain can be easily solved in the UI. With good UI design, certain validations in the backend domain layer can be eliminated as we are guaranteed consistency from the frontend.

### Fallacy of Perfectionism

Eric Evans advised that there is no perfect design. Many people will get stuck in DDD trying to come up with the perfect design and every iteration it never seems good enough. Let's be realistic, imperfection is fine, don't get stuck.

### Advantage of DDD

In a nutshell, DDD simplified the adding of new feature to a software.
