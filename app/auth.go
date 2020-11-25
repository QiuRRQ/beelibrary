package app

import (
	t "city/mytoken"
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
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in

		type MyCustomClaims struct {
			jwt.StandardClaims
			Email   string `json:"email"`
			Session string `json:"session"`
		}

		// token, err := jwt.ParseWithClaims(tokenPart, &t.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte("B33FR33"), nil
		// })

		token, err := jwt.ParseWithClaims(tokenPart, &t.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if jwt.SigningMethodHS256 != token.Method {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			secret := []byte("B33FR33")
			return secret, nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response = u.Message(false, err.Error())
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		claims := &t.CustomClaims{}
		if claim, ok := token.Claims.(*t.CustomClaims); ok && token.Valid {
			//log.Printf("%v %v %v %v \n", claim.Dbname, claim.Dbhost, claim.User_id, claim.StandardClaims.ExpiresAt)
			claims = claim
		} else {
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		ctx := context.WithValue(r.Context(), "email", claims.Email)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
