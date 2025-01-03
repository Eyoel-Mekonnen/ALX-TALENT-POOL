package models
import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)
type Application struct {
    ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    JobID           primitive.ObjectID `json:"jobid,omitempty" bson:"jobid,omitempty"`
    JobSeekerID     primitive.ObjectID `json:"jobseekerid,omitempty" bson:"jobseekerid,omitempty"`
    FilePath        string             `json:"filepath" bson:"filepath"`
    CreatedAt       time.Time          `json:"createdat" bson:"createdat"`
}
