package storage

func (s *Storage) JWT() (string, error) {
	return s.aux.JWT()
}

func (s *Storage) SaveAnalytics(a interface{}) error {
	return s.aux.SaveAnalytics(a)
}

func (s *Storage) LoadAnalytics(a interface{}) error {
	return s.aux.LoadAnalytics(a)
}

func (s *Storage) InstallID() string {
	return s.aux.InstallID()
}
