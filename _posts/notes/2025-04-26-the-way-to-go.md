---
title: "Notes for: The Way to Go"
source_title: "The Way to Go: A Thorough Introduction to the Go Programming Language"
source_author: "Ivo Balbaert"
source_published: "2012"
source_edition: 1
ISBN: "978-1469769175"
categories:
  - notes
tags:
  - golang
toc: true
classes: wide
published: true
---

> [!info]
> title: {{ page.source_title }}
> author: {{ page.source_author }}
> published: {{ page.source_published }}
> edition: {{ page.source_edition }}
> ISBN: {{ page.ISBN }}

I already have some experiences working on Golang codebase, so I will only highlight the important points from this book.

# Design objectives of Golang

- it aims to be the "C for the 21st century"
- software needs to be built quickly (for developer productivity)
	- a rigid, clean, fast dependency analysis allows fast compilation
	- explicit dependencies
- language should run well on modern multi-core hardware
- language should work well in a networked environment
- language should be a pleasure to use
- efficacy, speed and safety (type and memory safe), strongly and statically compiled, ease of programming like a dynamic language
- garbage collected to avoid memory problems (therefore not suitable for real-time software)
- for backward-compatibility, it is able to run C code as well.

what was removed from usual languages to maintain simplicity:

- no function or operator overloading
- no implicit conversions
- no classes and type inheritance
- no variant types
- no dynamic code loading
- no dynamic libraries
- no generics (back then)
- no exceptions (but has panic and recover)
- no assertions (runtime assertions that throws error. not the same as type assertions)
- no immutable variables

The spirit of Golang is to follow convention and idioms when writing code, so that it is easy to navigate between code bases.

# Install and Runtime Environment

- Golang provides compiler for different OS and architecture (32/64 bits)
	- a word may be 32-bits (4 bytes) or 64-bits (8 bytes) depending on the architecture
- Go compiler supports cross-compiling, developing a go program on a certain host architecture but compile it for a different target architecture
- with source code, you can also just compile the code in your own OS and architecture, as long as the OS and architecture is supported, your Go program will build correctly
- Go compiler generates native executable code.
	- but it will link a Go runtime code to every Go package
	- the Go runtime is somewhat comparable to VM used by Java/.Net languages
	- Go runtime is responsible for handling memory allocation, garbage collection, stack handling, goroutines, channels, and other managing memory reference data types.
	- therefore Go compiled code will be much bigger than source code, due to the runtime embedded but deployment is much easier since there is no need to link any external files

# Best Practices and Conventions

- naming should be short, concise, evocative, simplicity
	- package name should be lowercase, act as namespace
	- functions/methods/variables don't need to contain indication of the package name
	- getters can be named as a noun e.g. `person.Name()` instead of `person.GetName()`
	- setters can be prefixed with verb e.g. `person.SetName()`
- there are some universal method names used in Go standard library, we are encouraged to adopt those names when implement similar functionality for our own programs e.g. `Open(), Read(), Write(), String()`
- names of interfaces should have `-er` suffix, e.g. `Reader, Writer`
	- interfaces are short, usually max 3 methods
- **always recover panic in your own package**. panics should never cross package boundary, for a good developer experience for someone using your package.
- return errors as error value to the callers of your package
- documentation comments should be provided at the start of the package (the main entry, one of the file). this will be extracted by godoc.
- **Golang has no tail recursion optimization (TCO)**
	- a function is tail recursive, if the recursive call is the final line in the function with no further computation
	- this means the current recursion can be cleaned from the stack
	- compilers that implements TCO will be able to optimize stack memory
	- however Golang encourage you to use loop instead of recursion
	- this keeps compiler and runtime implementation simple. it also generally leads to more readable code (arguable?)

# Golang Constructs Basics

- only 25 keywords. this allows you to keep all required knowledge in your mind while working with the language
- building a Go program
	- it must contain one package called `main`
	- main package can depend on other packages (by convention one directory per package)
	- each package is compiled as a unit
	- every piece of code is compiled only once
	- if package A depends on B depends on C:
		- C is compiled first into C object file
		- B is compiled, compiler will link to C object file
		- A is compiled, but compiler will only link to B object file (this speeds up build at a large scale)
	- dependencies are explicitly declared by package imports
