# Douyu Client

Read Douyu's barrage data

Import
----
```go
import "github.com/35233/barrage-kit/douyuclient"
```

Demo
----
```go
package main
import (
	"fmt"
	"github.com/35233/barrage-kit/douyuclient"
	"time"
)

func main() {
    fmt.Println("started")
    client := douyuclient.New("openbarrage.douyutv.com:8601", 50)
    client.AddRoom("196")
    client.AddRoom("52004")
    msgChannel, err := client.Start()
    if err != nil {
        fmt.Println(err)
        return
    }
    go func() {
        time.Sleep(20 * time.Second)
        client.AddRoom("252140")
        time.Sleep(20 * time.Second)
        client.Stop()
    }()
    for msg := range msgChannel {
        fmt.Println("msg", msg.Text())
    }
    
    fmt.Println("restart after 5s")
    
    time.Sleep(5 * time.Second)
    msgChannel, err = client.Start()
    if err != nil {
        fmt.Println(err)
        return
    }
    go func() {
        time.Sleep(20 * time.Second)
        client.Stop()
    }()
    for msg := range msgChannel {
        fmt.Println("msg2", msg.Text())
    }
    
    fmt.Println("exit after 5s")
    
    time.Sleep(5 * time.Second)
    fmt.Println("end main")
}

```