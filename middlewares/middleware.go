package middlewares


import (
    "net/http"
    //"github.com/Eyoel-Mekonnen/ALX-TALENT-POOL/authentication"
)
/*Here the first function will be authentication function
 *The second function will be Role checker
 *Then do the thing that should be done
*/

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc{
     for _, m := range middlewares {
         f = m(f)
     }
     return f
}
