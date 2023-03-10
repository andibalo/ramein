// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/andibalo/ramein/orion/ent/predicate"
	"github.com/andibalo/ramein/orion/ent/template"
	"github.com/google/uuid"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeTemplate = "Template"
)

// TemplateMutation represents an operation that mutates the Template nodes in the graph.
type TemplateMutation struct {
	config
	op            Op
	typ           string
	id            *uuid.UUID
	name          *string
	_type         *string
	template      *string
	created_by    *string
	created_at    *time.Time
	updated_by    *string
	updated_at    *time.Time
	deleted_by    *string
	deleted_at    *time.Time
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Template, error)
	predicates    []predicate.Template
}

var _ ent.Mutation = (*TemplateMutation)(nil)

// templateOption allows management of the mutation configuration using functional options.
type templateOption func(*TemplateMutation)

// newTemplateMutation creates new mutation for the Template entity.
func newTemplateMutation(c config, op Op, opts ...templateOption) *TemplateMutation {
	m := &TemplateMutation{
		config:        c,
		op:            op,
		typ:           TypeTemplate,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withTemplateID sets the ID field of the mutation.
func withTemplateID(id uuid.UUID) templateOption {
	return func(m *TemplateMutation) {
		var (
			err   error
			once  sync.Once
			value *Template
		)
		m.oldValue = func(ctx context.Context) (*Template, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Template.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withTemplate sets the old Template of the mutation.
func withTemplate(node *Template) templateOption {
	return func(m *TemplateMutation) {
		m.oldValue = func(context.Context) (*Template, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m TemplateMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m TemplateMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Template entities.
func (m *TemplateMutation) SetID(id uuid.UUID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *TemplateMutation) ID() (id uuid.UUID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *TemplateMutation) IDs(ctx context.Context) ([]uuid.UUID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uuid.UUID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Template.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *TemplateMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *TemplateMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Template entity.
// If the Template object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TemplateMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *TemplateMutation) ResetName() {
	m.name = nil
}

// SetType sets the "type" field.
func (m *TemplateMutation) SetType(s string) {
	m._type = &s
}

// GetType returns the value of the "type" field in the mutation.
func (m *TemplateMutation) GetType() (r string, exists bool) {
	v := m._type
	if v == nil {
		return
	}
	return *v, true
}

// OldType returns the old "type" field's value of the Template entity.
// If the Template object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TemplateMutation) OldType(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldType is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldType requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldType: %w", err)
	}
	return oldValue.Type, nil
}

// ResetType resets all changes to the "type" field.
func (m *TemplateMutation) ResetType() {
	m._type = nil
}

// SetTemplate sets the "template" field.
func (m *TemplateMutation) SetTemplate(s string) {
	m.template = &s
}

// Template returns the value of the "template" field in the mutation.
func (m *TemplateMutation) Template() (r string, exists bool) {
	v := m.template
	if v == nil {
		return
	}
	return *v, true
}

// OldTemplate returns the old "template" field's value of the Template entity.
// If the Template object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TemplateMutation) OldTemplate(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTemplate is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTemplate requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTemplate: %w", err)
	}
	return oldValue.Template, nil
}

// ResetTemplate resets all changes to the "template" field.
func (m *TemplateMutation) ResetTemplate() {
	m.template = nil
}

// SetCreatedBy sets the "created_by" field.
func (m *TemplateMutation) SetCreatedBy(s string) {
	m.created_by = &s
}

// CreatedBy returns the value of the "created_by" field in the mutation.
func (m *TemplateMutation) CreatedBy() (r string, exists bool) {
	v := m.created_by
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedBy returns the old "created_by" field's value of the Template entity.
// If the Template object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TemplateMutation) OldCreatedBy(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreatedBy is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreatedBy requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedBy: %w", err)
	}
	return oldValue.CreatedBy, nil
}

// ResetCreatedBy resets all changes to the "created_by" field.
func (m *TemplateMutation) ResetCreatedBy() {
	m.created_by = nil
}

// SetCreatedAt sets the "created_at" field.
func (m *TemplateMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *TemplateMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the Template entity.
// If the Template object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TemplateMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *TemplateMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedBy sets the "updated_by" field.
func (m *TemplateMutation) SetUpdatedBy(s string) {
	m.updated_by = &s
}

// UpdatedBy returns the value of the "updated_by" field in the mutation.
func (m *TemplateMutation) UpdatedBy() (r string, exists bool) {
	v := m.updated_by
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedBy returns the old "updated_by" field's value of the Template entity.
// If the Template object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TemplateMutation) OldUpdatedBy(ctx context.Context) (v *string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdatedBy is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdatedBy requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdatedBy: %w", err)
	}
	return oldValue.UpdatedBy, nil
}

// ClearUpdatedBy clears the value of the "updated_by" field.
func (m *TemplateMutation) ClearUpdatedBy() {
	m.updated_by = nil
	m.clearedFields[template.FieldUpdatedBy] = struct{}{}
}

// UpdatedByCleared returns if the "updated_by" field was cleared in this mutation.
func (m *TemplateMutation) UpdatedByCleared() bool {
	_, ok := m.clearedFields[template.FieldUpdatedBy]
	return ok
}

// ResetUpdatedBy resets all changes to the "updated_by" field.
func (m *TemplateMutation) ResetUpdatedBy() {
	m.updated_by = nil
	delete(m.clearedFields, template.FieldUpdatedBy)
}

// SetUpdatedAt sets the "updated_at" field.
func (m *TemplateMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *TemplateMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the Template entity.
// If the Template object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TemplateMutation) OldUpdatedAt(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUpdatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUpdatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdatedAt: %w", err)
	}
	return oldValue.UpdatedAt, nil
}

// ClearUpdatedAt clears the value of the "updated_at" field.
func (m *TemplateMutation) ClearUpdatedAt() {
	m.updated_at = nil
	m.clearedFields[template.FieldUpdatedAt] = struct{}{}
}

// UpdatedAtCleared returns if the "updated_at" field was cleared in this mutation.
func (m *TemplateMutation) UpdatedAtCleared() bool {
	_, ok := m.clearedFields[template.FieldUpdatedAt]
	return ok
}

// ResetUpdatedAt resets all changes to the "updated_at" field.
func (m *TemplateMutation) ResetUpdatedAt() {
	m.updated_at = nil
	delete(m.clearedFields, template.FieldUpdatedAt)
}

// SetDeletedBy sets the "deleted_by" field.
func (m *TemplateMutation) SetDeletedBy(s string) {
	m.deleted_by = &s
}

// DeletedBy returns the value of the "deleted_by" field in the mutation.
func (m *TemplateMutation) DeletedBy() (r string, exists bool) {
	v := m.deleted_by
	if v == nil {
		return
	}
	return *v, true
}

// OldDeletedBy returns the old "deleted_by" field's value of the Template entity.
// If the Template object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TemplateMutation) OldDeletedBy(ctx context.Context) (v *string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDeletedBy is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDeletedBy requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDeletedBy: %w", err)
	}
	return oldValue.DeletedBy, nil
}

