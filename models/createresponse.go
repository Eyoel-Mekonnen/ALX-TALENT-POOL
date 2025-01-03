package models
import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

 type CreateResponse struct {
    ID                  primitive.ObjectID `json:"id"`                   // This will be returned from the database
    Name                string             `json:"name"`                 // Student's name
    JobSeekerLocation   string             `json:"jobseekerlocation"`    // Job seeker location
    Telegram            string             `json:"telegram"`             // Student's Telegram username
    Github              string             `json:"github"`               // Student's GitHub link
    Linkedin            string             `json:"linkedin"`             // Student's LinkedIn link
    Skills              string             `json:"skills"`               // List of skills
    Languages           string             `json:"languages"`            // List of languages
    ProfessionalSummary string             `json:"professionalsummary"`  // Professional summary
    PortfolioLink       string             `json:"portfoliolink"`        // Portfolio link
    Education           string             `json:"education"`            // Education details
}

