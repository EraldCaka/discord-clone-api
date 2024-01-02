package db

const (
	DBNAME = "discord-clone-api"
	DBURI  = "mongodb://localhost:27017"
)

type Store struct {
	User   UserStore
	Server ServerStore
}
