package smsproxy

type failingRepository struct {
	saveErrors   map[MessageID]error
	updateErrors map[MessageID]error
	repository   repository
}

func newFailingRepository() failingRepository {
	return failingRepository{
		saveErrors:   make(map[MessageID]error),
		updateErrors: make(map[MessageID]error),
		repository:   newRepository(),
	}
}

func (m *failingRepository) saveError(id MessageID, err error) *failingRepository {
	m.saveErrors[id] = err
	return m
}

func (m *failingRepository) updateError(id MessageID, err error) *failingRepository {
	m.updateErrors[id] = err
	return m
}

func (m *failingRepository) save(id MessageID) error {
	if err, ok := m.saveErrors[id]; ok {
		return err
	}
	return m.repository.save(id)
}

func (m *failingRepository) update(id MessageID, newStatus MessageStatus) error {
	if err, ok := m.updateErrors[id]; ok {
		return err
	}
	return m.repository.update(id, newStatus)
}

func (m *failingRepository) get(id MessageID) (MessageStatus, error) {
	return m.repository.get(id)
}

func (m *failingRepository) build() repository {
	return repository(m)
}
