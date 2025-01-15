package seed

import (
	"context"
	"fmt"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/services"
)

var usernames = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
	"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
	"quinn", "rachel", "steve", "trent", "ursula", "victor", "wendy", "xander",
	"yvonne", "zack", "amber", "brian", "carol", "doug", "eric", "fiona",
	"george", "hannah", "ian", "jessica", "kevin", "lisa", "mike", "natalie",
	"oliver", "peter", "queen", "ron", "susan", "tim", "uma", "vicky",
	"walter", "xenia", "yasmin", "zoe",
}

// Users populates the user repository with a predefined number of user entities.
func Users(svc *services.UserService) {
	ctx := context.Background()
	users := generateUsers(100)
	for _, user := range users {
		if _, err := svc.Register(ctx, user); err != nil {
			panic(err)
		}
	}
}

// generateUsers creates a slice of user entities with the specified number of users.
func generateUsers(num int) []*entities.User {
	users := make([]*entities.User, num)
	for i := 0; i < num; i++ {
		users[i] = &entities.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Password: "secret123",
		}
	}

	return users
}
