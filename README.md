go-amazon-paapi
===============

Amazon PAAPI client for golang

## Example

``` go
package main

import (
	"net/url"

	"github.com/m0t0k1ch1/go-amazon-paapi"
)

func main() {
	c := paapi.NewClient(
		"your access key",
		"your secret access key",
		"your associate tag",
	)

	var result []byte
	var err error

	result, err = c.ItemLookup("B003GQSYJO")
	if err != nil {
		panic(err)
	}
	println(string(result))

	result, err = c.ItemSearchByKeyword(
		"All",
		url.QueryEscape("clockwork orange"),
		"Small,Images",
	)
	if err != nil {
		panic(err)
	}
	println(string(result))
}
```
