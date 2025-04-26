package routers

import (
	"context"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glebkasilov/authorization/internal/domain/models"
	"github.com/glebkasilov/authorization/internal/domain/requests"
	"github.com/glebkasilov/authorization/internal/domain/responces"
	"github.com/glebkasilov/authorization/pkg/middlewares/authorization"
)

type Service interface {
	Login(ctx context.Context, user requests.Login) (string, error)
	Register(ctx context.Context, user requests.Register) error
	SetRole(ctx context.Context, id string, role string) error
	GetUser(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthorizationRouter struct {
	router  *gin.RouterGroup
	service Service
}

func Register(routerGroup *gin.RouterGroup, service Service) {
	router := AuthorizationRouter{
		router:  routerGroup,
		service: service,
	}
	router.init()
}

func (r *AuthorizationRouter) register(ctx *gin.Context) {
	var user requests.Register
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.Password != user.RepeatPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "passwords do not match",
		})
		return
	}

	if userFound, _ := r.service.GetUserByEmail(ctx, user.Email); userFound != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "user with this email already exists",
		})
		return
	}

	if err := r.service.Register(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "user registered",
	})
}

func (r *AuthorizationRouter) login(ctx *gin.Context) {
	var user requests.Login
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	token, err := r.service.Login(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (r *AuthorizationRouter) setRole(ctx *gin.Context) {
	var user requests.SetRole
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	payload, err := authorization.FromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if payload.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to get this user",
		})
		return
	}

	if err := r.service.SetRole(ctx, user.ID, user.Role); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "role set",
	})
}

func (r *AuthorizationRouter) getUser(ctx *gin.Context) {
	id := ctx.Param("id")

	payload, err := authorization.FromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if payload.Role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to get this user",
		})
		return
	}

	user, err := r.service.GetUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	userResponse := responces.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Fullname:     user.Fullname,
		Role:         user.Role,
		CurrentStage: user.CurrentStage,
		Points:       user.Points,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": userResponse,
	})
}

func (r *AuthorizationRouter) me(ctx *gin.Context) {
	payload, err := authorization.FromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := r.service.GetUser(ctx, payload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	userResponse := responces.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		Fullname:     user.Fullname,
		Role:         user.Role,
		CurrentStage: user.CurrentStage,
		Points:       user.Points,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": userResponse,
	})
}

func (r *AuthorizationRouter) init() {
	key := GetKey("./keys/public_key.pem")
	group := r.router.Group("/authorization")

	group.POST("/login", r.login)
	group.POST("/register", r.register)
	group.POST("/set-role", authorization.MiddlwareJWT(key), r.setRole)
	group.GET("/:id", authorization.MiddlwareJWT(key), r.getUser)
	group.GET("/me", authorization.MiddlwareJWT(key), r.me)
}
