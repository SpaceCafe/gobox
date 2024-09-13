package jsonapi

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spacecafe/gobox/gin-rest/types"
	"gorm.io/gorm/schema"
)

// Response represents the entire JSON:API response.
type Response[T any] struct {
	Meta     Meta    `json:"meta"`
	JSONAPI  JSONAPI `json:"jsonapi"`
	Data     *[]Data `json:"data"`
	Included *[]Data `json:"included,omitempty"`
	Links    Links   `json:"links"`
}

// JSONAPI contains details to the JSON:API specification.
type JSONAPI struct {

	// Version indicates the version of the JSON:API spec that the document complies with.
	Version string `json:"version"`
}

// Meta represents the meta information about the primary data in JSON:API.
type Meta struct {

	// Page is the current page number.
	Page int `json:"page"`

	// PageSize is the number of items per page.
	PageSize int `json:"page_size"`

	// Total is the total number of items.
	Total int `json:"total"`

	// TotalPages is the total number of pages.
	TotalPages int `json:"total_pages"`
}

// Data represents a resource object in JSON:API.
// It includes the type, ID, attributes, and relationships of the resource.
type Data struct {

	// Type specifies the model or resource name.
	Type string `json:"type"`

	// ID is the unique identifier of the resource object, usually a UUID.
	ID json.RawMessage `json:"id"`

	// Attributes contains the key-value pairs representing the resource's fields.
	Attributes map[string]json.RawMessage `json:"attributes,omitempty"`

	// Relationships maps relationship names to their corresponding relationship data.
	Relationships map[string]Relationship `json:"relationships,omitempty"`
}

// Relationship represents the relationships of a resource in JSON:API.
// It includes the data about the relationship, which provides information about the related resource.
type Relationship struct {

	// Data about the relationship.
	Data *[]RelationshipData `json:"data"`
}

// RelationshipData represents the data about a relationship in JSON:API.
// It includes the type and ID of the related resource, which are used to identify
// the resource that is related to the primary data in the JSON:API document.
type RelationshipData struct {

	// Type specifies the model or resource name.
	Type string `json:"type"`

	// ID is the unique identifier of the resource object, usually a UUID.
	ID json.RawMessage `json:"id"`
}

// Links represents the links related to a resource in JSON:API.
// It includes links for navigation and related resources.
type Links struct {

	// Self contains the link to the current resource.
	Self string `json:"self"`

	// First contains the link to the first page of resources.
	First string `json:"first,omitempty"`

	// Last contains the link to the last page of resources.
	Last string `json:"last,omitempty"`

	// Previous contains the link to the previous page of resources.
	Previous string `json:"prev,omitempty"`

	// Next contains the link to the next page of resources.
	Next string `json:"next,omitempty"`
}

// NewResponse creates and returns a new Response instance with default values.
func NewResponse[T any]() *Response[T] {
	return &Response[T]{
		Meta: Meta{
			Page:       1,
			PageSize:   1,
			Total:      1,
			TotalPages: 1,
		},
		JSONAPI: JSONAPI{
			Version: "1.1",
		},
		Data:     &[]Data{},
		Included: &[]Data{},
	}
}

// NewResponseFromEntity creates a new Response object from single entity and view options.
func NewResponseFromEntity[T any](resource types.Resource[T], entity *T) *Response[T] {
	response := NewResponse[T]()
	if entity != nil {
		response.Links = Links{Self: resource.BasePath().JoinPath(resource.PrimaryValue(entity)).String()}
		response.addData(context.Background(), resource, resource.Schema(), reflect.ValueOf(entity), response.Data)
	}
	return response
}

// NewResponseFromEntities creates a new Response object from a list of entities and view options.
func NewResponseFromEntities[T any](resource types.Resource[T], entities *[]T, options *types.ViewOptions) *Response[T] {
	response := NewResponse[T]()
	response.Meta = Meta{
		Page:       options.GetPage(),
		PageSize:   options.GetPageSize(),
		Total:      options.Total,
		TotalPages: options.GetTotalPages(),
	}
	response.Links = Links{Self: resource.BasePath().String()}
	response.setLinks(resource, options)
	if entities != nil {
		for i := range *entities {
			response.addData(context.Background(), resource, resource.Schema(), reflect.ValueOf(&(*entities)[i]), response.Data)
		}
	}
	return response
}

// Render writes the Response as JSON to the provided http.ResponseWriter, implements render.Render interface.
func (r *Response[T]) Render(writer http.ResponseWriter) error {
	r.WriteContentType(writer)
	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(r); err != nil {
		return err
	}
	return nil
}

// WriteContentType sets the Content-Type header to application/json with UTF-8 charset, implements render.Render interface.
func (r *Response[T]) WriteContentType(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/vnd.api+json; charset=utf-8")
}

