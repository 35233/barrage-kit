# STT serialize

Douyu's STT is a method of serializing structured data

1. Keys and values are directly split using '@='
1. Array uses '/' to separate elements 
1. '/' in key or value is encoded as '@S'
1. '@' in key or value is encoded as '@A'

Example：
 
    (1) key-value：key1@=value1/key2@=value2/key3@=value3/ 
    (2) Array：value1/value2/value3/ 
    
Import
----
```go
import "github.com/35233/barrage-kit/stt"
```
 
 Demo
 ----
 ```go
package main

import (
	"fmt"
	"github.com/35233/barrage-kit/stt"
)

func main() {
	fmt.Println(stt.Decode("a@A=b@Sc@A=d@S/d/v/"))
	fmt.Println(stt.Decode("a@=b/c@=d/"))
    fmt.Println(stt.Encode(map[string]interface{}{
		"type": "joingroup",
		"rid":  "123456",
		"gid":  "-9999",
	}))
}
```