package main

import (
	"fmt"
	"os"

	"server/common/models"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	// 1. INJECT POSTGIS EXTENSION HERE
	fmt.Println(`CREATE EXTENSION IF NOT EXISTS "postgis";`)

	// 2. Load your GORM models
	stmts, err := gormschema.New("postgres").Load(
		&models.User{},
		&models.Team{},
		&models.Challenge{},
		&models.Sport{},
		&models.Invitation{},
		&models.Location{},
		&models.Notification{},
		&models.UserSettings{},
		&models.Message{},
		&models.Chat{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	// 3. Print the rest of the schema
	fmt.Print(stmts)
}
