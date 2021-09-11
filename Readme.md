# JSON to XML

A simple JSON to XML converter written in Go.

Usage:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pgundlach/json2xml"
)

func dothings() error {
	f, err := os.Open("myfile.json")
	if err != nil {
		return err
	}
	str, err := json2xml.ToXML(f)
	if err != nil {
		return err
	}
	fmt.Println(str)
	return nil
}

func main() {
	err := dothings()
	if err != nil {
		log.Fatal(err)
	}
}
```

The (formatted) result is:

```xml
<data>
    <map>
        <array key="whatever">
            <entry>foo</entry>
            <entry>3.45</entry>
            <entry>bar</entry>
            <entry>1</entry>
        </array>
        <map key="something">
            <entry key="another">object</entry>
            <array key="and an">
                <entry>array</entry>
            </array>
        </map>
    </map>
</data>
```
