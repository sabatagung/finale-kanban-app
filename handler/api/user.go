package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserAPI interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("email or password is empty"))
		return
	}
	var userLogin = entity.User{
		Email:    user.Email,
		Password: user.Password,
	}
	id, err := u.userService.Login(r.Context(), &userLogin)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	cookie := &http.Cookie{
		Name:    "user_id",
		Value:   fmt.Sprintf("%d", id),
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 5),
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": id,
		"message": "login success",
	})
	return
}

func (u *userAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}
	// check request data
	if user.Fullname == "" || user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("register data is empty"))
		return
	}

	var userRegister = entity.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}
	userRegister, err = u.userService.Register(r.Context(), &userRegister)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userRegister.ID,
		"message": "register success",
	})
	return

}

func (u *userAPI) Logout(w http.ResponseWriter, r *http.Request) {
	var sessions = map[string]entity.User{}
	token, _ := r.Cookie("user_id")

	sessionToken := token.Value

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "logout success",
	})
	return

}

func (u *userAPI) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("user_id is empty"))
		return
	}

	deleteUserId, _ := strconv.Atoi(userId)

	err := u.userService.Delete(r.Context(), int(deleteUserId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "delete success"})
}
