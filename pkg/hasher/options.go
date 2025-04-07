package hasher

// Option -.
type Option func(hasher *BcryptPasswordHasher)

// Cost -.
func Cost(cost int) Option {
	return func(s *BcryptPasswordHasher) {
		s.cost = cost
	}
}
