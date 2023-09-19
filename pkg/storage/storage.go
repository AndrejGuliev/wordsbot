package storage

type Storage interface {
	AddUser(userID int64) error
	GetCurrentWord(userID int64) (int, string, string, error)
	EncCurrentWordNum(userID int64) error
	GetCurrentAnswNum(userID int64) (int, error)
	EncCurrentAnswNum(userID int64) error
	SetTest(userID int64, test string) error
	GetCurrentTest(userID int64) (string, error)
	EndTest(userID int64) error
	SetPosition(userID int64, pos int) error
	GetCurrentPosition(userID int64) (int, error)
	MakeTestsList(userID int64) ([]string, error)
	AddNewPair(userID int64, word [][]string) error
	AddNewTestName(userID int64, testName string) error
	ValidateName(userID int64, testName string) (bool, error)
	GetTestIdRange(userID int64, testName string) (int, int, error)
	SetRandomWord(userID int64, wordID int) error
	DeletePackage(userID int64, testName string) error
}
