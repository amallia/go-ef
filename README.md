# go-ef
A Go implementation of the Elias-Fano encoding

### Example
```go
package main
import (
    "fmt"
    "github.com/amallia/go-ef"
)

func main() {
    array := []uint64{1,5,10}
    size := len(array)
    max := array[size-1]
    obj := ef.New(max, size)

    obj.Compress(array)

    v := obj.Next()
    fmt.Println(v) // 1

    obj.Next()
    fmt.Println(obj.Value()) // 5
}
```