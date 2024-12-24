---
title: "Structure and Interpretation of Computer Program"
toc: true
toc_label: "Chapters"
---

**Structure and Interpretation of Computer Program** *The Big Ideas Behind Reliable, Scalable, and Maintainable Systems* - Gerald Jay Sussman, Hal Abelson, 2017
{: .notice--primary}

*to be written*

___

# Preface
Here are my favorite quotes from the preface, and the reason why I picked up this book.

*Alan J. Perlis, the first recipient of the Turing Award*
> It doesn’t matter much what the programs are about or what applications they serve. What does matter is how well they perform and how smoothly they fit with other programs in the creatio of still greater programs.


*Gerald Jay Sussman and Harold Abelson, professors, MIT*
> First, we want to establish the idea that a computer language is not just a way of getting a computer to perform operations but rather that it is a novel formal medium for expressing ideas about methodology. Thus, programs must be written for people to read, and only incidentally for machines to execute. Second, we believe that the essential material to be addressed by a (computer-science) subject at this level is not the syntax of particular programming-language constructs, nor clever algorithms for computing particular functions efficiently, nor even the mathematical analysis of algorithms and the foundations of computing, but rather the techniques used to control the intellectual complexity of large software systems.


# Chapter 1: Building Abstractions with Procedure

> The contrast between function and procedure is a reflection of the general distinction between describing properties of things and describing how to do things, or, as it is sometimes referred to, the distinction between declarative knowledge and imperative knowledge. In mathematics we are usually concerned with declarative (what is) descriptions, whereas in computer science we are usually concerned with imperative (how to) descriptions.

The example provided in the book uses the square-root function as an example. We can define square-root at a function of input variable X, that is greater than or equals to zero, and will be equal to value of X when squared. This is a concise definition to check if the square-root of X is correct, but is totally useless when we try to compute the square-root.

> So a procedure definition should be able to suppress detail. The users of the procedure may not have written the procedure themselves, but may have obtained it from another programmer as a black box. A user should not need to know how the procedure is implemented in order to use it

This actually ties in with Clean Code. Because cleanly named procedures are good abstractions that promotes the construction of even more complex procedures by other engineers.

```javascript
function factorial(n) {
  // this is a recursive PROCEDURE that is running a recursive PROCESS
  if (n === 1) {
    return 1
  }
  return n * factorial(n-1)
}

function factorial_v2(n) {
  function factorial_iterate(product, count, max) {
    // this is a recursive PROCEDURE that is running an iterative PROCESS
    if (count > max) {
      return product
    }
    return factorial_iterate(count * product, count + 1, max)
  }

  return factorial_iterate(1, 1, n)
}
```

Recursive procedures are procedures with definitions containing reference to the procedure itself (merely a lexical description).

A recursive process is characterized by a chain of deferred operations. The interpreter must keep track of a chain of operations it needs to perform in future, while resolving each procedure call in the chain.

In contrast, an iterative process is characterized by having the state of the process captured in the input parameter to the procedure. The interpreter need not keep track of additional state information, since the input to each procedure call sufficiently describe the entire current state of the process.

**Tail recursion** has long been known as a compiler optimization trick to compute recursive procedures as iterative process, which helps to limit the growing memory demand.

