package storage

type user struct {
	competiton    string
	lastWord      string
	correctAnsNum int
}

var (
	testLesson map[string]string
	testUsers  map[int64]user
)

func MakeLesson() map[string]string {
	lesson := map[string]string{
		"Привет": "Hello",
		"Пока":   "Goodbye",
		"День":   "Day",
		"Ночь":   "Night",
	}
	return lesson
}

func AddUser(userID int64) {
	var user user
	user.competiton = ""
	user.lastWord = ""
	user.correctAnsNum = 0
	testUsers[userID] = user
}
