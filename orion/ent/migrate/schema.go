// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// TemplatesColumns holds the columns for the "templates" table.
	TemplatesColumns = []*schema.Column{
		{Name: "oid", Type: field.TypeUUID},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "type", Type: field.TypeString},
		{Name: "template", Type: field.TypeString, SchemaType: map[string]string{"postgres": "text"}},
		{Name: "created_by", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_by", Type: field.TypeString, Nullable: true},
		{Name: "updated_at", Type: field.TypeTime, Nullable: true},
		{Name: "deleted_by", Type: field.TypeString, Nullable: true},
		{Name: "deleted_at", Type: field.TypeTime, Nullable: true},
	}
	// TemplatesTable holds the schema information for the "templates" table.
	TemplatesTable = &schema.Table{
		Name:       "templates",
		Columns:    TemplatesColumns,
		PrimaryKey: []*schema.Column{TemplatesColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		TemplatesTable,
	}
)

func init() {
}
