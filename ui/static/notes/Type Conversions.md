## Convert an int ascii to its character value


Use a [string conversion](https://golang.org/ref/spec#Conversions_to_and_from_a_string_type) to convert an ASCII numeric value to a string containing the ASCII character:

```go
fmt.Println(string(49)) // prints 1
```

The `go vet` command warns about the `int` to `string` conversion in the this code snippet because the conversion is commonly thought to create a decimal representation of the number. To squelch the warning, use a `rune` instead of an `int`:

```go
fmt.Println(string(rune(49))) // prints 1
```

This works for any rune value, not just the ASCII subset of runes.

Another option is to create a slice of bytes with the ASCII value and convert the slice to a string.

```go
b := []byte{49}
fmt.Println(string(b))  // prints 1
```

A variation on the previous snippet that works on all runes is:

```go
b := []rune{49}
fmt.Println(string(b))  // prints 1
```
