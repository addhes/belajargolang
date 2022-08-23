package handler

import (
	"fmt"
	"net/http"
	"tesa/auth"
	"tesa/helper"
	"tesa/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {

	// step 1: get input from request user
	// step 2: map input dari user ke struct RegisterUserInput
	// step 3: Struct di atas kita passing sebagai parameter service

	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {

		errors := helper.FormatError(err)
		erroMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", erroMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newuser, err := h.userService.RegisterUser(input)
	if err != nil {
		errors := helper.FormatError(err)
		erroMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", erroMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newuser.ID)
	if err != nil {
		errors := helper.FormatError(err)
		erroMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", erroMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newuser, token)

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		erroMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", erroMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinuser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedinuser.ID)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := user.FormatUser(loggedinuser, token)
	response := helper.APIResponse("Login success", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

	// step 1: user memasukan input email dan password
	// step 2: input ditangkap handler
	// step 3: mapping dari input user ke input struct
	// step 4: input struct passing service
	// step 5: di service mencari dg bantuan repository user dengan email
	// step 6: mencocokan password
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		erroMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Check email failed", http.StatusUnprocessableEntity, "error", erroMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	IsEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Check email failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{"is_available": IsEmailAvailable}

	metaMessage := "Email is not available"

	if IsEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

	// step 1: ada input email dari user
	// step 2: input email dimapping ke struct input
	// step 3: struct input dipassing ke service
	// step 4: service akan memanggil repository email untuk mencari sudah ada atau belum
	// step 5: jika sudah ada, maka responsenya false
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	// c.SaveUploadedFile(file, "./uploads/avatar.png")
	// step 1: get input dari user
	// step 2: Simpan gambarnya di folder "images/"
	// step 3: di service kita panggil repo
	// step 4: jwt (sementara harcode, seakan2 user yang login id=1)
	// step 5: repo ambil user dengan id=1
	// step 6: repo update data user simpan lokasi file

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaed": false}
		response := helper.APIResponse("Upload avatar failed", http.StatusUnprocessableEntity, "error", data)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	// path := "images/" + file.Filename
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaed": false}
		response := helper.APIResponse("Upload avatar failed", http.StatusUnprocessableEntity, "error", data)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaed": false}
		response := helper.APIResponse("Upload avatar failed", http.StatusUnprocessableEntity, "error", data)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{"is_uploaed": true}
	response := helper.APIResponse("Upload avatar success", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}
