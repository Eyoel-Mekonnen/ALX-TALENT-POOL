package handlers

import (
     "net/http"
     "encoding/json"
     "context"
     "go.mongodb.org/mongo-driver/mongo"
     "go.mongodb.org/mongo-driver/bson"
     "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/models"
     "github.com/gorilla/mux"
     "go.mongodb.org/mongo-driver/bson/primitive"
)
/*
var ContentType string = "Content-Type"
var Content string = "application/json"
*/
type UpdatedProfile struct {
    Message     string      `json:"message"`
    ID          string      `json:"ID"`
}

// @Summary Update Profile of Student
// @Description This Endpoint allows Student to update their profile
// @Tags Student
// @Accept json
// @Produce json
// @Param id path string true "Student Profile id"
// @Param body body models.Student true "Student Profile"
// @Success 200 {object} UpdatedProfile
// @Failure 404 {object} StatusNotFound
// @Failure 500 {object} ErrorResponse
// @Router /updateProfile/{id} [Put]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
    db, Ok := r.Context().Value("db").(*mongo.Database)
    //id := r.URL.Query().Get("id")
    vars := mux.Vars(r)
    id := vars["id"]
    if !Ok {
	response := map[string]interface{} {
	    "error": "There was an error retrieveing db",
	    "message": "There is an Internal Error",
	}
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	/*
        w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{} {
	      "error": "There was an error retrieveing db",
	      "message": "There is an Internal Error",
	})
	*/
	return
    }

    ctx, Ok := r.Context().Value("ctx").(context.Context)
    if !Ok {
        w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{} {
	    "error": "There was an error retrieving CTX",
	    "message": "There is an Internal Error",
	})
	return
    }
    jobid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        response := map[string]interface{} {
            "error": "There was an error when Converting to mongodb ID",
            "message": "There was an Error When Converting to mongodb Primitive",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response) 
	return
    }
    var student models.Student
    json.NewDecoder(r.Body).Decode(&student)

    result, err := db.Collection("students").UpdateOne(ctx, bson.M{"_id": jobid},bson.M{"$set": student})
    if err != nil {
        response := map[string]interface{} {
            "error": "Error while UpdatingOne",
            "message": "There was an error retrieving DB",
        }
	JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	/*
        w.Header().Set("Contenty-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{} {
	     "error": "There was an error retrieving Result from DB",
	     "message": "Error while retrieving Student from DB",
	})
	*/
	return
    }
    if result.MatchedCount == 0 {
        response := map[string]interface{} {
            "error": "Error while UpdatingOne",
            "message": "No document found for the specified ID.",
        }
        JsonResponse(w, ContentType, Content, http.StatusNotFound, response)
	return
    }
    if result.ModifiedCount == 0 {
        response := map[string]interface{} {
            "message": "No Changes Were Made to the Document",
        }
	JsonResponse(w, ContentType, Content, http.StatusOK, response)
	return
    }
    response := map[string]interface{} {
        "message": "Update Student User Successfully",
	    "ID": jobid,
    }

    JsonResponse(w, ContentType, Content, http.StatusOK, response)
    return
}
