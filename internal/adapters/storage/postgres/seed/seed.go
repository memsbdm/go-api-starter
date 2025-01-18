package seed

import (
	"context"
	"errors"
	"fmt"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/services"
	"sync"
)

const (
	defaultBatchSize = 50
	defaultPassword  = "secret123"
)

// UserGenerator handles the creation of test user data.
// It maintains a list of predefined usernames and provides methods
// to generate test users in batches efficiently and in parallel.
type UserGenerator struct {
	usernames []string
	password  string
	service   *services.UserService
}

// NewUserGenerator returns a new instance of UserGenerator
func NewUserGenerator(svc *services.UserService) *UserGenerator {
	return &UserGenerator{
		usernames: []string{
			"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "heidi",
			"ivan", "judy", "karl", "laura", "mallory", "nina", "oscar", "peggy",
			"quinn", "rachel", "steve", "trent", "ursula", "victor", "wendy", "xander",
			"yvonne", "zack", "amber", "brian", "carol", "doug", "eric", "fiona",
			"george", "hannah", "ian", "jessica", "kevin", "lisa", "mike", "natalie",
			"oliver", "peter", "queen", "ron", "susan", "tim", "uma", "vicky",
			"walter", "xenia", "yasmin", "zoe",
		},
		password: defaultPassword,
		service:  svc,
	}
}

// GenerateUsersOptions helps to provide options for generation
type GenerateUsersOptions struct {
	Count     int
	BatchSize int
}

// GenerateUsers creates and registers a specified number of test users in the system.
// It uses parallel batch processing to optimize performance.
// Returns an error if generation fails or if parameters are invalid.
func (g *UserGenerator) GenerateUsers(ctx context.Context, opts GenerateUsersOptions) error {
	if opts.Count <= 0 {
		return fmt.Errorf("count must be positive, got %d", opts.Count)
	}

	if opts.BatchSize <= 0 {
		opts.BatchSize = defaultBatchSize
	}

	var wg sync.WaitGroup
	errChan := make(chan error, opts.Count)

	for i := 0; i < opts.Count; i += opts.BatchSize {
		batchSize := opts.BatchSize
		if i+batchSize > opts.Count {
			batchSize = opts.Count - i
		}

		wg.Add(1)
		go func(start, size int) {
			defer wg.Done()
			if err := g.generateBatch(ctx, start, size, defaultPassword); err != nil {
				errChan <- err
			}
		}(i, batchSize)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to generate users: %v", errors.Join(errs...))
	}

	return nil
}

// generateBatch creates and registers a specific batch of users.
// This function is intended to be run in a separate goroutine.
func (g *UserGenerator) generateBatch(ctx context.Context, start, size int, password string) error {
	users := make([]*entities.User, size)
	for i := 0; i < size; i++ {
		idx := start + i
		username := fmt.Sprintf("%s%d", g.usernames[idx%len(g.usernames)], idx)
		users[i] = &entities.User{
			Username: username,
			Name:     g.usernames[idx%len(g.usernames)],
			Password: password,
		}
	}

	for _, user := range users {
		if _, err := g.service.Register(ctx, user); err != nil {
			return fmt.Errorf("failed to register user %s: %w", user.Username, err)
		}
	}

	return nil
}
