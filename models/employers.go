package models
import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Employer struct{
	ID          primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
    EmployerID  primitive.ObjectID  `json:"employerid,omitempty" bson:"employerid,omitempty"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"createdat" bson:"createdat"`
}
