package database

type Redis interface {
	Save() error
	Check() bool
}
