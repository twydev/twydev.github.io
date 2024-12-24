---
title: "Refactoring"
toc: true
toc_label: "Chapters"
---

A general approach to refactoring our code. It attempts to classify issues in our code into categories, and demonstrate the solutions to overcoming those issues. The list examples are not exhaustive, but it is sufficient for everyone to get started. 

Half of the book is made up of a reference catalog for you to quickly access the refactoring techniques that deals with the specific problem you are facing. Overall I felt that Clean Code by Uncle Bob was an easier read and may be more relevant to our daily work.

**Refactoring** *Improving the Design of Existing Code* - Martin Fowler, Kent Beck, John Brant, William Opdyke, Don Roberts, 2002
{: .notice--primary}

Everyone aims to write efficient software. However, writing legible software is equally, if not more, important. Clearly and cleanly structured code makes the software much easier to maintain and enhance, delivering more value in the long run.

## What is refactoring?

### Chapter 1: Refactoring, A First Example
This chapter showcase a small program that is part of a larger software, responsible for computing rental charges and printing customer statements. The initial version of the program cannot easily incorporate new features or changes to business logic.

**Approach**

1. **Refactor first, extend later**. If the code is not structured in a convenient way to add new changes, refactor it first, before making the changes. The extra effort is worth it in the long run.
2. **Set up tests**. Establish a set of tests to ensure that refactored code generates the same result as the existing code.
3. **Incremental small changes**. Incrementally make small changes and test often, to catch bugs more easily.
4. **Draw a simple UML** for the current program and update it as you make changes help to clarify your thoughts.
5. **Decompose large chunks** of code into smaller methods.
6. **Rename variables** and methods to make your code self-documenting.
7. **Replace temporary variables** with queries (e.g. getters). (Temp variables are typically problematic, causing too many parameters to be passed around. Problem worsens with long methods.)
8. **Make reusable**. Queries are reusable by the class. Changes in logic only need to be made in one place, which is the query implementation.
9. **Tradeoff between extensible code and optimized code**. Worry about optimization later. Refactor to make the code easier to maintain is the top priority.
10. **Replacing conditionals with polymorphism**. If the different conditions are subclasses of a main interface, then overriding methods in each subclasses will provide you the desired conditional logic while keeping the code easy to read and maintain.

```java
// Pseudo code in main object

class Main {
	method businessLogic() {
		if (condition A) {
			doSomethingA();
		} else if (condition B) {
			doSomethingB();
		}
	}
}

// Refactor using polymorphism

class Category {
	interface doSomething();
}

class CatA extends Category {
	method doSomething() {
		/* do something for condition A */
	}
}

class CatB extends Category {
	method doSomething() {
		/* do something for condition B */
	}
}

class Main {

	Category cat;
	/* in constructor create CatA or CatB according to logic */
	
	method businessLogic() {
		cat.doSomething(); // the correct override method will be resolved.
	}
}
```

8. **State Pattern.** *Referenced from Design Patterns by Gang of Four*. If a main object needs to change its behavior according to changes in its internal state, the behavior should not be implemented in this main object class.
	- Define separate object classes that encapsulate state-specific behavior. The main object class will call these separate objects to perform the required functions.
	- If states are described by separate state objects, those behavior can also be implemented in state objects, and main object will call them to perform the required functions.
9. **Rhythm of Refactoring**: build test > small changes > run test > small changes > run test > ...


### Chapter 2: Principles in Refactoring
This chapter talks about the definitions, benefits, and characteristics of refactoring.

> Refactoring = make software easier to understand and modify, without changing the observable behavior of the software.

You should be aware of which development activities you are performing at any point in time, and never mix the two up:
1. Add function to software - characterized by adding new tests, but no changes to existing code.
2. Refactoring - characterized by only using existing tests, and restructuring existing code.

Why some programs are hard to modify and work with:
- codes are hard to read
- codes have duplicated logic
- codes use complex conditional logic
- codes require additional behavior outside of the program

Benefits of Refactoring:
- Improves design of software
- Makes software easier to understand
- Identify bugs
- Speed up development as it is easier to add new functions to refactored code

Indirection (a.k.a dereferencing) = ability to act on a value through its name, container, reference and etc. instead of directly on the value itself.
- Enables sharing of logic
- Explains intention and implementation separately (self documenting through method and variable naming)
- Isolate changes (implement changes through subclasses)
- Encode conditional logic (via polymorphism)

However, indirection have costs. If there are any intermediate methods or components that used to serve a purpose or are supposed to be polymorphic but has become redundant, they should be taken out.

When to refactor:
- You have identified duplication of codes
- You want to add a new feature
- You need to fix a bug (it is likely that the code was not clear enough for you to spot the bug)
- You are doing code review (delivers more value when everyone can understand the code fully)

