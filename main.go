package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/mitchellh/mapstructure"
	"io"
	"log"
	"net/http"
	"time"
)

// User details
type User struct {
	EmailAddress string `json:"emailaddress,omitempty"`
	Password     string `json:"password,omitempty"`
	JwtToken	 string `json:"token,omitempty"`
}

// Message to return from end points
type Message struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// JwtToken for validating sessions
type JwtToken struct {
	Token string `json:"token"`
}


var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
var users []User


// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/register", CreateUser).Methods("POST")
	router.HandleFunc("/account", ProtectedEndpoint).Methods("GET")
	router.HandleFunc("/healthcheck", HealthCheckHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

// The handler for the login endpoint
func Login(w http.ResponseWriter, r *http.Request) {
	var newUser User
	var accountFound bool
	params := mux.Vars(r)
	newUser.EmailAddress = params["emailAddress"]
	newUser.Password = params["password"]

	for i := 0; i < len(users); i++ {
		if users[i].EmailAddress == newUser.EmailAddress && users[i].Password == newUser.Password {
			accountFound = true
			tokenString := generateTokenString(newUser)
			users[i].JwtToken = tokenString
			json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
		}
	}
	if accountFound == false {
		json.NewEncoder(w).Encode(Message{"1", "Could not find account"})
	}
}

// Handler for the register endpoint
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	var isNewUser = true
	params := mux.Vars(r)
	newUser.EmailAddress = params["emailAddress"]
	newUser.Password = params["password"]

	for i := 0; i < len(users); i++ {
		if users[i].EmailAddress == newUser.EmailAddress {
			isNewUser = false
		}
	}
	if isNewUser == true {
		users = append(users, newUser)
		json.NewEncoder(w).Encode(Message{"0", "New User created"})
	} else {
		json.NewEncoder(w).Encode(Message{"1", "That user already exists"})
	}

}

// ProtectedEndpoint can only be seen with a valid JWT token
func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	tokenValue := r.URL.Query()["token"][0]
	token, _ := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})

	claims := token.Claims.(jwt.MapClaims)

	//if the claims are found to have been expired but are tied to a user account and have been created within 7 days
	//refresh the token
	if claims.VerifyExpiresAt(time.Now().Add(time.Hour * 24 * 7 ).Unix(), true){
		for i := 0; i < len(users); i++ {
			if users[i].JwtToken == tokenValue {
				token = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"username": users[i].EmailAddress,
					"password": users[i].Password,
					"exp":      time.Now().Add(time.Hour * 24).Unix(),
				})
			}
		}
	}
	if token.Valid {
		var user User
		mapstructure.Decode(claims, &user)
		json.NewEncoder(w).Encode(user)
	} else {

		json.NewEncoder(w).Encode(Message{"1", "Invalid authorization token"})
	}
}


// Healthcheck endpoint
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
    // A very simple health check.
    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/json")
    io.WriteString(w, `{"alive": true}`)
}

// Generate a jwt token string for a user
func generateTokenString(user User) (string){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.EmailAddress,
		"password": user.Password,
		"exp": time.Now().Add(time.Hour * 24 ).Unix(),
	})

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
    return tokenString
}