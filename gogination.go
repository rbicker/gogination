package gogination

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// Builder is a small helper for pagination related to mongodb.
type Builder struct {
	structIdField string
	mongoIdField  string
}

// NewBuilder is a factory method for creating a new builder.
func NewBuilder(opts ...func(*Builder)) (*Builder, error) {
	b := &Builder{
		structIdField: "Id",
		mongoIdField:  "_id",
	}
	return b, nil
}

// NextFilter determines the filter to query the next document after the one given as obj.
// It takes the given filter and orderBy bson documents into consideration.
// The resulting filter will get returned as a bson document.
func (builder *Builder) NextFilter(obj interface{}, filter bson.D, sort bson.D) (bson.D, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Struct {
		return nil, ErrExpectedStruct
	}
	id := v.FieldByName(builder.structIdField)
	if !id.IsValid() {
		return nil, ErrNoIdField
	}
	idFilter := bson.E{
		Key: builder.mongoIdField,
		Value: bson.E{
			Key:   "$gt",
			Value: id.Interface(),
		},
	}
	// build pagination filter which makes sure that the
	// next page starts with the document coming after the given one
	// (depending on how the documents are / will be sorted)
	var pageFilter bson.D
	var next, exact bson.D
	// if sorting is requested
	if sort != nil {
		for _, e := range sort {
			name := e.Key
			var op string
			// determine if ascending / descending
			if v, ok := e.Value.(int); ok {
				switch v {
				case 1:
					op = "$gt"
				case -1:
					op = "$lt"
				default:
					return nil, ErrInvalidOrderBy
				}
			} else {
				return nil, ErrInvalidOrderBy
			}
			// handle id field which is special
			if name == builder.mongoIdField {
				// like id filter but considering
				// the operator (which might be $lt)
				flt := bson.E{
					Key: builder.mongoIdField,
					Value: bson.E{
						Key:   op,
						Value: id.Interface(),
					},
				}
				// if the only sort string is id
				if len(sort) == 1 {
					pageFilter = bson.D{flt}
					// no need to do more
					// as id is the only field
					break
				}
				// for id, we know that it is unique
				// therefore there cannot be exact
				// matches so we only set next and continue
				next = append(next, flt)
				continue
			}
			// get field value by it's name (camel case)
			f := v.FieldByName(strings.Title(name))
			if !f.IsValid() {
				// if field does not exist or has zero value
				// it cannot be used as a filter
				continue
			}
			value := f.Interface()
			// A) either the next document needs to have a
			// greater or smaller value for the given field
			// (depending on sort direction)
			next = append(next, bson.E{
				Key: name,
				Value: bson.E{
					Key:   op,
					Value: value,
				},
			})
			// B) or it needs to be an exact match
			// but the id needs to be greater, a criteria
			// we add after the loop
			exact = append(exact, bson.E{
				Key:   name,
				Value: value,
			})
		}
		// if only id was given as sort string
		// exact has a length of 0
		if len(exact) > 0 {
			// ensure the id is greater for exact matches
			exact = append(exact, idFilter)
			// put together page filter
			pageFilter = bson.D{
				{
					Key: "$or",
					Value: bson.A{
						next,
						exact,
					},
				},
			}
		}
	}
	// if no (valid) sorting is given, it will be sorted by
	// id, which makes the pagination filter quite simple
	if len(pageFilter) == 0 {
		pageFilter = bson.D{idFilter}
	}
	if len(filter) > 0 {
		// merge given and pagination filters
		filter = bson.D{
			{
				Key: "$and",
				Value: bson.A{
					filter,
					pageFilter,
				},
			},
		}
	} else {
		// no filter given as input
		// resulting filter equals the pagination filter
		filter = pageFilter
	}
	return filter, nil
}

// WithMongoIdField sets a custom mongodb id field (default: _id).
func WithMongoIdField(fieldName string) func(*Builder) error {
	return func(builder *Builder) error {
		builder.mongoIdField = fieldName
		return nil
	}
}

// WithStructIdField sets a custom id field (default: Id).
func WithStructIdField(fieldName string) func(*Builder) error {
	return func(builder *Builder) error {
		builder.structIdField = fieldName
		return nil
	}
}
