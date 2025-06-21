package migrations

import (
	"fmt"
	"log"

	"github.com/ziyadrw/faslah/config"
	podcast "github.com/ziyadrw/faslah/internal/modules/podcast/models"
	user "github.com/ziyadrw/faslah/internal/modules/user/models"
)

func Migrate() {
	db := config.GetDB()
	fmt.Println("Running Migrations...")
	err := db.AutoMigrate(
		&user.User{},
		&user.WatchHistory{},
		&podcast.Podcast{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations completed successfully!")
}
