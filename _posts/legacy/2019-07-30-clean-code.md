---
title: "Clean Code (PluralSight Course)"
toc: true
toc_label: "Chapters"
---

A very brief and general introduction to clean code by Cory House. Applicable to all languages. Touch on topics such as naming, conditionals, functions, classes, and comments.

**Clean Code** *Writing Code for Humans* - Cory House, accessed on PluralSight 2019
{: .notice--primary}

The main takeaway from this course:

{% include figure image_path="/assets/images/screenshots/clean-code-foundation-stack.png" alt="" caption="Clean Code Foundation Stack" %}

This pyramid provides a path for everyone to establish a strong foundation in software development.

## Naming

- Classes should indicate their responsibility.
- Methods should reflect their action / purpose.
  - if a name contains conjunction, you are likely violating Single Responsibility principle, and you should refactor.
- Avoid abbreviations.

## Conditionals

- Use positive boolean, instead of naming variables as "isNotTrue".
- Avoid "Stringly Typed" conditions. Use Enum to hold string values.
- Replace explicit boolean comparison with method calls that returns boolean.
- Use polymorphism to replace conditionals.
- Use table records (in DB or some storage) to determine condition.

## Functions

- Create functions to reuse code and avoid duplication.
  - a good sign is either spotting the exact same code block or similar code structure.
- Extract long methods into small functions of single responsibilities.
- Return early and fail fast.
  - Do validations and requirements checking upfront, before executing the core logic of the function.
  - This prevents the software from entering a corrupted state, and makes code easier to read.
- Avoid loose / short-lived local variables. Replace them with method calls.
- Use exceptions to catch problems.

## Classes

- Used to group attributes and methods of the same responsibility together.
- Should not have too much knowledge of other classes. 
  - Use inversion of control pattern to deal with such scenarios.

## Comments

- With proper naming and clean coding, comments are unnecessary most of the time.
- May be used to document design decisions, trade-offs, and design intent.
- Avoid apologies and warnings.
- Avoid commenting out codes, which becomes zombie codes.
- Avoid formatting characters, as they are nice to view but hard to maintain.
- Avoid tagging issue numbers, ticket numbers using comments, as that can already be handled with commit messages / branches using source control.

## Final Reminders

- If not broken, why fix?
- If broken, please fix. (avoid the Broken Windows Theory, which essentially means technical debt).
- Practice code reviews and pair programming.
- Follow boy scout rule (leave the code a little better than when you received it).
