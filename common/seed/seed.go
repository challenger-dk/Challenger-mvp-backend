package seed

import (
	"fmt"
	"log"
	"time"

	"server/common/config"
	"server/common/models"
	"server/common/models/types"

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

	log.Println("✅ Database seeding completed successfully!")
	return nil
}

func clearDatabase() error {
	log.Println("🧹 Clearing existing data...")

	// Delete in reverse order of dependencies
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

	users := []models.User{
		{
			Email:     "user1@challenger.dk",
			Password:  string(hashedPassword),
			FirstName: "Alice",
			LastName:  "Johnson",
			Bio:       "Tennisentusiast og weekendkriger",
			Age:       28,
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user2@challenger.dk",
			Password:  string(hashedPassword),
			FirstName: "Bob",
			LastName:  "Smith",
			Bio:       "Fodboldspiller på jagt efter konkurrencedygtige kampe",
			Age:       32,
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user3@challenger.dk",
			Password:  string(hashedPassword),
			FirstName: "Charlie",
			LastName:  "Brown",
			Bio:       "Basketballfanatiker",
			Age:       25,
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user4@challenger.dk",
			Password:  string(hashedPassword),
			FirstName: "Diana",
			LastName:  "Williams",
			Bio:       "Elsker at spille padel tennis og volleyball",
			Age:       30,
			Settings:  &models.UserSettings{},
		},
		{
			Email:     "user5@challenger.dk",
			Password:  string(hashedPassword),
			FirstName: "Eve",
			LastName:  "Davis",
			Bio:       "Løbe- og cykelentusiast",
			Age:       27,
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
			config.DB.Model(&teams[i]).Association("Users").Append(&users[3])
		}
		if i == 1 && len(users) > 2 {
			config.DB.Model(&teams[i]).Association("Users").Append(&users[1])
		}
		if i == 2 && len(users) > 4 {
			config.DB.Model(&teams[i]).Association("Users").Append(&users[4])
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
			Date:        now.AddDate(0, 0, 7),                              // Next week
			StartTime:   now.AddDate(0, 0, 7).Add(10 * time.Hour),          // 10 AM
			EndTime:     timePtr(now.AddDate(0, 0, 7).Add(12 * time.Hour)), // 12 PM
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
			Date:        now.AddDate(0, 0, 3),                              // In 3 days
			StartTime:   now.AddDate(0, 0, 3).Add(18 * time.Hour),          // 6 PM
			EndTime:     timePtr(now.AddDate(0, 0, 3).Add(20 * time.Hour)), // 8 PM
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
			Date:        now.AddDate(0, 0, 5),                              // In 5 days
			StartTime:   now.AddDate(0, 0, 5).Add(17 * time.Hour),          // 5 PM
			EndTime:     timePtr(now.AddDate(0, 0, 5).Add(19 * time.Hour)), // 7 PM
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
			Date:        now.AddDate(0, 0, 10),                              // In 10 days
			StartTime:   now.AddDate(0, 0, 10).Add(14 * time.Hour),          // 2 PM
			EndTime:     timePtr(now.AddDate(0, 0, 10).Add(18 * time.Hour)), // 6 PM
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
			config.DB.Model(&challenges[i]).Association("Users").Append(&users[3])
		}
		if i == 1 && len(users) > 1 {
			config.DB.Model(&challenges[i]).Association("Users").Append(&users[1])
		}
		if i == 2 && len(teams) > 2 {
			config.DB.Model(&challenges[i]).Association("Teams").Append(&teams[2])
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

	// Create some friendships
	config.DB.Model(&users[0]).Association("Friends").Append(&users[1], &users[2])
	config.DB.Model(&users[1]).Association("Friends").Append(&users[2])
	config.DB.Model(&users[3]).Association("Friends").Append(&users[4])

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
