# Parser Cababilities

## implemented

- classes
- methoddeclaration
- constructordeclaration
- adding a standardconstructor if there is none
- accesmodifier
- standard accessmodifier
- fields
- fields with expression
- lokale variablen declaration
- Variable Usage (Identifier => Typecheck needs to split between local or field)
- keep line and column number of tokens for error messages
- name resolution aka this.number or foo.bar()
- method calls
- differentiation between field, method or constructor
- constructor calls

## to be implemented

- for
- if
- if else
- while
- swicht case
- break, continue
- getter and setter members
- Post/Pre increment/decrement

=> Then TypeCheck , Bytecode etc.
This should be enough implementation to make a full compiler for a really basic c# language which can then be extended on
