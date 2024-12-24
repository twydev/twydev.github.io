---
title: "Introduction to TDD - TDD and the Terminator"
---

Webinar hosted by JetBrains that provides an introduction to TDD, the key concepts to write a failing test first on established interfaces, write minimal code to pass the test, and to follow a Red-Green-Refactor flow to gradually implement the software.

**TDD and The Terminator - An Introduction to Test Driven Development** - Layla Porter, accessed on JetBrains Webinar 2019
{: .notice--primary}

## Key Takeaways

If you remember nothing from the talk, remember these:

> 1. When you write tests after coding, you may be writing tests to fit your code and not your requirements!
> 2. If you need to make a private method public in order to test it, then its time to refactor. 
> 3. Write the least amount of code possible to make the test pass.

## How to do TDD?

1. Write a failing test. 
    - make use of interfaces (for strongly typed languages). Interfaces act as coding contracts, and empowers asynchronous development within the team.
    - write a test for a method of the interface.
    - Name of the test should give a human-readable indication of the requirements the method is trying to fulfill.
    - use assertion to indicate what the method should return.

    ```java
    // requirement, I want to know if my pet is a cat
    
    interface IAnimal {
      public boolean isCat();
    }

    @Test
    public void result_true_when_my_pet_is_a_cat() {
      IAnimal myPet = new Pet("persian");
      Assert.assertTrue(myPet.isCat()) // Fail!
    }

    class Pet {
      private String animalType;
      public Pet(String animalType) {
        this.animalType = animalType;
      }
    }
    ```

2. Follow the Red-Green-Refactor approach.
    - This approach follows a cycle of 
      1. writing a failing test  
      2. writing minimal code to pass the test and write more tests that fail
      3. refactor to pass all the tests so far

    - **Red** stage
      - at this stage, tests should fail
      - now implement the interface in the most minimal way to pass the test (you would likely return some hard-coded value that the assertion is expecting just to pass the test)
      - re-run the test, and it should pass (green!)

      ```java
      class Pet implements IAnimal {
        private String animalType;
        public Pet(String animalType) {
          this.animalType = animalType;
        }
        public boolean isCat() {
          // the most minimal code to pass test.
          return true;
        }
      }
      ```

    - **Green** stage
      - with only one test, the method implemented only covers a single use-case (with a hard-coded result!)
      - now write another test, that tests the same method, but for a negative result.
      - run the test, it should fail (red!).

      ```java
      @Test
      public void result_false_when_my_pet_is_a_dog() {
        IAnimal myPet = new Pet("husky");
        Assert.assertFalse(myPet.isCat()) // Fail!
      }
      ```

    - **Refactor** stage
      - refactor the implementation of the interface, write code that pass all the tests
      - **Tip** check out the testing framework you are using, and use test-case features to iteratively run those tests against different sets of inputs. More test inputs will reveal that the code below is not sufficient and more switch-cases will be needed to make the code more robust.
      
      ```java
      class Pet implements IAnimal {
        private String animalType;
        public Pet(String animalType) {
          this.animalType = animalType;
        }
        public boolean isCat() {
          switch(this.animalType) {
            case "persian":
              return true;
            case "husky":
              return false;
            default:
              return false;
          }
        }
      }
      ```

3. Handle new requirements.
    - With the introduction of new requirements, you may find that the current code may not perform well, and refactoring of the design may be necessary.
    - What should we do now? Let's use some Design Patterns! And let's start from the foundation, the SOLID principles. In Porter's opinion, the first 3 principles are the most important and useful for writing tests.
        1. **Single Responsibility Principle** ensures that your methods are small and easy to test.
        2. **Open/Closed Principle** wants the behavior of the code to be closed from modification, but the code is open to extending more rules/conditions that leads to the same behavior.
        3. **Liskov Substitution Principle**: through proper use of subtypes, the code can be expended (helps to achieve open/closed principle).
    
    ```java
    // new requirement, i only want to feed the pet if it is a cat
    // new requirement, i will have more pets in future

    // refactor, following SOLID principles
    // we want to be able to reuse the analysis code to determine if an animal is a cat (Feeder)
    // we want to be able to extend the analysis to include more animals without changing the behavior (Rules)

    @Test
    public void test() {
      ArrayList<IRule> rulesList = new ArrayList();
      rulesList.add(new PersianRule());
      rulesList.add(new HuskyRule());
      Feeder feeder = new Feeder(rulesList);
      
      IAnimal myPet = new Pet("husky");
      Assert.assertFalse(feeder.feed(myPet))
      
      IAnimal myPet2 = new Pet("persian");
      Assert.assertTrue(feeder.feed(myPet2))
    }

    // determine if animal is a cat is decoupled from the animal itself

    interface IAnimal {
      public String type();
    }

    class Pet implements IAnimal {
      private String animalType;
      public Pet(String animalType) {
        this.animalType = animalType;
      }
      public String type() {
        return this.animalType;
      }
    }

    // we use rules. if a rule matches an animal we will feed it.

    interface IRule {
      public boolean match(IAnimal animal);
      public boolean feed();
    }

    class PersianRule implements IRule {
      public boolean match(IAnimal animal) {
        return animal.type().equals("persian");
      } 
      public boolean feed() {
        return true;
      }
    }

    class HuskyRule implements IRule {
      public boolean match(IAnimal animal) {
        return animal.type().equals("husky");
      }
      public boolean feed() {
        return false;
      }
    }

    // finally the logic to run an animal against the rules

    class Feeder {
      private List<IRule> rulesList;
      public Feeder(List<IRule> rulesList) {
        this.rulesList = rulesList;
      }
      public boolean feed(IAnimal animal) {
        for (IRule rule : this.rulesList) {
          if (rule.match(animal)) {
            return rule.feed();
          }
        }
        return false;
      }
    }
    ``` 

## Why people fail at TDD?

- Underestimating the learning curve. 
  - if we were to start on this journey, we need to be patient with the team
  - we need to be empathetic when the team makes mistakes
  - it is going to be a long journey
- Confusing TDD with unit testing.
  - you need to follow the Red-Green-Refactor methodology.
- Thinking that unit testing is enough.
  - we need other forms of tests to support our development (integration tests, regression tests, ...)
- Not starting with failing tests.
- Not refactoring enough.
  - turning private methods public to test is bad bad practice and smelly smelly code.
- Not actually doing TDD. (it is possible to check all the checkboxes in TDD practices without having a TDD mindset)

## Is TDD suitable for you?

- It can be controversial and is a significant culture change.
- Initial drop in productivity can be disconcerting.
- Productivity will go up and reworks reduced, but only in the long term.
- TDD leads to increased understanding of requirements and their acceptance criteria.
