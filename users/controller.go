package users

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	userJSONDecoder := json.NewDecoder(r.Body)
	userJSONDecoder.Decode(&newUser)

	fmt.Println(newUser)
}
