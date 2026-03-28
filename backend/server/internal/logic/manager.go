package logic

type Manager struct {
	DataBase interface {
		Register(RegisterData) error
		Login(RegisterData) (bool, int, error)
		GetLeaderBoard() (LeaderBoard, error)
	}
}

func NewManager(dbM interface {
	Register(RegisterData) error
	Login(RegisterData) (bool, int, error)
	GetLeaderBoard() (LeaderBoard, error)
}) *Manager {
	return &Manager{
		DataBase: dbM,
	}
}

func (m *Manager) RegisterUser(data RegisterData) error {
	err := m.DataBase.Register(data)
	return err
}

func (m *Manager) LoginUser(data RegisterData) (bool, int, error) {
	isLogined, rank, err := m.DataBase.Login(data)
	return isLogined, rank, err
}

func (m *Manager) GetLeaderBoard() (LeaderBoard, error) {
	leaderBoard, err := m.DataBase.GetLeaderBoard()
	return leaderBoard, err
}
