package router

import (
	middleware "TA-management/internal/middlewares"
	coursecontroller "TA-management/internal/modules/course/controller"
	courseservice "TA-management/internal/modules/course/service"
	lookupcontroller "TA-management/internal/modules/lookup/controller"
	lookupservice "TA-management/internal/modules/lookup/service"
	studentcontroller "TA-management/internal/modules/student/controller"
	studentservice "TA-management/internal/modules/student/service"
	tadutycontroller "TA-management/internal/modules/ta_duty/controller"
	tadutyservice "TA-management/internal/modules/ta_duty/service"

	announcecontroller "TA-management/internal/modules/announce/controller"
	announceservice "TA-management/internal/modules/announce/service"

	authencontroller "TA-management/internal/modules/authen/controller"
	authenservice "TA-management/internal/modules/authen/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	clientID     string
	clientSecret string
	redirectURL  string
	jwtSecret    []byte
	cookieDomain string
)
var googleOAuthConfig *oauth2.Config

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// Allow all origins
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		// Allow specific methods
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		// Allow specific headers
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Bucket-Name, Folder-Path")
		// Allow credentials
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// func InitRouter() *gin.Engine {
// 	r := gin.Default()

// 	r.Use(corsMiddleware())

// 	db := config.ConnectDatabase()

// 	_ = godotenv.Load()

// 	clientID = os.Getenv("GOOGLE_CLIENT_ID")
// 	clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
// 	redirectURL = utils.GetenvDefault("GOOGLE_REDIRECT_URL", "http://localhost:8084/TA-management/auth/google/callback")
// 	jwtSecret = []byte(utils.GetenvDefault("JWT_SECRET", "change-me-please"))
// 	cookieDomain = os.Getenv("COOKIE_DOMAIN")

// 	// ====== OAuth2 config ======
// 	googleOAuthConfig = &oauth2.Config{
// 		ClientID:     clientID,
// 		ClientSecret: clientSecret,
// 		RedirectURL:  redirectURL,
// 		Scopes: []string{
// 			"openid",
// 			"email",
// 			"profile",
// 		},
// 		Endpoint: google.Endpoint,
// 	}

// 	courseRepo := courserepo.NewCourseRepository(db)
// 	courseService := courseservice.NewCourseService(courseRepo)

// 	authRepo := authenrepo.NewAuthenRepository(db)
// 	authService := authenservice.NewAuthenService(authRepo, googleOAuthConfig, jwtSecret)

// 	lookupRepo := lookuprepo.NewLookupRepository(db)
// 	lookupService := lookupservice.NewLookupService(lookupRepo)

// 	baseRouter := r.Group("/TA-management")

// 	authRouter := baseRouter.Group("/auth")
// 	{
// 		authencontroller.InitializeController(authService, googleOAuthConfig, authRouter)
// 	}

// 	publicRouter := baseRouter.Group("/public")
// 	{
// 		courseRouter := publicRouter.Group("/course")
// 		coursecontroller.InitializePublicController(courseService, courseRouter)
// 	}

// 	//authenticated routes Group
// 	authenticatedRouter := baseRouter.Group("")
// 	authenticatedRouter.Use(middleware.AuthMiddleware(jwtSecret))
// 	{
// 		courseRouter := authenticatedRouter.Group("/course")
// 		coursecontroller.InitializeController(courseService, courseRouter)

// 		lookupRouter := authenticatedRouter.Group("/lookup")
// 		lookupcontroller.InitializeController(lookupService, lookupRouter)
// 	}

// 	return r

// }
func InitRouter(
	authSvc authenservice.AuthenService,
	courseSvc courseservice.CourseService,
	lookupSvc lookupservice.LookupService,
	studentSvc studentservice.StudentService,
	tadutySvc tadutyservice.TaDutyService,
	announceSvc announceservice.AnnouncementService,
	oauthCfg *oauth2.Config,
	jwtSecret []byte,
) *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())

	baseRouter := r.Group("/TA-management")

	authRouter := baseRouter.Group("/auth")
	{
		authencontroller.InitializeController(authSvc, oauthCfg, authRouter)
	}

	publicRouter := baseRouter.Group("/public")
	{
		courseRouter := publicRouter.Group("/course")
		coursecontroller.InitializePublicController(courseSvc, courseRouter)
	}

	lookupPublicRouter := baseRouter.Group("/lookup")
	lookupcontroller.InitializePublicController(lookupSvc, lookupPublicRouter)

	authenticatedRouter := baseRouter.Group("")
	authenticatedRouter.Use(middleware.AuthMiddleware(jwtSecret))
	{
		courseRouter := authenticatedRouter.Group("/course")
		coursecontroller.InitializeController(courseSvc, courseRouter)

		lookupRouter := authenticatedRouter.Group("/lookup")
		lookupcontroller.InitializeProtectedController(lookupSvc, lookupRouter)

		studentRouter := authenticatedRouter.Group("/student")
		studentcontroller.InitializeController(studentSvc, studentRouter)

		tadutyRouter := authenticatedRouter.Group("/ta_duty")
		tadutycontroller.InitializeController(tadutySvc, tadutyRouter)

		announceRouter := authenticatedRouter.Group("/announce")
		announcecontroller.InitializeController(announceSvc, announceRouter)
	}

	return r
}
