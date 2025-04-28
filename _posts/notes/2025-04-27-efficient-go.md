---
title: "Notes for: Efficient Go"
source_title: "Efficient Go: Data-Driven Performance Optimization"
source_author: "Bartlomiej Plotka"
source_published: "2022"
source_edition: 1
ISBN: "978-1098105716"
categories:
  - notes
tags:
  - golang
toc: true
classes: wide
published: true
---

> title: {{ page.source_title }}
> author: {{ page.source_author }}
> published: {{ page.source_published }}
> edition: {{ page.source_edition }}
> ISBN: {{ page.ISBN }}

# What is Efficiency

- the word performance has a broad meaning and typically covers these elements
	- accuracy. how many errors the system produces
	- speed. how fast the system completes the required task
	- efficiency. how much extra resources are used (wasted) to complete the task
- common efficiency misconceptions
	- Optimised code is not readable. this only applies to extreme scenario, low level implementation, which is not required for most use cases. typically, inefficient code are less readable.
	- YAGNI rule should not be an excuse to skip simple and convenient efficient optimisation
	- Hardware is getting faster. Yes, but. Bad practices that causes inefficiencies tends to fill up all available hardware resources.
	- We can always scale horizontally to solve our problems. Yes, but horizontal scaling in distributed systems is order of magnitude more complex, more difficult, and more expensive, to get things right and to maintain, as compared to efficiency optimisation within one system.
	- Time to market is more important. Yes, but. It is also more tragic to go to market with a non-performant system. The reputation damage may not be recoverable.
- why should we try to be efficient?
	- it is harder to make efficient software slow
	- optimising for speed is fragile (due to other factors like network). often, being efficient with given resources is the only thing developers can control
	- speed is less portable (again, due to the different environments)

## Efficiency Optimisation

- the goal of efficiency optimisation, is to
	- modify code without changing functionality
	- so that overall execution is more efficient
	- or at least more efficient in categories we care about (which trades off by being worst in some other categories)
- reasonable optimisations
	- cutting away obvious wastes in our system
	- these could be legacies left behind by hasty implementations, or refactors, and are now obsolete but results in wasted resources
	- these are typically easy to optimise on, and provides significant immediate benefits, without affecting functionality
- deliberate optimisations
	- these are optimisations that are not obvious
	- it will require us to improve on the efficiency of one category/resource, at the sacrifice of another category/resource
	- this is a zero-sum game
	- but we should still pursue such optimisations if it makes sense for our use cases
- why is optimisation hard
	- we are bad at estimating which part of the system has performance problem
	- we are bad at estimating exact resource consumption
	- maintaining efficiency over time is hard
	- reliable verification of current performance is difficult (environment is complex)
	- optimising can impact other software qualities (tradeoffs)
	- for Golang, we don't have strict control over memory management
	- we did not define what is efficient enough

## Efficiency Requirements

- form a practice of stating efficiency requirements up front in a project
- encouraged to use Resource-Aware Efficiency Requirements (RAER)
	- state what is the operation we are interested in
	- the scope of input (size, shape of data etc.)
	- maximum latency of operation
	- resource consumption budget for this operation in this scope
- to be more pragmatic, just focus on the most critical operations, and focus on the boundaries (best and worst case scenarios)
- instead of absolute resource and latency measure, we can define them in relation to the input scope (e.g. use memory 2x the size of input data)
- we can always use a rough estimate to start, by applying napkin math calculation, and refine our requirements later
- tips to handling efficiency issues reported in a production environment
	- if there is already a workaround, direct users to that solution
	- if the user is using the system outside of functional scope, kindly explain to the user
	- if the issue reported is outside of RAER specs, kindly explain to the user
	- if the issue reported is within RAER specs, then the team should potentially try to work on it
	- but note, as usage and system evolves, efficiency requirements can change over time

## Latencies for Napkin Math Reference

Some notes
- KB refers to kilobyte, which is 1000 bytes per KB
- KiB refers to kibibyte, which is 1024 bytes per KiB, which is more precise and uses binary system
- numbers below are estimated based on year 2021, using x86 CPU from the Xeon family
- but most CPU independent numbers have been stable since 2005

