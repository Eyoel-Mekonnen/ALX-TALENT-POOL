package handlers

import (
    "net/http"
    "encoding/json"
    "context"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/models"
    "github.com/gorilla/mux"
)
type StatusNotFound struct {
    Error       string      `json:"error"`
    Message     string      `json:"message"`
}

type UpdatedJob struct {
    Message     string      `json:"message"`
    ID          string      `json:"ID"`
}

// @Update Job
// @Summary Update Existing Job
// @Description Update Existing Job
// @Tags Employer
// @Accept json
// @Produce json
// @Param body body models.Employer true "One of the fields are Updated"
// @Param id path string true "job id"
// @Success 200 {object} UpdatedJob
// @Failure 400 {object} BadRequest
// @Failure 404 {object} StatusNotFound
// @Failure 500 {object} ErrorResponse
// @Router /updatejob/{id} [Put]
func UpdateJob(w http.ResponseWriter, r *http.Request) {
    db, Ok := r.Context().Value("db").(*mongo.Database)
    //id := r.URL.Query().Get("id")
    vars := mux.Vars(r)
    id := vars["id"]
    if id == "" {
        response := map[string]interface{}{
	    "error": "There was an error retrieving id from Query",
	    "message": "There was an error retrieving ID from your Query",
	}
	JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
	return
    }
    if !Ok {
        response := map[string]interface{}{
	    "error": "There was an error retrieving db",
	    "message": "There was an Internal Error",
	}
	JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	return
    }


   ctx, Ok := r.Context().Value("ctx").(context.Context)
   if !Ok {
        response := map[string]interface{} {
	    "error": "There was an error retrieving Context",
	     "message": "There is an Internal Error",
	}
	JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	return
   }

   var employer models.Employer
   if err := json.NewDecoder(r.Body).Decode(&employer); err != nil {
       response := map[string]interface{}{
           "error": "Error Decoding struct request body",
	   "message": "There was an error decoding the struct body of request",
       }
       JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
       return
   }
   jobid, err := primitive.ObjectIDFromHex(id)
   if err != nil {
       response := map[string]interface{} {
           "error": "Error while Updating",
           "message": "There was an error Updating Employer Document",
       }
       JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
       return
   }
   result, err := db.Collection("employer").UpdateOne(ctx, bson.M{"_id": jobid}, bson.M{"$set": employer})
   if err != nil {
       response := map[string]interface{} {
           "error": "Error while Updating",
	   "message": "There was an error Updating Employer Document",
       }
       JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
       return
   }
   if result.MatchedCount == 0 {
       response := map[string]interface{}{
            "error": "Error while UpdatingOne",
	    "message": "No document found for the specified ID.",
       }
       JsonResponse(w, ContentType, Content, http.StatusNotFound, response)
       return
   }
   if result.ModifiedCount == 0 {
       response := map[string]interface{}{
           "message": "No Changes Were Made to the Document",
       }
       JsonResponse(w, ContentType, Content, http.StatusOK, response)
       return
   }
   response := map[string]interface{}{
       "message": "Update Employer User Successfully",
       "ID": id,
   }
   JsonResponse(w, ContentType, Content, http.StatusOK, response)
   return
}
