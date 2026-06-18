package domain

import "time"

type User struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Hash      []byte    `bson:"hash"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	LastLogin time.Time `bson:"last_login"`
}
