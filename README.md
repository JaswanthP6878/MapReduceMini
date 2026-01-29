### Map Reduce - mini

---
This project implements a **single-node, multi-threaded MapReduce framework**, inspired by the original Google MapReduce model, but designed to run entirely on a single machine using **concurrent execution**.

related blog post regarding implementation can be found [here](https://jaswanthp6878.github.io/blog/second-post/)

---

## Architecture Overview

Although the system runs on a **single node**, it follows a **logical master–worker architecture**:

### Master
- Runs as a dedicated goroutine
- Responsible for task scheduling, phase transitions, and synchronization
- Tracks progress of map and reduce workers
- Orchestrates the transition from **Map → Reduce** phase

### Workers
- Implemented as independent goroutines
- Execute either map or reduce tasks assigned by the master
- Communicate task completion and status back to the master
- Block or wait when required (e.g., at phase barriers)

This design mirrors a distributed MapReduce system while leveraging **Go’s goroutines and channels** for lightweight concurrency and coordination.

---

## Execution Model

### Map Phase
- Input data is split into logical chunks
- Each map worker goroutine:
  - Applies the user-defined `map` function
  - Emits intermediate key–value pairs
- Intermediate data is **partitioned into buckets** using a hash of the key
- Each bucket corresponds to a future reduce task

### Synchronization Barrier
- The master enforces a **global barrier** between phases
- Reduce workers do not start until **all map workers complete**
- Workers block or wait on synchronization primitives until the master signals phase completion

### Reduce Phase
- Reduce worker goroutines:
  - Read intermediate data from assigned buckets
  - Group values by key
  - Apply the user-defined `reduce` function
- Output is written to final result files
---

## Concurrency and Threading Model

- **Concurrency primitive**: Go goroutines
- **Synchronization**:
  - Phase barriers ensure correctness
  - Workers wait for master signals before transitioning phases
- **Thread safety**:
  - Shared metadata (task state, worker status) is protected using synchronization mechanisms
- **Parallelism**:
  - Multiple map tasks and reduce tasks execute concurrently on a single node

---


### Progress
- [x] Map part and Intermediate Data generation 
- [x] IR data partitioning into Buckets
- [x] Transitioning all workers from map phsae to reduce phase (wait on other map workers to complete task)    
- [x] Reduce phase 
- [x] Generating Output files
---

> Based on the original map-reduce [paper](http://nil.csail.mit.edu/6.824/2020/papers/mapreduce.pdf)