- structure of a Go program/package (common convention)
	- import statements
	- declare constants, variables, types
	- define `init()` function, which executes set up when this package is imported
	- define `main()` function (only if this is the `main` package)
	- define methods on types
	- define functions
	- methods and functions can be defined in the order they have been used. or alphabetical order if large numbers of functions exist
- execution of a Go program/package
	- packages in `main` are imported in the order as indicated, and recursively importing in every package (but each unique package is only imported once)
	- for every package (in reverse order) all constants and variables are evaluated, and then `init()` function are executed (LIFO)
		- initialization is always single threaded, to guarantee correct execution order
	- finally the same initialization happens for `main` package
	- then finally the `main()` function executes

## Types

| type | memory rep | mutability | zero val | alloc/init mem | 
| --- | --- | --- | --- | --- |
| int | value | immutable | 0 | new() |
| float | value | immutable | 0.0 | new() |
| bool | value | immutable | false | new() |
| string | value | immutable | "" | new() |
| array e.g. `[3]int` | value | **mutable** | `[0,0,0]` | new() |
| struct | value | **mutable** | all fields zero value | new() |
| slice e.g. `[]int` | reference | mutable | nil | make() |
| map | reference | mutable | nil | make() |
| channels | reference | mutable | nil | make() |
| interface | reference | mutable | nil | N.A. |

- when a **variable is declared, all memory is initialized, to the zero value for the corresponding type**
	- for value types, the variable contains actual value data
	- for reference types, the variable only contains a memory address to another location that contains the actual data structure
	- reference types are therefore very similar to pointers
	- and the zero value for pointers and reference types are `nil`
	- you are therefore discouraged by the IDE to not use a pointer of a reference type. just work directly with the reference type variable.
- `new()` allocates memory for value types
- `make()` is only for reference types, and it initializes the memory (or data structure) for that reference type so that it is usable
	- it typically accepts optional params, to initialize maps/slices of certain sizes
- `int`, `uint`, `uintptr` are architecture dependent, defaults to 32-bit or 64-bits
- strings are a sequence of UTF-8 characters with variable-width per character (will use 1 byte ASCII-code when possible, and expand the bytes used depending on the character) which allows Golang string to occupy less space as compared to C++ or Java.
	- using `for` loop over each byte will therefore produce undesirable results
	- using `for range` loop over a string will allow the correct parsing of each UTF-8 character

## Slice in memory

- memory structure has 3 fields
	- has a pointer to underlying array, at the starting index of the slice
	- has a length field
	- has a capacity field
- slice can be resliced up to the capacity limit (which is determined by the underlying array)
	- if a slice is shortened from the starting index, then there is no way to expand it back to the starting index, even if the underlying array still exists
- `append()` adds an element to the slice (and hence the underlying array is mutated)
	- if underlying array do not have sufficient capacity, a new slice (with new array) is allocated and returned. **The original slice and array will remain intact.**
- when slices are assigned to a new variable e.g. `slice2 := slice1`, the **memory structure of slice1 is copied to slice2**
	- so slice1 and slice2 points to the same underlying array
	- but slice1 and slice2 length field can now mutate independently of each other
- `copy(dst, src []T) int` copies elements from source slice to destination slice
	- it returns number of elements copied
	- it is limited by the minimum length between source and destination slices
- **slices can potentially reference a very big array with high impact on memory**
	- you can copy the relevant data to a new slice and return from your function
	- the new slice will be using a shorter array with only relevant data
	- this allows the underlying large array to be GC-ed

## Map in memory

- maps are good for key-value store and retrieval
	- but it is still 100x slower than direct indexing in an array or a slice
	- so take note for ultra high performance requirements
- map capacity grows by only 1, whenever an additional element needs to be added beyond the current capacity
	- **it is better for performance to set the capacity for large maps or map that keeps growing to reduce memory operations during runtime**

## Structs in memory

- All data in a struct, even if it contains another struct, are allocated in a continuous block of memory
	- this gives huge memory performance boosts
	- fields must therefore be defined in struct literal in the same order as their struct definition
	- and since fields are sequentially defined, anonymous fields are supported, as long as you define the values in sequence
	- anonymous fields will use the type as field name e.g. `myStruct.int, myStruct.float32`
	- therefore each struct can only have one anonymous field of the same type
- **if struct contains reference types, those are allocated with zero values**
	- those reference types must be initialized separately before they can be used
	- it may be a good idea to use a factory pattern to create a new struct
