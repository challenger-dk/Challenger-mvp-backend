package models

type ResourceType string

// Only allowed resource types and statuses
const (
	ResourceTypeTeam      ResourceType = "team"
	ResourceTypeFriend    ResourceType = "friend"
	ResourceTypeChallenge ResourceType = "challenge"
)
