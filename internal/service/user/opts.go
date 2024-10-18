package user

type opt func(s *Service)

// WithCryptCost set custom crypto cost
func WithCryptCost(cost int) opt {
	return func(s *Service) {
		s.bcryptCost = cost
	}
}
