package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rshlin/go-blog-api-assesment/api"
	"github.com/rshlin/go-blog-api-assesment/blog/repository"
	"github.com/rshlin/go-blog-api-assesment/blog/service"
	"github.com/rshlin/go-blog-api-assesment/server"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var technicalTestCmd = &cobra.Command{
	Use:   "technical",
	Short: "Technical Test: RESTful API for a blog platform",
	Long:  `This command starts a server for a RESTful API for a blog platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetInt("port")
		appConfigPath, _ := cmd.Flags().GetString("config")

		initialPostsPath, _ := cmd.Flags().GetString("initial-posts-path")

		serverConfig := server.LoadConfig(appConfigPath)

		postsCfg := repository.LoadInMemoryConfig(initialPostsPath)

		fmt.Printf("Starting server at %s:%d\n", address, port)

		repo := repository.NewInMemoryBlogRepository(
			repository.WithConfig(postsCfg),
		)
		svc := service.NewSimpleBlogService(repo)
		srv := server.NewServer(svc, serverConfig)

		authStore, err := server.NewAuthStore(serverConfig.AuthStore)
		if err != nil {
			log.Fatal("Failed to initialize auth store", err)
		}

		authenticator, err := server.NewAuthenticator(serverConfig.Authenticator, authStore)
		if err != nil {
			log.Fatal("Failed to initialize authenticator", err)
		}

		middleware := server.CreateMiddleware(authenticator, serverConfig)

		router := mux.NewRouter()

		handler := api.HandlerFromMux(srv, router)

		s := &http.Server{
			Handler: handler,
			Addr:    fmt.Sprintf("%s:%d", address, port),
		}

		router.Use(middleware...)

		log.Fatal(s.ListenAndServe())
	},
}

var address string
var port int

var appConfigPath string
var initialBlogPostsJsonPath string

func init() {
	rootCmd.AddCommand(technicalTestCmd)

	// server args
	technicalTestCmd.Flags().StringVarP(&address, "address", "a", "127.0.0.1", "bind address for the server")
	technicalTestCmd.Flags().IntVarP(&port, "port", "p", 8080, "port for the server")
	technicalTestCmd.Flags().StringVarP(&appConfigPath, "config", "c", "app.yaml", "path to server configuration file")

	// blog args
	technicalTestCmd.Flags().StringVar(&initialBlogPostsJsonPath, "initial-posts-path", "", "path to the initial blog posts")
}