| Operation | Latency | Throughput |
| --- | --- | --- |
| 3 Ghz CPU clock cycle | 0.3 ns | NA |
| CPU register access | 0.3 ns (1 CPU cycle) | NA |
| CPU L1 cache access | 0.9 ns (3 CPU cycles) | NA |
| CPU L2 cache access | 3 ns | NA |
| Sequential memory R/W (64 Bytes) | 5 ns | 10 GBps |
| CPU L3 cache access | 20 ns | NA |
| Hashing, not crypto-safe (64 Bytes) | 25 ns | 2 GBps |
| RAM R/W (64 Bytes) | 50 ns | 1 GBps |
| Mutex lock/unlock | 17 ns | NA |
| System call | 500 ns | NA |
| Hashing, crypto-safe (64 Bytes) | 500 ns | 200 MBps |
| Sequential SSD read (8 KB) | 1 us | NA |
| Context switch | 10 us | NA |
| Sequential SSD write, -fsync (8 KB) | 10 us | 1 GBps |
| TCP echo server (32 KiB) | 10 us | 4 GBps |
| Sequential SSD write, +fsync (8 KB) | 1 ms | 10 MBps |
| Sorting (64-bit integers) | NA | 200 MBps |
| Random SSD seek (8 KiB) | 100 us | 70 MBps |
| Compression | NA | 100 MBps |
| Decompression | NA | 200 MBps |
| Proxy: Envoy/ProxySQL/NGINX/HAProxy | 50 us | ? |
| Network within same region | 250 us | 100 MBps |
| MySQL, memcached, Redis query | 500 us | ? |
| Random HDD Seek (8 KB) | 10 ms | 0.7 MBps |
| Network NA East <-> West | 60 ms | 25 MBps |
| Network EU West <-> NA East | 80 ms | 25 MBps |
| Network NA West <-> SG | 180 ms | 25 MBps |
| Network EU West <-> SG | 160 ms | 25 MBps |

## Optimisation Technology Levels

The different parts that forms software execution:
1. system - The highest level of abstraction. It is made up of multiple processes, which can be distributed or not. Each process will be using modules
2. module - Encapsulates certain functionality behind an API. Modules are implemented using data structures and algorithms
3. code implementation - The same data structure and algorithm can be implemented differently at the code level. Also, depending on the language, we also depend on the compiler and runtime.
4. operating system - the interface that runs our software, and communicate with underlying machine
5. hardware - the actual machine and various components like CPU, memory, hard disk, network

recognising the levels matters
- optimisation within a level can achieve speedups with factors of 10 to 20
- to achieve better results beyond that, you will need to work on other levels
- our teams may be biased to only optimise within the level we are comfortable with, which may not be the easiest solution

## Efficiency-Aware Development Flow

Here is a recommended Test-Fix-Benchmark-Optimise workflow (TFBO):
1. at the start of the project, set functional and efficiency requirements
2. after development, test and assess that all functional requirements are met (top priority)
3. if not, we must fix functional issues first, then go back to step 2.
4. once functional test passed, benchmark and assess that efficiency requirements are also passing
5. if not, analyse and identify bottlenecks
6. choose a technology level to optimise, and implement reasonable/deliberate optimisation
7. go back to step 2.
8. only after both functional and efficiency tests have passed, then we release the feature.

# CPU Usage

- assembly instructions
	- we can use Go tooling to disassemble our compiled code into assembly dialects (depends on our target architecture)
	- assembly instructions are sequential
	- typically moving data in/out of registers
	- computing data from input registers and sending to output registers
- Go compiler
	- tokenises source code, builds Abstract Syntax Tree (AST)
	- optimises by removing dead code
	- performs escape analysis to determine what variables can allocate on the stack/heap
	- performs function inlining (which avoids heap memory usage)
	- AST is converted to Static Single Assignment form (SSA) for further machine-independent optimisation
	- SSA is then converted to target machine code for specific Instruction Set Architecture (ISA) and OS
	- apply further ISA and OS specific optimisation
	- finally package all the machine code and debug info into a single object file
	- debug info can be omitted from build, to shrink the size of output if necessary
	- we can also choose to build with compiler options that shows the optimisations made by the compiler

## CPU internals

- How CPU works
	- a CPU interacts with the main memory RAM and other I/O devices
	- a CPU consists of multiple physical core
	- each core contains an Arithmetic Logic Unit (ALU), registers, and a hierarchy of caches (L1, L2, L3)
	- registers are the fastest but smallest local storage and are only used for short-term variables or for internal working of the CPU
	- the caches are on-chip SRAM, which are closer and faster for ALU to access than the main memory RAM
	- the core runs in cycles, it can only perform one instruction per cycle (Single Instruction Single Data SISD)
		- some CPU can perform Single Instruction Multiple Data (SIMD)
		- parallel processing of the same instruction but on different data
	- but the bottleneck in modern computing is memory, not CPU cycles
