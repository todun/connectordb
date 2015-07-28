package operator

import (
	"connectordb/streamdb/datastream"
	"connectordb/streamdb/operator/adminoperator"
	"connectordb/streamdb/operator/authoperator"
	"connectordb/streamdb/operator/interfaces"
	"connectordb/streamdb/operator/messenger"
	"connectordb/streamdb/operator/plainoperator"
	"connectordb/streamdb/users"
)

type Database interface {
	GetUserDatabase() users.UserDatabase
	GetDatastream() *datastream.DataStream
	GetMessenger() *messenger.Messenger
}

/** Gets a general-purpose administrative operator without database-ruining
permissions.
**/
func NewOperator(db Database) Operator {
	ao := newAdminOperator(db)
	po := interfaces.PathOperatorMixin{&ao}
	return &po
}

func newAdminOperator(db Database) adminoperator.AdminOperator {
	op := newPlainOperator(db)
	ao := adminoperator.AdminOperator{&op}
	return ao
}

func newPlainOperator(db Database) plainoperator.PlainOperator {
	return plainoperator.NewPlainOperator(db.GetUserDatabase(), db.GetDatastream(), db.GetMessenger())
}

/*
NewUserOperator creates an operator that can act on behalf of the given user; returns an error
if the username does not exist.
**/
func NewUserOperator(db Database, username string) (Operator, error) {
	bootstrapOperator := NewOperator(db)
	op, err := authoperator.NewUserAuthOperator(bootstrapOperator, username)
	po := interfaces.PathOperatorMixin{op}
	return &po, err
}

/*
NewDeviceOperator creates an operator that keeps permissions contained to the
device at the given path.
*/
func NewDeviceOperator(db Database, devicepath string) (Operator, error) {
	bootstrapOperator := NewOperator(db)
	op, err := authoperator.NewDeviceAuthOperator(bootstrapOperator, devicepath)
	po := interfaces.PathOperatorMixin{op}
	return &po, err
}

/*
NewDeviceApiOperator creates an operator that keeps permissions contained to the
device at the given path. Additionally, it fails to be created if the given
apikey does not match the one for the specified device.
*/
func NewDeviceApiOperator(db Database, devicepath, apikey string) (Operator, error) {
	bootstrapOperator := NewOperator(db)
	op, err := authoperator.NewDeviceLoginOperator(bootstrapOperator, devicepath, apikey)
	po := interfaces.PathOperatorMixin{op}
	return &po, err

}

/*
NewDeviceIdOperator creates an operator that contains what it can do to the
scope of the device with the given id.
*/
func NewDeviceIdOperator(db Database, deviceID int64) (Operator, error) {
	bootstrapOperator := NewOperator(db)
	op, err := authoperator.NewDeviceIdOperator(bootstrapOperator, deviceID)
	po := interfaces.PathOperatorMixin{op}
	return &po, err
}

/*
NewUserLoginOperator creates an operator that contains what it can do to the
scope of the user with the given username and password.
*/
func NewUserLoginOperator(db Database, username, password string) (Operator, error) {
	bootstrapOperator := NewOperator(db)
	op, err := authoperator.NewUserLoginOperator(bootstrapOperator, username, password)
	po := interfaces.PathOperatorMixin{op}
	return &po, err

}
