package handlers

import (
    "fmt"
    "encoding/json"
    "context"
    "net/http"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/models"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/authentication"
    "go.mongodb.org/mongo-driver/mongo"
    //"golang.org/x/crypto/bcrypt"
    "go.mongodb.org/mongo-driver/bson"
    //"go.mongodb.org/mongo-driver/mongo/options"
)
/*
func HashPassword(password string) (string, error) {
     bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
     return string(bytes), err
}

func CheckPassWordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
*/

type SignedUp struct {
    Message         string      `json:"message"`
    OutputMessage   string      `json:"outputmessage"`
}

// @Summary SignUp 
// @Description Employer and Studnets will signout from this page
// @Tags Employer/Student
// @Accept json
// @Produce json
// @Param body body models.User true "User signup details"
// @Success 200 {object} SignedUp
// @Failure 400 {object} BadRequest
// @Failure 500 {object} ErrorResponse
// @Router /signup [Post]
func SignUp(db *mongo.Database, ctx context.Context) http.HandlerFunc{
    return func(w http.ResponseWriter, r *http.Request) {
        var user = models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
	     //fmt.Printf("I am the json body: %v and the error: %v \n", r.Body, err)
	     //http.Error(w, "Error decoding json/not correct format", http.StatusInternalServerError)

	     response := map[string]interface{}{
	        "error": "Invalid Json format Passed",
		    "message": "Ensure Invalid JSON is passed",
	     }
         JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
	     return
	}
	var password string = user.Password
	hashed, err := authentication.HashPassword(password)
	if err != nil {
	    response := map[string]interface{}{
	         "error": "Error Hashing Password",
		 "message": "Error is Internall Server Error",
	    }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
	}
	user.Password = hashed
	var result models.User
	fmt.Println(user)
        err1 := db.Collection("users").FindOne(ctx, bson.M{"email": user.Email}).Decode(&result)
	if err1 != nil && err1 != mongo.ErrNoDocuments {
	    response := map[string]interface{}{
	        "error": "There was an error on the find operation",
		    "message": "The error occured internally when retrieving value",
	    }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
	} else if err1 == nil {
	    message := fmt.Sprintf("A user with email %v Already Exists", result.Email)
	    response := map[string]interface{} {
	         "message->User Already Exists": message,
	    }
        JsonResponse(w, ContentType, Content, http.StatusOK, response)
	    return
	}
	Insertresult, err2 := db.Collection("users").InsertOne(ctx, user)
	if err2 != nil {
        response := map[string]interface{} {
            "Error": "There was an Error inserting into database",
            "message": "An error occurred when inserting to database",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
	} else {
	    message := fmt.Sprintf("name: %s email: %s password: %s DocumentID: %v", user.Username, user.Email, user.Password, Insertresult.InsertedID)
        response := map[string]interface{} {
            "Message": "Successfully Cretated",
            "OutputMessage": message,
        }
        JsonResponse(w, ContentType, Content, http.StatusOK, response)
	    return
	}
 }
}