// ClearDeletedBy clears the value of the "deleted_by" field.
func (m *TemplateMutation) ClearDeletedBy() {
	m.deleted_by = nil
	m.clearedFields[template.FieldDeletedBy] = struct{}{}
}

// DeletedByCleared returns if the "deleted_by" field was cleared in this mutation.
func (m *TemplateMutation) DeletedByCleared() bool {
	_, ok := m.clearedFields[template.FieldDeletedBy]
	return ok
}

// ResetDeletedBy resets all changes to the "deleted_by" field.
func (m *TemplateMutation) ResetDeletedBy() {
	m.deleted_by = nil
	delete(m.clearedFields, template.FieldDeletedBy)
}

// SetDeletedAt sets the "deleted_at" field.
func (m *TemplateMutation) SetDeletedAt(t time.Time) {
	m.deleted_at = &t
}

// DeletedAt returns the value of the "deleted_at" field in the mutation.
func (m *TemplateMutation) DeletedAt() (r time.Time, exists bool) {
	v := m.deleted_at
	if v == nil {
		return
	}
	return *v, true
}

// OldDeletedAt returns the old "deleted_at" field's value of the Template entity.
// If the Template object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *TemplateMutation) OldDeletedAt(ctx context.Context) (v *time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDeletedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDeletedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDeletedAt: %w", err)
	}
	return oldValue.DeletedAt, nil
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (m *TemplateMutation) ClearDeletedAt() {
	m.deleted_at = nil
	m.clearedFields[template.FieldDeletedAt] = struct{}{}
}

// DeletedAtCleared returns if the "deleted_at" field was cleared in this mutation.
func (m *TemplateMutation) DeletedAtCleared() bool {
	_, ok := m.clearedFields[template.FieldDeletedAt]
	return ok
}

// ResetDeletedAt resets all changes to the "deleted_at" field.
func (m *TemplateMutation) ResetDeletedAt() {
	m.deleted_at = nil
	delete(m.clearedFields, template.FieldDeletedAt)
}

