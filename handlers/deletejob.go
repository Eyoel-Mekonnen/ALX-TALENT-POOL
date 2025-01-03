package handlers

import (
     "net/http"
     "context"
     "go.mongodb.org/mongo-driver/mongo"
     "go.mongodb.org/mongo-driver/bson"
     "github.com/gorilla/mux"
     "go.mongodb.org/mongo-driver/bson/primitive"
)
type DeleteJobID struct {
    Message     string     `json:"message"`
    ID          string     `json:ID`
}
// @Summary Create a new student profile
// @Description This endpoint Deletes A job based on ID
// @Tags Employer
// @Accept json
// @Produce json
// @Param id path string true "Delete Job based on ID"
// @Success 200 {object} DeleteJobID
// @Failure 500 {object} ErrorResponse
// @Router /deletejob/{id} [delete]
func DeleteJob(w http.ResponseWriter, r *http.Request) {
     db, Ok := r.Context().Value("db").(*mongo.Database)
     if !Ok {
         response := map[string]interface{} {
             "error": "There was an error retrieving db",
	     "message": "There is an Internal Error",
	 }
	 JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	 return
     }

     //id := r.URL.Query().Get("id")
     vars := mux.Vars(r)
     id := vars["id"]
     if id == "" {
	 response := map[string]interface{} {
              "error": "Error Retrieving Id from Query",
	      "message": "There is an Error Retrieving ID",
	 }
	 JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
	 return
     }
     ctx, Ok := r.Context().Value("ctx").(context.Context)

     if !Ok  {
         response := map[string]interface{} {
             "error": "There was an error retrieving CTX",
             "message": "There is an Internal Error CTX",
         }
         JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
         return
     }
     jobid, err := primitive.ObjectIDFromHex(id)
     result, err := db.Collection("employer").DeleteOne(ctx, bson.M{"_id": jobid})
     if err != nil {
           response := map[string]interface{} {
             "error": "Error Deleting from DB",
             "message": "There was an Internal Error When Deleting from db",
         }
         JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
         return
     }
     if result.DeletedCount == 0 {
         response := map[string]interface{} {
             "error": "Document Not Found",
             "message": "No Document was found with specified Criteria",
         }
         JsonResponse(w, ContentType, Content, http.StatusNotFound, response)
         return
     }
     response := map[string]interface{} {
         "message": "Deleted Student User Successfully",
         "ID": id,
     }
     JsonResponse(w, ContentType, Content, http.StatusOK, response)
     return
}