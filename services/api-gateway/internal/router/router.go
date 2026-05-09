package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"

	"github.com/Kost0/internship-exchange/services/api-gateway/internal/clients"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/config"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/handler"
	"github.com/Kost0/internship-exchange/services/api-gateway/internal/middleware"
)

func New(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	authMW := middleware.NewAuthMiddleware(cfg.JWTSecret, rdb)

	authConn, err := clients.NewGRPCConn(cfg.AuthServiceAddr)
	if err != nil {
		panic(err)
	}

	profileConn, err := clients.NewGRPCConn(cfg.ProfileServiceAddr)
	if err != nil {
		panic(err)
	}

	listingConn, err := clients.NewGRPCConn(cfg.ListingServiceAddr)
	if err != nil {
		panic(err)
	}
	
	appConn, err := clients.NewGRPCConn(cfg.AppServiceAddr)
	if err != nil {
		panic(err)
	}

	authHandler := handler.NewAuthHandler(authConn)
	profileHandler := handler.NewProfileHandler(profileConn)
	listingHandler := handler.NewListingHandler(listingConn)
	appHandler := handler.NewApplicationHandler(appConn)

	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)
		})

		r.Get("/listings", listingHandler.GetListings)
		r.Get("/listings/{id}", listingHandler.GetListing)
		r.Get("/profile/company/{id}", profileHandler.GetCompanyProfile)
		r.Get("/companies/{id}", profileHandler.GetCompanyProfile)

		r.Group(func(r chi.Router) {
			r.Use(authMW.Authenticate)

			r.Route("/profile/student", func(r chi.Router) {
				r.Get("/", profileHandler.GetMyStudentProfile)
				r.Put("/", profileHandler.UpdateStudentProfile)
				r.Post("/avatar", profileHandler.UploadAvatar)
				r.Post("/resume", profileHandler.UploadResume)
				r.Post("/education", profileHandler.AddEducation)
				r.Put("/education/{id}", profileHandler.UpdateEducation)
				r.Delete("/education/{id}", profileHandler.DeleteEducation)
				r.Post("/experience", profileHandler.AddExperience)
				r.Put("/experience/{id}", profileHandler.UpdateExperience)
				r.Delete("/experience/{id}", profileHandler.DeleteExperience)
				r.Post("/projects", profileHandler.AddProject)
				r.Put("/projects/{id}", profileHandler.UpdateProject)
				r.Delete("/projects/{id}", profileHandler.DeleteProject)
				r.Get("/{id}", profileHandler.GetStudentProfile)
				r.Get("/{id}/resume", profileHandler.GetResumeURL)
			})

			r.Route("/profile/company", func(r chi.Router) {
				r.Get("/", profileHandler.GetMyCompanyProfile)
				r.Put("/", profileHandler.UpdateCompanyProfile)
				r.Post("/logo", profileHandler.UploadLogo)
			})

			r.Route("/listings", func(r chi.Router) {
				r.Get("/my", listingHandler.GetMyListings)
				r.Post("/", listingHandler.CreateListing)
				r.Put("/{id}", listingHandler.UpdateListing)
				r.Delete("/{id}", listingHandler.DeleteListing)
				r.Post("/{id}/publish", listingHandler.PublishListing)
				r.Post("/{id}/close", listingHandler.CloseListing)
				r.Get("/{id}/applications", appHandler.GetListingApplications)
			})

			r.Route("/applications", func(r chi.Router) {
				r.Post("/", appHandler.Apply)
				r.Get("/my", appHandler.GetMyApplications)
				r.Delete("/{id}", appHandler.Withdraw)
				r.Put("/{id}/status", appHandler.ChangeStatus)
				r.Get("/{id}/history", appHandler.GetHistory)
			})
		})
	})

	return r
}