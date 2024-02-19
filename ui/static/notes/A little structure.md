Structures group related values together, making it simpler and less error-prone to pass them around.

```go
type location struct {
	lat float64
	long float64
}
var spirit location
spirit.lat = -14.5684
spirit.long = 175.472636

var opportunity location
opportunity.lat = -1.9462
opportunity.long = 354.4734

fmt.Println(spirit, opportunity)
```

Output:

```
{-14.5684 175.472636}
{-1.9462 354.4734}
```

## Initializing structures

```go
type location struct {
	lat, long float64
}

ariane := location{lat: -1.9462, long: 354.4734}
fmt.Println(ariane)

atlas := location{4.5, 135.9}
fmt.Println(atlas)
```

Output:

```
{-1.9462 354.4734}
{4.5 135.9}
```

## Printing keys of structures

```go
curiosity := location{-4.5895, 137.4417}

fmt.Printf("%v\n", curiosity)
fmt.Printf("%+v\n", curiosity)
```

Output:

```
{-4.5895 137.4417}
{lat:-4.5895 long:137.4417}
```

## Encoding structures to JSON

The Marshal function from the json package is used to encode the data into JSON format. Marshal returns the JSON data as bytes, which can be sent over the wire or converted to a string for display. It may also return an error.


```go
package main
import (
	"encoding/json"
	"fmt"
	"os"
)
func main() {
	type location struct {
		Lat, Long float64
	}
	
	curiosity := location{-4.5895, 137.4417}
	bytes, err := json.Marshal(curiosity)
	exitOnError(err)
	
	fmt.Println(string(bytes))

	// exitOnError prints any errors and exits.
	func exitOnError(err error) {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
```

> Notice that the JSON keys match the field names of the location structure. For this to work, the json package requires fields to be exported. If Lat and Long began with a lower-case letter, the output would be {}.

## Customizing JSON with struct tags

Go’s json package requires that fields have an initial uppercase letter and multi-word field names use Camel-Case by convention. You may want JSON keys in snake_case, particularly when inter-operating with Python or Ruby. The fields of a structure can be tagged with the field names you want the json package to use.

```go
type location struct {
	Lat float64 `json:"latitude"`
	Long float64 `json:"longitude"`
}
```

To read more on Go's json standard library visit [[Go by Example JSON]]
