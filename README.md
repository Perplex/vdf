# vdf

A Parser for [Valves Data Format (known as vdf)](https://developer.valvesoftware.com/wiki/KeyValues) written in Go. 
Comments are not preserved during parsing.

## Installation

It is go gettable

```
$ go get github.com/perplex/vdf
```
   

## Usage

The parser will return a KeyValue type that can handle querying for submaps, objects, and values without the need to worry about type assertions.

```go
firstkey
{
	secondkey
	{
		"attr1" "val1"
		"attr2" "val2"
	}
}
```


In the above example firstkey would be a submap as it is of the type map[string]interface{}, while secondkey can either be a submap of a object. An object is defined as a map[string]string, so it is the final map structure before looping over values. Val1 and val2 are both values as there is no map structure below them. Querying for secondkey as a submap would look like this

```go
package main

import (
	"fmt"
	"github.com/perplex/vdf"
)

func main() {

	obj, err := vdf.ParseFile("path/to/example.vdf")
	if err != nil {
		panic(err)
	}
	
	sm, err := obj.GetSubmap("firstkey", "secondkey")
	if err != nil {
		panic(err)
	}
	
	fmt.Println(sm.GetKeys())
}

```
This would print attr1 and attr2.

## Inspiration

This is based on the original fork of this repo from [andygrunwald](https://github.com/andygrunwald/vdf), and 
[simple-vdf](https://github.com/rossengeorgiev/vdf-parser) which could handle duplicate keys and various other nuances 
in vdf files.
