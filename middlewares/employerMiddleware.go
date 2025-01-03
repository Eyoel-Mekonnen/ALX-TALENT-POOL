package middlewares

import (
    "fmt"
    "context"
    "net/http"
    "encoding/json"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/authentication"
    "go.mongodb.org/mongo-driver/mongo"
)

//type Middleware func(http.HandlerFunc) http.HandlerFunc
func VerifyJWT() Middleware {
     return func(checkrole http.HandlerFunc) http.HandlerFunc {
          return func(w http.ResponseWriter, r *http.Request) {
        
        fmt.Println("I am inside here")
		token := r.Header.Get("Authorization")
		token = token[len("Bearer "):]
		if token == "" {
		    w.Header().Set("Content-Type", "application/json")
		    w.WriteHeader(http.StatusUnauthorized)
		    json.NewEncoder(w).Encode(map[string]interface{}{
		        "error": "Unauthorized",
			    "message": "Missing Token",
		    })
		    return
		}
        fmt.Println("I am inside here")
		tokenString, role, id, err := authentication.VerifyToken(token)
		if err != nil && (role == "" || tokenString == ""){
		    /** Return Error of not being authorized */
		     w.Header().Set("Content-Type", "application/json")
		     w.WriteHeader(http.StatusUnauthorized)
		     json.NewEncoder(w).Encode(map[string]interface{} {
		        "error": "Unauthorized",
			    "message": "Token is Invalid or Missing Claim",
		     })
		     return
		}
		ctx := context.WithValue(r.Context(), "role", role)
        ctx = context.WithValue(ctx, "id", id)
		r = r.WithContext(ctx)
		checkrole(w,r) //DO something here with the data
	  }
     }
}

func CheckRole(Role string, db *mongo.Database, ctx context.Context) Middleware {
     return func(f http.HandlerFunc) http.HandlerFunc {
         return func(w http.ResponseWriter, r *http.Request) {
	        role, Ok := r.Context().Value("role").(string)
	        if !Ok || role != Role {
	            w.Header().Set("Content-Type", "application/json")
		        w.WriteHeader(http.StatusUnauthorized)
		        message := fmt.Sprintf("This is Protected site for %v", role)
		        json.NewEncoder(w).Encode(map[string]interface{}{
		            "error": "You are not allowed to go to this site",
		            "message": message,
		        })
		        return
	        }
	        ctx := context.WithValue(r.Context(), "db", db)
		    ctx = context.WithValue(ctx, "ctx", ctx)
		    r = r.WithContext(ctx)
            f(w, r)
	    }
     }
}