- The "memory wall" problem
	- reading data from main memory to the CPU is much slower compared to the number of instructions modern CPU can process in the same period of time
	- therefore, the CPU is often waiting for data to arrive, before it can process
- Hierarchical Cache System
	- data is read from main memory in contiguous blocks, beyond what is requested
	- these data are stored in the CPU L-caches
	- CPU will attempt to read the required data for the next cycle from the L-cache first
	- therefore, if data structure used in our program are stored sequentially in memory, the CPU will have a higher chance of a cache hit
- CPU Pipelining
	- executing one assembly instruction may actually take more than one cycle
	- there are a couple of stages to the instruction, e.g. fetching instruction, decoding instruction, executing, and then storing results, which takes 4 cycles
	- but each stage is physically handled by a different physical component of the CPU, and all physical components can be running within the same cycle
	- with pipelining, the CPU is smart enough to overlap the instructions, and process a different stage for a different instruction within one cycle
	- e.g. instruction one is currently in decoding stage, so the CPU can also perform the fetch stage for the second instruction
	- this improves throughput, which is why we said modern CPUs can process "one instruction per cycle"
- Out-of-Order execution
	- the CPU is also smart enough to use availability of input data to schedule some instructions before others when possible, with the help of an internal queue
	- the objective is to make the CPU doing as much work as possible with least amount of delay waiting for data
- Speculative Execution
	- some CPU also has the feature to guess the most likely branch of execution and perform that ahead of time
- Hyper-Threading
	- is Intel's proprietary name for Simultaneous Multithreading which other CPU makers also implements
	- allowing one physical CPU core to operate in a mode visible to the OS as multiple separate logical cores
	- the CPU core internally will help schedule instructions headed for the multiple logical cores
	- objective is to keep the core as busy as possible, as we know the main bottleneck is waiting for data
	- the CPU core will handle context switching internally

## Thread Scheduling

- Linux OS scheduler (Completely Fair Scheduler CFS)
	- each process has a dedicated memory space and unique pid
	- a process can have multiple threads, which are the smallest scheduling unit. threads do not share machine code sequences, but they share context with the parent process.
	- the Linux scheduler performs preemptive thread scheduling (freeze thread execution at any time)
	- context needs to be switched when running a different process on the CPU
	- CPU time is cut into slices, and allocated to each thread fairly
	- if a thread is waiting for external events which may be slow (I/O reads) it may voluntarily yield its CPU time
	- the more threads there are, the less time each thread will be allocated
	- for strict realtime requirements (threads that never wants to be preempted), your system might need a different OS with realtime scheduler
- Go Runtime Scheduler
	- Go runtime will make full use of allocated CPU time but switching between goroutines
	- from the perspective of OS scheduler, the Go program is a single process
	- this avoids expensive context switching by the actual OS scheduler
	- all goroutine scheduling happens at the Go runtime application level
	- goroutine have a flat hierarchy therefore all of them share the same memory context
	- Go runtime will multiplex goroutines onto any allocated OS threads
	- when an OS thread is preempted, Go runtime saves the goroutine's execution state, and is able to schedule that same goroutine to run on a potentially different OS thread in a multi-threaded OS environment

## Impact on our code

- if our code is too dynamic and does not allow the compiler to check some memory allocation statically, then the program will depend on runtime checks, which will have performance impact
- contiguous memory data structures are preferred as it empowers existing CPU cache hits and optimisations
- if our code has less conditional branching logic, the CPU is also better able to pipeline and optimise the program execution
- goroutines and concurrency makes our program more complex to maintain and debug and measure performance benchmark. it should be considered as a last deliberate optimisation option. most optimisation problems can be resolved without using concurrency.

# Memory Usage

- symptoms of memory issues
	- process crashes due to OOM errors
	- program running slower than usual while memory usage is higher than average (could be due to memory pressure from trashing or swapping)
	- program is running slower than usual, with high spikes in CPU utilisation (could be caused by excessive memory allocation and release)
- physical memory
	- usually DRAM chips, powered continuously, using memory cells (transistors) to store 1 bit
	- DRAM is cheaper and easier to produce than SRAM, but is slower
	- the most popular DRAM is from the SDRAM family, 5th gen called DDR4
	- the memory is byte addressable, which means each byte has an address represented by an integer (therefore on a 32-bit system, largest integer is 2^32, and these systems typically cannot handle RAM with more capacity than 4 GB)
	- since each byte is addressable, the memory supports random access (hence called RAM)
	- but the industry has been focusing on other improvements, so random access remains a slow operation
		- focuses on higher capacity to store more data
		- more bandwidth and lower latency for reading and writing large chunks of contiguous memory
		- lower voltage since all components in a machine are competing for power. lower power also leads to lower heat and better battery life for devices
		- production cost of RAM must be kept low since they are produced in large quantities and used in all systems