Potential problems of refactoring:
- **Database changes** are often hard to manage, especially schema changes. Having an intermediate layer of abstraction to between your program and the database helps minimize modification needed.
- **Interface changes**. All instances of interface needs to be changed if method signature changes. This is especially problematic if interface is published and some interface users are not accessible. You may mark the old interface as deprecated and let the old interface call the new interface. Understanding the trouble, it is better to not publish interfaces prematurely.

When not to refactor:
- when existing code is not even stable enough to work and effort is better spent rewriting the code.
- when you are super close to a deadline. Otherwise, you should always be refactoring.

Refactoring and Design:
- **Speculative design** = an attempt to put all good qualities into a system before any code is written. Like a waterfall model, this process is too easy to guess wrong. Refactoring along the way plug the holes in the initial software design. You will have less pressure to get things right on the first try.
- **Flexible programs** are often more complex to provide flexibility. Instead, implement a simple design that is easy to refactor to a flexible one if there ever is a need.

Refactoring and Performance:
- Well factored code allows changes more quickly, making your efforts to optimise the program more productive.
- Well factored code provides finer granularity for performance analysis.
- Focus on refactoring first, and then only optimise the code after having performance profiling. This will significantly reduce wasted effort to optimise portions of codes that only run once.

## Where should we refactor?

### Chapter 3: Bad Smells in Code
This chapter provides examples of smelly codes that should be refactored. However, the author also emphasize that we should develop the intuition to identify what needs refactoring as there are no hard rules that can be easily defined.

> If it stinks, change it.

1. **Duplicated Code** - duplicated code structure in two methods, or classes, or even having the same expression in two places.
2. **Long Method** - places that need comments are indication of semantic distance and should be refactored into a separate method for clarity. Even if it is just one line of code. Conditionals and loops are also indication of codes to be extracted to a new method.
3. **Large Class** - indicated by too many instance variables and too much code. Break down into smaller classes, subclass certain behavior, split into new and separate domain objects.
4. **Long Parameter List** - replace by passing an object or method that can be used to query for the necessary parameters. Exception is when you want to be totally independent from other objects and want to pass parameters by value.
5. **Divergent Change** - many different kinds of changes made to the software require code modification in the same method or class is an indication that the internals of that class/method should be refactored, so that different kinds of changes modify different classes/methods.
6. **Shotgun Surgery** - inversion of divergent change. If a single event requires code changes in 3 classes, it is a clear indication that some codes from these 3 classes can be extracted to form a new class/method.
7. **Feature Envy** - occurs when a method is more interested in getting and processing data from an object of a different class, instead of its own class. Keep data and methods that changes together in the same class.
8. **Data Clumps** - data fields that always appear together should be grouped in an object. A good test of relevance is to delete one field and see if the rest of the data fields still make sense.
9. **Primitive Obsession** - instead of using primitives, use classes to form the foundation Type used by other objects and methods, and also group data values within objects.
10. **Switch Statements** - use polymorphism to replace switch statements. 
11. **Parallel Inheritance Hierarchies** - edge case of shotgun surgery. Occurs when sub-classing a class requires you to make a separate subclass in another class, as if both classes mirror their class hierarchy. You may try to make objects of one hierarchy reference the objects of the other hierarchy.
12. **Lazy Class** - class that used to serve a purpose but is not so meaningful now and can be deleted or broken up and merged with other classes.
13. **Speculative Generality** - usually indicated by classes and methods that were overly complicated to serve a future use case that may or may not happen. Remove such abstract classes, methods, parameters to methods.
14. **Temporary Field** - usually occurs when object attributes or variables are only set for very specific purpose, but hangs around within the execution context doing nothing but confuse others. Refactor the fields to group them with the specific method that uses them into a method object.
15. **Message Chains** - an getting an object from another object, which internally gets it from yet another object. Observe the end result required and extract the necessary method to achieve this.
16. **Middle Man** - when delegation and encapsulation gets out of hand, and a class/object only house methods with the sole purpose of calling other objects/methods. You can call the underlying objects directly without the middle man and it would not make much of a difference.
17. **Inappropriate Intimacy** - when too classes are overly dependent on each other's internal structures or private fields/methods. Extract commonality between the two classes into a new class, or make one class the subclass of the other.
18. **Alternative Classes with Different Interfaces** - consolidate them. May need to create a super class.
19. **Incomplete Library Class** - extend or encapsulate external libraries that does not fully meet your needs.
20. **Data Class** - dumb data container class. Ensure there are no public access fields or collections. Move behaviors into the data class so that it gains more responsibilities besides holding data.
21. **Refused Bequest** - when a subclass does not want to use some of the methods/fields inherited from the superclass. This smelliness can still be accepted, but if the subclass is using interfaces that refuse the superclass, this is some kind of anti-pattern. The subclass should be refactored to a standalone class that delegates to the superclass for whatever necessary execution. Then the class will have no problem using an interface that does not support the superclass.
22. **Comments** - usually indicate smelly code that should be refactored. Before writing comments, try refactoring first. Useful comments usually **describe why** a code was design the way it is, instead of **explaining what** the code is currently doing.