- Methods are functions that accept a certain struct type as receiver
	- **the receiver type and method implementations must be declared within the same package**
	- you can define your method with either a pointer or value receiver
		- pointer receiver allows you to mutate the struct. it is also efficient to pass large struct data
		- **value receiver is a copy of struct data, therefore modifying value type fields will not work.**
		- **but a value receive will allow you to mutate reference type fields**, since the receiver is a copy of the reference type containing a pointer to the underlying memory data.
	- **Golang runtime is smart enough to allow pointer/value receiver methods to interoperate on either pointer or value inputs**
		- Golang will perform reference/dereference before calling the method
		- So if you method expect a value receiver, it will still receive a value receiver, regardless of a pointer being used at the calling line of code

### Struct Inheritance

- inheritance of fields and methods can achieved using embedding
	- "parent" is the inner struct
	- "subtype" is the outer struct
	- inner struct can be anonymous, so outer struct will inherit the fields and methods
- **(gotcha) however if the inner struct method accesses a struct field as part of its implementation, the inner struct field will be used even if you called the method on the outer struct**
	- the only way to achieve the desired override behavior is either for the outer struct to either define its own implementation of the same method
	- or update the inner struct with a new field value
	- this is expected, since methods are defined in the same package that defined the corresponding struct
		- if you subtype a third party package struct and expect that your outer struct shadow fields will override inner struct methods
		- then the compiler and Go runtime will incur a lot more complexity and overhead
	- this explicit behavior also made the language simpler
- embedding more inner structs allows multiple inheritance

## Interfaces

- interface can define a method set without implementation of methods (abstract)
- Go allows us to define a variable with the interface type (typescript also allows for this)
- the interface variable has the following structure in memory
	- the zero value of interface variable is `nil` (similar to pointer)
	- it contains a method table pointer field, which points to the actual implementation of methods matching this interface, for a provided concrete type.
	- it contains a receiver field, which holds the data of a concrete type variable
		- the receiver field can be assigned with a value, or a pointer
		- e.g. if a struct is assigned to the interface variable, then it is a copy of a struct (similar to using value receiver for a struct method)
		- e.g. if a struct pointer is assigned to the interface variable, then this field stores the memory address (similar to using pointer receiver for a struct method)
		- e.g. if any other data types are assigned, the behavior will depend on whether they are primitive (values) or reference types
- an interface can embed another interface to achieve inheritance, which is the same as explicitly enumerating all methods of the embedded interface
- type assertions
	- interface variable can be cast to concrete type using type assertion at runtime e.g. `t, ok := interfaceVar.(T)`
	- this can be used for type switches
	- a concrete value can also be tested to see if it implements an interface e.g. `ivar, ok := var.(MyInterface)`
- **when calling methods on an interface variable, the reference/dereference receiver behaviors are different from struct methods**
	- when interface variable contains a pointer
		- it can be called with pointer receiver methods
		- it can be called with value receiver methods (automatic dereference)
	- when interface variable contains a value
		- it can be called with value receiver methods
		- it **cannot** be called with pointer receiver methods
	- reason behind this behavior
		- when assigning to an interface, the method table pointer is determined upfront
		- when you assign only a value to the interface, the method table pointer will only point to the set of value receiver methods
		- when you assign a pointer to the interface, then the method table pointer will point to both value/pointer receiver methods
		- since interface variable by design needs to be dynamic, it cannot perform the same automatic reference that a concrete struct type can perform
		- this explicit behavior allows the code writer to be very intentional with which method set he is trying to assign when using the interface variable (he can choose to provide a value or a pointer)
- therefore when embedding inner interface as a pointer in an outer struct, it will allow the outer struct to be called with all methods (value/pointer receiver) of that interface. You will still need to implement the methods though, since interfaces are abstract.

### Empty Interface

- an empty interface `type Any interface{}` can be used to represent an "any" type since it has no method set
- (gotcha) the code snippet below produces an error:

```go
var dataslice []myType = //..some data..//
var anyslice []interface{} = dataslice

// error: cannot use dataslice (type []myType) as type []interface{} in assigment
```

- because the two variables `dataslice` and `anyslice` are two different concrete arrays, even though `anyslice` is an array containing "any" type of elements
- to overcome this problem, you can assign the elements one by one.

## Reflection