## OS memory management

- challenges for memory management from OS perspective
	- each process needs their dedicated memory space
	- external fragmentation of memory can occur, since many processes are dynamically reserving memory, which can leads to pockets of empty and unusable space over time
	- memory isolation and memory safety, to prevent processes from accessing other spaces or unauthorised spaces
	- memory usage needs to be efficient as processes do not use all the memory they request for, so the memory needs to be dynamically managed
- virtual memory
	- each process is given a simplified view of the RAM
	- the kernel will take care of actual physical memory operations, like coordinating with other processes for space, bin packing, defragmentation, security, limits and swaps etc.
	- memory is represented as fixed-size chunks called pages (for virtual memory) and frames (for physical memory)
	- pages are typically 4 KB, but can be changed to larger sizes to match specific CPU operations
	- the kernel maps pages to frames
	- details about a page (mapping, state, permissions, metadata) are stored as an entry in many hierarchical page tables
	- Memory Management Unit (MMU) uses a Translation Lookaside Buffer (TLB) which caches page tables. MMU sits inside the CPU, and it helps to translate CPU memory address (virtual) to physical memory address, to avoid expensive look up for the page tables in the RAM itself.
	- more virtual memory can be allocated to processes, than what is available in physical memory (memory overcommitment)
	- physical memory is only allocated when the process tries to access those virtual memory spaces (on-demand paging)
	- we can inspect the memory usage using `proc`, to look at `VSS` stats (virtual) and `RSS` stats (physical)
	- when the process first tries to access some of those unallocated physical memories, the MMU will trigger a hardware interrupt, and a couple of operations may happen
		- OS may allocate more RAM frames, if they are free and available
		- OS may de-allocate unused RAM frames (trashing)
		- if the OS is still OOM, it can use a swap disk partition to back up virtual memory pages (swapping)
		- the OS can also crash and reboot
		- it can also terminate lower priority processes

## Go memory management

- in the allocated memory space, it holds the following regions
	- program code (.text)
	- initialised data (.data)
	- space for uninitialised data (.bss)
	- heap (which can grow dynamically)
	- shared libraries
	- other explicit memory mappings (depending on whether we made syscall like `mmap` in our code)
	- stack (LIFO structure)
- goroutines will each have their own stack frames, but they all sit inside the stack region of the parent process
- Go is value oriented, which means each variable is allocated at a memory address in the stack frame, and at that memory location, it directly stores the value of that variable.
	- therefore simple assignment like `var1 = var2` will copy the value at memory address of `var2` to the location of `var1`
- during compilation, escape analysis will be run to determine if a variable can be allocated on the stack or on the heap
- during runtime, Go allocator will be invoked to allocate space on the heap, and dynamically manage the heap (bin packing, request for more pages from OS, avoid locking and prevents fragmentation)
- the Go garbage collector offers two configurable options
	- `GOGC` represents GC percentage, by default is 100 
	- GC will be run when heap size expands to 100% of the size it gas at the end of the last GC cycle (GC will estimate the time to run based on heap growth rate)
	- the `GOMEMLIMIT` option, disabled by default, will run the GC when heap memory usage reaches the limit.
	- This may lead to disruptive frequent triggers of GC if your program uses memory above the limit.
	- you can also manually trigger GC in the code. This is usually used for testing/benchmarking.
	- GC too often will be too disruptive to the program
	- GC not frequent enough will cause more memory to be used as old memory are not reclaimed
- the current Go GC implementation is a concurrent, non-generational, tricolour mark and sweep collector
	- first a stop the world event (STW) is triggered to inject a write barrier to all goroutines
	- then use 25% of CPU capacity given to the process to concurrently mark all objects in the heap that are still in use
	- finally remove write barriers on goroutines using another STW event
- when Go allocator wants to allocate new heap memory to a goroutine
	- it will sweep the heap to find unmarked regions that can be used
	- this sweeping task contributes to memory allocation latency, even though it is a functionality of garbage collection
- heap memory is managed using bucketed object pooling
	- maintaining buckets of memory of different sizes
	- buckets from a resource pool, that can be used by different goroutines

## Impact on our code

