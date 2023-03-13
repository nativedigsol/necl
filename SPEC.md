# NECL Syntax Specification

This is the specification of the syntax and semantics of NECL.

## Syntax Notation

This notation is intended for human consumption rather than machine consumption, with the following conventions:

- Double and single quotes (`"` and `'`) are used to mark literal character sequences, which may be either punctuation markers or keywords.
- The symbol `|` indicates that any one of its left and right operands may be present.
- Parentheses `(` and `)` are used to group items together to apply the `|` operator to them collectively.
- The symbol `=` is used to declare variables

## Lexical Elements

### Comments

Comments start with the // sequence and end with the next newline sequence. A line comment is considered equivalent to a newline sequence.

### Operators and Delimiters

The following character sequences represent operators, delimiters, and other special tokens:

```
+   {   ==  <   &&  
-   }   !=  >   ||  
*   [   =   <=  !
/   ]   :   >=  (
%   ${  ?   \   )
```

## Structural elements

The structural language consists of syntax representing the following constructs:

- Attributes, which assign a value to a specified name.
- Blocks, which create a child body.
- Body Content, which consists of a collection of attributes and blocks.

```
attribute = "value"
block {
    sub_block {
        body_content = "foo"
    }
}
```

Note: Blocks **MUST** have a name assigned to it

## Data Types

NECL supports the common data types:

- Number (assigned integers and floats): `number = 3.14` or `number = -10`
- String (a collection of characters): `string = "Hello World!`
- Multiline string (a collection of lines): 
```
multiline = "line1" \
              "line2" \
              "lineN" \
              "final line"
```
- Boolean (true of false values): `bool = false` or `bool = true`
- Array (collection of data) = `array = ["foo", "bar", 2023, false]`

## Expressions

### If

A "if" is a conditional construct to make an attribute based on a condition, applying it's value by using the `?` and `:` operators.

```
msg = "Hello World!"
has_hello = if contains(msg, "hello") ? true : false
has_hi = if contains(msg, "hi") ? true : false
// has_hello = true
// has_hi = false

// Note:
// This is only an example, in this case it would make more sense to simply do
// has_hello = contains(msg, "hello")
// has_hi = contains(msg, "hi")
```

### For

A "for loop" is a construct for constructing a collection by projecting the items from another collection.

Note that if using a "for loop" with an array, 2 variables need to be declared: the index and the value. `... for index, value in array ...`

```
months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
monthNumber = for i, month in months : i + 1
// monthNumber = [1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12]
```

### Operations

Operations apply a particular operator to either one or more expression terms.

#### Arithmetic operators
```
a + b   // sum 
a - b   // difference
a * b   // product
a / b   // quotient
```

Note 1: NECL does not support the remainder, exponentiation and floor division operators. These are offered via functions
Note 2: These operations can only be done to integers. If you try to do this operation with float values, it will return an error when parsing

#### Comparative operators

Note that these can only be applied to integers

```
a == b    // Equal
a != b    // Not equal
a < b     // less than
a <= b    // less than or equal to
a > b     // greater than
a >= b    // greater than or equal to
```

### Binary operators

We can easily build logic gates by combining operators

```
cond1 = true
cond2 = false

// Conjunction
and = cond1 == cond2 // false
or = if (cond1 == true || cond2 == true) ? true : false // true

// Denial
nand = !and // true
nor = !nor // false

// Exclusive
xor = if (cond1 == cond2) ? false : true // true
xnor = !xor // false
```

### Functions

The following functions come by default with the NECL interpreter:

#### Strings

- upper(str) // Uppercases a string
- lower(str) // Lowercases a string
- concat(str, val) // Adds a string to the end of another string
- contains(str, substr) // Checks if a string contains a substring
- len(str) // Checks the length of the string

#### Numeric

- power(number, power) // Perform an exponent arithmetic operation
- floor(quotient, dividend) // Performs a floor division
- remainder(quotient, dividend) // Gets the remainder of a division

#### Gate Logic

- and(cond1, cond2) // AND gate
- or(cond1, cond2) // OR gate
- nand(cond1, cond2) // NAND gate
- nor(cond1, cond2) // NOR gate
- xor(cond1, cond2) // XOR gate
- xnor(cond1, cond2) // XNOR gate
