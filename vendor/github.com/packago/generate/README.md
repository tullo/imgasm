Go package inspired by icza's great [SO answer](https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang). I made slight modification to include numbers and omit characters that can be difficult to differentiate (such as lI1).

Get:
```
go get -u github.com/packago/generate
```

Use:
```
package main

import (
  "github.com/packago/generate"
  "fmt"
)

func main() {
  fmt.Println(generate.RandomLettersNumbers(89))
  fmt.Println(generate.RandomLowercaseNumbers(6))
  fmt.Println(generate.RandomLowercaseNumbers(4))
}
```
