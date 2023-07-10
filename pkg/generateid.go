package pkg

import "github.com/google/uuid"

func GenerateUniqueId() string {
	// Generate a new UUID
	id := uuid.New()
	return id.String()
}
