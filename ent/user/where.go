// Code generated by ent, DO NOT EDIT.

package user

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/hedgehog125/project-reboot/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.User {
	return predicate.User(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.User {
	return predicate.User(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.User {
	return predicate.User(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.User {
	return predicate.User(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.User {
	return predicate.User(sql.FieldLTE(FieldID, id))
}

// Username applies equality check predicate on the "username" field. It's identical to UsernameEQ.
func Username(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldUsername, v))
}

// AlertDiscordId applies equality check predicate on the "alertDiscordId" field. It's identical to AlertDiscordIdEQ.
func AlertDiscordId(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldAlertDiscordId, v))
}

// AlertEmail applies equality check predicate on the "alertEmail" field. It's identical to AlertEmailEQ.
func AlertEmail(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldAlertEmail, v))
}

// Content applies equality check predicate on the "content" field. It's identical to ContentEQ.
func Content(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldContent, v))
}

// FileName applies equality check predicate on the "fileName" field. It's identical to FileNameEQ.
func FileName(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldFileName, v))
}

// Mime applies equality check predicate on the "mime" field. It's identical to MimeEQ.
func Mime(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldMime, v))
}

// Nonce applies equality check predicate on the "nonce" field. It's identical to NonceEQ.
func Nonce(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldNonce, v))
}

// KeySalt applies equality check predicate on the "keySalt" field. It's identical to KeySaltEQ.
func KeySalt(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldKeySalt, v))
}

// PasswordHash applies equality check predicate on the "passwordHash" field. It's identical to PasswordHashEQ.
func PasswordHash(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldPasswordHash, v))
}

// PasswordSalt applies equality check predicate on the "passwordSalt" field. It's identical to PasswordSaltEQ.
func PasswordSalt(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldPasswordSalt, v))
}

// HashTime applies equality check predicate on the "hashTime" field. It's identical to HashTimeEQ.
func HashTime(v uint32) predicate.User {
	return predicate.User(sql.FieldEQ(FieldHashTime, v))
}

// HashMemory applies equality check predicate on the "hashMemory" field. It's identical to HashMemoryEQ.
func HashMemory(v uint32) predicate.User {
	return predicate.User(sql.FieldEQ(FieldHashMemory, v))
}

// HashKeyLen applies equality check predicate on the "hashKeyLen" field. It's identical to HashKeyLenEQ.
func HashKeyLen(v uint32) predicate.User {
	return predicate.User(sql.FieldEQ(FieldHashKeyLen, v))
}

// UsernameEQ applies the EQ predicate on the "username" field.
func UsernameEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldUsername, v))
}

// UsernameNEQ applies the NEQ predicate on the "username" field.
func UsernameNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldUsername, v))
}

// UsernameIn applies the In predicate on the "username" field.
func UsernameIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldUsername, vs...))
}

// UsernameNotIn applies the NotIn predicate on the "username" field.
func UsernameNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldUsername, vs...))
}

// UsernameGT applies the GT predicate on the "username" field.
func UsernameGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldUsername, v))
}

// UsernameGTE applies the GTE predicate on the "username" field.
func UsernameGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldUsername, v))
}

// UsernameLT applies the LT predicate on the "username" field.
func UsernameLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldUsername, v))
}

// UsernameLTE applies the LTE predicate on the "username" field.
func UsernameLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldUsername, v))
}

// UsernameContains applies the Contains predicate on the "username" field.
func UsernameContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldUsername, v))
}

// UsernameHasPrefix applies the HasPrefix predicate on the "username" field.
func UsernameHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldUsername, v))
}

// UsernameHasSuffix applies the HasSuffix predicate on the "username" field.
func UsernameHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldUsername, v))
}

// UsernameEqualFold applies the EqualFold predicate on the "username" field.
func UsernameEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldUsername, v))
}

// UsernameContainsFold applies the ContainsFold predicate on the "username" field.
func UsernameContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldUsername, v))
}

// AlertDiscordIdEQ applies the EQ predicate on the "alertDiscordId" field.
func AlertDiscordIdEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldAlertDiscordId, v))
}

// AlertDiscordIdNEQ applies the NEQ predicate on the "alertDiscordId" field.
func AlertDiscordIdNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldAlertDiscordId, v))
}

// AlertDiscordIdIn applies the In predicate on the "alertDiscordId" field.
func AlertDiscordIdIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldAlertDiscordId, vs...))
}

// AlertDiscordIdNotIn applies the NotIn predicate on the "alertDiscordId" field.
func AlertDiscordIdNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldAlertDiscordId, vs...))
}

