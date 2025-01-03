package handlers


import (
     "net/http"
)
type SignedOut struct {
    Message     string      `json:"message"`
}

// @Summary SignOut
// @Description This endpoints allows users to signou
// @Tags Employer/Student
// @Accept json
// @Produce json
// @Success 200 {object} SignedOut
// @Router /signout [Get]
func SignOut(w http.ResponseWriter, r *http.Request) {
    response := map[string]interface{} {
        "Message": "You SignedOut Successfully",
    }
    JsonResponse(w, ContentType, Content, http.StatusOK, response)
    return
}
