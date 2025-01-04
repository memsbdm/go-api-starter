package seed

import (
	"context"
	"fmt"
	"go-starter/internal/adapters/storage/postgres/repositories"
	"go-starter/internal/domain/entities"
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

func Seed(repo *repositories.UserRepository) {
	ctx := context.Background()
	users := generateUsers(100)
	for _, user := range users {
		if _, err := repo.Create(ctx, user); err != nil {
			panic(err)
		}
	}
}

func generateUsers(num int) []*entities.User {
	users := make([]*entities.User, num)
	for i := 0; i < num; i++ {
		users[i] = &entities.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
		}
	}

	return users
}
