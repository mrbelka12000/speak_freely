package user

type opt func(s *Service)

func WithCryptCost(cost int) opt {
	return func(s *Service) {
		s.bcryptCost = cost
	}
}