// AlertDiscordIdGT applies the GT predicate on the "alertDiscordId" field.
func AlertDiscordIdGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldAlertDiscordId, v))
}

// AlertDiscordIdGTE applies the GTE predicate on the "alertDiscordId" field.
func AlertDiscordIdGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldAlertDiscordId, v))
}

// AlertDiscordIdLT applies the LT predicate on the "alertDiscordId" field.
func AlertDiscordIdLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldAlertDiscordId, v))
}

// AlertDiscordIdLTE applies the LTE predicate on the "alertDiscordId" field.
func AlertDiscordIdLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldAlertDiscordId, v))
}

// AlertDiscordIdContains applies the Contains predicate on the "alertDiscordId" field.
func AlertDiscordIdContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldAlertDiscordId, v))
}

// AlertDiscordIdHasPrefix applies the HasPrefix predicate on the "alertDiscordId" field.
func AlertDiscordIdHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldAlertDiscordId, v))
}

// AlertDiscordIdHasSuffix applies the HasSuffix predicate on the "alertDiscordId" field.
func AlertDiscordIdHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldAlertDiscordId, v))
}

// AlertDiscordIdEqualFold applies the EqualFold predicate on the "alertDiscordId" field.
func AlertDiscordIdEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldAlertDiscordId, v))
}

// AlertDiscordIdContainsFold applies the ContainsFold predicate on the "alertDiscordId" field.
func AlertDiscordIdContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldAlertDiscordId, v))
}

// AlertEmailEQ applies the EQ predicate on the "alertEmail" field.
func AlertEmailEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldAlertEmail, v))
}

// AlertEmailNEQ applies the NEQ predicate on the "alertEmail" field.
func AlertEmailNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldAlertEmail, v))
}

// AlertEmailIn applies the In predicate on the "alertEmail" field.
func AlertEmailIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldAlertEmail, vs...))
}

// AlertEmailNotIn applies the NotIn predicate on the "alertEmail" field.
func AlertEmailNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldAlertEmail, vs...))
}

// AlertEmailGT applies the GT predicate on the "alertEmail" field.
func AlertEmailGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldAlertEmail, v))
}

// AlertEmailGTE applies the GTE predicate on the "alertEmail" field.
func AlertEmailGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldAlertEmail, v))
}

// AlertEmailLT applies the LT predicate on the "alertEmail" field.
func AlertEmailLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldAlertEmail, v))
}

// AlertEmailLTE applies the LTE predicate on the "alertEmail" field.
func AlertEmailLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldAlertEmail, v))
}

// AlertEmailContains applies the Contains predicate on the "alertEmail" field.
func AlertEmailContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldAlertEmail, v))
}

// AlertEmailHasPrefix applies the HasPrefix predicate on the "alertEmail" field.
func AlertEmailHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldAlertEmail, v))
}

// AlertEmailHasSuffix applies the HasSuffix predicate on the "alertEmail" field.
func AlertEmailHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldAlertEmail, v))
}

// AlertEmailEqualFold applies the EqualFold predicate on the "alertEmail" field.
func AlertEmailEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldAlertEmail, v))
}

// AlertEmailContainsFold applies the ContainsFold predicate on the "alertEmail" field.
func AlertEmailContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldAlertEmail, v))
}

// ContentEQ applies the EQ predicate on the "content" field.
func ContentEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldContent, v))
}

// ContentNEQ applies the NEQ predicate on the "content" field.
func ContentNEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldContent, v))
}

// ContentIn applies the In predicate on the "content" field.
func ContentIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldIn(FieldContent, vs...))
}

// ContentNotIn applies the NotIn predicate on the "content" field.
func ContentNotIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldContent, vs...))
}

// ContentGT applies the GT predicate on the "content" field.
func ContentGT(v []byte) predicate.User {
	return predicate.User(sql.FieldGT(FieldContent, v))
}

// ContentGTE applies the GTE predicate on the "content" field.
func ContentGTE(v []byte) predicate.User {
	return predicate.User(sql.FieldGTE(FieldContent, v))
}

// ContentLT applies the LT predicate on the "content" field.
func ContentLT(v []byte) predicate.User {
	return predicate.User(sql.FieldLT(FieldContent, v))
}

// ContentLTE applies the LTE predicate on the "content" field.
func ContentLTE(v []byte) predicate.User {
	return predicate.User(sql.FieldLTE(FieldContent, v))
}

// FileNameEQ applies the EQ predicate on the "fileName" field.
func FileNameEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldFileName, v))
}

