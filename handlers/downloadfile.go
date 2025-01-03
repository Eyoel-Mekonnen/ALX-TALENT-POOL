package handlers

import (
     //"bytes"
     "context"
     "fmt"
     //"strings"
     //"strconv"
     //"io"
     "os"
     //"net/url"
     //"strings"
     "cloud.google.com/go/storage"
     "google.golang.org/api/option"
     "github.com/gorilla/mux"
     "go.mongodb.org/mongo-driver/mongo"
     "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/models"
     "go.mongodb.org/mongo-driver/bson/primitive"
     "go.mongodb.org/mongo-driver/bson"
     "time"
     "net/http"
)
type DownloadLinkURL struct {
    Message  string   `json:"message"`
    URLLINK  string  `json:"URLLINK"`
}
type BadRequest struct {
    Error    string   `json:"error"`
    Message  string   `json:"message"`
}
// @Summary Download a file by providing link
// @Description This Endpoint allows Employers to download a resume
// @Tags Employer
// @Accept json
// @Produce json
// @Param jobid path string true "job id"
// @Param applicationid path string true "application id"
// @Success 200 {object} DownloadLinkURL
// @Failure 400 {object} BadRequest
// @Failure 500 {object} ErrorResponseProfile
// @Router /download/{jobid}/applications/{applicationid} [Get]
func DownloadFile(w http.ResponseWriter, r *http.Request) {
      //In this retrieve objectName which is file path
      //bucketName which is the name of the bucket in th
    db, Ok := r.Context().Value("db") .(*mongo.Database)
    if !Ok {
        response := map[string]interface{} {
            "error": "There was an Error retrieveing db",
            "message": "There is an Internal Error",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    }

    vars := mux.Vars(r)
    jobid := vars["jobid"]
    if jobid == "" {
        response := map[string]interface{} {
            "error": "There was an error retrieving jobid from Query",
            "message": "There was an error retrieving JobID from Query",
        }
        JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
        return
    }
    applicationid := vars["applicationid"]
    if applicationid == "" {
       response := map[string]interface{} {
            "error": "There was an error retrieving applicationid from Query",
            "message": "There was an error retrieving ApplicationID from Query",
        }
        JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
        return
    }
    ctx, Ok := r.Context().Value("ctx").(context.Context)
    if !Ok {
       response := map[string]interface{} {
            "error": "There was an Error retrieving CTX",
            "message": "There is an Internal Error",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    }
    employerid, Ok := r.Context().Value("id").(primitive.ObjectID)
    if !Ok {
        response := map[string]interface{} {
            "error": "There was an error Retrieving ID",
            "message": "There is an Internal Error",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    }
    id, err := primitive.ObjectIDFromHex(jobid)
    if err != nil {
        response := map[string]interface{} {
            "error": "Error When converting to Primitive ObjectID",
            "message": "There was An Error Converting to PrimitiveObjectID",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    }
    var employer models.Employer
    err1 := db.Collection("employer").FindOne(ctx, bson.M{"_id": id, "employerid": employerid}).Decode(&employer)
    if err1 != nil && err1 != mongo.ErrNoDocuments {
        response := map[string]interface{} {
            "error": err1.Error(),
            "message": "The Error occured Internally When retrieving Value",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    } else if err1 == mongo.ErrNoDocuments {
        response := map[string]interface{} {
            "error": "There No user with id and jobid",
            "message": "Error Finding id and jobid",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    }
    applicationID, err := primitive.ObjectIDFromHex(applicationid)
    if err != nil {
        response := map[string]interface{} {
            "error": "There was an Error converting to Primitive ObjectID",
            "message" : "There was An Error Converting to PrimitiveObjectID",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    }
    var application models.Application
    err2 := db.Collection("applications").FindOne(ctx, bson.M{"_id": applicationID, "jobid": id}).Decode(&application)
    if err2 != nil && err2 != mongo.ErrNoDocuments {
        response := map[string]interface{} {
            "error": "There was an error on the find operation",
            "message": "The error occured internally when retrieving Value",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
    } else if err2 == mongo.ErrNoDocuments {
        
        fmt.Println(applicationID)
        
        fmt.Println(id)
        response := map[string]interface{} {
            "error": "There is no Application with that id and applicationid",
            "message": "No applicationID or jobID was found with those credentials",
        }
        JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
        return
    }
    bucketName := "gethired-2cb65.appspot.com"
    objName := application.FilePath
    fmt.Println(objName)
    fmt.Println(bucketName)
    link, err := downloadFromFirebase(bucketName, objName)
    if err != nil {
        response := map[string]interface{} {
            "error": err.Error(),
            "message": "URL LINK was not retrieved Successfully",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    }
    //decodedURL := strings.ReplaceAll(link, "\\u0026", "&")
    fmt.Println(link)
    /*
    response := map[string]interface{} {
        "message": "Application URL successfully Retrieved",
        "URLLINK":link,
    }
    JsonResponse(w, ContentType, Content, http.StatusOK, response)
    return
    */
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)

    // Write the URL directly to the response
    fmt.Fprintln(w, link)
}

func downloadFromFirebase(bucketName string, objName string) (string, error) {
     ctx := context.Background()
     filePath := os.Getenv("FIREBASE1")
     /* I am creating a New client that allows me to interact with Firebase */
     client, err := storage.NewClient(ctx, option.WithCredentialsFile(filePath))

     if err != nil {
         fmt.Println("Error is inside here the client")
         fmt.Println(err.Error())
         return "", fmt.Errorf("failed to create storage client: %v", err)
     }
     defer client.Close()
     /* I am getting the bucket with the specified from the client */
     //bucket := client.Bucket(bucketName)
     /*Am getting the reference of the file path */
     //obj := bucket.Object(objName)
     /* Creating Expriation for it */
     expiration := time.Now().Add(1 * time.Hour)
     /* Singing the refrecne with public timed URL to be downloaded */
     urL, err := client.Bucket(bucketName).SignedURL(objName, &storage.SignedURLOptions {
          Method: "GET",
   	      Expires: expiration,
     })
     if err != nil {
         fmt.Println("The error url generation")
         fmt.Println(err.Error())
         return "", fmt.Errorf("failed to generate signed URL: %v", err)
     }
     /*Returning this url */
    
     if err != nil {
        return "", fmt.Errorf("failed to generate signed URL: %v", err)
     }
     //decodedURL := urL
     //decodedURL = strings.ReplaceAll(urL, "\\u0026", "&")
     /*
     decodedURL = strings.TrimSpace(decodedURL)
     decodedURL = strings.Trim(decodedURL, `"`)
     */
     //fmt.Println(urL)
     return urL, nil
}