- observing virtual memory size is not useful, due to on-demand paging by the OS
- it is impossible to tell how much memory a process has used precisely in a given time
- OS memory usage expands to use all available RAM, due to the "lazy" approach used by memory management, even if your program did not use so much memory
- tail latency of our program memory access can be much slower than DRAM latency, if we are hit with unfortunate cache misses, MMU misses, and page swaps/trashing
- slow program could be simply caused by high RAM usage, and CPU is mostly under-utilised
- for Go programs, looking at our heap size is usually a good start to optimise memory
- usually there is no need to explicitly assign `nil` to a variable to free memory, but if your goroutine is long running and the variable remains unused for a long time before next assignment, then early `nil` assignment can indicate to GC to reclaim the memory
- having more objects on heap to GC, means that CPU load will be higher during GC and the collection will be slower
- GC is also destructive to the hierarchical cache system, causing more cache misses that slows down the program overall
- if GC is not fast enough to catch up to the allocation, then it can result in "memory leaks" (memory usage grows)
- the best practice is simply to allocate less memory

# Observability

- auto or manual instrumentation
	- tools automatically generate metrics outside of your process
	- or you need to manually add calls to the tools within your process
- raw events or aggregated data
	- e.g. full HTTP requests responses
	- or count of success/failure HTTP requests
- pulled or pushed from process
	- a centralised remote process periodically collects data from your application process
	- or your process must periodically sends data to a centralised tool

## Common Tools

- Logging
	- simplest tool we can use to get primitive latency of a operation, by measuring duration between start and end within our own process and logging it out
	- sometimes it is more valuable to aggregate the total latency and total number of operations to find the average latency
	- Go benchmark tool from standard library helps to measure that average latency (helpful for very fast operation)
	- at a micro-benchmark level (focusing on our process), logs can provide aggregated stats
	- at a macro-benchmark level (distributed system), logs providing raw events are better for further analysis across the system
	- can consider using third party API or `log/slog` for structured logging
- Tracing
	- will need to depend on third party libraries, like OpenTelemetry
	- uses a root span to trace an entire transaction, and children spans for different parts of the transaction
	- the value comes from context propagation, even across different processes
	- but it is harder to maintain, has risk of vendor lock-in, may be expensive, and it is still challenging to observe very fast/short duration transactions within a process
- Metrics
	- will need to depend on third party libraries
	- aggregated numerical values for some stats of our system
	- this signal is useful for most efficiency analysis, as we can compare numbers before and after optimisations
	- aim to have low metric cardinality to keep the overall instrumentation maintainable

## Metrics Semantics

- metric can be defined by 2 things
	- semantics. what does the number mean, what do we measure, what unit, how do we call it
	- granularity. how details is the information, per operation or per goroutine, or over a time period
