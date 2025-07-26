# Antlr4-C3 / ANTLR-NG – API Reference

> **Package entry point:** `import * as c3 from "antlr4-c3";`
>
> All names shown below are re-exported via the package **index.ts** and can therefore be imported directly from the main package entry.
>
> ```ts
> import {
>   CodeCompletionCore,
>   SymbolTable,
>   ScopedSymbol,
>   ClassSymbol,
>   VariableSymbol,
>   TypeAlias,
>   ArrayType,
>   ReferenceKind,
>   MemberVisibility,
> } from "antlr4-c3";
> ```
>
> The documentation is split into three logical groups:
>
> 1. **C3 Engine** – the grammar-agnostic code-completion algorithm (`CodeCompletionCore`).
> 2. **Symbol & Type System** – a lightweight language-independent symbol table implementation.
> 3. **Shared Utility Types & Errors** – enumerations, helper interfaces and exceptions.

---

## 1. C3 Engine

### 1.1 `class CodeCompletionCore`
Provides grammar-agnostic collection of completion candidates for a given caret position. Operates directly on the ATN that powers every ANTLR/ANTLR-NG parser.

**Constructor**
```ts
new CodeCompletionCore(parser: Parser);
```
* `parser` – an *already configured* parser instance (token stream attached). The input does **not** have to be parsed yet.

**Important public fields**
* `ignoredTokens: Set<number>` – token types that must never appear in the candidate set (e.g. operators).
* `preferredRules: Set<number>` – parser rules that should replace the plain token candidate (e.g. `functionRef` → ask symbol-table for functions).
* `translateRulesTopDown = false` – when `true`, the rule list returned for a candidate is ordered from outer-most to inner-most rule.
* Debug flags: `showResult`, `showDebugOutput`, `debugOutputWithTransitions`, `showRuleStack`.

**Method**
```ts
collectCandidates(caretTokenIndex: number, context?: ParserRuleContext): CandidatesCollection;
```
* `caretTokenIndex` – token index that visually contains the caret.
* `context` *(optional)* – restricts the search space to a specific parser rule (huge speed-up for large grammars).

**Return type – `CandidatesCollection`**
```ts
class CandidatesCollection {
  tokens: Map<number, number[]>;             // <tokenType, followingTokens[]>
  rules:  Map<number, {
            startTokenIndex: number;
            ruleList: number[];              // Rule call-stack
          }>;
}
```

### 1.2 Quick Example
```ts
import { CodeCompletionCore } from "antlr4-c3";

const parser = new ExprParser(new CommonTokenStream(new ExprLexer(CharStreams.fromString(source))));

const core = new CodeCompletionCore(parser);
core.ignoredTokens = new Set([ ExprLexer.PLUS, ExprLexer.MINUS ]);
core.preferredRules = new Set([ ExprParser.RULE_variableRef, ExprParser.RULE_functionRef ]);

const caretIndex   = 17;                     // obtain from IDE/editor
const candidates   = core.collectCandidates(caretIndex);

// Convert lexer tokens → user visible text
const keywords = [...candidates.tokens.keys()].map(t => parser.vocabulary.getDisplayName(t));
```

---

## 2. Symbol & Type System
Lightweight but powerful symbol-table implementation that works for any programming-language grammar. All symbol classes ultimately derive from **`BaseSymbol`** and many also implement **`IType`**.

### 2.1 Core Infrastructure
| Class | Responsibility |
|-------|----------------|
| `BaseSymbol` | Name, visibility, modifier set, parent/child traversal & resolution helpers. |
| `ScopedSymbol` (implements `IScopedSymbol`) | A `BaseSymbol` that can own *children* (namespace, class, method …). Provides `addSymbol`, `resolve`, `getAllSymbols` etc. |
| `SymbolTable` (implements `ISymbolTable`) | Top-level scope. Manages dependencies to other tables and convenience helpers such as `addNewNamespaceFromPath`. |
| `DuplicateSymbolError` | Thrown when `ISymbolTableOptions.allowDuplicateSymbols` is `false` and a duplicate is detected. |

