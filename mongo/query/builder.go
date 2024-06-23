package query

import (
	"github.com/snowmerak/template/mongo/query/option"
	"go.mongodb.org/mongo-driver/bson"
)

type Builder struct {
	field      string
	conditions bson.D
}

func New(field string) *Builder {
	return &Builder{
		field: field,
	}
}

func (b *Builder) Build() bson.D {
	return bson.D{{Key: b.field, Value: b.conditions}}
}

func (b *Builder) Gt(value interface{}) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$gt", Value: value})
	return b
}

func (b *Builder) Gte(value interface{}) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$gte", Value: value})
	return b
}

func (b *Builder) Lt(value interface{}) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$lt", Value: value})
	return b
}

func (b *Builder) Lte(value interface{}) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$lte", Value: value})
	return b
}

func (b *Builder) Eq(value interface{}) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$eq", Value: value})
	return b
}

func (b *Builder) Ne(value interface{}) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$ne", Value: value})
	return b
}

func (b *Builder) In(value ...interface{}) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$in", Value: value})
	return b
}

func (b *Builder) Nin(value ...interface{}) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$nin", Value: value})
	return b
}

func And(conditions ...bson.D) bson.D {
	return bson.D{{Key: "$and", Value: conditions}}
}

func Or(conditions ...bson.D) bson.D {
	return bson.D{{Key: "$or", Value: conditions}}
}

func Not(conditions bson.D) bson.D {
	return bson.D{{Key: "$not", Value: conditions}}
}

func Nor(conditions ...bson.D) bson.D {
	return bson.D{{Key: "$nor", Value: conditions}}
}

func (b *Builder) Exists(value bool) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$exists", Value: value})
	return b
}

const (
	TypeDouble              = 1
	TypeString              = 2
	TypeObject              = 3
	TypeArray               = 4
	TypeBinaryData          = 5
	TypeUndefined           = 6
	TypeObjectId            = 7
	TypeBoolean             = 8
	TypeDate                = 9
	TypeNull                = 10
	TypeRegex               = 11
	TypeDBPointer           = 12
	TypeJavaScript          = 13
	TypeSymbol              = 14
	TypeJavaScriptWithScope = 15
	TypeInt32               = 16
	TypeTimestamp           = 17
	TypeInt64               = 18
	TypeDecimal128          = 19
	TypeMinKey              = -1
	TypeMaxKey              = 127
)

func (b *Builder) Type(value int) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$type", Value: value})
	return b
}

func TextSearch(value string, language option.Option[string], caseSensitive option.Option[bool], diacriticSensitive option.Option[bool]) bson.D {
	query := bson.D{{Key: "$search", Value: value}}
	if language.IsSome() {
		query = append(query, bson.E{Key: "$language", Value: language.Unwrap()})
	}
	if caseSensitive.IsSome() {
		query = append(query, bson.E{Key: "$caseSensitive", Value: caseSensitive.Unwrap()})
	}
	if diacriticSensitive.IsSome() {
		query = append(query, bson.E{Key: "$diacriticSensitive", Value: diacriticSensitive.Unwrap()})
	}

	return bson.D{{Key: "$text", Value: query}}
}

const (
	GeometryPoint        = "Point"
	GeometryMultiPoint   = "MultiPoint"
	GeometryLineString   = "LineString"
	GeometryMultiLine    = "MultiLineString"
	GeometryPolygon      = "Polygon"
	GeometryMultiPolygon = "MultiPolygon"
	GeometryGeometryColl = "GeometryCollection"
)

func Geometry(geoType string, coordinates ...[2]float64) bson.E {
	c := make([][]float64, len(coordinates))
	for i := range coordinates {
		c[i] = coordinates[i][:]
	}
	return bson.E{
		Key: "$geometry",
		Value: bson.D{
			{Key: "type", Value: geoType},
			{Key: "coordinates", Value: c},
		},
	}
}

func GeometryCollection(geometries ...bson.E) bson.E {
	return bson.E{Key: "$geometry", Value: bson.D{{Key: "type", Value: GeometryGeometryColl}, {Key: "geometries", Value: geometries}}}
}

func Box(bottomLeft, topRight [2]float64) bson.E {
	return bson.E{Key: "$box", Value: [][]float64{bottomLeft[:], topRight[:]}}
}

func Center(center [2]float64, radius float64) bson.E {
	return bson.E{Key: "$center", Value: []interface{}{center[:], radius}}
}

func (b *Builder) GeoIntersects(value bson.D) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$geoIntersects", Value: value})
	return b
}

func (b *Builder) GeoWithin(value bson.D) *Builder {
	b.conditions = append(b.conditions, bson.E{Key: "$geoWithin", Value: value})
	return b
}

func (b *Builder) Near(value bson.E, minDistance option.Option[float64], maxDistance option.Option[float64]) *Builder {
	query := bson.D{value}
	if minDistance.IsSome() {
		query = append(query, bson.E{Key: "$minDistance", Value: minDistance.Unwrap()})
	}
	if maxDistance.IsSome() {
		query = append(query, bson.E{Key: "$maxDistance", Value: maxDistance.Unwrap()})
	}

	b.conditions = append(b.conditions, bson.E{Key: "$near", Value: query})
	return b
}

func (b *Builder) NearSphere(value bson.E, minDistance option.Option[float64], maxDistance option.Option[float64]) *Builder {
	query := bson.D{value}
	if minDistance.IsSome() {
		query = append(query, bson.E{Key: "$minDistance", Value: minDistance.Unwrap()})
	}
	if maxDistance.IsSome() {
		query = append(query, bson.E{Key: "$maxDistance", Value: maxDistance.Unwrap()})
	}

	b.conditions = append(b.conditions, bson.E{Key: "$nearSphere", Value: query})
	return b
}
