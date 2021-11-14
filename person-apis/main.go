package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	loggermiddleware "github.com/meateam/api-gateway/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	pb "github.com/sternomer/person-api/proto"
	"github.com/sternomer/person-api/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	ParamId        = "id"
	ParamBirthDate = "birthdate"
	ParamFirstName = "firstname"
	ParamLastName  = "lastname"
)

type PersonItem struct {
	Id        string `json:"id"`
	Birthdate string `json:"birthdate"`
	FirstName string `json:"firstName"`
	Lastname  string `json:"lastname"`
}
type Router struct {
	client pb.PersonServiceClient
	pb.UnimplementedPersonServiceServer
	logger *logrus.Logger
}

func initClientConnection() pb.PersonServiceClient {

	address := viper.GetString(utils.PersonServiceAddress)
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("failed to get mongo connection parameters")
	}

	client := pb.NewPersonServiceClient(conn)

	return client
}

func corsRouterConfig() cors.Config {

	corsConfig := cors.DefaultConfig()
	corsConfig.AddExposeHeaders("x-uploadid")
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowWildcard = true
	corsConfig.AllowOrigins = strings.Split("http://localhost*", ",")
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders(
		"x-content-length",
		"authorization",
		"cache-control",
		"x-requested-with",
		"content-disposition",
		"content-range",
		"destination",
		"fileID",
	)

	return corsConfig
}
func (r *Router) CreatePerson(c *gin.Context) {
	var PersonFilter PersonItem
	if bindErr := c.Bind(&PersonFilter); bindErr != nil {
		c.String(http.StatusBadRequest, "Create person method is failed", bindErr)
		return
	}
	request := &pb.Person{
		Id:        PersonFilter.Id,
		Birthdate: PersonFilter.Birthdate,
		FirstName: PersonFilter.FirstName,
		LastName:  PersonFilter.Lastname,
	}
	res, err := r.client.CreatePerson(c, request)
	if err != nil {
		httpStatusCode := gwruntime.HTTPStatusFromCode(status.Code(err))
		loggermiddleware.LogError(r.logger, c.AbortWithError(httpStatusCode, err))
		return
	}
	c.JSON(http.StatusOK, res)
}
func (r *Router) ReadPerson(c *gin.Context) {
	request := &pb.ReadPersonReq{Id: c.Param(ParamId)}
	res, err := r.client.ReadPerson(c, request)
	if err != nil {
		httpStatusCode := gwruntime.HTTPStatusFromCode(status.Code(err))
		loggermiddleware.LogError(r.logger, c.AbortWithError(httpStatusCode, err))
		return
	}
	c.JSON(http.StatusOK, res)

}
func (r *Router) ListPersons(c *gin.Context) {
	request := &pb.ListPersonsReq{}
	res, err := r.client.ListPersons(c, request)
	if err != nil {
		httpStatusCode := gwruntime.HTTPStatusFromCode(status.Code(err))
		loggermiddleware.LogError(r.logger, c.AbortWithError(httpStatusCode, err))
		return
	}
	c.JSON(http.StatusOK, res)

}
func (r *Router) DeletePerson(c *gin.Context) {
	request := &pb.DeletePersonReq{Id: c.Param(ParamId)}
	res, err := r.client.DeletePerson(c, request)
	if err != nil {
		httpStatusCode := gwruntime.HTTPStatusFromCode(status.Code(err))
		loggermiddleware.LogError(r.logger, c.AbortWithError(httpStatusCode, err))
		return
	}
	c.JSON(http.StatusNoContent, res)
}
func (r *Router) UpdatePerson(c *gin.Context) {
	request := &PersonItem{}
	c.BindJSON(request)

	res, err := r.client.UpdatePerson(c, &pb.Person{Id: request.Id,Birthdate: request.Birthdate,FirstName: request.FirstName, LastName: request.Lastname})
	if err != nil {
		httpStatusCode := gwruntime.HTTPStatusFromCode(status.Code(err))
		loggermiddleware.LogError(r.logger, c.AbortWithError(httpStatusCode, err))
		return
	}
	c.JSON(http.StatusOK, res)
}
func main() {

	// Loading dotenv file parameters
	err := utils.LoadConfig()
	if err != nil {
		fmt.Println("cannot load config:", err)
		return
	}
	routerPort := viper.GetString(utils.GrpcRouterPort)

	r := &Router{}
	r.client = initClientConnection()

	mainRouter := gin.Default()
	mainRouter.Use(
		cors.New(corsRouterConfig()),
	)

	mainRouter.POST("/api/person", r.CreatePerson)
	mainRouter.GET("/api/person/:id", r.ReadPerson)
	mainRouter.GET("/api/persons", r.ListPersons)
	mainRouter.DELETE("/api/person/:id", r.DeletePerson)
	mainRouter.PUT("/api/person", r.UpdatePerson)

	err = mainRouter.Run(":" + routerPort)
	if err != nil {
		fmt.Println("failed to run api gateway. \nrouter error: ", err)
		return
	}
}