#### Quick Start
```ts
import {
  SymbolTable, ClassSymbol, MethodSymbol,
  VariableSymbol, FundamentalType, ReferenceKind, MemberVisibility,
} from "antlr4-c3";

const table = new SymbolTable("main", { allowDuplicateSymbols: false });

// ① Add a class
const myClass = new ClassSymbol("ExampleClass");
table.addSymbol(myClass);

// ② Add a field and a method inside the class
const intType = new FundamentalType("int", ReferenceKind.Instance, TypeKind.Integer);

const field = new VariableSymbol("count", intType);
field.visibility = MemberVisibility.Private;
myClass.addSymbol(field);

const method = new MethodSymbol("increment");
method.visibility = MemberVisibility.Public;
myClass.addSymbol(method);

// ③ Resolve symbols
console.log(table.resolveSync("ExampleClass")?.qualifiedName());        // → "ExampleClass"
console.log(myClass.resolveSync("count") === field);                    // → true
```

### 2.2 Pre-defined Symbol Flavours
* `BlockSymbol` – anonymous scoped block (e.g. `{ ... }`).
* `NamespaceSymbol` – namespace or module container.
* `ClassSymbol`, `InterfaceSymbol` – OO types.
* `FieldSymbol`, `VariableSymbol`, `ParameterSymbol` – data holders.
* `MethodSymbol`, `RoutineSymbol` – callable executables.
* `LiteralSymbol` – constant literal values.
* `TypedSymbol` – generic symbol carrying an `IType` *(base for most concrete symbols).*  
* `TypeAlias` – language constructs like `typedef` or `using`.

### 2.3 Type Helpers
* `FundamentalType` – built-in primitives (`int`, `string`, …).
* `ArrayType` – typed array (`elementType`, optional fixed `size`).
* `IType` (interface) – implemented by every type class; provides `kind`, `reference` and `baseTypes`.
* Enumerations: `TypeKind`, `ReferenceKind`.

---

## 3. Shared Utility Types & Errors
| Name | Description |
|------|-------------|
| `MemberVisibility` | Encodes public/protected/private/… semantics across languages. |
| `Modifier` | Flag-set for modifier keywords (`static`, `final`, `abstract`, …). |
| `ISymbolTableOptions` | Currently only `allowDuplicateSymbols`. |
| `utils.longestCommonPrefix` | Returns longest common prefix of two arrays. |

---

## Error Handling
Almost every `addSymbol` operation can throw `DuplicateSymbolError` when the *target scope* already contains a symbol with the same name and duplicates are *not* allowed.

```ts
try {
  table.addSymbol(new VariableSymbol("x", intType));
  table.addSymbol(new VariableSymbol("x", intType)); // ← boom
} catch (e) {
  if (e instanceof DuplicateSymbolError) {
    console.error(e.message);
  }
}
```

---

## Advanced Example – Integrating C3 With SymbolTable
```ts
// 1. Build/maintain a SymbolTable while walking the parse-tree (listener or visitor pattern) …
//    assume the variable `table` now contains your in-memory model.

// 2. In your editor integration, after every keystroke:
const caretTokenIndex = getTokenIndexForCaret();
const context         = findEnclosingContext(caretTokenIndex);   // optional

const core = new CodeCompletionCore(parser);
core.preferredRules = new Set([
  MyLangParser.RULE_varRef,
  MyLangParser.RULE_funcRef,
]);

const candidates = core.collectCandidates(caretTokenIndex, context);

// 3. Translate tokens → keywords and rules → symbols
const keywordSuggestions = [...candidates.tokens.keys()]
  .map(t => parser.vocabulary.getDisplayName(t));

const symbolSuggestions: string[] = [];
for (const [ruleIndex, { ruleList }] of candidates.rules) {
  if (ruleIndex === MyLangParser.RULE_varRef) {
    for (const v of table.getAllSymbolsSync(VariableSymbol, false)) {
      symbolSuggestions.push(v.name);
    }
  }
  // … handle other rules
}

return [...keywordSuggestions, ...symbolSuggestions];
```

---

## Versioning & Compatibility
The public API follows **semantic-versioning**. Breaking changes will only happen in a new **major** version.

---

### See Also
* [Project README](../readme.md) – in-depth discussion of C3 design & theory.
* [Unit tests](../tests/) – concrete, executable examples.
* [Release notes](../release-notes.md) – detailed change log.

---

© Mike Lischke – MIT License