// FileNameNEQ applies the NEQ predicate on the "fileName" field.
func FileNameNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldFileName, v))
}

// FileNameIn applies the In predicate on the "fileName" field.
func FileNameIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldFileName, vs...))
}

// FileNameNotIn applies the NotIn predicate on the "fileName" field.
func FileNameNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldFileName, vs...))
}

// FileNameGT applies the GT predicate on the "fileName" field.
func FileNameGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldFileName, v))
}

// FileNameGTE applies the GTE predicate on the "fileName" field.
func FileNameGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldFileName, v))
}

// FileNameLT applies the LT predicate on the "fileName" field.
func FileNameLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldFileName, v))
}

// FileNameLTE applies the LTE predicate on the "fileName" field.
func FileNameLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldFileName, v))
}

// FileNameContains applies the Contains predicate on the "fileName" field.
func FileNameContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldFileName, v))
}

// FileNameHasPrefix applies the HasPrefix predicate on the "fileName" field.
func FileNameHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldFileName, v))
}

// FileNameHasSuffix applies the HasSuffix predicate on the "fileName" field.
func FileNameHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldFileName, v))
}

// FileNameEqualFold applies the EqualFold predicate on the "fileName" field.
func FileNameEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldFileName, v))
}

// FileNameContainsFold applies the ContainsFold predicate on the "fileName" field.
func FileNameContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldFileName, v))
}

// MimeEQ applies the EQ predicate on the "mime" field.
func MimeEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldMime, v))
}

// MimeNEQ applies the NEQ predicate on the "mime" field.
func MimeNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldMime, v))
}

// MimeIn applies the In predicate on the "mime" field.
func MimeIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldMime, vs...))
}

// MimeNotIn applies the NotIn predicate on the "mime" field.
func MimeNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldMime, vs...))
}

// MimeGT applies the GT predicate on the "mime" field.
func MimeGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldMime, v))
}

// MimeGTE applies the GTE predicate on the "mime" field.
func MimeGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldMime, v))
}

// MimeLT applies the LT predicate on the "mime" field.
func MimeLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldMime, v))
}

// MimeLTE applies the LTE predicate on the "mime" field.
func MimeLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldMime, v))
}

// MimeContains applies the Contains predicate on the "mime" field.
func MimeContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldMime, v))
}

// MimeHasPrefix applies the HasPrefix predicate on the "mime" field.
func MimeHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldMime, v))
}

// MimeHasSuffix applies the HasSuffix predicate on the "mime" field.
func MimeHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldMime, v))
}

// MimeEqualFold applies the EqualFold predicate on the "mime" field.
func MimeEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldMime, v))
}

// MimeContainsFold applies the ContainsFold predicate on the "mime" field.
func MimeContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldMime, v))
}

// NonceEQ applies the EQ predicate on the "nonce" field.
func NonceEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldNonce, v))
}

// NonceNEQ applies the NEQ predicate on the "nonce" field.
func NonceNEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldNonce, v))
}

// NonceIn applies the In predicate on the "nonce" field.
func NonceIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldIn(FieldNonce, vs...))
}

// NonceNotIn applies the NotIn predicate on the "nonce" field.
func NonceNotIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldNonce, vs...))
}

// NonceGT applies the GT predicate on the "nonce" field.
func NonceGT(v []byte) predicate.User {
	return predicate.User(sql.FieldGT(FieldNonce, v))
}

// NonceGTE applies the GTE predicate on the "nonce" field.
func NonceGTE(v []byte) predicate.User {
	return predicate.User(sql.FieldGTE(FieldNonce, v))
}

// NonceLT applies the LT predicate on the "nonce" field.
func NonceLT(v []byte) predicate.User {
	return predicate.User(sql.FieldLT(FieldNonce, v))
}

// NonceLTE applies the LTE predicate on the "nonce" field.
func NonceLTE(v []byte) predicate.User {
	return predicate.User(sql.FieldLTE(FieldNonce, v))
}

// KeySaltEQ applies the EQ predicate on the "keySalt" field.
func KeySaltEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldKeySalt, v))
}

// KeySaltNEQ applies the NEQ predicate on the "keySalt" field.
func KeySaltNEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldKeySalt, v))
}

// KeySaltIn applies the In predicate on the "keySalt" field.
func KeySaltIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldIn(FieldKeySalt, vs...))
}

// KeySaltNotIn applies the NotIn predicate on the "keySalt" field.
func KeySaltNotIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldKeySalt, vs...))
}

