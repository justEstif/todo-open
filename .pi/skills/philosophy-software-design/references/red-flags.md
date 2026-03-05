# Software Design Red Flags Reference

Quick-reference for detecting design problems. Ordered by frequency of occurrence in practice.

## Shallow module

**Signal**: Interface is nearly as complex as implementation. Many small classes/methods that individually do very little.

**Detection**: Count the public API surface vs. lines of implementation. If ratio is close to 1:1, the module is shallow.

**Example**: A class with 6 getter/setter methods wrapping 6 private fields, adding no logic.

**Fix**: Merge with related functionality to create a deeper module. Ask "what complexity is this hiding?"

## Pass-through method

**Signal**: A method that does almost nothing except invoke another method with a similar or identical signature.

**Detection**: Method body is 1-3 lines, mostly delegating to another object's method with same parameters.

```
// Red flag
class OrderService {
    createOrder(items, userId) {
        return this.orderRepository.createOrder(items, userId);
    }
}
```

**Fix**: Expose the lower layer directly, redistribute functionality so the method adds real value, or merge the layers.

## Pass-through variable

**Signal**: A variable passed through a long chain of methods but only used deep in the call stack.

**Detection**: Parameter appears in 3+ method signatures but is only consumed by the deepest method.

**Fix**: Use context objects, dependency injection, or restructure to reduce the chain.

## Information leakage

**Signal**: Two or more modules depend on the same piece of knowledge (file format, protocol, data structure).

**Detection**: Changing one internal detail requires changes in multiple modules. Look for:
- Shared format/protocol knowledge
- Interface types that mirror internal structures
- Constructor parameters exposing implementation choices

**Fix**: Consolidate the shared knowledge into a single module. One of the two modules should own the knowledge entirely.

## Temporal decomposition

**Signal**: Code organized by execution order rather than logical grouping. Common in read-process-write pipelines where reading, processing, and writing are in separate modules even though they share data structure knowledge.

**Detection**: Module names reflect sequence ("first", "then", "after") or pipeline stages rather than responsibilities.

**Fix**: Group by information — all code that deals with a particular data format or concept should live together.

## Overexposed configuration

**Signal**: Configuration parameters that most users don't need, pushed to callers instead of using sensible defaults.

**Detection**: Functions with 5+ parameters where most callers pass the same values. Constructor requiring many options that could have defaults.

```
// Red flag: pushing complexity to every caller
createServer(port, host, backlog, keepAlive, timeout, maxHeaders,
             headerTimeout, requestTimeout, maxConnections)

// Better: defaults + optional overrides
createServer({ port: 3000 })
```

**Fix**: Use sensible defaults. Only expose configuration callers actually vary.

## Conjoined methods

**Signal**: Can't understand method A without reading method B. Methods are intellectually joined even if structurally separate.

**Detection**: Reading one function requires constantly jumping to another to understand behavior. Shared mutable state between methods with no clear contract.

**Fix**: Make each method self-contained with a clear interface contract, or merge them if they're truly one logical operation.

## Special-general mixture

**Signal**: A general-purpose module contains special-case code for a specific use case.

**Detection**: `if` branches or parameters that serve only one caller. Comments like "this is for the billing page".

**Fix**: Keep the general mechanism clean. Let the specific use case implement its specialization externally.

## Repetition

**Signal**: Same pattern of code repeated in multiple places with minor variations.

**Detection**: Copy-pasted blocks with small differences. Bug fixes that need to be applied in multiple places.

**Fix**: Extract common pattern into a shared abstraction — but only if the pattern is stable and repeated 3+ times. Premature abstraction of 2 occurrences often isn't worth it.

## Non-obvious code

**Signal**: Reader cannot quickly understand what code does or why.

**Detection**: You need to trace through multiple files or hold many things in memory to understand a single operation. Generic containers used instead of named types.

```
// Non-obvious
const result: [string, [number, boolean]] = process(data);

// Obvious
interface ProcessResult {
    userId: string;
    score: number;
    isActive: boolean;
}
const result: ProcessResult = process(data);
```

**Fix**: Better naming, extract named types, add comments explaining "why", simplify control flow.

## Vague naming

**Signal**: Names that are too generic to convey meaning (e.g., `data`, `result`, `handle`, `process`, `manager`, `info`, `tmp`).

**Detection**: Reading the name doesn't tell you what specific thing it represents. Multiple variables could swap names without confusion.

**Fix**: Choose names that are precise about what the thing IS or DOES. If it's hard to name, the abstraction might be unclear.

---

## Common Mistakes

Quick-reference of the most frequent design mistakes and their fixes.

| Mistake | Why it fails | Fix |
|---------|-------------|-----|
| Too many small classes | Each boundary adds cognitive overhead without adding depth | Merge related shallow classes into deeper modules |
| Splitting by execution order | Forces shared knowledge across phases (read/process/write) | Organize by information, not by when things happen |
| Exposing implementation in interfaces | Callers depend on internals, changes propagate everywhere | Design interfaces around abstractions, hide format/protocol details |
| Treating comments as optional | Design intent and assumptions are lost | Write interface comments first, maintain them with the code |
| Configuration parameters for everything | Each parameter pushes a decision to the caller | Use sensible defaults, minimize required configuration |
| Quick-and-dirty tactical fixes | Each shortcut adds complexity that compounds over time | Invest 10-20% extra in good design per change |
| Pass-through methods everywhere | Adds interface surface without adding depth | Merge into caller or callee, or redistribute real functionality |
| Designing for one specific use case | Special-purpose interfaces accumulate special cases and bloat | Ask "what is the simplest interface that covers all current needs?" |
