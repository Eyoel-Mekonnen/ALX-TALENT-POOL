package handlers

import (
    "fmt"
    "bytes"
    "net/http"
    "net/http/httptest"
    "github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/models"
    "testing"
)
/*
dd - > to delete a line
yy -> To copy
P/p -> to paste it
d -> cut - > P/p to paste it
shift + A -> to go to the end of the line and in insert mode
or $ -> to just go to the end of the line
*/

type Test struct {
    TestName        string
    data            models.User
    StatusCode      int
    ExpectedOutput  map[string]interface{} 
}
var tests = []Test {
    {
        TestName: "Invalid Json Format",
        data: models.User{Username: "", Email: "", Password: "", Role: ""}
        StatusCode: http.StatusBadRequest
        ExpectedOutput: map[string]interface{} {
                "error": "Invalid Json format Passed",
                "message": "Ensure Invalid JSON is passed", 
            }
    },
    {
        TestName: "Successful Insertion",
        data: models.User{Username: "Eyoel Mekonnen", Email: "eyoelakatommyshellby@gmail.com", Password: "eyoel123@alx", Role: "Student"},
        StatusCode: http.StatusOK,
        ExpectedOutput: map[string]interface{} {
                "Message": "Successfully Created",
                "OutputMessage": {"name":user.Name, "email": user.Email, "password": user.Password},
            }
    },
    {
        TestName: "User Already Exists",
        data: models.User{Username: "Eyoel Mekonnen", Email: "eyoelakatommyshellby@gmail.com", Password: "eyoel123@alx", Role: "Student"},
        StatusCode: http.StatusOK,
        ExpectedOutput: map[string]interface{} {
            }
    }
}

func SignUpTest() {
    /* Here httptest.NewRequest Does not accept JSON so we need to use byte */
    var user models.User = {"Eyoel Mekonnen", "eyoelmekonnenbogale@gmail.com", "Student"}
    jsonBytes, err := json.Marshal(user)
    req := httptest.NewRequest("POST", "/signup", bytes.NewReader(jsonBytes))
    w := httptest.NewRecorder()
    req.Header.Set("Content-Type", "application/json")
    SignUp(req, w)
    var result map[string]interface{}
    resp := w.Result()
    /*
        if not json you can use body, _ := io.ReadAll(resp.Body)

    */
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        t.Fatalf("Failed to Decode response: %v", err)
    }
    if w.Code != httpStatusOK {
        t.Fatalf("Expected status OK but got %v: ", w.Code)
    }
}