- reflection first cast the input as an empty interface, then retrieves the various attributes of the input
	- `reflect.TypeOf(x)` gives the type of the input
	- `reflect.ValueOf(x)` gives the raw value of the input
	- `reflect.ValueOf(x).Kind()` gives the type, based on the raw value
	- `reflect.ValueOf(x).Interface()` gives the interface value, based on the raw value
- reflection allows you to change the value of the input
	- `reflect.ValueOf(&x)` must be provided with a pointer (**this is a recurring design of Golang. use pointer to allow mutable value**)
	- `v := reflect.ValueOf(&x).Elem()` indirects through the pointer and returns the underlying element
	- `v.CanSet()` will be true
	- `v.SetFloat(1.234)` will succeed
- `Printf(format string, args ... interface{})` uses reflection extensively to figure out what are the inputs and what format strings are supported.
- I have once encountered a function in production that attempts to mutate all string fields in an input struct
	- the function made heavy use of reflection
	- one developer used the function but provided the in-built Error type instead of struct to the function call
	- the Error was mutated, but it resulted in memory corruption
	- later attempts by the program to log that Error object crashed the entire Go program (it was not even a panic that can be recovered. it was a fatal crash)
	- my hypothesis of what could have happened:
		- the function uses reflection and recursively mutated any string fields in a struct (and any nested structs)
		- however, Error contains an un-exported struct that contains the string message and a length field
		- the function mutated the Error string (making the string shorter), but it failed to update the length field
		- this violated the memory alignment assumptions
		- when the Go runtime later tries to access the Error object in memory, it reads the unmodified length, and reads that segment of memory, which includes addresses that are already invalid, this resulted in a segmentation fault.

# Standard Library

just some notes on standard library

- `net/rpc` provides utility for RPC
- `encoding/gob` (stands for Go Binary) provides de/serialization to binary format. it can be considered as an alternative for data representation (to JSON which is popular). But it can only be used by Go processes since it uses Go reflection.
- `netchan` (stands for network channel) has been moved out of Golang v1
- `websocket` has been moved out of Golang v1

I think standard library must be studied by all users of Go, to avoid reinventing the wheel.

# Golang OOP Equivalence

- Encapsulation
	- package scope: objects defined within their own package with a name that starts with lowercase
	- exported: objects defined within their own package with a name that starts with uppercase
	- a type can only have methods defined in its own package
- Inheritance (by composition)
	- embedding types with desired behaviors
	- supports multiple inheritance through multiple embedded types
- Polymorphism
	- a variable of a type can be assigned to a variable of an interface, if type implements the interface
	- types and interfaces are loosely coupled
	- multiple inheritance is supported by types implementing multiple interfaces

# JSON Handling

- JSON handling will be common for web application development
- Golang to JSON data type conversions
	- bool marshals to JSON booleans
	- float64 marshals to JSON numbers
	- string marshals to JSON strings
	- nil marshals to JSON null
	- maps can be marshalled to JSON if keys of the maps are string
	- channel, complex, function types cannot be marshalled
	- cyclical data structure cannot be marshalled as it will cause the `Marshal()` function call to loop infinitely
	- pointers are marshalled as the values they are pointing to
	- only exported fields of a struct will be marshalled

# Error Handling

- you can create a new error object with custom message using `errors.New()`
- you can also create a new error object with format string message using `fmt.Errorf()`
- you can also implement your own custom error by having a type implement the built-in error interface
	- typically used if you need to load more data into the error object or you need to implement custom methods on the error object
- if a runtime problem occurs, a panic will be triggered (value of interface type `runtime.Error`)
- flow of control when panic occurs in a nested function call (also called panicking)
	- current function execution stops immediately
	- `defer` functions executes in LIFO order
	- control is passed back to function caller
	- caller either handles the panic, or continue execution of defer functions and bubble panic upwards
	- at top level, program crashes and error condition is reported to the CLI using the value in the panic
- `recover` can only be used inside a deferred function
	- retrieves error value passed through the call of panic
	- some standard libraries uses recover so that errors are explicitly returned to caller, instead of causing a panic to bubble from the API
- (gotcha) `defer` expression are evaluated at the line they are called, but the function will be executed later in the deferred phase
	- e.g. `defer myFn(420+69)` the summation will be evaluated on the spot

# Testing

