package models

type Session struct {
	Id    uint64  // Session Id
	User  User    // User struct
	Sheet []Sheet // Sheet struct
}
