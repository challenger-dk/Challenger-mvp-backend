package seed

import (
	"fmt"
	"log"
	"time"

	"server/common/config"
	"server/common/models"
	"server/common/models/types"
	"server/common/services"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedDatabase seeds the database with test data
func SeedDatabase() error {
	log.Println("🌱 Starting database seeding...")

	// Clear existing data (optional - comment out if you want to keep existing data)
	if err := clearDatabase(); err != nil {
		return fmt.Errorf("failed to clear database: %w", err)
	}

	// Seed sports (if not already seeded)
	if err := config.SeedSports(); err != nil {
		return fmt.Errorf("failed to seed sports: %w", err)
	}

	// Get sports for associations
	var sports []models.Sport
	if err := config.DB.Find(&sports).Error; err != nil {
		return fmt.Errorf("failed to fetch sports: %w", err)
	}

	// Seed users
	users, err := seedUsers()
	if err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// Seed locations
	locations, err := seedLocations()
	if err != nil {
		return fmt.Errorf("failed to seed locations: %w", err)
	}

	// Seed teams
	teams, err := seedTeams(users, sports, locations)
	if err != nil {
		return fmt.Errorf("failed to seed teams: %w", err)
	}

	// Seed challenges
	challenges, err := seedChallenges(users, teams, locations, sports)
	if err != nil {
		return fmt.Errorf("failed to seed challenges: %w", err)
	}

	// Seed friendships
	if err := seedFriendships(users); err != nil {
		return fmt.Errorf("failed to seed friendships: %w", err)
	}

	// Seed invitations
	if err := seedInvitations(users, teams, challenges); err != nil {
		return fmt.Errorf("failed to seed invitations: %w", err)
	}

	// Seed conversations
	if err := seedConversations(users, teams); err != nil {
		return fmt.Errorf("failed to seed conversations: %w", err)
	}

	log.Println("✅ Database seeding completed successfully!")
	return nil
}

func clearDatabase() error {
	log.Println("🧹 Clearing existing data...")

	// Delete in reverse order of dependencies
	if err := config.DB.Exec("DELETE FROM messages").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM conversation_participants").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM conversations").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM user_challenges").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM challenge_teams").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM user_friends").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM user_teams").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM user_favorite_sports").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM team_sports").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM invitations").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM notifications").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM challenges").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM teams").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM locations").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM user_settings").Error; err != nil {
		return err
	}
	if err := config.DB.Exec("DELETE FROM users").Error; err != nil {
		return err
	}

	log.Println("✅ Database cleared")
	return nil
}