- standard library `testing` provides automated testing, logging, error reporting, benchmarking
- test programs must sit in the same package as source programs
- test programs must be named `*_test.go`
- test programs are not compiled by the normal Go compiler (so they are not deployed). only `gotest` compiles all programs
- test functions should be named with the Test prefix and have a particular header form `TestXXX(t *testing.T)`
- the options `-cpuprofile` and `-memprofile` can output the respective resources profile into a file
- `pprof` is a runtime library that needs to be imported into your source program to output more profiling data for analysis

# Concurrency

## Goroutines

- Golang provides goroutines and channels for concurrency, which requires support from the language, compiler and runtime.
- The spirit of Go concurrency: **Do not communicate by sharing memory. Share memory by communicating**.
- Communication forces coordination
- concurrency and parallelism
	- a single process on a machine run in its own address space
	- it can be made up of multiple OS threads, running concurrently on 1 processor core (multi-threaded)
	- only if the same application process runs on multiple core, then it is considered parallelized
- the classic concurrency model in other programing languages
	- synchronize different threads by locking data
	- this is complex and error-prone
	- race conditions lead to unpredictable results
- Go uses the Communicating Sequential Processes or Message Passing Model
	- Go provides goroutines to run concurrent computations
	- goroutines do not correspond one-to-one to OS threads
	- the Go runtime component, goroutine-scheduler, manages the multiplexing of goroutine to OS threads
	- goroutines can run on multiple threads or within one thread only, and these are abstracted from us and well managed by Go runtime
	- goroutines share memory address space, but synchronizing data is discouraged (but provided by `sync` standard library). 
	- channels should be used to communicate between goroutines
- goroutine footprint
	- it is lightweight since it is a virtual abstraction on top of actual OS threads
	- they are created with a 4K memory stack-space on the heap
	- segmented stack can grow or shrink dynamically according to usage, all managed by Go runtime
	- therefore a large number of goroutines can be created on the fly (~100K in the same address space)
	- memory are freed when goroutines exit (no need for GC). goroutines exit silently without notifying the caller.
	- similarly, **when `main()` returns or exits, the Go program terminates. It doesn't wait for goroutines to finish**.
		- `main()` can be designed to receive signals from goroutines
		- `main()` can attempt to wait for all goroutines to complete execution before exiting (can consider using `sync.WaitGroup`)
	- `runtime.Gosched()` called within goroutine function will yield the processor (can be used if computation is intensive), but there is no need to explicitly or manually resume the goroutine, the Go runtime will handle that.
- **by focusing on writing good concurrent programs, we let the underlying Go runtime decide if the goroutines will be run in parallel using multiple cores or not. because a concurrent program that runs well can also be parallelized well.**
- number of cores to use can be set using `runtime.GOMAXPROCS(n)`
	- experiments suggests best performance is obtained by setting n = one less than number of cores on the machine
	- higher n does not necessarily means better performance due to communication overhead
	- if n is set to 1 then only 1 core will be used, and there will be no parallel executions.
	- in later versions of Golang, Go runtime by default uses all available core, if `GOMAXPROCS` is not set
- goroutines may run during `init()` function of `main()`

## Channels

- channel is like a conduit (pipe) through which goroutines can send data (typed values) between each other
	- it is like a message queue
	- data follows a FIFO order
- **only one goroutine has access to a data item at any given time, data races therefore cannot occur by design**
- **channel sending and receiving are atomic**, they always complete without interruption
- Go compiler is able to detect some deadlock scenarios using static analysis, if it knows that some goroutines will be blocked due to missing send/receive implementation on channels
- by default, communication are synchronous and unbuffered
	- sender will be blocked until a receiver is ready to receive
	- and similarly, receiver is blocked until a data is sent
	- and if more sender or receiver tries to access the channel, they will be waiting behind the current blocked sender/receiver
- asynchronous communication can be achieved through use of buffered channels `make(chan type, bufSize)`
	- this prevents blocking the sender (unless buffer is full) or receiver (unless buffer is empty)
	- in a multiprocessor set up, using a buffered channel for producer-consumer pattern can result in goroutines that never block, since the producer and consumer may run in parallel
	- however, the bigger the channel buffer, the more memory will be used. benchmarking can help us to figure out the optimal performance and memory tradeoff and setting the appropriate buffer size.
