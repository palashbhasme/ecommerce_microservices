package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

type MyObjectID string
type User struct {
	ID        MyObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string     `bson:"first_name" json:"first_name"`
	LastName  string     `bson:"last_name" json:"last_name"`
	Email     string     `bson:"email" json:"email"`
	DOB       time.Time  `bson:"dob" json:"dob"`
	Phone     string     `bson:"phone" json:"phone"`
	Addresses []Address  `bson:"addresses" json:"addresses"`
	Account   Account    `bson:"account" json:"account"`
	CreatedAt string     `bson:"created_at" json:"created_at"`
	UpdatedAt string     `bson:"updated_at" json:"updated_at"`
}

type Address struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	Line1     string `bson:"line1" json:"line1"`           // Street address
	Line2     string `bson:"line2,omitempty" json:"line2"` // Optional second line
	City      string `bson:"city" json:"city"`
	State     string `bson:"state" json:"state"`
	Country   string `bson:"country" json:"country"`
	ZipCode   string `bson:"zip_code" json:"zip_code"`
	IsDefault bool   `bson:"is_default" json:"is_default"` // Marks the default address
}

type Account struct {
	ID           string `bson:"_id,omitempty" json:"id"`
	Role         string `bson:"role" json:"role"`
	Username     string `bson:"username" json:"username"`
	PasswordHash string `bson:"password_hash" json:"password_hash"` // Store hashed passwords only
	IsActive     bool   `bson:"is_active" json:"is_active"`
}

func (id MyObjectID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	objectID, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return bson.TypeNull, nil, err
	}
	return bson.MarshalValue(objectID)
}
