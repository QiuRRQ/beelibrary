package app

import (
	u "city/utils"
	"context"
	"fmt"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/login", "/api/user"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                           //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		fmt.Println(len(splitted))
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in

		type MyCustomClaims struct {
			Dbname  string `json:"dbname"`
			Dbhost  string `json:"dbhost"`
			User_id string `json:"user_id"`
			jwt.StandardClaims
		}

		token, err := jwt.ParseWithClaims(tokenPart, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("B33FR33"), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		claims := &MyCustomClaims{}
		if claim, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			//log.Printf("%v %v %v %v \n", claim.Dbname, claim.Dbhost, claim.User_id, claim.StandardClaims.ExpiresAt)
			claims = claim
		} else {
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		ctx := context.WithValue(r.Context(), "dbname", claims.Dbname)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