- channel variables can also be declared to be omnidirectional
	- `chan<-` can only send data
	- `<-chan` can only receive data
	- misusing them after declaration will therefore cause compile errors
	- and read-only channels cannot be closed. **the convention is to have the sender close the channel to signal end of processing**.
- the convention is usually to close the channel using defer right after channel creation
	- receiver can test whether the receive failed because channel has been closed using `v, ok := <-ch`
	-  if receiver is using `for v := range ch {}` construct, the loop will automatically ends when channel has been closed
- select statement (looks similar to switch-case) can be used to select between different channels to send/receive
	- if all channels are blocked, it waits until one channel can proceed
	- if multiple channels can proceed, it chooses one at random
	- therefore it is usually used in a loop, to implement a listener pattern. the loop will be broken by an explicit break statement when conditions are met
	- if default clause is present, then it will be executed if no channels can proceed

### Channel Usage Patterns

- `time.Ticker` and `time.Tick()` returns a channel that sends a tick at specified time intervals, which you can use to implement time driven tasks (like periodic logging, or rate limit, since receiving will block until the next tick)
- select statement can be used to implement some interesting patterns to avoid blocking
	- listener
		- select between multiple different channels
		- one of the channel can be dedicated to sending the quit signal
		- if all channels are blocked, and quit signal is received, the loop can terminate and the goroutine can clean up and exit
	- timeout
		- select between receiving from a data channel, and a time based channel
		- if data comes first, move on
		- if time based channel receives first, then a timeout has occurred
	- fastest result (similar to JavaScript `Promise.race()`)
		- create a common buffered channel (this is a minor optimization, in case the main function is still spinning up goroutines, but one of the goroutine has already completed execution. having buffered channel allows sending data without waiting for a receiver)
		- main function spins up goroutines
		- the main function blocks to receive the result from the common channel
		- each goroutine queries data from a different external source
		- each goroutine implements a select statement with only one case: to send the queried data to the common channel
		- but the select statement also has an empty default clause
		- so all goroutine will query the data first before attempting to send the data to the channel
		- when it is time to send data, the fastest goroutine will succeed first.
		- the slower goroutine will be blocked by the channel.
		- but since they are attempting to send data in a select case statement, they will fallback to execute the empty default clause
		- and then the slower goroutine will exit, without blocking and without sending data.
		- this ensure that we only receive the fastest result, and all other goroutines will automatically clean up
- implement semaphores
	- a buffered channel with n capacity
	- it can hold n number of placeholder data which represents semaphores
	- goroutines "check out" the semaphore from the channel before accessing shared memory
	- goroutines "check in" the semaphore back to the channel after they are done
- implement lazy evaluation generators
	- since sender will block unless receivers read the data
	- and if sender is in a loop
	- then the sender only computes the next result in each loop after sending succeeds
- implement futures pattern (similar to JavaScript `Promise`)
	- the concept is similar to lazy evaluation
	- you can pass the reference for future around, while sending a separate goroutine to compute the results
	- you can then "await" for the future when you finally need the result
- exclusive data access via "background task"
	- set up a single goroutine to execute functions sent over a channel (reference to a function can be passed)
	- only that goroutine can execute those functions, which exclusively access shared data in memory
	- other goroutines/functions will dispatch functions to the channel
	- this is simply another implementation of semaphores
- concurrent computation of sequential tasks
	- a typical single-threaded approach could be to loop over the data, and in each iteration, perform tasks A, B, C in sequence
	- we can re-write this by having task A communicating output to task B, and then to task C, using channels, forming a pipeline
	- each step in the pipeline (A, B, C) can be run on a different goroutine, since these goroutines now communicate exclusively using channels
	- the only loop we need is to iterate over the data, and sending the data into the pipeline
	- we delegated the scheduling of concurrent computation to the Go runtime (this can also be parallelized potentially, depending on architecture and set up of your Go runtime, but it is abstracted from us)

### My understanding on Channels

- Golang essentially offered a queue-like data structure with exclusive access for goroutines, to achieve the same synchronization that we typically use with locks and shared memory
- exclusive access is the key feature here, as it essentially removes deadlock from poor implementation occurred when manually managing locks
- the same channel behavior can definitely be implemented using traditional locks and share memory, nothing special here. so Golang is basically having the channel API, and encouraging us to use this API and adopt a different paradigm to manage our concurrency

# Common Pitfalls

