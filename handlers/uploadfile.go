package handlers

import (
    "bytes"
    "context"
    "fmt"
    "io"
    "strings"
    //"mime/multipart"
    "net/http"
    //"strings"
    "cloud.google.com/go/storage"
    "google.golang.org/api/option"
    "os"
    "github.com/gorilla/mux"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/models"
    "go.mongodb.org/mongo-driver/bson/primitive"
)
type UploadedFile struct {
    Message     string              `json:"message"`
    ID          primitive.ObjectID  `json:"ID"`
}

// @Summary Upload File
// @Description This Endpoint allows student to upload Resume
// @Tags Student
// @Accept json
// @Produce json
// @Param id path string true "Upload file to a jobID"
// @Success 200 {object} UploadedFile 
// @Failure 400 {object} BadRequest
// @Failure 500 {object} ErrorResponse
// @Router /uploadfile/{id} [Post]
func UploadFile(w http.ResponseWriter, r *http.Request) {
    db, Ok := r.Context().Value("db").(*mongo.Database)
    if !Ok {
        response := map[string]interface{} {
            "error": "There was an error retrieveing db",
            "message": "There is an Internal Error",
	    }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    }
    vars := mux.Vars(r)
    id := vars["id"]
    id = strings.TrimSpace(id)
    /*The above is the job id that is passed in the url */
    if id == "" {
        fmt.Println("I am inside here because it is empty")
        response := map[string]interface{} {
            "error": "There was an error retrieving id from Query",
            "message": "There was an error retieving ID of job from Query",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    }
    ctx, Ok := r.Context().Value("ctx").(context.Context)
    if !Ok {
        response := map[string]interface{} {
            "error": "There was an error retrieving CTX",
	        "message": "There is an Internal Error",
         }
         JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
         return
    }
    /*Job seeker id is the one below here */
    jobseekerid , Ok := r.Context().Value("id").(primitive.ObjectID)
    if !Ok {
        response := map[string]interface{} {
            "error": "There was an error retrieving ID",
            "message": "There is an Internal Error",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
    }
     /*
      * Here am parsing the file that is passed and checking the size limit 
      * Doing the above populates the retreived File int the r.FileForm("file")
      */
    err := r.ParseMultipartForm(10 << 20)

    if err != nil {
        response := map[string]interface{} {
            "error": "Unable To Parse Form",
            "message": "There was An Error Parsing Form",
        }
        JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
	    return
    }
     /*
      Here the file that is passed in the file from frontend in assigned to the variable file
     */
    file, _, err := r.FormFile("file")
    if err != nil {
        response := map[string]interface{} {
            "error": "Unable To Parse Form",
            "message": "There was An Error Parsing Form",
        }
        JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
        return
     }
     /* 
       We close this file when the function finishes
     */
     defer file.Close()
     /* G
      * Here I am creating a dynamic buffer that cna hold of type bytes
     */
    buf := new(bytes.Buffer)
     /* Here am copying the content of file to the buffer i created */
    _, err  = io.Copy(buf, file)
    if err != nil {
        response := map[string]interface{} {
            "error": "Unable To Parse Form",
            "message":"There was an Error Parsing Form",
        }
        JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
        return
     }
     /*I am creating an object file name that will be used name the file*/
     objName := "cvs/" + r.FormValue("filename")
     fmt.Println(objName)
     /*Calling the uploadFileToFireBase passing the name the populated buf.Bytes*/
     //err = uploadFileToFireBase(objName, buf.Bytes())
     /*
     if err != nil {
        response := map[string]interface{} {
            "error": err,
            "message": "owhh!!There was an Error to Upload file to Firebase",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
	    return 
     }
     */
     /* I need to check if the job id exists right */
     jobId, err2 := primitive.ObjectIDFromHex(id)
     fmt.Println(jobId)
     if err2 != nil {
        response := map[string]interface{} {
            "Error": "ID could not be converted",
            "message": "ID of job is not being converted correctly",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
     }
     fmt.Println(id)
     var job models.Employer
     err1 := db.Collection("employer").FindOne(ctx, bson.M{"_id": jobId}).Decode(&job)
     if err1 != nil {
        response := map[string]interface{} {
            "error": err1.Error(),
            "message": "The job you are applying to does not exist",
        }
        JsonResponse(w, ContentType, Content, http.StatusBadRequest, response)
        return
     }
     err = uploadFileToFireBase(objName, buf.Bytes())
     if err != nil {
        response := map[string]interface{} {
           "error": err.Error(), 
           "message": "Error uplaoding to firebase",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
     }
     var application models.Application
     application.CreatedAt = time.Now()
     application.JobSeekerID = jobseekerid
     application.JobID = jobId
     //jobId, err := primitive.ObjectIDFromHex(id)
     application.FilePath = objName
     result, err := db.Collection("applications").InsertOne(ctx, application)
     if err != nil {
        response := map[string]interface{} {
            "error": "There was an Error Inserting in database",
            "message": "Data was not inserted Successfully",
        }
        JsonResponse(w, ContentType, Content, http.StatusInternalServerError, response)
        return
     }
     response := map[string]interface{} {
         "message": "Application Added Successuflly to DB",
         "ID": result.InsertedID,
     }
     JsonResponse(w, ContentType, Content, http.StatusOK, response)
     return
}

func uploadFileToFireBase(objName string, fileContent []byte) error {
    ctx := context.Background()
    filePath := os.Getenv("FIREBASE1")
    fmt.Println(filePath)
    /* I want you to create a client which will communicate with the credentials I gave you */
    client, err := storage.NewClient(ctx, option.WithCredentialsFile(filePath))
    fmt.Println("FIREBASE Environment Variable:", os.Getenv("FIREBASE"))

    if err != nil {
         fmt.Println("I ma the client error")
         return fmt.Errorf("failed to create storage client: %v", err)
    }
    /* We will close this client when we are done with this function */
    defer client.Close()
    /* I am constructing my bucket name into a variable */
    bucketName := "gethired-2cb65.appspot.com"
    /*I am passing this name and retreiving the bucket from firebase and storing it into bucket */
    bucket := client.Bucket(bucketName)
    fmt.Println("Bucket is : %v", bucket)
    /* Here based on the retrieved bucket am creating a reference and name where the file will be stroed
    * and then writing it to that referece */
    obj := bucket.Object(objName)
    /* Creating a new Writer That will be used to write the file to the specified Object */
    writer := obj.NewWriter(ctx)
    /* The line below is acutally writing the file to the bucket */
    fmt.Println("Writer initialized:", writer)
    _, err = writer.Write(fileContent)
   if err != nil {
       fmt.Println("______________")
       fmt.Println(err)
       fmt.Println("_______________")
       fmt.Println("I am the error is when using writer.Write %v", err)
       return fmt.Errorf("failed to upload file: %v", err)
   }
   err = writer.Close()
   if err != nil {
       fmt.Println("I am the error when trying to close", err)
       return fmt.Errorf("failed to upload file: %v", err)
   }
   fmt.Println("File uploaded successfully to Firebase Storage!")
   return nil
}
