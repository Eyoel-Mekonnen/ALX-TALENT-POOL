package handlers

import (
    "encoding/json"
    "net/http"
    "go.mongodb.org/mongo-driver/mongo"
    "context"
    "fmt"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/models"
)
type ErrorResponseProfile struct {
    Error   string    `json:"error"`
    Message string    `json:"message"`
}

var ContentType string = "Content-Type"

var Content string = "application/json"

// @Summary Create a new student profile
// @Description This endpoint creates a new student profile with personal and professional details
// @Tags Student
// @Accept json
// @Produce json
// @Param body body models.Student true "Student profile details"
// @Success 200 {object} models.CreateResponse
// @Failure 500 {object} ErrorResponseProfile
// @Router /createProfile [post]
func CreateProfile(w http.ResponseWriter, r *http.Request) {
    db, Ok := r.Context().Value("db").(*mongo.Database)
    if !Ok  {
        w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{} {
	    "Error": "There was an error retriveing db",
	    "Message": "There is an Internal Error",
	})
	return
    }
    ctx, Ok := r.Context().Value("ctx").(context.Context)
    if !Ok {
	response := map[string]interface{} {
	     "Error": "There was an error retriveing context",
	     "Message": "There is an Internal Error",
	}

	JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
    /*
        w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{} {
	    "error": "There was an error retriveing context",
	    "message": "There is an Internal Error",
	})
    */
	return
    }
    var student models.Student
    json.NewDecoder(r.Body).Decode(&student)
    fmt.Println(student)
    result, err := db.Collection("students").InsertOne(ctx, student)
    if err != nil {
	    response := (map[string]interface{} {
            "Error": "There was an error inserting in database",
	        "Message": "Data was not inserted Successfully",
	    })
     JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	 return
    } 
    response := map[string]interface{} {
	"id": result.InsertedID,
    "name": student.Name,
	"jobseekerlocation": student.JobSeekerLocation,
	"telegram": student.Telegram,
	"github": student.Github,
	"linkedin": student.Linkedin,
	"skills": student.Skills,
	"languages": student.Languages,
	"professionalsummary": student.ProfessionalSummary,
	"portfolioLink": student.PortfolioLink,
	"education": student.Education,
    }
    JsonResponse(w, ContentType, Content, http.StatusOK, response)
    return
}
