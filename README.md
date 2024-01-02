# `glox`

![glox logo](https://i.ibb.co/gw9rJ7k/glox-logo.png)

Go implementation of the Lox programming language from [Crafting Interpreters](https://craftinginterpreters.com/) by [@munificent](https://github.com/munificent).

## Run glox
```sh
# Build the interpreter.
make all
# Run the interperter on an example file.
./bin/glox examples/fibonacci.g
```

Alternatively, you can run `make example` to build the glox interpreter and run it on `example/fibonacci.g`. 

## Example Script
```lox
// Recursive Fibonnaci implementation.
print "Printing the 10-th Fibonnaci number";

fun fib(n) {
  if (n <= 1) {
    return n;
  }
  return fib(n-1) + fib(n-2);
}

print fib(10);
```


## Grammar
From [Crafting Interpreters](https://craftinginterpreters.com/appendix-i.html).

### Syntax Grammar

```ebnf
program        → declaration* EOFToken ;
declaration    → classDecl
               | funDecl
               | varDecl
               | statement ;

classDecl      → "class" IDENTIFIER ( "<" IDENTIFIER )?
                 "{" function* "}" ;
funDecl        → "fun" function ;
varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
statement      → exprStmt
               | forStmt
               | ifStmt
               | printStmt
               | returnStmt
               | whileStmt
               | block ;

exprStmt       → expression ";" ;
forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
                           expression? ";"
                           expression? ")" statement ;
ifStmt         → "if" "(" expression ")" statement
                 ( "else" statement )? ;
printStmt      → "print" expression ";" ;
returnStmt     → "return" expression? ";" ;
whileStmt      → "while" "(" expression ")" statement ;
block          → "{" declaration* "}" ;
expression     → assignment ;

assignment     → ( call "." )? IDENTIFIER "=" assignment
               | logic_or ;

logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → equality ( "and" equality )* ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;

unary          → ( "!" | "-" ) unary | call ;
call           → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
primary        → "true" | "false" | "nil" | "this"
               | NUMBER | STRING | IDENTIFIER | "(" expression ")"
               | "super" "." IDENTIFIER ;
function       → IDENTIFIER "(" parameters? ")" block ;
parameters     → IDENTIFIER ( "," IDENTIFIER )* ;
arguments      → expression ( "," expression )* ;
```

### Lexical Grammar

```
NUMBER         → DIGIT+ ( "." DIGIT+ )? ;
STRING         → "\"" <any char except "\"">* "\"" ;
IDENTIFIER     → ALPHA ( ALPHA | DIGIT )* ;
ALPHA          → "a" ... "z" | "A" ... "Z" | "_" ;
DIGIT          → "0" ... "9" ;
```

## Caveats

* Assumes ASCII source text

## Miscellaneous

* [Go Style](https://google.github.io/styleguide/go/)
* [Aglet Mono Typeface](https://fonts.adobe.com/fonts/aglet-mono)
