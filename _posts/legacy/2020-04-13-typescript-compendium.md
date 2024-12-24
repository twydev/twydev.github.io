---
title: "My TypeScript Compendium"
toc: true
toc_label: "Chapters"
---

My notes are derived mostly from [TypeScript Deep Dive](https://basarat.gitbook.io/typescript/) and [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/). But before diving into TypeScript, it is best to start from having a strong understanding in JavaScript as TypeScript is ultimately transpiled into JavaScript.

## Project Organization

- `tsconfig.json` is used to define the project, files to include and compiler options.
- Code splitting can be achieved with Webpack by preserving dynamic imports of TypeScript instead of transpiling (can be configured in `module` setting).
- Global types can be controlled under compiler options as well.
- `lib` option controls what default type support we will be having from TypeScript compiler (e.g. `esnext`, `dom`). Another option, `target`, controls the JavaScript version we will transpile to.
  - Supported types are declared in `lib.d.ts` that ships with the compiler.
- Avoid using `outFile` option since it introduces higher chance of error and slower compilation.
- Type Declaration Space consists of `class`, `interface`, `type` declarations. Among the 3, only `class` declaration creates a variable.
- Importing module with relative path look up follows this sequence to look for a `.ts`, `.d.ts` or `.js` file:
  - Look for file in `node_modules` of current directory, and recursively look in parent directory uptil root of file system, for
  - a file name matches the import name or
  - a folder name matches the import name, and contains a file `importname/index.ts` or
  - a folder name matches the import name, contains `importname/package.json`, and a file is specified in `types` or `main` of the `package.json` file.
- `declare module` will make your module available in your program's global namespace. (equivalent to populating types in `global.d.ts` if your module only export types).
- If we are importing a file just for types, and this file does not directly use the import, then it may transpile into an empty file. Use the following approach to ensure import

```javascript
import mod = require('mod');
import mod2 = require('mod2');
const guaranteeImport: any = mod && mod2;
```

_`namespace` in TypeScript allows grouping of functions within a file, but using file based modules will achieve the same objective with cleaner code._

## JS Migration

- Understand that all JavaScript is legit TypeScript.
- Add `tsconfig.json`.
- Change `.js` files to `.ts`. Suppress errors with `any` type.
- Add new code with strong typing.
- Gradually refactor old code.

## Ambient Declaration

The **Definitely Typed** community provides many type definitions for popular JavaScript libraries. They can all be imported via npm, and they begins with `@type/` prefix.

But even if the definitions for a third-party library cannot be found, we can create our own ambient declaration, so that using this pure JavaScript library will not cause TypeScript compiler to throw error.

```javascript
// thirdpartylib.d.ts
declare function process(x: number): Promise<number>; // must use declare keywords
export default process;
```

For declaring variables type definition, we can consider using `interface` so that it can be extended easily in future.

## Coding (TypeScript)

### Class Access Modifiers

- `private` only accessable within the class definition.
- `protected` only accessable within the class definition, and child class definition. (even the constructor)
- `public` is the default. Accessible from anywhere, even the class instances.
- `abstract` class cannot be initiated, can only be `extended`. `abstract` members must be implemented by child classes.
- `static` members can only be accessed by fully qualified class name access. (I believe under the hood this has the same delegation property as ES6 static modifier).

### Type Alias

Essentially using `type alias = validTypeAnnotation`. (Therefore, it is important to note that `type` is only an alias keyword, the true definition of types happens in all the annotations)

### Interface

Interface is the main building block of TypeScript type system.

```javascript
interface myOwnInterface {
  prop: number,
  optionalProp?: number,
  readonly readonlyProp: number,
}

interface myOwnInterface { // open-ended
  (): string, // makes implementation of this interface callable as a function
  (input: number): number, // overloads the above function with a different signature
  namedFunc(): number,
  new(): string, // allows function to be called with NEW keyword
  explicitThis(this: myOwnInterface): void, // this function implementation uses THIS. Explicit declaration prevents implicit ANY for THIS.
}

interface myIndexable extends myOwnInterface {
  [index: string]: myOwnInterface, // this makes interface indexable, and also declare nested structure
}
```

- TypeScript interfaces are open-ended, so subsequent interface definition with the same interface name will be merged as one.
- It is pure syntactic sugar with no impact on JavaScript runtime.
- Classes can `implement` interface. (But does not mean you have to, all depends on use cases).
- Interface can `extend` classes as well, only inheriting the declaration of members without implmentation. If the class contains private members, then only a direct child class of this parent can implement such an interface extended from the parent.
- `readonly` can be used on any property in a class, interface or in any types. Using `Readonly<type>` marks all properties of input type as readonly automatically and returns a new type. Readonly is not failsafe, as the values can still be mutated through aliasing.

```javascript
var x = { readonlyProp: 1 } as myOwnInterface // dirty shortcut assertion

function mutateReadonly(aliasParam = { readonlyProp: number }) {
  aliasParam.readonlyProp = 2;
}

mutateReadonly(x); // mutates readonly property!
```

### Functions

- Supports parameter and output annotations.
- Supports optional parameter `?`.
- Support default parameter value.
- Supports **Function Overloading**

### Tuple

TypeScript has extended the capabilities of JavaScript array by introducing a way to declare tuples.

### Enum

- Enum extends beyond numbers, allowing us to create string enums, or even heterogeneous enums.
- Enum declarations are open-ended.
- If the numeric value of the first enum member is declared, value of subsequent members will increment from this first value. (Value declaration is optional).
- Declaring constant values for enums improves performance.
- Obtaining all enum values can be done using `keyof typeof [enum]`
- Value can be reversed mapped to key using `enum[value]`

### Type Assertion vs Type Casting

Type assertion can be achieved by using `var x = {} as myType`, so that no error will be thrown even if some properties of that type is initially missing. This is purely compile time checking, whereas type casting implies runtime conversion.

To assert to any given type simply requires a double assertion to first assert to `any` then to desired type. Assertion is generally harmful if used wrongly since it undermines type checking.

### Freshness (Strict Object Literal Check)

Only applies to **object literal** because it has a higher chance to suffer from typo errors or from misusing of APIs. Essentially if a function parameter has been annotated, then calling the function with object literal as parameter must strictly match the parameter declaration.

### Literal Types

Using JavaScript primitives _values_ as type, meaning that the value of those instances must be exactly the same as the primitive value stated in the type definition. This seems to only be useful when creating union types to restrict values of certain instances.

```JavaScript
type diceRoll = 1 | 2 | 3 | 4 | 5 | 6;
type direction = 'North' | 'South' | 'East' | 'West';
```

### Union and Intersection Types

- Union types are created using `type | type`.
  - Only properties common to all types can be accessed on a variable with union type.
- Intersection types can be created using `type & type`.
  - Can be used to implement mixin pattern.

### Discriminate Type Union

This comes down to using switch cases to check the value of certain literal member of this union type, that will be unique of each of the constituting types. A tip from the book _Typscript Deep Dive_ is to assign the input that falls through all the cases (therefore it is an unidentified type) to `never` to automatically throw an error.

Union types can be used to support backward compatibility, such as implementing interface versioning, and we can use the discriminate approach mentioned above to perform specific processing.

### Generics

Generics are supported using syntax like `function myFunc<T>(inputList: T[])`. This helps to constrain our inputs and outputs, provide type support while preventing us from reimplementing the same functionalities for different types.

```javascript
// constraining between two params that are related
function getProperty<T, K extends keyof T>(obj: T, key: K) {
  return obj[key];
}

// declaring parameter to be class constructor
function create<T>(c: { new (): T }): T {
  return new c();
}
```

### Type Compatibility

Type compatibility (assign instance of type A to reference for type B) depends on a number of factors:

- Child class are compatible with Parent class (polymorphism)
- Structurally similar types are compatible (same properties in object)
- Function types are compatible if they have:
  - sufficient information in return type (type A returns more information than type B) (extra output are ignored).
  - accepts fewer parameters (type A accepts less parameters than type B) (extra parameters are ignored).
  - optional and rest parameters are compatible.
  - parameter types are compatible.
- Enums are compatible with numbers.
- Classes are compatible by only comparing members and methods (ignores static and constructor).
  - Also, any `private` or `protected` members must originate from the same parent class, in order for two classes to be compatible.

### Never

`never` type is assignable to function that never returns or function that always throws. A `never` type can only be assigned to a `never` type, and not anything else.

### Error Handling (TypeScript)

Error subtype provided by JavaScript that we can use are: `RangeError,ReferenceError, SyntaxError, TypeError, URIError`.

The book _Typescript Deep Dive_ do not encourage throwing errors, but instead passing the errors around (perhaps through callback functions), so that errors can be annotated as optional function return type. This is an interesting approach since I really like to have potential errors tracked by our type system, but I also do not want to lose the ability to interrupt the execution by throwing.

### Non-Null Assertion

We can assert that a variable is non-null, non-undefined by using a suffix `!` on the variable. Using `!` suffix in declaration just signals to the compiler that the property being declared will contain a valid value before it is accessed (and it is the coder's responsibility to do so). (Again these are dangerous augmentations that should be avoided for cleaner code).

### Predefined Conditional Types

- `Exclude<T, U>` — Exclude from T those types that are assignable to U.
- `Extract<T, U>` — Extract from T those types that are assignable to U.
- `NonNullable<T>` — Exclude null and undefined from T.
- `ReturnType<T>` — Obtain the return type of a function type.
- `InstanceType<T>` — Obtain the instance type of a constructor function type.

## Note on Style

This is a list of recommendation from _Typescript Deep Dive_

- `camelCase` for function, variables, inner members, and filenames.
- `PascalCase` for namespace, classes, types, and interfaces (no legacy `I` prefix).
- `PascalCase` for enum and enum members.
- No explicit use of `undefined` and `null`. Use `!= null` or `== null` to guard against both `null` and `undefined`.
- `tsfmt` ships with compiler, and is useful for automatically formatting code.
- Prefer single quotes. Use backticks if we need to escape single/double quotes.
- Prefer 2 spaces, no tabs.
- Use semicolon to end statements.
- Prefer primitive types instead of native objects of primitives.
- Use `type` if we want to perform union or intersection. Use `interface` when we want to extend or implement.

## TSCompiler

TODO: read https://basarat.gitbook.io/typescript/overview/program

---

## Additional Readings

- https://www.joelonsoftware.com/2003/10/08/the-absolute-minimum-every-software-developer-absolutely-positively-must-know-about-unicode-and-character-sets-no-excuses/
- https://tsherif.wordpress.com/2013/08/04/constructors-are-bad-for-javascript/
- https://blog.izs.me/2013/08/designing-apis-for-asynchrony
- https://jakearchibald.com/2015/tasks-microtasks-queues-and-schedules/
- http://asmjs.org/
- TypeScript Handbook Reference - https://www.typescriptlang.org/docs/handbook/advanced-types.html
- https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#tests
- https://egghead.io/courses/introduction-to-reactive-programming
- https://gist.github.com/staltz/868e7e9bc2a7b8c1f754
