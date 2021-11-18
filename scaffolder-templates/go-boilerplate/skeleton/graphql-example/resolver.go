package graphql_example

import (
	"context"
	"errors"
	"log"
	"time"

	"git.xenonstack.com/lib/golang-boilerplate/auth"
	cauth "git.xenonstack.com/util/golang-boilerplate/graphql-example/src/auth"

	"git.xenonstack.com/util/golang-boilerplate/graphql-example/src/methods"
	"git.xenonstack.com/util/golang-boilerplate/graphql-example/src/models"

	"github.com/99designs/gqlgen/graphql"
	"github.com/jinzhu/gorm"
	ot "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

// Resolver is a basic resolver structure can be modified according to service needs
// in this db client is saved
type Resolver struct {
	DB *gorm.DB
}

// NewRootResolvers is a method to intialise resolvers configuration
func NewRootResolvers(db *gorm.DB) Config {
	// save db client
	c := Config{
		Resolvers: &Resolver{
			DB: db,
		},
	}

	//schema directive
	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		log.Println(ctx)
		token := ctx.Value(cauth.CtxKey{})
		log.Println("=-=-=-=-=-=", token)
		if token != nil {
			// validate the Token
			mapd, err := auth.ValidateTokenString(token.(string))
			if err != nil {
				zap.S().Error("Token err....", err)
				return nil, err
			}
			//add claims to context
			for key, value := range mapd {
				ctx = context.WithValue(ctx, key, value)
			}
			return next(ctx)
		} else {
			return nil, errors.New("You are not authorised to perform this action")
		}
	}
	return c
}

// Mutation is a method to execute graphql mutations
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query is a method to execute graphql queries
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) SignupUser(ctx context.Context, input NewUser) (*models.UserDetail, error) {
	//start span
	span, _ := ot.StartSpanFromContext(ctx, "signup user")
	//Finish span after all operations
	defer span.Finish()
	span.SetTag("event", "signup user in resolver")

	span.LogKV("task", "validating password")
	flag := methods.ValidatePassword(input.Password)
	zap.S().Info(flag)
	if flag == 1 {
		return &models.UserDetail{}, errors.New("Password must have minimum eight characters, at least one uppercase letter, at least one lowercase letter, at least one number and at least one special character")
	}
	span.LogKV("task", "hashing password")
	pass := methods.HashForNewPassword(input.Password)
	span.LogKV("task", "create new user and save in database")
	user := models.UserDetail{
		Name:      *input.Name,
		Email:     input.Email,
		Password:  pass,
		Role:      "user",
		CreatedAt: time.Now(),
	}
	err := r.Resolver.DB.Create(&user).Error
	if err != nil {
		log.Println("err...", err)
		return &models.UserDetail{}, err
	}
	span.LogKV("task", "generate jwt token")
	mapd := make(map[string]interface{})
	mapd["id"] = user.ID
	mapd["email"] = user.Email
	token, err := auth.NewToken(mapd)
	if err != nil {
		zap.S().Error("===--=-=-=", err)
		return &models.UserDetail{}, err
	}
	user.Token = token
	zap.S().Info(token)
	span.LogKV("task", "send final output")
	return &user, nil
}
func (r *mutationResolver) AddAddress(ctx context.Context, input NewAddress) (*models.UserDetail, error) {
	//start span
	span, _ := ot.StartSpanFromContext(ctx, "get user info")
	//Finish span after all operations
	defer span.Finish()
	span.SetTag("event", "add address to db")
	zap.S().Info(ctx.Value("id"), "====", ctx.Value("email"))

	span.LogKV("task", "check user exists")
	//data type conversion of id
	id := ctx.Value("id").(float64)
	log.Println(id)
	//fetch users from db
	user := []models.UserDetail{}
	r.Resolver.DB.Where("id=?", id).Find(&user)
	if len(user) == 0 {
		return &models.UserDetail{}, errors.New("No user exists with this id")
	} else if len(user) > 1 {

	}
	span.LogKV("task", "add address to db")
	address := models.Address{
		UserID:  int(id),
		Country: input.Country,
		State:   input.State,
		Zip:     input.Zip,
	}
	err := r.Resolver.DB.Create(&address).Error
	if err != nil {
		zap.S().Error(err)
		return &models.UserDetail{}, err
	}

	user[0].Address = address
	span.LogKV("task", "send final outpout")
	return &user[0], nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) UserInfo(ctx context.Context, id *int) ([]*models.UserDetail, error) {
	//start span
	span, _ := ot.StartSpanFromContext(ctx, "get user info")
	//Finish span after all operations
	defer span.Finish()
	span.SetTag("event", "fetch user information from db")
	span.LogKV("task", "fetching information")
	// initialize users slice
	users := make([]*models.UserDetail, 0)
	//fetch user details from db
	if id == nil {
		r.Resolver.DB.Find(&users)
		return users, nil
	}
	r.Resolver.DB.Where("id=?", *id).Find(&users)
	span.LogKV("task", "send final output")
	for i := 0; i < len(users); i++ {
		address := make([]models.Address, 0)

		r.Resolver.DB.Where("user_id=?", users[i].ID).Find(&address)
		log.Println(address)
		if len(address) != 0 {
			users[i].Address = address[0]
		}
	}
	return users, nil
}
