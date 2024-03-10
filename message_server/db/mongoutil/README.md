# Usage

Copy in a separated package called `mongoutil` the content of the `.go` file, then use it as follows

``` go
import (
  "path/of/mongoutil" 
)

type example struct {
  Field1 int             `bson:"field_1"`
  Field2 *mongoutil.UUID `bson:"field_2"`
}
```

## Usage with Goa (v1)

This is the proper way to use the UUID type in Goa V1, a similar way can be used with other versions.
Remember to copy the `func_helpers.go` file in the `design` directory of your project.

``` go
package design

// ... 

var myType = Type("myType", func() {
    BSONMember("field_1", Integer)
    BSONMember("field_2", UUID, func() {
        Metadata("struct:field:type", "*mongoutil.UUID", "path/of/mongoutil")
        // if you use Goa v2, use Meta(...) instead of Metadata
    })
})
```