// addData adds data to the response, including attributes and relationships, based on the provided schema and value.
func (r *Response[T]) addData(ctx context.Context, resource types.Resource[T], schema *schema.Schema, value reflect.Value, dest *[]Data) {
	fields := append(schema.Fields[:0:0], schema.Fields...)
	relationships := schema.Relationships.Relations

	data := &Data{
		Type:          schema.Table,
		ID:            r.serialize(ctx, schema.PrioritizedPrimaryField, value),
		Attributes:    make(map[string]json.RawMessage, len(fields)),
		Relationships: make(map[string]Relationship, len(relationships)),
	}

	// Process each relationship in the schema.
	for _, relationship := range relationships {
		fields = removeByRef(fields, relationship.Field)
		if !relationship.Field.Readable {
			continue
		}

		// Remove foreign keys from fields.
		for _, reference := range relationship.References {
			if reference.ForeignKey != nil {
				fields = removeByRef(fields, reference.ForeignKey)
			}
		}

		// Determine the field name for the relationship.
		fieldName := resource.NamingStrategy().ColumnName(schema.Table, relationship.Name)
		if jsonTag, ok := relationship.Field.Tag.Lookup("json"); ok {
			fieldName = strings.Split(jsonTag, ",")[0]
			if fieldName == "-" {
				continue
			}
		}

		// Get the value of the relationship field.
		fieldValue := relationship.Field.ReflectValueOf(ctx, value)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		// Handle different kinds of field values (struct, slice, array).
		switch fieldValue.Kind() {
		case reflect.Struct:
			if id := r.serialize(ctx, relationship.FieldSchema.PrioritizedPrimaryField, fieldValue); id != nil {
				data.Relationships[fieldName] = Relationship{Data: &[]RelationshipData{{
					Type: relationship.FieldSchema.Table,
					ID:   id,
				}}}
				r.addData(ctx, resource, relationship.FieldSchema, fieldValue, r.Included)
			}
		case reflect.Slice, reflect.Array:
			var relData []RelationshipData
			for i := 0; i < fieldValue.Len(); i++ {
				relData = append(relData, RelationshipData{
					Type: relationship.FieldSchema.Table,
					ID:   r.serialize(ctx, relationship.FieldSchema.PrioritizedPrimaryField, fieldValue.Index(i)),
				})
				r.addData(ctx, resource, relationship.FieldSchema, fieldValue.Index(i), r.Included)
			}
			data.Relationships[fieldName] = Relationship{Data: &relData}
		default:
		}
	}

	// Process each readable field in the schema.
	for _, field := range fields {
		if field.Readable && field.DBName != "" && !field.PrimaryKey {
			data.Attributes[field.DBName] = r.serialize(ctx, field, value)
		}
	}

	*dest = append(*dest, *data)
}

// removeByRef removes a specific field reference from a slice of field pointers.
// It replaces the removed element with the last element and returns the shortened slice.
func removeByRef(slice []*schema.Field, ref *schema.Field) []*schema.Field {
	for i, v := range slice {
		if v == ref {

			// Move the last element to the position of the element to remove
			slice[i] = slice[len(slice)-1]

			// Return the slice without the last element
			return slice[:len(slice)-1]
		}
	}
	return slice
}

// serialize transforms the given field value into a json.RawMessage.
// It handles special cases like time.Time and uses a generic JSON serializer for other types.
func (r *Response[T]) serialize(ctx context.Context, field *schema.Field, value reflect.Value) json.RawMessage {
	fieldValue, zero := field.ValueOf(ctx, value)
	if zero {
		return nil
	}

	switch v := fieldValue.(type) {
	case time.Time:
		out, _ := json.Marshal(v.Format(time.RFC3339))
		return out
	default:
		out, _ := schema.JSONSerializer{}.Value(ctx, field, value, fieldValue)
		if outMsg, ok := out.(string); ok {
			return []byte(outMsg)
		}
	}
	return nil
}

// setLinks sets pagination links.
func (r *Response[T]) setLinks(resource types.Resource[T], options *types.ViewOptions) {
	if options.Total <= 0 {
		return
	}

	link := *resource.BasePath()
	query := link.Query()
	query.Set("page[size]", strconv.Itoa(options.PageSize))

	// Current page
	query.Set("page[number]", strconv.Itoa(options.Page))
	link.RawQuery = query.Encode()
	r.Links.Self = link.String()

	// First page
	query.Set("page[number]", "1")
	link.RawQuery = query.Encode()
	r.Links.First = link.String()

	// Last page
	query.Set("page[number]", strconv.Itoa(r.Meta.TotalPages))
	link.RawQuery = query.Encode()
	r.Links.Last = link.String()

	// Previous page
	if options.Page > 1 {
		query.Set("page[number]", strconv.Itoa(options.Page-1))
		link.RawQuery = query.Encode()
		r.Links.Previous = link.String()

	}

	// Next page
	if options.Page < r.Meta.TotalPages {
		query.Set("page[number]", strconv.Itoa(options.Page+1))
		link.RawQuery = query.Encode()
		r.Links.Next = link.String()
	}
}
