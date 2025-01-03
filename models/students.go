package models

import (
     "go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
    ID	                primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
    Name 	        string `json:"name" bson:"name"`
    JobSeekerLocation   string `json:"jobseekerlocation" bson:"jobseekerlocation"`
    Telegram	        string `json:"telegram" bson:"telegram"`
    Github              string `json:"github" bson:"github"`
    Linkedin            string `json:"linkedin" bson:"linkedin"`
    Skills              string `json:"skills" bson:"skills"`
    Languages           string `json:"languages" bson:"languages"`
    ProfessionalSummary string `json:"professionalsummary" bson:"professionalsummary"`
    PortfolioLink       string `json:"porfoliolink" bson:"portfoliolink"`
    Education           string `json:"education" bson:"education"`
}
