package dependency_injection

import (
	"database/sql"
	domain_repository "main/internal/Domain/Repository"
	infra_repository "main/internal/Infrastructure/Repository"
	"os"
	"sync"
	_ "github.com/lib/pq"
)

type Container struct {
	DB             *sql.DB
	PostRepository domain_repository.PostRepository
}

var lock = sync.Mutex{}
var container *Container

func GetContainer() *Container {
	if container == nil {
		lock.Lock()
		defer lock.Unlock()
		db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			panic(err)
		}

		container = &Container{
			DB:             db,
			PostRepository: infra_repository.NewPostRepository(db),
		}
	}
	return container
}