// KeySaltGT applies the GT predicate on the "keySalt" field.
func KeySaltGT(v []byte) predicate.User {
	return predicate.User(sql.FieldGT(FieldKeySalt, v))
}

// KeySaltGTE applies the GTE predicate on the "keySalt" field.
func KeySaltGTE(v []byte) predicate.User {
	return predicate.User(sql.FieldGTE(FieldKeySalt, v))
}

// KeySaltLT applies the LT predicate on the "keySalt" field.
func KeySaltLT(v []byte) predicate.User {
	return predicate.User(sql.FieldLT(FieldKeySalt, v))
}

// KeySaltLTE applies the LTE predicate on the "keySalt" field.
func KeySaltLTE(v []byte) predicate.User {
	return predicate.User(sql.FieldLTE(FieldKeySalt, v))
}

// PasswordHashEQ applies the EQ predicate on the "passwordHash" field.
func PasswordHashEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldPasswordHash, v))
}

// PasswordHashNEQ applies the NEQ predicate on the "passwordHash" field.
func PasswordHashNEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldPasswordHash, v))
}

// PasswordHashIn applies the In predicate on the "passwordHash" field.
func PasswordHashIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldIn(FieldPasswordHash, vs...))
}

// PasswordHashNotIn applies the NotIn predicate on the "passwordHash" field.
func PasswordHashNotIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldPasswordHash, vs...))
}

// PasswordHashGT applies the GT predicate on the "passwordHash" field.
func PasswordHashGT(v []byte) predicate.User {
	return predicate.User(sql.FieldGT(FieldPasswordHash, v))
}

// PasswordHashGTE applies the GTE predicate on the "passwordHash" field.
func PasswordHashGTE(v []byte) predicate.User {
	return predicate.User(sql.FieldGTE(FieldPasswordHash, v))
}

// PasswordHashLT applies the LT predicate on the "passwordHash" field.
func PasswordHashLT(v []byte) predicate.User {
	return predicate.User(sql.FieldLT(FieldPasswordHash, v))
}

// PasswordHashLTE applies the LTE predicate on the "passwordHash" field.
func PasswordHashLTE(v []byte) predicate.User {
	return predicate.User(sql.FieldLTE(FieldPasswordHash, v))
}

// PasswordSaltEQ applies the EQ predicate on the "passwordSalt" field.
func PasswordSaltEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldEQ(FieldPasswordSalt, v))
}

// PasswordSaltNEQ applies the NEQ predicate on the "passwordSalt" field.
func PasswordSaltNEQ(v []byte) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldPasswordSalt, v))
}

// PasswordSaltIn applies the In predicate on the "passwordSalt" field.
func PasswordSaltIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldIn(FieldPasswordSalt, vs...))
}

// PasswordSaltNotIn applies the NotIn predicate on the "passwordSalt" field.
func PasswordSaltNotIn(vs ...[]byte) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldPasswordSalt, vs...))
}

// PasswordSaltGT applies the GT predicate on the "passwordSalt" field.
func PasswordSaltGT(v []byte) predicate.User {
	return predicate.User(sql.FieldGT(FieldPasswordSalt, v))
}

// PasswordSaltGTE applies the GTE predicate on the "passwordSalt" field.
func PasswordSaltGTE(v []byte) predicate.User {
	return predicate.User(sql.FieldGTE(FieldPasswordSalt, v))
}

// PasswordSaltLT applies the LT predicate on the "passwordSalt" field.
func PasswordSaltLT(v []byte) predicate.User {
	return predicate.User(sql.FieldLT(FieldPasswordSalt, v))
}

// PasswordSaltLTE applies the LTE predicate on the "passwordSalt" field.
func PasswordSaltLTE(v []byte) predicate.User {
	return predicate.User(sql.FieldLTE(FieldPasswordSalt, v))
}

// HashTimeEQ applies the EQ predicate on the "hashTime" field.
func HashTimeEQ(v uint32) predicate.User {
	return predicate.User(sql.FieldEQ(FieldHashTime, v))
}

// HashTimeNEQ applies the NEQ predicate on the "hashTime" field.
func HashTimeNEQ(v uint32) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldHashTime, v))
}

// HashTimeIn applies the In predicate on the "hashTime" field.
func HashTimeIn(vs ...uint32) predicate.User {
	return predicate.User(sql.FieldIn(FieldHashTime, vs...))
}

// HashTimeNotIn applies the NotIn predicate on the "hashTime" field.
func HashTimeNotIn(vs ...uint32) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldHashTime, vs...))
}

// HashTimeGT applies the GT predicate on the "hashTime" field.
func HashTimeGT(v uint32) predicate.User {
	return predicate.User(sql.FieldGT(FieldHashTime, v))
}

