package gogination

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

type Person struct {
	Id   string
	Name string
	Age  int
}

var peter = Person{
	Id:   "abc",
	Name: "Peter Piper",
	Age:  25,
}

func TestBuilder_NextFilter(t *testing.T) {
	type args struct {
		obj     interface{}
		filter  bson.D
		orderBy bson.D
	}
	tests := []struct {
		name    string
		args    args
		want    bson.D
		wantErr bool
	}{
		{
			name: "invalid object",
			args: args{
				obj:     "I am a string",
				filter:  nil,
				orderBy: nil,
			},
			wantErr: true,
		},
		{
			name: "missing id",
			args: args{
				obj: struct {
					Name string
				}{
					Name: "Peter Piper",
				},
				filter:  nil,
				orderBy: nil,
			},
			wantErr: true,
		},
		{
			name: "simple",
			args: args{
				obj:     peter,
				filter:  nil,
				orderBy: nil,
			},
			want: bson.D{
				bson.E{
					Key: "_id",
					Value: bson.E{
						Key: "$oid",
						Value: bson.E{
							Key:   "$gt",
							Value: "abc",
						},
					},
				},
			},
		},
		{
			name: "filtered",
			args: args{
				obj: peter,
				filter: bson.D{
					bson.E{
						Key:   "name",
						Value: 25,
					},
				},
				orderBy: nil,
			},
			want: bson.D{
				bson.E{
					Key: "$and",
					Value: bson.A{
						bson.D{
							bson.E{
								Key:   "name",
								Value: 25,
							},
						},
						bson.D{
							bson.E{
								Key: "_id",
								Value: bson.E{
									Key: "$oid",
									Value: bson.E{
										Key:   "$gt",
										Value: "abc",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "invalid order",
			args: args{
				obj:    peter,
				filter: nil,
				orderBy: bson.D{
					bson.E{
						Key:   "Age",
						Value: 3,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "order by id desc",
			args: args{
				obj:    peter,
				filter: nil,
				orderBy: bson.D{
					bson.E{
						Key:   "_id",
						Value: -1,
					},
				},
			},
			want: bson.D{
				bson.E{
					Key: "_id",
					Value: bson.E{
						Key: "$oid",
						Value: bson.E{
							Key:   "$lt",
							Value: "abc",
						},
					},
				},
			},
		},
		{
			name: "ordered by age",
			args: args{
				obj:    peter,
				filter: nil,
				orderBy: bson.D{
					bson.E{
						Key:   "age",
						Value: 1,
					},
				},
			},
			wantErr: false,
			want: bson.D{
				bson.E{
					Key: "$or",
					Value: bson.A{
						bson.D{
							bson.E{ // greater match for age
								Key: "age",
								Value: bson.E{
									Key:   "$gt",
									Value: 25,
								},
							},
						},
						bson.D{
							bson.E{ // exact match for age...
								Key:   "age",
								Value: 25,
							},
							bson.E{ // ... but greater id
								Key: "_id",
								Value: bson.E{
									Key: "$oid",
									Value: bson.E{
										Key:   "$gt",
										Value: "abc",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &Builder{
				structIdField: "Id",
				mongoIdField:  "_id",
			}
			got, err := builder.NextFilter(tt.args.obj, tt.args.filter, tt.args.orderBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextFilter() got = %+v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBuilder(t *testing.T) {
	type args struct {
		opts []func(*Builder)
	}
	tests := []struct {
		name    string
		args    args
		want    *Builder
		wantErr bool
	}{
		{
			name: "new",
			args: args{},
			want: &Builder{
				structIdField: "Id",
				mongoIdField:  "_id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBuilder(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBuilder() got = %v, want %v", got, tt.want)
			}
		})
	}
}
