---
name: philosophy-software-design
description: Apply "A Philosophy of Software Design" (John Ousterhout) principles to code review and architectural design. Dual mode — reviews existing code for design quality AND advises on new module/API/abstraction design. Use when reviewing code, refactoring modules, designing APIs or abstractions, discussing architecture decisions, or when explicitly asked to apply software design principles. Triggers on code reviews, refactoring tasks, module/API design, architecture discussions, and explicit requests like "review software design", "check design quality", "apply philosophy of software design", "review design complexity", "check module depth".
---

# Software Design

Apply John Ousterhout's "A Philosophy of Software Design" principles pragmatically — focus on the highest-impact issues, not dogmatic compliance.

## Mode Selection

Determine mode from context:

- **Review mode**: User asks to review, audit, or refactor existing code
- **Design mode**: User is designing new modules, APIs, abstractions, or architecture
- **Both**: Large tasks involving redesign of existing code

## Review Mode

### 1. Read the code under review

Read all relevant files. Build a mental model of module boundaries, interfaces, and information flow.

### 2. Evaluate against core principles

For each module/class/function, assess (see [references/principles.md](references/principles.md) for details):

| Principle | Key question |
|-----------|-------------|
| **Depth** | Does this module hide significant complexity behind a simple interface? |
| **Information hiding** | Are implementation details leaking across module boundaries? |
| **Cognitive load** | How much must a developer know to use or modify this code? |
| **Change amplification** | Would a simple change require edits in many places? |
| **Obviousness** | Can a reader understand behavior quickly without deep investigation? |
| **Interface comments** | Do comments describe the abstraction, not repeat the code? |
| **Design in reviews** | Are code reviews evaluating design quality, not just correctness? |
| **Strategic investment** | Is ~10-20% of dev time going toward design improvements? |

### 3. Scan for red flags

Check for red flags from the book (see [references/red-flags.md](references/red-flags.md) for complete list):

- Shallow modules (complex interface, little functionality)
- Pass-through methods/variables
- Information leakage between modules
- Temporal decomposition (split by when, not by what)
- Conjoined methods (can't understand one without reading another)
- Overexposed configuration parameters

### 4. Report findings

For each issue found:

```
**[Principle violated]** — severity: high|medium|low
Location: file:line
Problem: What's wrong and why it matters
Suggestion: Concrete improvement with brief rationale
```

Prioritize by impact. Skip low-severity issues unless explicitly asked for exhaustive review. Always explain WHY something is a problem in terms of complexity consequences (cognitive load, change amplification, unknown unknowns).

### 5. Score the design (0-10)

End every review with an overall score:

- **9-10**: Deep modules, clean abstractions, excellent information hiding, comments capture design intent
- **7-8**: Solid design with minor issues — a few shallow modules or small leakage
- **5-6**: Mixed — some good abstractions but notable complexity problems
- **3-4**: Significant issues — shallow modules, widespread leakage, tactical shortcuts
- **0-2**: Pervasive complexity — unknown unknowns everywhere, no clear abstractions

Always state the current score and the **specific changes needed to reach the next level**. Example: "Current: 6/10. To reach 8: merge the 3 shallow service wrappers into a single deep module, and encapsulate the JSON format knowledge that leaks between Parser and Writer."

## Design Mode

### 1. Understand the requirement

Clarify what the module/API/system needs to do. Identify the key use cases.

### 2. Apply "Design it twice"

Propose at least two alternative designs. For each, evaluate:

- **Interface simplicity**: How many concepts must users understand?
- **Depth**: How much complexity does the interface hide?
- **Information hiding**: What decisions are encapsulated vs. exposed?
- **Generality**: Is it somewhat general-purpose without being over-engineered?

### 3. Guide toward deep modules

Steer designs toward:

- **Simple interfaces** with powerful implementations
- **Fewer, more capable classes** over many tiny ones
- **Pulling complexity downward** — make the module's life harder so callers' lives are easier
- **Defining errors out of existence** — choose semantics that eliminate error conditions
- **General-purpose interfaces** with specialized implementations when appropriate

### 4. Validate with red flags

Before finalizing, check the proposed design against [references/red-flags.md](references/red-flags.md). If any red flags appear, iterate.

## Pragmatic Application

These are guidelines, not laws. When reviewing:

- A shallow method that improves readability at call sites is fine
- Comments explaining "why" are valuable; don't demand them for obvious code
- Small helper functions for repeated 3-line patterns are acceptable
- Perfect information hiding is not always worth the refactoring cost

Always weigh the **cost of change** against the **complexity reduction**. Recommend changes only when the net benefit is clear.
