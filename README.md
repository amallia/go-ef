# go-ef [![Build Status](https://travis-ci.org/amallia/go-ef.svg?branch=master)](https://travis-ci.org/amallia/go-ef)
A Go implementation of the Elias-Fano encoding

### Example
```go
package main
import (
    "fmt"
    "github.com/amallia/go-ef"
    "os"
)

func main() {
    array := []uint64{1,5,10}
    size := len(array)
    max := array[size-1]
    obj := ef.New(max, size)

    obj.Compress(array)

    v, err := obj.Next()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println(v) // 1

    obj.Next()
    fmt.Println(obj.Value()) // 5
}
```