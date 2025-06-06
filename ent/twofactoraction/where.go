// Code generated by ent, DO NOT EDIT.

package twofactoraction

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/hedgehog125/project-reboot/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldLTE(FieldID, id))
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEQ(FieldType, v))
}

// Version applies equality check predicate on the "version" field. It's identical to VersionEQ.
func Version(v int) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEQ(FieldVersion, v))
}

// ExpiresAt applies equality check predicate on the "expiresAt" field. It's identical to ExpiresAtEQ.
func ExpiresAt(v time.Time) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEQ(FieldExpiresAt, v))
}

// Code applies equality check predicate on the "code" field. It's identical to CodeEQ.
func Code(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEQ(FieldCode, v))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldNotIn(FieldType, vs...))
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldGT(FieldType, v))
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldGTE(FieldType, v))
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldLT(FieldType, v))
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldLTE(FieldType, v))
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldContains(FieldType, v))
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldHasPrefix(FieldType, v))
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldHasSuffix(FieldType, v))
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEqualFold(FieldType, v))
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldContainsFold(FieldType, v))
}

// VersionEQ applies the EQ predicate on the "version" field.
func VersionEQ(v int) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEQ(FieldVersion, v))
}

// VersionNEQ applies the NEQ predicate on the "version" field.
func VersionNEQ(v int) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldNEQ(FieldVersion, v))
}

// VersionIn applies the In predicate on the "version" field.
func VersionIn(vs ...int) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldIn(FieldVersion, vs...))
}

// VersionNotIn applies the NotIn predicate on the "version" field.
func VersionNotIn(vs ...int) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldNotIn(FieldVersion, vs...))
}

// VersionGT applies the GT predicate on the "version" field.
func VersionGT(v int) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldGT(FieldVersion, v))
}

// VersionGTE applies the GTE predicate on the "version" field.
func VersionGTE(v int) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldGTE(FieldVersion, v))
}

// VersionLT applies the LT predicate on the "version" field.
func VersionLT(v int) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldLT(FieldVersion, v))
}

// VersionLTE applies the LTE predicate on the "version" field.
func VersionLTE(v int) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldLTE(FieldVersion, v))
}

// ExpiresAtEQ applies the EQ predicate on the "expiresAt" field.
func ExpiresAtEQ(v time.Time) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEQ(FieldExpiresAt, v))
}

// ExpiresAtNEQ applies the NEQ predicate on the "expiresAt" field.
func ExpiresAtNEQ(v time.Time) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldNEQ(FieldExpiresAt, v))
}

// ExpiresAtIn applies the In predicate on the "expiresAt" field.
func ExpiresAtIn(vs ...time.Time) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldIn(FieldExpiresAt, vs...))
}

// ExpiresAtNotIn applies the NotIn predicate on the "expiresAt" field.
func ExpiresAtNotIn(vs ...time.Time) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldNotIn(FieldExpiresAt, vs...))
}

// ExpiresAtGT applies the GT predicate on the "expiresAt" field.
func ExpiresAtGT(v time.Time) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldGT(FieldExpiresAt, v))
}

// ExpiresAtGTE applies the GTE predicate on the "expiresAt" field.
func ExpiresAtGTE(v time.Time) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldGTE(FieldExpiresAt, v))
}

// ExpiresAtLT applies the LT predicate on the "expiresAt" field.
func ExpiresAtLT(v time.Time) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldLT(FieldExpiresAt, v))
}

// ExpiresAtLTE applies the LTE predicate on the "expiresAt" field.
func ExpiresAtLTE(v time.Time) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldLTE(FieldExpiresAt, v))
}

// CodeEQ applies the EQ predicate on the "code" field.
func CodeEQ(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEQ(FieldCode, v))
}

// CodeNEQ applies the NEQ predicate on the "code" field.
func CodeNEQ(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldNEQ(FieldCode, v))
}

// CodeIn applies the In predicate on the "code" field.
func CodeIn(vs ...string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldIn(FieldCode, vs...))
}

// CodeNotIn applies the NotIn predicate on the "code" field.
func CodeNotIn(vs ...string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldNotIn(FieldCode, vs...))
}

// CodeGT applies the GT predicate on the "code" field.
func CodeGT(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldGT(FieldCode, v))
}

// CodeGTE applies the GTE predicate on the "code" field.
func CodeGTE(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldGTE(FieldCode, v))
}

// CodeLT applies the LT predicate on the "code" field.
func CodeLT(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldLT(FieldCode, v))
}

// CodeLTE applies the LTE predicate on the "code" field.
func CodeLTE(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldLTE(FieldCode, v))
}

// CodeContains applies the Contains predicate on the "code" field.
func CodeContains(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldContains(FieldCode, v))
}

// CodeHasPrefix applies the HasPrefix predicate on the "code" field.
func CodeHasPrefix(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldHasPrefix(FieldCode, v))
}

// CodeHasSuffix applies the HasSuffix predicate on the "code" field.
func CodeHasSuffix(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldHasSuffix(FieldCode, v))
}

// CodeEqualFold applies the EqualFold predicate on the "code" field.
func CodeEqualFold(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldEqualFold(FieldCode, v))
}

// CodeContainsFold applies the ContainsFold predicate on the "code" field.
func CodeContainsFold(v string) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.FieldContainsFold(FieldCode, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.TwoFactorAction) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.TwoFactorAction) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.TwoFactorAction) predicate.TwoFactorAction {
	return predicate.TwoFactorAction(sql.NotPredicates(p))
}