- latency measures
	- Go `time` uses nanoseconds, but it may make sense to standardise your metrics to seconds to compare across different tools and processes
	- short and very fast latencies should be measured using average latency (like using Go benchmark)
	- take note, Go runtime is prone to leap second problem (adding extra second to be in sync with the Earth's rotation) and virtual machine suspension (process goes to sleep) which leads to wrong measurement
	- we can measure latency at different levels of granularity
		- end-to-end from client side
		- HTTP server end-to-end
		- only application logic
		- only specific function in the process
	- each level of granularity has its value and helps identify different bottlenecks, so it may be helpful to measure a couple of them
	- using percentiles are helpful, as averages may be misleading
- CPU usage
	- we can use `proc` on linux to view the stats
	- CPU cycles. total clock cycles used to execute our program
	- CPU instructions. total instructions executed for our program.
	- these measures are good, as they are independent of thread scheduling, and latency of memory fetches
	- CPU time divided by CPU capacity = CPU usage.
		- CPU time is split into user time and kernel time
	- low CPU time could mean a lot of waiting for other slower events like I/O
	- CPU usage can tell you if the program is reaching CPU saturation, which may be an issue if the program cannot handle additional spike in usage depending on your use case
- memory usage
	- use Go `runtime/metrics` to collect insights about GC, memory allocation, heap etc.
	- use `proc` to look at OS thread memory (but be skeptical due to the complexity and laziness of OS memory management)
		- VSS stands for virtual set size, number of pages allocated to process
		- RSS is residential set size, is number of pages resident in RAM
		- PSS is proportional set size, is number of shared memory pages divided equally among all users
		- WSS is working set size, estimating number of pages currently used to perform work by our process

# Efficiency Assessment

- complexity analysis can be a preliminary estimate of our function's behaviour
	- time and space asymptotic complexity as a function of some N input
	- upper bounds, average, and lower bounds are useful to help us estimate the limits of our function
	- we can easily have a quick estimate of how much resources is needed for any given input size N
	- we can also easily figure out which is the critical bottleneck of our algorithm
	- if requirements are simple but complexity is high, it is an indication that our algorithm might be bad
	- if complexity scales faster than inputs, then we might have scaling problems in the future
	- if measured and observed resource usage deviates significantly from complexity analysis, then we might be having efficiency issues (memory leaks, other poor implementations)
- benchmarking is made up of 4 components (also called stress or load testing)
	- number of iterations
	- experiment on the system for each iteration
	- measure and observe resource stats for each iteration
	- finally compare results from all iterations
- benchmarking is harder than functional testing
	- we need to have many different test cases and test data
	- the performance is often nondeterministic (due to the complex environment)
	- it is expensive to maintain the tests
	- it is difficult to assert what is correct or good enough, and requirements evolve over time
- common challenges to reliable benchmark tests
	- human errors
		- keep our SOP simple
		- track what version of system and benchmarking tools we are running.
		- keep our work well documented and organised
		- be skeptical about results that appears too good to be true
	- reproducing production 
		- it is difficult to reproduce production conditions and workload, due to cost concerns and feasibility
		- we can aim to only simulate key characteristics from production in our benchmark tests
		- we can also focus only on the use cases that our users truly care about
	- performance nondeterminism
		- the complex environment we are operating in can often affect our measures
		- if variance of our results is low, perhaps this is not such a big issue
		- ensure stable state of machine we are benchmarking on (no background threads, thermal scaling, power management)
		- if running on share infra, be skeptical of noisy neighbours
		- cheap CI cloud platform runners are usually intended to run functional tests, not to produce reliable benchmarks
		- benchmark machines may be having different limits from production
		- run experiments longer to amortise the impact of overheads and noise
		- expire older benchmark results, as system conditions is evolving over time
- benchmarking levels
	- benchmarking in production
		- the most accurate way but it is challenging
		- we must not impact real users
		- feedback loop is long (fully deploy the system before we can test it, and then make changes to deploy again)
	- macro-benchmark (or testing in a sandbox/staging/UAT environment that is quite similar to production)
		- can be reliable and effective, and doesn't impact production
		- this may be expensive and difficult to maintain to keep in sync with production
		- feedback loop is also long
	- micro-benchmark (testing an isolated part of our process)
		- fastest feedback loop
		- easy and cheap to implement, like a unit test
		- may not be able to identify all bottlenecks
- it is recommended to start small with micro-benchmark, and move to higher level tests if necessary for your requirements

# Benchmarking

## Micro-benchmark

- should be kept micro. consider other levels if you violates some of the assumptions
	- testing single functionality
	- short and fast
	- not too many goroutines are created
	- uses resources that are within the limits of development machine

```go
// filename must end with _test.go suffix

// test function must start with Benchmark prefix
// test function must accept (b *testing.B) as the only parameter
func BenchmarkOurFn(b *testing.B) {
	b.ReportAllocs() // trace memory allocations, equivalent to -benchmem
	
	InitOurFn()
	
	b.ResetTimer() // reset timer to ignore above overhead init
	
	// b.N is an optimal number that benchmark tool will decide
	// therefore, we should avoid using variable i within our loop
	for i := 0; i< b.N; i++ { 
		OurFn()
	}
	
	// can use b.ReportMetric to emit custom metrics
}

// results from test table can also be analysed by benchstat
func BenchmarkOurFnTestTable(b *testing.B) {
	for _, tcase := range []struct { input int } {
		{ input: 1 },
		{ input: 99 },
		{ input: 2048 }
	} {
		b.Run(fmt.Sprintf("input-%d", tcase.input), func(b *testing.B)) {
			InitOurFn(tcase.input)
			for i := 0; i< b.N; i++ { 
				OurFn()
			}
		}
	}
}
```

```sh
export ver=v1
go test -run '^$' \ # regex matches nothing
	-bench '^BenchmarkOurFn$' \ # targets my particular benchmark
	-benchtime 10s \ # run the benchmark test for 10s
	-count 5 \ # repeats the benchmark test 5 times
	-cpu 4 \ # sets the GOMAXPROCS
	-benchmem \ # trace memory allocations
	-memprofile=${ver}.mem.pprof \ # output mem profile
	-cpuprofile=${ver}.cpu.pprof \ # output CPU profile (may have impact on ultra fast tests)
	| tee ${ver}.txt # save to temporary file

# recommended to use community library `benchstat` to further analyse the results
benchstat v1.txt
benchstat v1.txt v2.txt # compares results
```

- relating back to TBFO workflow
	- execute micro-benchmark and save results for current performance
	- analyse and figure out bottlenecks, implement improvements in a new branch
	- make sure new branch passes functional tests
	- re-run benchmark and compare results
	- compile all notes and analysis in pull request (but discard raw results, as results should expire!)
	- benchmark code can be committed
	- provide code comment or document in the PR description, how to replicate the same conditions used to run the benchmark (so that someone else may repeat the experiment in the future)
- micro-benchmark is not able to reveal insights about memory efficiency, particularly:
	- GC latency
	- maximum memory usage
- but it is good enough to let us know the number of memory allocations and we can already start improving on that
- compiler optimisation may lead to unexpected results for our benchmark
	- the same function is run in a loop, compiler may inline the function
	- some inputs are constant and doesn't change, the compiler may cache the results instead of performing the work at runtime
	- the outputs from our function are often discarded in the benchmark, the compiler may deem this as unused and not even run the function at all
- there are some tricks to trick the compiler to not be optimal or lazy
	- using a global exported variable to assign a constant value. the compiler must allocate memory, since the variable may change value (by a goroutine for example) and the compiler cannot predict that runtime dynamics
	- consuming the global variable in our function as input parameter
	- assigning output of our function to another global exported variable
- but these tricks conflict with benchmark principle: to be as close to production as possible. it is a delicate balance.

## Macro-benchmark

- testing our product in a deployed environment similar to production
- ways to handle dependencies (e.g. database)
	- use a realistic dependency, this is the best approach
	- use a fake, this is hard to simulate and maintain
	- use a substitute, like local storage. it is not as accurate but may be sufficient depending on your requirements
- an observability tool, third party platform, will need to be used (no Go built-in tool for this scale)
	- collect metrics from load tester
	- collect metrics from our product
- we also need a load tester to send requests to our product
	- consider the k6 open source project
- we also need a framework to help us orchestrate the test
	- the book recommends its own package https://github.com/efficientgo/e2e
- metrics to look out for
	- server-side latency
	- CPU time
	- memory, whether there are memory leaks and impact of GC
- common practices
	- load test the product at a target RPS and observe how much resources is required to maintain a certain level of p90 latency
	- run load tester from a different location to simulate realistic client
	- deploy to remote servers instead of running on local machine
	- use realistic dependencies
	- scale and test your product on multiple replica to observe if load balancing works
- typical workflow
	- plan what components to use for macro-benchmarking
	- commit all details of the set up in code repo (test framework is likely sitting on a different repo)
	- clearly document experiment details, like environment conditions and software versions (use Google doc for example)
	- perform benchmark test
	- confirm there are no functional errors
	- save load tester results to the same document
	- compile insights from various metrics and stats
	- compile profiling results
	- analyse and discover bottleneck, and implement improvements in a new branch of the system
	- run the workflow again, and compare benchmark results
	- merge branch and release if results are good
- since macro-benchmark are more expensive and intensive, use it when required.
	- recommended as an acceptance test against RAER specs of the entire system after a major feature/release
	- when debugging and optimising regressions or incidents that trigger efficiency problems

# Bottleneck Analysis

## Profiling in Go

- using `pprof` has the advantage of using common representation, file format, and visualisation for different resource profile data
- `runtime/pprof` can be used out of the box to profile the runtime
- `google/pprof` can be used to read the output profile
	- can also display the profiles in the web browser in graphical form
- profiling can be as significant as the other 3 signals (log, trace, metric) when analysing efficiency
- we tend to capture profiles during benchmarking
- there are 3 main patterns to capture profiles 
	- instrument directly in the program code, supports custom profile
	- triggers profile output when running Go benchmark (limited by what Go provides, which are CPU and memory)
	- using HTTP handlers (make your Go program accept a request for profile data)
- the Go standard library `net/http/pprof` already provides HTTP handlers you can easily add to your server
- you can connect the `pprof` web browser tool to your server to view the profiles
- common profilers to note
	- heap profiler
	- goroutine profiler
	- CPU profiler (this is expensive to gather, so it needs to be explicitly started and stopped. expect latency when requesting from HTTP handler. and the data will be sampled)
	- Off-CPU time (block profiler) shows time our program spent waiting without utilising the CPU
- advanced tips
	- profiles can be exported and shared to other team members
	- continuous profiling by having a separate platform to periodically capture profiles from your server using the HTTP handlers (this can be done in production, or only during macro-benchmark tests)
	- `pprof` format supports comparing and aggregating profiles
		- one profile (e.g. tested using input A) can be subtracted from another profile (tested using input A+B)
		- shows delta in different data across two profiles
		- profiles can also be merged (this is not natively supported by `go tool pprof`)

# Optimisation Tips

- optimise and benchmark one improvement at a time
- standard functions from Go may not be designed for your use case, and it can be a source of inefficiency
- standard functions may generally work on multiple types, and you may want to reimplement them for only a specific type that is relevant to your use case to be more efficient
- you can potentially use `unsafe` to perform deliberate optimisation, but we must be careful and understand the tradeoff (memory safety) and whether it truly meets our efficiency requirements
- stream data and process instead of allocating the entire space in memory
- when using concurrency to distribute work to improve latency, take note and observe communication overhead
	- a lot of time may be spent sending data over channels and queuing for channels
	- sometimes, a distributed algorithm that guarantees no conflict without using channel communication may be a better choice to distribute work
- sometimes, we can simply think out of the box instead of performing hard optimisation (e.g. cache the results that are frequently requested. this amortises any initial resource costs over repeated operations)

# Common Patterns

- Do less work
	- skip unnecessary logic (e.g. double validating the same property)
	- do things once (e.g. reuse memory in place, instead of allocating new memory)
	- leverage math to do less (e.g. distributed algorithm instead of communication over channel)
	- use pre-computed information (e.g. request for stats in API instead of recalculating stats)
	- but be strategic, if you are just lazy, then you may incur more work in the future with tech debts
- Trading functionality for efficiency
- Trading space for time
	- precompute results in a look up table
	- caching
	- augment data structure to contain more stats and metadata
	- decompressed data
- Trading time for space
	- usually the exact opposite of the above patterns
- The Three Rs Optimisation Method
	- Reduce Allocation (e.g. prevent escape to heap)
	- Reuse Memory (e.g. use already allocated objects for loops)
	- Recycle (this is handled by GC, but we can prompt it)
		- optimise allocated struct (with less internal pointers)
		- GC tuning (using the two options to GC more often, but is very brittle as our program evolves)
		- manually trigger GC (generally not encouraged)
		- allocating memory off-heap (using syscall like `mmap`)
- avoid resource leaks
	- we use an unbounded resource for the same amount of load
	- eventually the resource runs out
	- this can easily happen with a slice with short length but large underlying array
	- this can easily happen with forgotten long-lived goroutines holding reference to large data in the heap
- control goroutine lifecycle to prevent resource leaks
	- know when and how the goroutine will terminate
	- will the goroutine run indefinitely? in what scenarios?
	- should caller function wait for all goroutines to finish? then use `sync.WaitGroup`
	- if you can't answer some of these questions, then there is a potential for resource leak
	- there are libraries like `goleak` to help test for resource leaks
- reliably close resources
	- always read docs. if the resource needs to be closed, always close it. no linters to help with this.
	- using `defer` to close resource may fail. we should be notified using some kind of error capture or logger.
	- sometimes, we need to pass resource handlers around and we cannot `defer` to close them within the same function that created them. We can append any resource closers to a list on creation, and before program exits we can loop through the list and close each resource.
	- some resources don't have closer, but requires us to exhaust/drain the resource to release it
- pre-allocate when possible
	- this allows our program to operate on bounded resources
	- we may potentially avoid dynamic resize/memory allocation during the lifetime of the program
	- even for implementing linked list, where each node points to the next (requires the use of pointers), we can still optimise it
		- for example, pre-allocating an array of to store the nodes, to act like a resource pool
		- nodes can still be created to point to another node, but all nodes are assigned in the array
		- there will not be any nodes floating in the heap, but we are bounded by the size of the array
		- memory will be contiguous
- overusing memory in arrays
	- pre-allocating memory can also lead to overuse, if the array is large but our slice only use a small part of the memory
	- we may need to implement methods that helps to dynamically copy data to a smaller array
- memory reuse and pooling
	- resource pooling can help avoid repeated allocations (which improves allocation latency)
	- `sync.Pool` provides a form of memory pooling that is thread-safe
	- it is expected that all values in the pool are functionally identical. therefore a memory pool is not a cache where values are different
	- however, the `sync.Pool` is intended for short duration usage, and GC may potentially wiped out all unused resources in the pool, and requesting from the pool will result in new allocations.