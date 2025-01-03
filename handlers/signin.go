package handlers


import (
    "fmt"
    "encoding/json"
    "net/http"
    "context"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/models"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/authentication"
    //"golang.org/x/crypto/bcrypt"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")
func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
        jwt.MapClaims{
	    "username": username,
	    "exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
	    return "", err
	}
	return tokenString, nil
}
/* The jwt.Parse is used to separte the token return the two body parts sign it with
 * the provided secret Key from anonymous function and then recreate the jwt and compare it
 * with the original one
 */


func verifyToken(tokenString string) error {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	/* Here the token *jwt.Token is used if we want to some modification */
        return secretKey, nil
    })
    if err != nil {
        return err
    }
    if !token.Valid {
	/* Creates and return an error object with the message below */
        return fmt.Errorf("Invalid Token")
    }
    return nil
}

/*
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPassWordHash(password, hash string) bool {
    err := bcrypt.CompareHashPassword([]byte(hash), []byte(password))
    return err == nil
}
*/
type SuccessMessage struct{
    Token       string      `json:"token"`
    Message     string      `json:"message"`
}
type Unauthorized struct {
    Error       string      `json:"error"`
    Message     string      `json:"message"`
}
type SigninBody struct {
    Email        string      `json:"email"`
    Password     string      `json:"password"`
}

// @Summary Signin
// @Description Allows Students and Employer to Signin
// @Tags Employer/Student
// @Accept json
// @Produce json
// @param body body SigninBody true "User signin detail"
// @Success 200 {object} SuccessMessage
// @Failure 400 {object} BadRequest 
// @Failure 401 {object} Unauthorized
// @Failure 500 {object} ErrorResponseProfile
// @Router /signin [POST]
func SignIn(db *mongo.Database, ctx context.Context) http.HandlerFunc {

    return func (w http.ResponseWriter, r *http.Request) {
        var user = models.User{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
	    w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusBadRequest)
	    json.NewEncoder(w).Encode(map[string]interface{} {
	        "error": "Invalid Json format Passed",
		"message": "Ensure Invalid JSON is passed",
	    })
	    return
	}
	var result models.User
	err := db.Collection("users").FindOne(ctx, bson.M{"email": user.Email}).Decode(&result)
	if err != nil && err != mongo.ErrNoDocuments {
	    w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusInternalServerError)
	    json.NewEncoder(w).Encode(map[string]interface{}{
	        "error": "There was an error on the find operation",
		"message": "The error occured internally when retrieving value",
	    })
	    return
	} else if err == mongo.ErrNoDocuments {
	    w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusUnauthorized)
	    json.NewEncoder(w).Encode(map[string]interface{} {
	        "error": "There is no user with That Email",
		"message": "Unauthorized User",
	    })
	    return
	}
	var passwordHashed string = result.Password //Here it is a hashed password right
	hashedPassword, err1 := authentication.HashPassword(user.Password)
	if err1 != nil || hashedPassword == ""{
	     w.Header().Set("Content-Type", "application/json")
	     w.WriteHeader(http.StatusInternalServerError)
	     json.NewEncoder(w).Encode(map[string]interface{}{
	         "error": "There was an error creating hash",
		 "message": "An Internal Error when creating hash",
	     })
	     return
	}
	//fmt.Println(result.Email, result.Password)
	check := authentication.CheckPasswordHash(user.Password, passwordHashed)
	//fmt.Fprintf(w, "hashedPassword %v \n PasswordHashed(retreived fromDB) %v", hashedPassword, passwordHashed)
	if check != true {
	    w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusUnauthorized)
	    json.NewEncoder(w).Encode(map[string]interface{}{
	        "error": "Password is Incorrect",
		    "message": "Unauthorized User",
	    })
	    return
	}
	token, err := authentication.CreateToken(result.Username, result.ID, result.Role)
	if err != nil  || token == ""{
	    w.Header().Set("Content-Type", "application/json")
	    w.WriteHeader(http.StatusInternalServerError)
	    json.NewEncoder(w).Encode(map[string]interface{}{
	        "error": "Token Was Not created ",
		    "message": "Internal Error Creating Token",
	    })
	    return
	}
	tokenValue := map[string]interface{} {
	    "token": token,
	    "message": "Token Created Successfully",
	}
	json.NewEncoder(w).Encode(tokenValue)
	return
    }
}
