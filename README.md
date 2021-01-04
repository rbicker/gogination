gogination
==========

mongodb pagination, written in golang

# usage
```go
package example

import (
	"context"
	"time"

	"github.com/rbicker/gogination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name  string `bson:"name"`
	Age   int    `bson:"age"`
	State string `bson:"state"`
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// connect to mongodb server
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	// set db & collection
	db := client.Database("example")
	col := db.Collection("people")

	// filter, sorting & limit
	filter := bson.D{
		bson.E{
			Key:   "state",
			Value: "Washington",
		},
	}
	sort := bson.D{
		bson.E{
			Key:   "age",
			Value: -1,
		},
	}
	opts := options.Find()
	opts.SetSort(sort)
	opts.SetLimit(10)

	// query
	cur, err := col.Find(ctx, filter, opts)
	if err != nil {
		panic(err)
	}
	var people []Person
	err = cur.All(ctx, people)
	if err != nil {
		panic(err)
	}

	// determine filter for next page
	builder, err := NewBuilder()
	if err != nil {
		panic(err)
	}
	last := people[len(people)-1]
	nextFilter, err := builder.NextFilter(last, filter, sort)
	if err != nil {
		panic(err)
	}

	// todo: encode filter and add to api response
	// ...
}

```