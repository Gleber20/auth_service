package http

import (
	"auth_service/internal/domain"
	"auth_service/internal/errs"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SignUpRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignUp
// @Summary Регистрация
// @Description Создать новый аккаунт
// @Tags Auth
// @Consume json
// @Produce json
// @Param request_body body SignUpRequest true "информация о новом аккаунте"
// @Success 201 {object} CommonResponse
// @Failure 422 {object} CommonError
// @Failure 400 {object} CommonError
// @Failure 404 {object} CommonError
// @Failure 500 {object} CommonError
// @Router /auth/sign-up [post]
func (s *Server) SignUp(c *gin.Context) {
	var input SignUpRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		s.handleError(c, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	if err := s.uc.UserCreator.CreateUser(c, domain.User{
		FullName: input.FullName,
		Username: input.Username,
		Password: input.Password,
	}); err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, CommonResponse{Message: "User created successfully!"})
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// SignIn
// @Summary Вход
// @Description Войти в аккаунт
// @Tags Auth
// @Consume json
// @Produce json
// @Param request_body body SignInRequest true "логин и пароль"
// @Success 200 {object} TokenPairResponse
// @Failure 400 {object} CommonError
// @Failure 404 {object} CommonError
// @Failure 500 {object} CommonError
// @Router /auth/sign-in [post]
func (s *Server) SignIn(c *gin.Context) {
	var input SignInRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		s.handleError(c, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	userID, err := s.uc.Authenticator.Authenticate(c, domain.User{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		s.handleError(c, err)
		return
	}

	accessToken, refreshToken, err := s.generateNewTokenPair(userID)
	if err != nil {
		s.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
