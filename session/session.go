package session

type Sessioner interface {
	Get(key string) (value any, ok bool)
	Set(key string, value any)
	Delete(key string)
	GetAllKeys() (keys []string)
	CurrentState() string
	SetState(state string)
	ClearData()
}

type SessionManager[K comparable] interface {
	GetOrCreate(id K) (Sessioner, error)
	Set(id K, session Sessioner) error
	Delete(id K) error
}
