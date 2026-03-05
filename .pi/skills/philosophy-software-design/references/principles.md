# Software Design Principles Reference

## Table of Contents

1. [Complexity and its symptoms](#complexity-and-its-symptoms)
2. [Deep vs shallow modules](#deep-vs-shallow-modules)
3. [Information hiding and leakage](#information-hiding-and-leakage)
4. [General-purpose modules](#general-purpose-modules)
5. [Different layer, different abstraction](#different-layer-different-abstraction)
6. [Pull complexity downward](#pull-complexity-downward)
7. [Define errors out of existence](#define-errors-out-of-existence)
8. [Design it twice](#design-it-twice)
9. [Comments as design tool](#comments-as-design-tool)
10. [Strategic vs tactical programming](#strategic-vs-tactical-programming)
11. [Consistency](#consistency)
12. [Obviousness](#obviousness)

---

## Complexity and its symptoms

Complexity is anything that makes software hard to understand or modify. It manifests as:

- **Change amplification**: A simple change requires modifying many places
- **Cognitive load**: A developer must learn many things to make a change
- **Unknown unknowns**: It's not obvious what needs to change (worst symptom)

Complexity accumulates incrementally. Each "small hack" adds up. Fight it continuously.

**Detection**: If a developer asks "what else do I need to change?" and the answer isn't obvious, there are unknown unknowns.

## Deep vs shallow modules

A module's **depth** = functionality provided / interface complexity.

**Deep module**: Simple interface, hides significant implementation complexity.

```
// Deep: Unix file I/O
// Interface: open, read, write, close (4 functions)
// Hides: disk layout, buffering, caching, permissions, locking, journaling
fd = open("/path/file", O_RDONLY);
read(fd, buffer, size);
close(fd);
```

**Shallow module**: Interface nearly as complex as implementation.

```
// Shallow: adds a layer without hiding complexity
class FileReader {
    constructor(path: string, encoding: string, bufferSize: number,
                retryCount: number, timeout: number) { ... }
    read(offset: number, length: number, callback: Function): void { ... }
}
```

**Guideline**: Prefer fewer, deeper modules over many shallow ones. A class with one method that does something substantial is often better than splitting it into five trivial classes.

## Information hiding and leakage

**Information hiding**: Each module encapsulates design decisions that other modules don't need to know.

**Information leakage**: Implementation details escape module boundaries. Forms:

- Two modules that both depend on a file format (shared knowledge)
- Interface parameters that expose internal data structures
- Back-door leakage through shared state or global variables

```
// Leakage: caller must know internal serialization format
function saveUser(user: User, format: "json" | "msgpack"): void { ... }

// Hidden: module decides serialization internally
function saveUser(user: User): void { ... }
```

**Fix**: If two modules share knowledge, consider merging them or extracting the shared knowledge into a single module.

## General-purpose modules

Build modules that are **somewhat general-purpose**: the interface serves current needs but is general enough to support other uses naturally.

```
// Too specific: only works for undo
class UndoHistory {
    addUndoAction(action: UndoAction): void
    undo(): void
}

// Somewhat general: text module handles any editing, undo uses it
class TextDocument {
    insert(position: number, text: string): void
    delete(range: Range): void
}
```

**Test**: "What is the simplest interface that covers all my current needs?" If it also covers future needs naturally, even better — but don't design for speculative requirements.

## Different layer, different abstraction

Each layer in a system should provide a different abstraction from the layers above and below.

**Red flag**: Pass-through methods — a method that does little except invoke another method with a similar signature.

```
// Bad: TextDocument.insert just calls TextArea.insert
class TextDocument {
    insert(pos, text) { this.textArea.insert(pos, text); }
}
```

**Fixes**:
- Expose the lower layer directly
- Redistribute functionality so each layer does something distinct
- Merge the layers if they don't provide different abstractions

## Pull complexity downward

When complexity is inevitable, push it into the module's implementation rather than its interface. Make the module developer's life harder so every caller's life is easier.

```
// Complexity pushed to caller (bad)
const config = loadConfig("db.yml");
const pool = createPool(config.host, config.port, config.maxConn);
const conn = await pool.acquire(config.timeout);

// Complexity pulled down (good)
const conn = await db.connect(); // module handles config, pooling, retries
```

**Guideline**: Configuration parameters are a form of pushing complexity upward. Only expose them when callers actually need control. Use sensible defaults.

## Define errors out of existence

Reduce exception handling complexity by redefining operations so error conditions can't occur.

```
// Error-prone: delete fails if file doesn't exist
function deleteFile(path: string): void {
    if (!exists(path)) throw new FileNotFoundError();
    // ...
}

// Error-free: delete ensures file doesn't exist (idempotent)
function deleteFile(path: string): void {
    if (!exists(path)) return; // goal already achieved
    // ...
}
```

**Techniques**:
- Make operations idempotent (repeating has no additional effect)
- Use default values instead of throwing on missing data
- Mask exceptions at low levels when possible (e.g., TCP retransmissions hide packet loss)

**Exception**: Errors that genuinely cannot be handled internally must still be reported.

## Design it twice

Before implementing, sketch at least two fundamentally different approaches. Compare:

- Which interface is simpler?
- Which provides better information hiding?
- Which handles edge cases more cleanly?
- Which is easier to evolve?

Even if the first approach seems obvious, the comparison often reveals improvements or hybrid solutions.

## Comments as design tool

Comments should capture information **not obvious from the code**:

- **Interface comments** (what): Describe what a module/function does from the caller's perspective, including preconditions, postconditions, and side effects
- **Implementation comments** (why): Explain non-obvious design decisions, tricky algorithms, or subtle constraints
- **Cross-module comments** (dependencies): Document non-obvious relationships between modules

**Write comments first**: Describe the interface before writing implementation. If the comment is hard to write, the design may be too complex.

**Bad comments** (don't write these):
```
// Increment counter
counter++;

// Returns the user's name
function getUserName(): string { ... }
```

**Good comments**:
```
// Retry with exponential backoff because the payment gateway
// rate-limits burst requests and returns 429 after 5 rapid calls.
async function processPayment(order: Order): Promise<Receipt> { ... }
```

## Strategic vs tactical programming

- **Tactical**: Get it working as fast as possible. Creates complexity debt.
- **Strategic**: Invest in good design. Each piece of code should be clean and contribute to good system design.

Invest ~10-20% of development time in design improvements. This pays compound returns.

**Tactical tornado**: A developer who produces code fast but leaves complexity everywhere. Their output velocity is an illusion — they create more work for everyone else.

## Consistency

Use consistent patterns for names, coding conventions, design patterns, and invariants.

- Name similar things similarly, different things differently
- Use the same pattern for similar operations throughout the codebase
- Document conventions and enforce them

**Consistency reduces cognitive load** — developers apply prior knowledge instead of re-learning each module.

## Obviousness

Code is obvious when a reader can understand it quickly without significant effort.

**Techniques for obviousness**:
- Good naming (specific, consistent)
- Consistent coding style
- Judicious use of whitespace to separate logical sections
- Comments for anything non-obvious

**Enemies of obviousness**:
- Event-driven programming (hard to follow control flow)
- Generic containers (e.g., `Pair<Pair<String, Int>, List<String>>` instead of a named type)
- Code that requires reading other code to understand
