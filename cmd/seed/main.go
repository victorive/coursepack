package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/victorive/klaro-backend/database/seeders"
	"github.com/victorive/klaro-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	all := flag.Bool("all", false, "Run all seeders in registration order")
	name := flag.String("seeder", "", "Run a specific seeder by name")
	flag.Parse()

	if !*all && *name == "" {
		fmt.Fprintln(os.Stderr, "usage: seed -all | -seeder=<name>")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	if *all {
		for _, s := range seeders.All() {
			runOne(db, s)
		}
		return
	}

	s := seeders.Find(*name)
	if s == nil {
		log.Fatalf("seeder %q not found", *name)
	}
	runOne(db, s)
}

func runOne(db *gorm.DB, s seeders.Seeder) {
	fmt.Printf("seeder %s ... ", s.Name())
	if err := s.Run(db); err != nil {
		fmt.Printf("FAILED: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("OK")
}
