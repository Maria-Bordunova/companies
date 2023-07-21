package entity

type EventType string

const (
	CompanyCreated EventType = "created"
	CompanyUpdated EventType = "updated"
	CompanyDeleted EventType = "deleted"
)