func seedUsers() ([]models.User, error) {
	log.Println("👥 Seeding users...")

	// Default password for all test users: "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password12"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedPasswordStr := string(hashedPassword)

	users := []models.User{
		{
			Email:     "user1@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Alice",
			LastName:  "Johnson",
			Bio:       "Tennisentusiast og weekendkriger",
			BirthDate: time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user2@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Bob",
			LastName:  "Smith",
			Bio:       "Fodboldspiller på jagt efter konkurrencedygtige kampe",
			BirthDate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user3@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Charlie",
			LastName:  "Brown",
			Bio:       "Basketballfanatiker",
			BirthDate: time.Date(1992, 1, 1, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user4@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Diana",
			LastName:  "Williams",
			Bio:       "Elsker at spille padel tennis og volleyball",
			BirthDate: time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user5@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Eve",
			LastName:  "Davis",
			Bio:       "Løbe- og cykelentusiast",
			BirthDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user6@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Frank",
			LastName:  "Hansen",
			Bio:       "Badminton spiller og tennisfan",
			BirthDate: time.Date(1998, 3, 15, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user7@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Grace",
			LastName:  "Nielsen",
			Bio:       "Svømning og vandpolo elsker",
			BirthDate: time.Date(1994, 7, 22, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user8@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Henry",
			LastName:  "Andersen",
			Bio:       "Fodbold og basketball på weekenderne",
			BirthDate: time.Date(1996, 11, 8, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user9@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Iris",
			LastName:  "Petersen",
			Bio:       "Yoga og pilates instruktør",
			BirthDate: time.Date(1991, 5, 30, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user10@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Jack",
			LastName:  "Larsen",
			Bio:       "Håndbold spiller og fitness entusiast",
			BirthDate: time.Date(1999, 9, 12, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user11@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Karen",
			LastName:  "Christensen",
			Bio:       "Tennis og padel tennis spiller",
			BirthDate: time.Date(1993, 2, 18, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user12@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Lars",
			LastName:  "Sørensen",
			Bio:       "Fodbold træner og løbeentusiast",
			BirthDate: time.Date(1988, 12, 5, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user13@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Maria",
			LastName:  "Rasmussen",
			Bio:       "Basketball og volleyball spiller",
			BirthDate: time.Date(1997, 4, 25, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user14@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Niels",
			LastName:  "Jørgensen",
			Bio:       "Cykling og mountainbike kører",
			BirthDate: time.Date(2001, 8, 14, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user15@challenger.dk",
			Password:  &hashedPasswordStr,
			FirstName: "Olivia",
			LastName:  "Madsen",
			Bio:       "Tennis, badminton og squash spiller",
			BirthDate: time.Date(1995, 6, 20, 0, 0, 0, 0, time.UTC),
			Settings:  &models.UserSettings{},
		},
	}

	// Create users
	for i := range users {
		if err := config.DB.Create(&users[i]).Error; err != nil {
			return nil, fmt.Errorf("failed to create user %s: %w", users[i].Email, err)
		}
	}

	// Associate favorite sports
	tennisSport := findSportByName("Tennis", config.DB)
	footballSport := findSportByName("Football", config.DB)
	basketballSport := findSportByName("Basketball", config.DB)
	padelSport := findSportByName("PadelTennis", config.DB)
	volleyballSport := findSportByName("Volleyball", config.DB)
	runningSport := findSportByName("Running", config.DB)
	bikingSport := findSportByName("Biking", config.DB)
	badmintonSport := findSportByName("Badminton", config.DB)
	swimmingSport := findSportByName("Swimming", config.DB)
	handballSport := findSportByName("Handball", config.DB)
	squashSport := findSportByName("Squash", config.DB)

	// Original users (0-4)
	if tennisSport != nil {
		config.DB.Model(&users[0]).Association("FavoriteSports").Append(tennisSport)
		config.DB.Model(&users[3]).Association("FavoriteSports").Append(tennisSport, padelSport)
	}
	if footballSport != nil {
		config.DB.Model(&users[1]).Association("FavoriteSports").Append(footballSport)
	}
	if basketballSport != nil {
		config.DB.Model(&users[2]).Association("FavoriteSports").Append(basketballSport)
	}
	if volleyballSport != nil {
		config.DB.Model(&users[3]).Association("FavoriteSports").Append(volleyballSport)
	}
	if runningSport != nil && bikingSport != nil {
		config.DB.Model(&users[4]).Association("FavoriteSports").Append(runningSport, bikingSport)
	}

	// New users (5-14)
	if len(users) > 5 && badmintonSport != nil && tennisSport != nil {
		config.DB.Model(&users[5]).Association("FavoriteSports").Append(badmintonSport, tennisSport)
	}
	if len(users) > 6 && swimmingSport != nil {
		config.DB.Model(&users[6]).Association("FavoriteSports").Append(swimmingSport)
	}
	if len(users) > 7 && footballSport != nil && basketballSport != nil {
		config.DB.Model(&users[7]).Association("FavoriteSports").Append(footballSport, basketballSport)
	}
	if len(users) > 8 {
		// User 8 (Iris) - Yoga/Pilates (no specific sports in DB, skip)
	}
	if len(users) > 9 && handballSport != nil {
		config.DB.Model(&users[9]).Association("FavoriteSports").Append(handballSport)
	}
	if len(users) > 10 && tennisSport != nil && padelSport != nil {
		config.DB.Model(&users[10]).Association("FavoriteSports").Append(tennisSport, padelSport)
	}
	if len(users) > 11 && footballSport != nil && runningSport != nil {
		config.DB.Model(&users[11]).Association("FavoriteSports").Append(footballSport, runningSport)
	}
	if len(users) > 12 && basketballSport != nil && volleyballSport != nil {
		config.DB.Model(&users[12]).Association("FavoriteSports").Append(basketballSport, volleyballSport)
	}
	if len(users) > 13 && bikingSport != nil {
		config.DB.Model(&users[13]).Association("FavoriteSports").Append(bikingSport)
	}
	if len(users) > 14 && tennisSport != nil && badmintonSport != nil && squashSport != nil {
		config.DB.Model(&users[14]).Association("FavoriteSports").Append(tennisSport, badmintonSport, squashSport)
	}

	log.Printf("✅ Created %d users", len(users))
	return users, nil
}

func seedLocations() ([]models.Location, error) {
	log.Println("📍 Seeding locations...")

	locations := []models.Location{
		{
			Address:     "Fælledparken Tennisbaner",
			Coordinates: types.Point{Lat: 55.6908, Lon: 12.5704},
			PostalCode:  "2100",
			City:        "København",
			Country:     "Danmark",
		},
		{
			Address:     "Fælledparken Fodboldbaner",
			Coordinates: types.Point{Lat: 55.6915, Lon: 12.5712},
			PostalCode:  "2100",
			City:        "København",
			Country:     "Danmark",
		},
		{
			Address:     "Nørrebro Basketballbaner",
			Coordinates: types.Point{Lat: 55.6892, Lon: 12.5514},
			PostalCode:  "2200",
			City:        "København",
			Country:     "Danmark",
		},
		{
			Address:     "Islands Brygge Strandvolleyball",
			Coordinates: types.Point{Lat: 55.6667, Lon: 12.5833},
			PostalCode:  "2300",
			City:        "København",
			Country:     "Danmark",
		},
		{
			Address:     "Vesterbro Idrætsanlæg",
			Coordinates: types.Point{Lat: 55.6719, Lon: 12.5500},
			PostalCode:  "1650",
			City:        "København",
			Country:     "Danmark",
		},
	}

	for i := range locations {
		if err := config.DB.Create(&locations[i]).Error; err != nil {
			return nil, fmt.Errorf("failed to create location %s: %w", locations[i].Address, err)
		}
	}

	log.Printf("✅ Created %d locations", len(locations))
	return locations, nil
}

func seedTeams(users []models.User, sports []models.Sport, locations []models.Location) ([]models.Team, error) {
	log.Println("👥 Seeding teams...")

	if len(users) < 2 || len(locations) < 1 {
		return nil, fmt.Errorf("not enough users or locations to create teams")
	}

	tennisSport := findSportByName("Tennis", config.DB)
	footballSport := findSportByName("Football", config.DB)
	basketballSport := findSportByName("Basketball", config.DB)

	// Build a lookup from user ID to user pointer so we do not rely on
	// implicit assumptions like "ID starts at 1 and increments by 1".
	userByID := make(map[uint]*models.User, len(users))
	for i := range users {
		user := &users[i]
		userByID[user.ID] = user
	}

	teams := []models.Team{
		{
			Name:        "Ace Tennis Klub",
			Description: stringPtr("En venlig tennisklub for alle færdighedsniveauer"),
			CreatorID:   users[0].ID,
			LocationID:  &locations[0].ID,
		},
		{
			Name:        "Nørrebro Ballers",
			Description: stringPtr("Konkurrencedygtigt basketballhold"),
			CreatorID:   users[2].ID,
			LocationID:  &locations[2].ID,
		},
		{
			Name:        "Weekend Warriors FC",
			Description: stringPtr("Afslappet fodboldhold"),
			CreatorID:   users[1].ID,
			LocationID:  &locations[1].ID,
		},
	}

	for i := range teams {
		if err := config.DB.Create(&teams[i]).Error; err != nil {
			return nil, fmt.Errorf("failed to create team %s: %w", teams[i].Name, err)
		}

		// Add creator to team using ID lookup (avoids slice index out of range)
		if creator, ok := userByID[teams[i].CreatorID]; ok {
			config.DB.Model(&teams[i]).Association("Users").Append(creator)
		}

		// Add sports
		if i == 0 && tennisSport != nil {
			config.DB.Model(&teams[i]).Association("Sports").Append(tennisSport)
		} else if i == 1 && basketballSport != nil {
			config.DB.Model(&teams[i]).Association("Sports").Append(basketballSport)
		} else if i == 2 && footballSport != nil {
			config.DB.Model(&teams[i]).Association("Sports").Append(footballSport)
		}

		// Add additional members
		if i == 0 && len(users) > 3 {
			// Ace Tennis Klub: Diana (3), Frank (5), Karen (10), Olivia (14)
			config.DB.Model(&teams[i]).Association("Users").Append(&users[3])
			if len(users) > 5 {
				config.DB.Model(&teams[i]).Association("Users").Append(&users[5])
			}
			if len(users) > 10 {
				config.DB.Model(&teams[i]).Association("Users").Append(&users[10])
			}
			if len(users) > 14 {
				config.DB.Model(&teams[i]).Association("Users").Append(&users[14])
			}
		}
		if i == 1 && len(users) > 2 {
			// Nørrebro Ballers: Bob (1), Henry (7), Maria (12)
			config.DB.Model(&teams[i]).Association("Users").Append(&users[1])
			if len(users) > 7 {
				config.DB.Model(&teams[i]).Association("Users").Append(&users[7])
			}
			if len(users) > 12 {
				config.DB.Model(&teams[i]).Association("Users").Append(&users[12])
			}
		}
		if i == 2 && len(users) > 4 {
			// Weekend Warriors FC: Eve (4), Henry (7), Lars (11)
			config.DB.Model(&teams[i]).Association("Users").Append(&users[4])
			if len(users) > 7 {
				config.DB.Model(&teams[i]).Association("Users").Append(&users[7])
			}
			if len(users) > 11 {
				config.DB.Model(&teams[i]).Association("Users").Append(&users[11])
			}
		}
	}

	log.Printf("✅ Created %d teams", len(teams))
	return teams, nil
}

func seedChallenges(users []models.User, teams []models.Team, locations []models.Location, sports []models.Sport) ([]models.Challenge, error) {
	log.Println("🏆 Seeding challenges...")

	if len(users) < 1 || len(locations) < 1 {
		return nil, fmt.Errorf("not enough users or locations to create challenges")
	}

	// Build a lookup from user ID to user pointer so we do not rely on
	// implicit assumptions like "ID starts at 1 and increments by 1".
	userByID := make(map[uint]*models.User, len(users))
	for i := range users {
		user := &users[i]
		userByID[user.ID] = user
	}

	now := time.Now()
	challenges := []models.Challenge{
		{
			Name:        "Weekend Tennis Match",
			Description: "Søger doublepartnere til et venligt match",
			Sport:       "Tennis",
			LocationID:  locations[0].ID,
			CreatorID:   users[0].ID,
			IsIndoor:    false,
			IsPublic:    true,
			IsCompleted: false,
			Status:      models.ChallengeStatusOpen,
			Date:        now.AddDate(0, 0, 7), // Next week
			StartTime:   now.AddDate(0, 0, 7),
			EndTime:     timePtr(now.AddDate(0, 0, 7).Add(2 * time.Hour)),
			TeamSize:    intPtr(4),
		},
		{
			Name:        "Basketball Pickup Game",
			Description: "Afslappet pickup game, alle færdighedsniveauer er velkomne",
			Sport:       "Basketball",
			LocationID:  locations[2].ID,
			CreatorID:   users[2].ID,
			IsIndoor:    false,
			IsPublic:    true,
			IsCompleted: false,
			Status:      models.ChallengeStatusOpen,
			Date:        now, // Today
			StartTime:   now,
			EndTime:     timePtr(now.Add(2 * time.Hour)),
			TeamSize:    intPtr(10),
		},
		{
			Name:        "Fodbold Træningssession",
			Description: "Holdtræningssession",
			Sport:       "Football",
			LocationID:  locations[1].ID,
			CreatorID:   users[1].ID,
			IsIndoor:    false,
			IsPublic:    false,
			IsCompleted: false,
			Status:      models.ChallengeStatusOpen,
			Date:        now, // Today
			StartTime:   now,
			EndTime:     timePtr(now.Add(2 * time.Hour)),
			Comment:     stringPtr("Medbring dit eget udstyr"),
		},
		{
			Name:        "Volleyball Stranddag",
			Description: "Strandvolleyballturnering",
			Sport:       "Volleyball",
			LocationID:  locations[3].ID,
			CreatorID:   users[3].ID,
			IsIndoor:    false,
			IsPublic:    true,
			IsCompleted: false,
			Status:      models.ChallengeStatusOpen,
			Date:        now.AddDate(0, 0, 10), // In 10 days
			StartTime:   now,
			EndTime:     timePtr(now.AddDate(0, 0, 10).Add(2 * time.Hour)),
			HasCost:     true,
			PlayFor:     stringPtr("Sjov"),
		},
	}

	for i := range challenges {
		if err := config.DB.Create(&challenges[i]).Error; err != nil {
			return nil, fmt.Errorf("failed to create challenge %s: %w", challenges[i].Name, err)
		}

		// Add creator to challenge using ID lookup (avoids slice index out of range)
		if creator, ok := userByID[challenges[i].CreatorID]; ok {
			config.DB.Model(&challenges[i]).Association("Users").Append(creator)
		}

		// Add some additional participants
		if i == 0 && len(users) > 3 {
			// Weekend Tennis Match: Diana (3), Frank (5), Karen (10)
			config.DB.Model(&challenges[i]).Association("Users").Append(&users[3])
			if len(users) > 5 {
				config.DB.Model(&challenges[i]).Association("Users").Append(&users[5])
			}
			if len(users) > 10 {
				config.DB.Model(&challenges[i]).Association("Users").Append(&users[10])
			}
		}
		if i == 1 && len(users) > 1 {
			// Basketball Pickup Game: Bob (1), Henry (7), Maria (12)
			config.DB.Model(&challenges[i]).Association("Users").Append(&users[1])
			if len(users) > 7 {
				config.DB.Model(&challenges[i]).Association("Users").Append(&users[7])
			}
			if len(users) > 12 {
				config.DB.Model(&challenges[i]).Association("Users").Append(&users[12])
			}
		}
		if i == 2 && len(teams) > 2 {
			// 5v5 Football Tournament: Add team and individual users
			config.DB.Model(&challenges[i]).Association("Teams").Append(&teams[2])
			if len(users) > 11 {
				config.DB.Model(&challenges[i]).Association("Users").Append(&users[11])
			}
		}
	}

	log.Printf("✅ Created %d challenges", len(challenges))
	return challenges, nil
}

func seedFriendships(users []models.User) error {
	log.Println("🤝 Seeding friendships...")

	if len(users) < 3 {
		return nil
	}

	// Create a network of friendships to enable friend suggestions
	// Alice (0) is friends with Bob (1), Charlie (2), Frank (5)
	config.DB.Model(&users[0]).Association("Friends").Append(&users[1], &users[2])
	if len(users) > 5 {
		config.DB.Model(&users[0]).Association("Friends").Append(&users[5])
	}

	// Bob (1) is friends with Alice (0), Charlie (2), Henry (7)
	config.DB.Model(&users[1]).Association("Friends").Append(&users[2])
	if len(users) > 7 {
		config.DB.Model(&users[1]).Association("Friends").Append(&users[7])
	}

	// Diana (3) is friends with Eve (4), Karen (10)
	config.DB.Model(&users[3]).Association("Friends").Append(&users[4])
	if len(users) > 10 {
		config.DB.Model(&users[3]).Association("Friends").Append(&users[10])
	}

	// Frank (5) is friends with Alice (0), Grace (6), Olivia (14)
	if len(users) > 6 {
		config.DB.Model(&users[5]).Association("Friends").Append(&users[6])
	}
	if len(users) > 14 {
		config.DB.Model(&users[5]).Association("Friends").Append(&users[14])
	}

	// Henry (7) is friends with Bob (1), Maria (12)
	if len(users) > 12 {
		config.DB.Model(&users[7]).Association("Friends").Append(&users[12])
	}

	// Karen (10) is friends with Diana (3), Olivia (14)
	if len(users) > 14 {
		config.DB.Model(&users[10]).Association("Friends").Append(&users[14])
	}

	// Lars (11) is friends with Jack (9), Niels (13)
	if len(users) > 13 {
		config.DB.Model(&users[11]).Association("Friends").Append(&users[9], &users[13])
	}

	log.Println("✅ Created friendships")
	return nil
}

func seedInvitations(users []models.User, teams []models.Team, challenges []models.Challenge) error {
	log.Println("📨 Seeding invitations...")

	if len(users) < 3 || len(teams) < 1 || len(challenges) < 1 {
		return nil
	}

	invitations := []models.Invitation{
		{
			InviterId:    users[0].ID,
			InviteeId:    users[3].ID,
			ResourceType: models.ResourceTypeTeam,
			ResourceID:   teams[0].ID,
			Status:       models.StatusPending,
			Note:         "Bliv medlem af vores tennisklub!",
		},
		{
			InviterId:    users[1].ID,
			InviteeId:    users[2].ID,
			ResourceType: models.ResourceTypeFriend,
			ResourceID:   0, // Not applicable for friend requests
			Status:       models.StatusPending,
		},
		{
			InviterId:    users[2].ID,
			InviteeId:    users[4].ID,
			ResourceType: models.ResourceTypeChallenge,
			ResourceID:   challenges[1].ID,
			Status:       models.StatusAccepted,
			Note:         "Håber du kan komme!",
		},
	}

	for i := range invitations {
		if err := config.DB.Create(&invitations[i]).Error; err != nil {
			return fmt.Errorf("failed to create invitation: %w", err)
		}
	}

	log.Printf("✅ Created %d invitations", len(invitations))
	return nil
}

// Helper functions
func findSportByName(name string, db *gorm.DB) *models.Sport {
	var sport models.Sport
	if err := db.Where("name = ?", name).First(&sport).Error; err != nil {
		return nil
	}
	return &sport
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func intPtr(i int) *int {
	return &i
}

func seedConversations(users []models.User, teams []models.Team) error {
	log.Println("💬 Seeding conversations...")

	if len(users) < 3 {
		return nil
	}

	// Create direct conversations using the service function
	// This ensures participants are created with proper joined_at timestamps
	directConv1, err := services.CreateDirectConversation(users[0].ID, users[1].ID)
	if err != nil {
		return fmt.Errorf("failed to create direct conversation 1: %w", err)
	}

	directConv2, err := services.CreateDirectConversation(users[0].ID, users[2].ID)
	if err != nil {
		return fmt.Errorf("failed to create direct conversation 2: %w", err)
	}

	directConv3, err := services.CreateDirectConversation(users[1].ID, users[2].ID)
	if err != nil {
		return fmt.Errorf("failed to create direct conversation 3: %w", err)
	}

	directConversations := []*models.Conversation{directConv1, directConv2, directConv3}

	// Create a group conversation using the service function
	if len(users) >= 4 {
		participantIDs := []uint{users[0].ID, users[1].ID, users[2].ID, users[3].ID}
		groupConv, err := services.CreateGroupConversation(users[0].ID, participantIDs, "Weekend Sports Group")
		if err != nil {
			return fmt.Errorf("failed to create group conversation: %w", err)
		}

		// Add some messages to the group
		_, err = services.SendMessage(groupConv.ID, users[0].ID, "Hey everyone! Ready for this weekend?")
		if err != nil {
			return fmt.Errorf("failed to send message 1: %w", err)
		}

		time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamps

		_, err = services.SendMessage(groupConv.ID, users[1].ID, "Absolutely! What time are we meeting?")
		if err != nil {
			return fmt.Errorf("failed to send message 2: %w", err)
		}

		time.Sleep(10 * time.Millisecond)

		_, err = services.SendMessage(groupConv.ID, users[2].ID, "I'm thinking 10 AM at Fælledparken")
		if err != nil {
			return fmt.Errorf("failed to send message 3: %w", err)
		}
	}

	// Create team conversations using SyncTeamConversationMembers
	for i := range teams {
		// Get all team members
		var teamMembers []models.User
		config.DB.Model(&teams[i]).Association("Users").Find(&teamMembers)

		// Extract member IDs
		memberIDs := make([]uint, len(teamMembers))
		for j := range teamMembers {
			memberIDs[j] = teamMembers[j].ID
		}

		// Sync team conversation (creates conversation and participants)
		err := services.SyncTeamConversationMembers(teams[i].ID, memberIDs)
		if err != nil {
			return fmt.Errorf("failed to sync team conversation for team %d: %w", teams[i].ID, err)
		}
	}

	// Add some messages to the first direct conversation using SendMessage
	if len(directConversations) > 0 {
		_, err := services.SendMessage(directConversations[0].ID, users[1].ID, "Hey! Want to play tennis this weekend?")
		if err != nil {
			return fmt.Errorf("failed to send message 1: %w", err)
		}

		time.Sleep(10 * time.Millisecond)

		_, err = services.SendMessage(directConversations[0].ID, users[0].ID, "Sure! Saturday morning works for me")
		if err != nil {
			return fmt.Errorf("failed to send message 2: %w", err)
		}

		time.Sleep(10 * time.Millisecond)

		_, err = services.SendMessage(directConversations[0].ID, users[1].ID, "Perfect! See you at 9 AM at Fælledparken?")
		if err != nil {
			return fmt.Errorf("failed to send message 3: %w", err)
		}
	}

	log.Printf("✅ Created conversations with messages")
	return nil
}
