To easily read and parse CSV files in Go, you can use two methods of [`encoding/csv` package:

- [`csv.Reader.ReadAll()`](https://pkg.go.dev/encoding/csv#Reader.ReadAll) to read and parse the entire file at once. Note however that a very large file may not fit into the memory.
- [`csv.Reader.Read()`](https://pkg.go.dev/encoding/csv#Reader.Read) to read the CSV file line by line.

> See also our example of [how to write data to a CSV file in Go](https://gosamples.dev/write-csv)

> In the examples below, we use `data.csv` file:
> 
> ```bash
> vegetables,fruits
> carrot,banana
> potato,strawberry
> ```

## Read the entire CSV file at once

In this example, we open the CSV file, initialize [`csv.Reader`](https://pkg.go.dev/encoding/csv#NewReader) and read all the data into a `[][]string` slice where the first index is the file’s line number and the second is an index of the comma-separated value in this line. We can do something with these data, for example, convert to an array of structs.

```go
package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
)

type ShoppingRecord struct {
    Vegetable string
    Fruit     string
}

func createShoppingList(data [][]string) []ShoppingRecord {
    var shoppingList []ShoppingRecord
    for i, line := range data {
        if i > 0 { // omit header line
            var rec ShoppingRecord
            for j, field := range line {
                if j == 0 {
                    rec.Vegetable = field
                } else if j == 1 {
                    rec.Fruit = field
                }
            }
            shoppingList = append(shoppingList, rec)
        }
    }
    return shoppingList
}

func main() {
    // open file
    f, err := os.Open("data.csv")
    if err != nil {
        log.Fatal(err)
    }

    // remember to close the file at the end of the program
    defer f.Close()

    // read csv values using csv.Reader
    csvReader := csv.NewReader(f)
    data, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal(err)
    }

    // convert records to array of structs
    shoppingList := createShoppingList(data)

    // print the array
    fmt.Printf("%+v\n", shoppingList)
}
```

Output:

```bash
[{Vegetable:carrot Fruit:banana} {Vegetable:potato Fruit:strawberry}]
```

## Read a CSV file line by line

Reading line by line is similar to [reading the whole file at once](https://gosamples.dev/read-csv/#read-the-entire-csv-file-at-once), but in this case, we use [`csv.Reader.Read()`](https://pkg.go.dev/encoding/csv#Reader.Read) method to read the next line of data in the infinite loop. This loop is exited when no more data are available, i.e., [`io.EOF`](https://pkg.go.dev/io#EOF) error occurs.

```go
package main

import (
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "os"
)

func main() {
    // open file
    f, err := os.Open("data.csv")
    if err != nil {
        log.Fatal(err)
    }

    // remember to close the file at the end of the program
    defer f.Close()

    // read csv values using csv.Reader
    csvReader := csv.NewReader(f)
    for {
        rec, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }
        // do something with read line
        fmt.Printf("%+v\n", rec)
    }
}
```

Output:

```bash
[vegetables fruits]
[carrot banana]
[potato strawberry]
```
