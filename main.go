package main

import (
    "net/http"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/handlers"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/utils"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/middlewares"
   _ "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/docs"
    httpSwagger "github.com/swaggo/http-swagger"
    //_"github.com/swaggo/swag/example/basics/docs"
)
// @title ALX Talent Pool API
// @version 1.0
// @description This is the API documentation for ALX Talent Pool system
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:80
// @BasePath /

func main() {
     r := mux.NewRouter()

     client, ctx, err := utils.DatabaseConnection()
     if err != nil {
         fmt.Println("Error: ", err)
     }
     db := client.Database("ALX-Talent-Pool-DB")
     fmt.Println("I am the db: ", db)
     if err != nil {
          fmt.Println(err)
     } else {
	  fmt.Println("Connection Successfull ", client)
     }
     r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
     r.HandleFunc("/", handlers.LandingPage).Methods("GET")
     r.HandleFunc("/signup", handlers.SignUp(db, ctx)).Methods("POST")
     r.HandleFunc("/signin", handlers.SignIn(db, ctx)).Methods("POST")
     r.HandleFunc("/signout", handlers.SignOut).Methods("GET")
     r.HandleFunc("/protectedEmployer", middlewares.Chain(handlers.ProtectedEmployer, middlewares.CheckRole("Employer", db, ctx), middlewares.VerifyJWT())).Methods("GET")
     r.HandleFunc("/protectedStudent", middlewares.Chain(handlers.ProtectedStudent, middlewares.CheckRole("Student", db, ctx), middlewares.VerifyJWT())).Methods("GET")
     r.HandleFunc("/createProfile", middlewares.Chain(handlers.CreateProfile, middlewares.CheckRole("Student", db, ctx), middlewares.VerifyJWT())).Methods("POST")
     r.HandleFunc("/updateProfile/{id}", middlewares.Chain(handlers.UpdateProfile, middlewares.CheckRole("Student", db, ctx), middlewares.VerifyJWT())).Methods("PUT")
     r.HandleFunc("/deleteProfile/{id}", middlewares.Chain(handlers.DeleteProfile, middlewares.CheckRole("Student", db, ctx), middlewares.VerifyJWT())).Methods("DELETE")
     r.HandleFunc("/createJob", middlewares.Chain(handlers.CreateJob, middlewares.CheckRole("Employer", db, ctx), middlewares.VerifyJWT())).Methods("POST")
     r.HandleFunc("/updatejob/{id}", middlewares.Chain(handlers.UpdateJob, middlewares.CheckRole("Employer", db, ctx), middlewares.VerifyJWT())).Methods("PUT")
     r.HandleFunc("/deletejob/{id}", middlewares.Chain(handlers.DeleteJob, middlewares.CheckRole("Employer", db, ctx), middlewares.VerifyJWT())).Methods("DELETE")
     r.HandleFunc("/uploadfile/{id}", middlewares.Chain(handlers.UploadFile, middlewares.CheckRole("Student", db, ctx), middlewares.VerifyJWT())).Methods("POST")
     //router.get('/jobs/:jobId/applications/:applicationId/download-cv', EmployerMiddleWare.getEmployer, Employer.getLinkDownloadCv);
     r.HandleFunc("/download/{jobid}/applications/{applicationid}", middlewares.Chain(handlers.DownloadFile, middlewares.CheckRole("Employer", db, ctx), middlewares.VerifyJWT())).Methods("GET")
     http.ListenAndServe(":80", r)
}