## How to refactor

### Chapter 4: Building Tests
Frequent testing after every change is made improves overall productivity. Debugging is much faster when you can narrow down the bugs to the previous change.

- Use a unit test framework that makes it easy to add new tests, group tests together, and to conveniently run all tests.
- The tests should indicate and differentiate between failures and errors.
- Bugs identified by functional testing should be supplement by new unit tests that exposes the bug. If the logic is complex, try to write multiple unit tests, each with a smaller scope to narrow down the failure. Once the bug has been fixed, all these unit tests should provide you assurance that no other bugs that cause the same error falls through the crack.
- To avoid writing too many tests, evaluate the risks of each test scope, and write test for those that poses high risk.
- Write tests for boundary cases, and special cases.
- Write tests to ensure that exceptions and expected errors occur properly.
- When dealing with inheritance and polymorphism, consider the 80/20 rule. Test to make sure that the most probable cases and combinations of objects are working correctly.

## Guest Contributions

### Chapter 13: Refactoring, Reuse, Reality - William Opdyke
Developers are reluctant to refactor their code because of one of the following reasons:

- they do not understand how to refactor.
- they will no longer be with the project to enjoy the long-term benefits.
- refactoring is an overhead activity and they prioritize adding new features.
- refactoring may break the program.

The author goes on to address each of these reasons and how we may overcome them. We can learn how to refactor by reading this book and practicing. We have to understand that refactoring has short-term benefits too, making debugging faster and improving overall development time immediately. We can reduce refactoring overhead by using tools and technologies to assist us, and making refactoring part of our development regime to adopt it as a habit. Lastly, refactoring can be made safe with the use of tools, test suites, and code reviews.


### Chapter 14: Refactoring Tools - Don Roberts, John Brant
A good refactoring tool should satisfy the following requirements:
- ability to search for program entities across the entire program, and be able to differentiate functions and variables of the same name.
- must preserve the behavior of a method when the internal of the method is being refactored (analysed under the hood using a Parse Tree).
- speed and accuracy (safety) tradeoff. Some tools do not provide all refactoring capabilities so that they can maintain certain level of speed and safety.
- ability to undo refactoring changes, or even version control.
- integration with IDE.


### Chapter 15: Putting It All Together - Kent Beck
Here is how you may pick up the skill of refactoring:

1. Get used to setting a goal (specific and small issue to address by refactoring the code).
2. Stop when you are unsure.
3. Backtrack when necessary.
4. Refactor with a partner.

You know you are mastering the skill when you feel absolute confidence that no matter how screwed up the code was handed to you, you can make the code better, or make it good enough to keep development going. Also, you need to gain the ability to stop refactoring, with confidence.

Eventually, when you feel that the design of any system is fluid and moldable, and you are able to perceive the whole design at once, that is when you have mastered the skill.

## Reference Catalog
Please refer to the book if you need detailed examples for each type of refactoring.

### Chapter 6: Composing Methods
Most of the time, problems come from methods that are too long. This chapter deals with refactoring methods: how to handle temporary variables, input parameters, good naming conventions (for variables and methods) and how to replace algorithm implementation.

### Chapter 7: Moving Features Between Objects
A key decision to make when refactoring is in object design, and where to assign responsibilities. This chapter talks about strategy to move methods and fields between objects, creation of new classes, delegation, extension, and considering the possibility of future changes when making such decisions.

### Chapter 8: Organizing Data
One advantage of object languages is to the ability to group data with objects. This chapter addresses encapsulating data structure, association between objects (MVC, uni/bi-directional), and handling of type codes.

### Chapter 9: Simplifying Conditional Expressions
This chapter talks about decomposing and consolidating conditionals, using polymorphism, strategy for using null objects and assertions.

### Chapter 10: Making Method Calls Simpler
This chapter once again deals with refactoring methods, but focusing on method signatures, parameters, separation between queries and modifiers, using factory methods to replace constructors, casting, and throwing exceptions.

### Chapter 11: Dealing with Generalization
Generalization would will make the program more extensible and flexible. This refactoring would require re-organizing the inheritance hierarchy, which is the focus of this chapter.

### Chapter 12: Big Refactoring
This chapter provides the bigger strategy and overview of refactoring, keeping in mind the purpose we are trying to achieve.

- **Tease Apart Inheritance** - about tangled inheritance that confuses people.
- **Convert Procedural Design to Objects**
- **Separate Domain from Presentation** - separate business logic from UI.
- **Extract Hierarchy** - reorganizing classes.
