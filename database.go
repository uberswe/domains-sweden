package baseproject

import (
	"fmt"
	"github.com/uberswe/domains-sweden/config"
	"github.com/uberswe/domains-sweden/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

func connectToDatabase(c config.Config) (db *gorm.DB, err error) {
	return connectLoop(c, 0)
}

func connectLoop(c config.Config, count int) (db *gorm.DB, err error) {
	db, err = attemptConnection(c)
	if err != nil {
		if count > 300 {
			return db, fmt.Errorf("could not connect to database after 300 seconds")
		}
		time.Sleep(1 * time.Second)
		return connectLoop(c, count+1)
	}
	return db, err
}

func attemptConnection(c config.Config) (db *gorm.DB, err error) {
	if c.Database == "sqlite" {
		// In-memory sqlite if no database name is specified
		dsn := "file::memory:?cache=shared"
		if c.DatabaseName != "" {
			dsn = fmt.Sprintf("%s.db", c.DatabaseName)
		}
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	} else if c.Database == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.DatabaseUsername, c.DatabasePassword, c.DatabaseHost, c.DatabasePort, c.DatabaseName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if c.Database == "postgres" {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", c.DatabaseHost, c.DatabaseUsername, c.DatabasePassword, c.DatabaseName, c.DatabasePort)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	} else {
		return db, fmt.Errorf("no database specified: %s", c.Database)
	}
	return db, err
}

func migrateDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.Session{},
		&models.Domain{},
		&models.Nameserver{},
		&models.Release{},
		&models.Fetch{},
		&models.NameserverAggregate{},
		&models.Sitemap{},
		&models.Parse{},
		&models.ParseEvent{},
		&models.Migration{},
	)

	parseDropContentScreenshotBlurredScreenshot(db)
	domainChangeLongtextToText(db)

	return err
}

func parseDropContentScreenshotBlurredScreenshot(db *gorm.DB) {
	m := models.Migration{
		Key: "01_parse_drop_content_screenshot_blurred_screenshot",
	}
	res := db.Where(m).First(&m)
	if res.Error == gorm.ErrRecordNotFound {
		res = db.Save(&m)
		if res.Error != nil {
			log.Println(res.Error)
			return
		}
	} else {
		return
	}

	err := db.Migrator().DropColumn(&models.Parse{}, "content")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.Migrator().DropColumn(&models.Parse{}, "screenshot")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.Migrator().DropColumn(&models.Parse{}, "blurred_screenshot")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Migrated", m.Key)
}

func domainChangeLongtextToText(db *gorm.DB) {
	m := models.Migration{
		Key: "02_domain_change_longtext_to_text",
	}
	res := db.Where(m).First(&m)
	if res.Error == gorm.ErrRecordNotFound {
		res = db.Save(&m)
		if res.Error != nil {
			log.Println(res.Error)
			return
		}
	} else {
		return
	}

	err := db.Migrator().AlterColumn(&models.Domain{}, "Host")
	if err != nil {
		return
	}

	err = db.Migrator().CreateIndex(&models.Domain{}, "Host")
	if err != nil {
		return
	}

	log.Println("Migrated", m.Key)
}
