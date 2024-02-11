package main

import (
	"flag"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jittakal/todo-api-hexagonal/api/proto"
	"github.com/jittakal/todo-api-hexagonal/docs"
	rpc "github.com/jittakal/todo-api-hexagonal/internal/adapters/handler/grpc"
	"github.com/jittakal/todo-api-hexagonal/internal/adapters/handler/rest"
	"github.com/jittakal/todo-api-hexagonal/internal/adapters/repository/inmemory"
	"github.com/jittakal/todo-api-hexagonal/internal/domain"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Define command-line flags for handler and repository types
	handlerType := flag.String("handler", "all", "Type of handler (http/grpc/all)")
	repoType := flag.String("repository", "in-memory", "Type of repository (in-memory)")

	flag.Parse()

	// Set up Logrus logger
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// Initialize repositories and services based on flags
	var taskRepository domain.TaskRepository
	switch *repoType {
	case "in-memory":
		taskRepository = inmemory.NewTaskRepository()
	default:
		logrus.Info("Unsupported repository type: ", *repoType)
		return
	}

	// Initialize handler based on the selected type
	switch *handlerType {
	case "http":
		startHTTPServer(taskRepository)
	case "grpc":
		startGRPCServer(taskRepository)
	case "all":
		go startHTTPServer(taskRepository)
		startGRPCServer(taskRepository)
	default:
		logrus.Info("Unsupported handler type:", *handlerType)
		return
	}
}

func startHTTPServer(taskRepository domain.TaskRepository) {

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Todo API"
	docs.SwaggerInfo.Description = "API documentation for the Todo API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	// Create a new Gin router
	router := gin.Default()

	// use ginSwagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define health check route
	router.GET("/todo/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Create a new REST handler
	taskHandler := rest.NewGinTaskHandler(taskRepository)

	// Group routes for tasks
	tasksGroup := router.Group("/todo/v1/tasks")
	{
		tasksGroup.POST("", taskHandler.CreateTask)
		tasksGroup.GET("", taskHandler.GetAllTasks)
		tasksGroup.GET("/:id", taskHandler.GetTask)
		tasksGroup.PUT("/:id", taskHandler.UpdateTask)
		tasksGroup.DELETE("/:id", taskHandler.DeleteTask)
		tasksGroup.POST("/:id/done", taskHandler.MarkTaskDone)
	}

	// Create an instance of the HTTP server
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the HTTP server and handle errors
	logrus.Info("HTTP server listening on port 8080")
	err := httpServer.ListenAndServe()
	if err != nil {
		logrus.Errorf("Error starting HTTP server: %v", err)
	}
}

func startGRPCServer(taskRepository domain.TaskRepository) {
	grpcServer := grpc.NewServer()

	// Create a new gRPC handler
	grpcHandler := rpc.NewGRPCTaskHandler(taskRepository)
	proto.RegisterTaskServiceServer(grpcServer, grpcHandler)

	// Start the gRPC server
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logrus.Fatalf("Failed to listen: %v", err)
	}
	logrus.Info("gRPC server listening on :8081")
	err = grpcServer.Serve(lis)
	if err != nil {
		logrus.Errorf("Error starting gRPC server: %v", err)
	}
}
