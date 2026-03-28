package main

import (
	"TA-management/config"
	"TA-management/internal/logs"
	"TA-management/internal/modules/announce/discord"
	announcerepo "TA-management/internal/modules/announce/repository"
	announceservice "TA-management/internal/modules/announce/service"
	authenrepo "TA-management/internal/modules/authen/repository"
	authenservice "TA-management/internal/modules/authen/service"
	courserepo "TA-management/internal/modules/course/repository"
	courseservice "TA-management/internal/modules/course/service"
	lookuprepo "TA-management/internal/modules/lookup/repository"
	lookupservice "TA-management/internal/modules/lookup/service"
	studentrepo "TA-management/internal/modules/student/repository"
	studentservice "TA-management/internal/modules/student/service"
	tadutyrepo "TA-management/internal/modules/ta_duty/repository"
	tadutyservice "TA-management/internal/modules/ta_duty/service"
	router "TA-management/internal/routers"
	"TA-management/internal/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config

func main() {

	_ = godotenv.Load()
	log := logs.InitializeLogger()
	defer logs.SyncLogger(log)

	db := config.ConnectDatabase()
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	// calendarKey := os.Getenv("GOOGL_CALENDAR_KEY")
	BOTKey := os.Getenv("BOT_HOLIDAYS_KEY")
	BOTURL := os.Getenv("BOT_HOLIDAYS_URL")

	redisHost := utils.GetenvDefault("REDIS_HOST", "localhost")
	redisPort := utils.GetenvDefault("REDIS_PORT", "6379")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisClient := config.ConnectRedis(redisHost, redisPort, redisPassword)
	defer redisClient.Close()

	// ====== OAuth2 config ======
	googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  utils.GetenvDefault("GOOGLE_REDIRECT_URL", "http://localhost:8084/TA-management/auth/google/callback"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	// discord credential
	discordClientUrl := os.Getenv("DISCORD_CLIENT_URL")
	guilID := os.Getenv("GUILD_ID")

	//Initialize Repositories & Services
	authenRepo := authenrepo.NewAuthenRepository(db.DB)
	authenSvc := authenservice.NewAuthenService(authenRepo, googleOAuthConfig, jwtSecret)

	courseRepo := courserepo.NewCourseRepository(db.DB)
	courseSvc := courseservice.NewCourseService(courseRepo, redisClient)

	lookupRepo := lookuprepo.NewLookupRepository(db.DB)
	lookupSvc := lookupservice.NewLookupService(lookupRepo)

	studentRepo := studentrepo.NewStudentRepository(db)
	studentSvc := studentservice.NewStudentService(studentRepo)

	tadutyRepo := tadutyrepo.NewTaDutyRepository(db.DB)
	tadutySvc := tadutyservice.NewTaDutyServiceImplementation(tadutyRepo, log)

	discordClient := discord.NewDiscordClient(discordClientUrl, guilID)

	announceRepo := announcerepo.NewAnnouncementRepository(db.DB)
	announceSvc := announceservice.NewAnnouncementService(announceRepo, discordClient)

	//start CRONS JOB
	c := cron.New(cron.WithLocation(time.FixedZone("ICT", 7*3600)))
	c.AddFunc("@weekly", func() {
		lookupSvc.SyncOfficialHoliday(BOTKey, BOTURL)
	})
	c.AddFunc("0 0 * * *", func() {
		log.Info("🗂️ Running semester cleanup cron...")
		if err := courseSvc.SoftDeleteExpiredData(); err != nil {
			log.Info("❌ Semester cleanup failed: " + err.Error())
		} else {
			log.Info("✅ Semester cleanup completed successfully")
		}
	})
	c.Start()

	go func() {
		log.Info("🚀 Startup Sync: Initializing holiday data...")
		err := lookupSvc.SyncOfficialHoliday(BOTKey, BOTURL)
		if err != nil {
			log.Info("❌ Initial sync holidays failed: %v", err)
		} else {
			log.Info("✅ Initial sync holidays successful")
		}
	}()

	routes := router.InitRouter(authenSvc, courseSvc, lookupSvc, studentSvc, tadutySvc, announceSvc, googleOAuthConfig, jwtSecret)

	port := 8084
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: routes,
	}

	fmt.Printf("🚀 TA-Management Server started on :%d\n", port)
	log.Fatal(server.ListenAndServe())
}