- never use `goto` with a preceding label (usage of `goto` is intended to skip forward)
- when implementing `String()` method, don't use other functions like `fmt.Print` that itself depends on `String()` (this causes infinite recursive loop)
- always use `Flush()` to terminate buffered writing
- never change the `for` loop counter variable inside the loop itself. use the language construct.
- use factory pattern to create your structs (if it contains reference types that needs to be initialized)
- use pointer receiver if you need to mutate the receiver
- reference types do not need to be passed explicitly as pointers to functions (interface type is also a reference type!)
- do not misuse `defer` inside a loop to clean up resources
	- `defer` will finally execute in LIFO order after function returns, NOT at the end of the loop
	- it may lead to a large number of unclosed resources if the function is large and does not return
	- directly close the resources in each loop iteration is simpler
- **prefer to pass value to functions when it is small (cheap to copy). pass pointer if mutability or size dictates**
	- passing by value to the function, means the value lives on the stack
	- stack allocation is fast (only bumping the pointer)
	- stack is auto-cleaned when function pops from the stack
	- stack is cache friendly (CPU cache reads from contiguous stack memory faster)
	- does not use the heap, which leads to less memory fragmentation in the heap over time
	- heap memory requires GC
	- Go compiler uses **escape analysis** to decide if a variable can live on the stack or the heap
- **always start with a simple single-threaded implementation. only use goroutines and channels when concurrency becomes a concern**
	- it is more important to get the program right

## Misusing Closures and Goroutines

```go
func main() {
	// version A
	for ix := range values {
		go func(){
			fmt.Print(ix)
		}()
	}

	// version B
	for ix := range values {
		go func(ix interface{}){
			fmt.Print(ix)
		}(ix)
	}

	// version C
	for ix := range values {
		val := values[ix]
		go func(){
			fmt.Print(val)
		}()
	}
}
```

- version A will not work as intended
	- `ix` is a variable in the `main()` function scope
	- all goroutines have shared access to a single `ix` variable
	- the goroutine executions are scheduled by the runtime, and may only begin at the end of the main loop
	- there is a chance that all goroutine is reading the same value for `ix`
- version B works 
	- because the `ix` value is copied to goroutine stack on invocation
	- in each loop iteration, a unique value is copied
	- even if goroutines are scheduled to run at the end of the loop, all goroutines have received a different value in their function argument
- version C also works
	- because `val` is a variable that only exists in the scope of each iteration (not shared between iteration)
	- each goroutine is invoked with a closure over a unique `val` variable

# Common Language Specific Patterns

- comma, ok pattern
- defer pattern (to close resources, to recover from panic)
- visibility pattern (exported objects)
- factory pattern (to initialize un-exported fields)
	- naturally you will make a struct private, but export the factory function from your package to create the struct correctly
- operator pattern
	- implement operators as methods 
	- the methods are internally responsible to switch between types and operate on the values correctly
	- this pattern allows chaining of operations
	- an interface can be used to describe this polymorphism

```go
func (x *SubTypeA) Add (y MainType) MainType

func (x *SubTypeB) Add (y MainType) MainType

subtypeA.Add(subtypeB).Add(subtypeC)

// each subtype implements the same interface
type Operable interface {
	Add(x MainType) MainType
	Subtract(x MainType) MainType
}
```

## Performance Advice

- convert string to byte slice to manipulate characters in a string
- use `utf8.RuneCountInString(str)` to count number of characters correctly (because characters uses variable number of bytes in Golang)
- use `bytes.Buffer` instead of string concatenation to accumulate string content (fastest and most memory efficient)
- when using goroutine for concurrent computation
	- **work done inside the goroutine must be much higher than the overhead of creating goroutine and sending data through channels to gain performance benefit**
	- buffered channels can help increase throughput (at the cost of memory)
	- channels can become bottlenecks, consider passing references to shared memory instead (e.g. pass pointers to array, and unpack the values in the receiver instead)
- use slices when possible instead of arrays
- use `for-range` loop over a slice if you only need the value, not the index, of the elements (this is slightly faster as it does not perform individual look up for elements by index)
- using a `map` instead of a sparse `array` can reduce memory footprint
- specify initial capacity for maps
- use pointer receiver for methods of structs
- using constants or flags to extract constant values from the code (hardcoded literal built into the compiled code, no memory allocation is required at runtime)
- use caching when large amount of memory are being allocated