// Where appends a list predicates to the TemplateMutation builder.
func (m *TemplateMutation) Where(ps ...predicate.Template) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the TemplateMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *TemplateMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Template, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *TemplateMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *TemplateMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Template).
func (m *TemplateMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *TemplateMutation) Fields() []string {
	fields := make([]string, 0, 9)
	if m.name != nil {
		fields = append(fields, template.FieldName)
	}
	if m._type != nil {
		fields = append(fields, template.FieldType)
	}
	if m.template != nil {
		fields = append(fields, template.FieldTemplate)
	}
	if m.created_by != nil {
		fields = append(fields, template.FieldCreatedBy)
	}
	if m.created_at != nil {
		fields = append(fields, template.FieldCreatedAt)
	}
	if m.updated_by != nil {
		fields = append(fields, template.FieldUpdatedBy)
	}
	if m.updated_at != nil {
		fields = append(fields, template.FieldUpdatedAt)
	}
	if m.deleted_by != nil {
		fields = append(fields, template.FieldDeletedBy)
	}
	if m.deleted_at != nil {
		fields = append(fields, template.FieldDeletedAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *TemplateMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case template.FieldName:
		return m.Name()
	case template.FieldType:
		return m.GetType()
	case template.FieldTemplate:
		return m.Template()
	case template.FieldCreatedBy:
		return m.CreatedBy()
	case template.FieldCreatedAt:
		return m.CreatedAt()
	case template.FieldUpdatedBy:
		return m.UpdatedBy()
	case template.FieldUpdatedAt:
		return m.UpdatedAt()
	case template.FieldDeletedBy:
		return m.DeletedBy()
	case template.FieldDeletedAt:
		return m.DeletedAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *TemplateMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case template.FieldName:
		return m.OldName(ctx)
	case template.FieldType:
		return m.OldType(ctx)
	case template.FieldTemplate:
		return m.OldTemplate(ctx)
	case template.FieldCreatedBy:
		return m.OldCreatedBy(ctx)
	case template.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case template.FieldUpdatedBy:
		return m.OldUpdatedBy(ctx)
	case template.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	case template.FieldDeletedBy:
		return m.OldDeletedBy(ctx)
	case template.FieldDeletedAt:
		return m.OldDeletedAt(ctx)
	}
	return nil, fmt.Errorf("unknown Template field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *TemplateMutation) SetField(name string, value ent.Value) error {
	switch name {
	case template.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case template.FieldType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
		return nil
	case template.FieldTemplate:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTemplate(v)
		return nil
	case template.FieldCreatedBy:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedBy(v)
		return nil
	case template.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case template.FieldUpdatedBy:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedBy(v)
		return nil
	case template.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case template.FieldDeletedBy:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDeletedBy(v)
		return nil
	case template.FieldDeletedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDeletedAt(v)
		return nil
	}
	return fmt.Errorf("unknown Template field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *TemplateMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *TemplateMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *TemplateMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Template numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *TemplateMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(template.FieldUpdatedBy) {
		fields = append(fields, template.FieldUpdatedBy)
	}
	if m.FieldCleared(template.FieldUpdatedAt) {
		fields = append(fields, template.FieldUpdatedAt)
	}
	if m.FieldCleared(template.FieldDeletedBy) {
		fields = append(fields, template.FieldDeletedBy)
	}
	if m.FieldCleared(template.FieldDeletedAt) {
		fields = append(fields, template.FieldDeletedAt)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *TemplateMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *TemplateMutation) ClearField(name string) error {
	switch name {
	case template.FieldUpdatedBy:
		m.ClearUpdatedBy()
		return nil
	case template.FieldUpdatedAt:
		m.ClearUpdatedAt()
		return nil
	case template.FieldDeletedBy:
		m.ClearDeletedBy()
		return nil
	case template.FieldDeletedAt:
		m.ClearDeletedAt()
		return nil
	}
	return fmt.Errorf("unknown Template nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *TemplateMutation) ResetField(name string) error {
	switch name {
	case template.FieldName:
		m.ResetName()
		return nil
	case template.FieldType:
		m.ResetType()
		return nil
	case template.FieldTemplate:
		m.ResetTemplate()
		return nil
	case template.FieldCreatedBy:
		m.ResetCreatedBy()
		return nil
	case template.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case template.FieldUpdatedBy:
		m.ResetUpdatedBy()
		return nil
	case template.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case template.FieldDeletedBy:
		m.ResetDeletedBy()
		return nil
	case template.FieldDeletedAt:
		m.ResetDeletedAt()
		return nil
	}
	return fmt.Errorf("unknown Template field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *TemplateMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *TemplateMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *TemplateMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *TemplateMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *TemplateMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *TemplateMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *TemplateMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Template unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *TemplateMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Template edge %s", name)
}