// HashTimeGTE applies the GTE predicate on the "hashTime" field.
func HashTimeGTE(v uint32) predicate.User {
	return predicate.User(sql.FieldGTE(FieldHashTime, v))
}

// HashTimeLT applies the LT predicate on the "hashTime" field.
func HashTimeLT(v uint32) predicate.User {
	return predicate.User(sql.FieldLT(FieldHashTime, v))
}

// HashTimeLTE applies the LTE predicate on the "hashTime" field.
func HashTimeLTE(v uint32) predicate.User {
	return predicate.User(sql.FieldLTE(FieldHashTime, v))
}

// HashMemoryEQ applies the EQ predicate on the "hashMemory" field.
func HashMemoryEQ(v uint32) predicate.User {
	return predicate.User(sql.FieldEQ(FieldHashMemory, v))
}

// HashMemoryNEQ applies the NEQ predicate on the "hashMemory" field.
func HashMemoryNEQ(v uint32) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldHashMemory, v))
}

// HashMemoryIn applies the In predicate on the "hashMemory" field.
func HashMemoryIn(vs ...uint32) predicate.User {
	return predicate.User(sql.FieldIn(FieldHashMemory, vs...))
}

// HashMemoryNotIn applies the NotIn predicate on the "hashMemory" field.
func HashMemoryNotIn(vs ...uint32) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldHashMemory, vs...))
}

// HashMemoryGT applies the GT predicate on the "hashMemory" field.
func HashMemoryGT(v uint32) predicate.User {
	return predicate.User(sql.FieldGT(FieldHashMemory, v))
}

// HashMemoryGTE applies the GTE predicate on the "hashMemory" field.
func HashMemoryGTE(v uint32) predicate.User {
	return predicate.User(sql.FieldGTE(FieldHashMemory, v))
}

// HashMemoryLT applies the LT predicate on the "hashMemory" field.
func HashMemoryLT(v uint32) predicate.User {
	return predicate.User(sql.FieldLT(FieldHashMemory, v))
}

// HashMemoryLTE applies the LTE predicate on the "hashMemory" field.
func HashMemoryLTE(v uint32) predicate.User {
	return predicate.User(sql.FieldLTE(FieldHashMemory, v))
}

// HashKeyLenEQ applies the EQ predicate on the "hashKeyLen" field.
func HashKeyLenEQ(v uint32) predicate.User {
	return predicate.User(sql.FieldEQ(FieldHashKeyLen, v))
}

// HashKeyLenNEQ applies the NEQ predicate on the "hashKeyLen" field.
func HashKeyLenNEQ(v uint32) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldHashKeyLen, v))
}

// HashKeyLenIn applies the In predicate on the "hashKeyLen" field.
func HashKeyLenIn(vs ...uint32) predicate.User {
	return predicate.User(sql.FieldIn(FieldHashKeyLen, vs...))
}

// HashKeyLenNotIn applies the NotIn predicate on the "hashKeyLen" field.
func HashKeyLenNotIn(vs ...uint32) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldHashKeyLen, vs...))
}

// HashKeyLenGT applies the GT predicate on the "hashKeyLen" field.
func HashKeyLenGT(v uint32) predicate.User {
	return predicate.User(sql.FieldGT(FieldHashKeyLen, v))
}

// HashKeyLenGTE applies the GTE predicate on the "hashKeyLen" field.
func HashKeyLenGTE(v uint32) predicate.User {
	return predicate.User(sql.FieldGTE(FieldHashKeyLen, v))
}

// HashKeyLenLT applies the LT predicate on the "hashKeyLen" field.
func HashKeyLenLT(v uint32) predicate.User {
	return predicate.User(sql.FieldLT(FieldHashKeyLen, v))
}

// HashKeyLenLTE applies the LTE predicate on the "hashKeyLen" field.
func HashKeyLenLTE(v uint32) predicate.User {
	return predicate.User(sql.FieldLTE(FieldHashKeyLen, v))
}

// HasLoginAttempts applies the HasEdge predicate on the "loginAttempts" edge.
func HasLoginAttempts() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, LoginAttemptsTable, LoginAttemptsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasLoginAttemptsWith applies the HasEdge predicate on the "loginAttempts" edge with a given conditions (other predicates).
func HasLoginAttemptsWith(preds ...predicate.LoginAttempt) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := newLoginAttemptsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.User) predicate.User {
	return predicate.User(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.User) predicate.User {
	return predicate.User(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.User) predicate.User {
	return predicate.User(sql.NotPredicates(p))
}
