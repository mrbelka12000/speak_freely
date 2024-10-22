package linguo_sphere_backend

type AuthMethod int

const (
	AuthMethodWeb AuthMethod = iota + 1
	AuthMethodTG
)
