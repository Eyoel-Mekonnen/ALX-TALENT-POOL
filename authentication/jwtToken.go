package authentication


import (
     "fmt"
     "time"
     "github.com/golang-jwt/jwt/v5"
     "go.mongodb.org/mongo-driver/bson/primitive"
)

var secretKey = []byte("secret-key")
func CreateToken(username string, id primitive.ObjectID, role string) (string, error) {
     token := jwt.NewWithClaims(jwt.SigningMethodHS256,
     jwt.MapClaims{
         "username": username,
         "id": id.Hex(), 
	     "role": role,
	     "exp": time.Now().Add(time.Hour * 24).Unix(),
     })
     tokenString, err := token.SignedString(secretKey)
     if err != nil {
         return "", err
     }
     return tokenString, nil
}

func VerifyToken(tokenString string) (string, string, primitive.ObjectID, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
        //Here i need to retireve the role right
	return secretKey, nil
    })
    if err != nil {
        return "", "", primitive.NilObjectID, err
    }
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        role, ok := claims["role"].(string)
	    if !ok {
	         return "", "", primitive.NilObjectID, fmt.Errorf("role claim is not a string or missing")
	    }
        idString, ok := claims["id"].(string)
        if !ok {
            return "", "", primitive.NilObjectID, fmt.Errorf("role claim is not a string or missing")
        }
        id, err := primitive.ObjectIDFromHex(idString)
        if err != nil {
            return "", "", primitive.NilObjectID, fmt.Errorf("invalid ObjectID in token")
        }
	    return token.Raw, role, id, nil
    } else {
        return "", "", primitive.NilObjectID, err
    }
}
