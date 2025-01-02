---
title: "Clean Architecture (PluralSight Course)"
toc: true
toc_label: "Chapters"
published: false
---

Extremely high level course on clean architecture, briefly covering Domain-Centric vs Database-Centric designs, Component Cohesions, Microservices, Testing, and extensibility of the architecture.

**Clean Architecture** *Patterns, Practices, and Principles* - Matthew Renze, accessed on PluralSight 2019
{: .notice--primary}

I have skipped over CQRS, which is a topic covered in greater details in other online courses. This course is an introduction to the concept of Clean Architecture. Many of the good principles covered in this course overlap with Domain Driven Design and feels more like a high-level refresher to modern architecture design approach than any thing new.

## Domain-Centric vs Database-Centric Architecture
From Uncle Bob: the first concern of an architect should be the usability of the software, not the implementation details. This thinking puts into perspective what parts of software development are essential vs what are merely details.

In the domain-centric approach, domain and use-cases are deemed essential, and presentation and persistence are merely implementation details.

In the database-centric approach, database and persistence are deemed essential, vs all other components.

Some examples of domain-centric architecture:
- Hexagonal Architecture, by Alistair Cockburn
- Onion Architecture, by Jeffrey Palermo
- Clean Architecture, by Uncle Bob (incorporate BCE model by Ivar Jacobson)

{% include figure image_path="/assets/images/screenshots/layered-dependency-inversion.png" alt="" caption="Inversion of Control + Dependency Inversion in Modern Architecture" %}

Abstraction should not depend on details. Details should depend on abstraction. Hence dependency of Persistence and Infrastructure layer is inverted. However, flow of control is passed through the layers, down to the Persistence and Infrastructure layer, following Inversions of Control pattern.

We can achieve this by using Interfaces to abstract away the actual implementation of each layers / services (Dependency Injection pattern).

## Functional Cohesion
Organizing codes into folders and namespaces that follows business use-cases (e.g. Orders, Sales, Products) instead of software components (Models, Views, Controllers).

Projects can still be categorized at a high level by layers (e.g. Presentation, Application, Domain, Persistence, Infrastructure, Common), however within each layer, the classes and methods should be sub-categorized by functions.

This approach to organizing software code provides spatial locality, ease navigation, and avoid locking the code into specific implementation technologies, however we would lose all advantages of other framework conventions and possible automatic scaffolding.

## Microservices
Each Bounded Context in the business domain may be implemented as distinct Microservices. This approach will force the separation of contexts physically through the implementation.

Advantage of Microservices is independence of each services, which allows separate development and scaling. However, this approach requires higher up-front cost and managing a distributed system also introduces another set of complexities.

## Testable Architecture
The course makes an accurate observation on the current state of testing in our industry. I have personally made the same observation.

Development teams do very little testing, ineffective and inefficient testing, which still lead to high bug count in production. The top reasons being:
- insufficient time
- "not my job" perception
- too difficult

**Test-Driven Development in a nutshell**

1. Create a failing test
2. Get the test to pass
3. Refactor the code to improve quality
4. (cycle!)

**Advantages of TDD**

By creating a comprehensive suite of tests, we ensure that all important codes are covered. In order to run those tests, we will design our methods and classes to be more testable. Coincidentally, testable components are also more maintainable. Falling back on a suite of tests eliminates our fear of breaking stuff when we refactor/extend our system.

The following diagram shows the state of tests in typical development.

{% include figure image_path="/assets/images/screenshots/test-automation-pyramid.png" alt="" caption="Test Automation Pyramid" %}

Acceptance Tests are important to verify the functionality of features being delivered. Acceptance Tests capture the criteria of completeness. However, if we require full tests involving manual testing and UI testing to complete Acceptance Tests, the up-front cost is too high for our development to remain efficient and effective.

If we could instead, by designing our architecture in a testable way, use service tests in place of full tests for Acceptance testing, we will speed up the process significantly.

{% include figure image_path="/assets/images/screenshots/service-test-as-acceptance.png" alt="" caption="Implement Acceptance Testing using Service Test only" %}

## Evolving Architecture

**Last Responsible Moment**

Deferring implementation decisions to the last responsible moment can be advantageous as it provides decision makers more information without causing too much delay. (from Lean Software Development: An Agile Toolkit, by Mary and Tom Poppendieck). 

{% include figure image_path="/assets/images/screenshots/last-responsible-moment.png" alt="" caption="Decision making time window" %}

**Evolutionary Architecture**

Even after implementing the software, an evolutionary design coupled with agile development approach provides the flexibility for the software to continually improve and adopt to changing market, changing technology, and changing preferences.

## Further Readings

**Books**
- Patterns of Enterprise Application Architecture, by Martin Fowler
- Clean Architecture, by Robert C. Martin (Uncle Bob)
- Domain Driven Design, by Eric Evans
- Dependency Injection in .NET, by Mark Seaman

**PluralSight Courses**
- Domain-Driven Design Fundamentals
- Domain-Driven Design in Practice
- Modern Software Architecture
- Microservices Architecture
- Dependency Injection On-Ramp

**Websites**
- http://martinfowler.com
- https://goodenoughsoftware.net
- http://udidahan.com
- www.matthewrenze.com
