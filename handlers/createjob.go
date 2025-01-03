package handlers

import (
    "net/http"
    "encoding/json"
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/models"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type CreateJobSuccessResponse struct {
    Message      string             `json:"message"`
    ID           primitive.ObjectID `json:"ID"`
    Title        string             `json:"title"`
    EmployerID   primitive.ObjectID `json:"EmployerID"`
    Description  string             `json:"description"`
    CreatedAt    time.Time          `json:"createdat"`
}
type ErrorResponse struct {
    Error   string  `json:"Error"`
    Message string  `json:"message"`
}
// CreateJob godoc
// @Summary Create a new employer job
// @Description Create a new job entry for an employer in the database
// @Tags Employer
// @Accept json
// @Produce json
// @Param body body models.Employer true "Employer Job Details"
// @Success 200 {object} CreateJobSuccessResponse
// @Failure 500 {object} ErrorResponse
//@Router /createjob [post]
func CreateJob(w http.ResponseWriter, r *http.Request) {
     db, Ok := r.Context().Value("db").(*mongo.Database)
     if !Ok {
         response := map[string]interface{}{
		 "Error": "There was an error retreiving DB",
		 "message": "There was an Internal Error",
	 }
	 JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	 return
     }
     ctx, Ok := r.Context().Value("ctx").(context.Context)
     if !Ok {
          response := map[string]interface{}{
	      "Error": "There was an error retrieving context",
	      "message": "There is an Internal Error",
	  }
	  JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	  return
     }
     employerid, Ok := r.Context().Value("id").(primitive.ObjectID)
     if !Ok {
        response := map[string]interface{} {
            "Error": "There was an Error retrieving id",
            "message": "There was an Error retrieving User ID",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
     }
     var employer models.Employer
     if err := json.NewDecoder(r.Body).Decode(&employer); err != nil {
         response := map[string]interface{}{
            "Error": "There was an Error Decoding",
	        "message": "There was an Error Decoding Struct Body of the Employer",
	 }
	 JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	 return
     }
     employer.CreatedAt = time.Now()
     employer.EmployerID = employerid
     result, err := db.Collection("employer").InsertOne(ctx, employer)
     if err != nil {
         response := map[string]interface{} {
             "Error": "There was an Error Inserting in database",
	        "message": "Data was not inserted Successfully",
	 }
	 JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	 return
     }
     response := map[string]interface{} {
	   "message": "Job Created Successfully",
       "ID": result.InsertedID,
	   "title": employer.Title,
       "employerID": employer.EmployerID,
	   "description": employer.Description,
	   "createdat": employer.CreatedAt,
     }
     JsonResponse(w, ContentType, Content, http.StatusOK, response)
     return
}
// JobCreatedResponse defines the structure for the job creation response
