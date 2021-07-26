# Ghost's Parser

Ghost uses a recursive decent parser. It starts from the top or outermost grammar rule (expression) and works its way down into the nested subexpressions before finally reaching the leaves of the syntax tree.

A recursive decent parser is a literal translation of the grammar's rules straight into imperative code. Each rule becomes a function. The decent is described as "recursive" because when a grammar rule refers to itself -- directly or indirectly -- that translates to a recursive function call. Recursive decent parsers are fast, robust, and can handle sophisticated error reporting while being easy to code and maintain.

## Precedence order

| Rule | Forwards |
|------|----------|
| `parse` | `-> declaration` |
| `declaration` | `-> classDeclaration, letDeclaration, functionDeclaration, statement` |
| `classDeclaration` | `-> "class" IDENTIFIER "{" function* "}"` |
| `functionDeclaration` | `-> "function" IDENTIFIER "(" parameters? ")" block` |
| `letDeclaration` | `-> "let" IDENTIFIER ( "=" expression )` |
| `statement` | `-> expressionStatement, ifStatement, whileStatement, printStatement, blockStatement` |
| `expressionStatement` | `-> expression` |
| `expression` | `-> or` |
| `or` | `-> and ( "or" and )` |
| `and` | `-> ternary ( "and" ternary )` |
| `ternary` | `-> equality ( equality "?" expression ":" expression )` |
| `equality` | `-> comparison ( ( "!=", "==" ) comparison )*` |
| `comparison` | `-> term ( ( ">", ">=", "<", "<=" ) term )*` |
| `term` | `-> factor ( ( "-", "+" ) factor )*` |
| `factor` | `-> unary ( ( "/", "*" ) unary )*` |
| `unary` | `-> ( "!", "-" ) unary, call` |
| `call` | `-> primary ( "(" arguments? ")" )` |
| `arguments` | `-> expression ( "," expression )` |
| `primary` | `-> NUMBER, STRING, "true", "false", "null", "(" expression ")" ` |