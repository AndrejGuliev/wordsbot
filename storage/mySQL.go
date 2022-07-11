package storage

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type WordsBotStorage struct {
	db *sql.DB
}

func NewWordsBotStorage(db *sql.DB) *WordsBotStorage {
	return &WordsBotStorage{db: db}
}

func (s *WordsBotStorage) AddUser(userID int64) error {
	_, err := s.db.Exec("INSERT INTO users (telegram_id) VALUES(?)", userID)
	return err
}

func (s *WordsBotStorage) CurrentWord(userID int64) (int, string, string, error) {
	var wordID int
	var word, translation string
	row := s.db.QueryRow("SELECT u.current_word, w.translation, w.word FROM users u inner join words w on u.current_word = w.word_id where telegram_id =?", userID)
	err := row.Scan(&wordID, &translation, &word)
	if err == sql.ErrNoRows {
		return 0, "", "", nil
	}
	return wordID, word, translation, err
}

func (s *WordsBotStorage) EncCurrentWordNum(userID int64, currentWord int) error {
	_, err := s.db.Exec("UPDATE users SET current_word = (SELECT word_id from words WHERE word_id >? AND test = (SELECT test FROM words WHERE word_id = ?) LIMIT 1) WHERE telegram_id = ?", currentWord, currentWord, userID)
	return err
}

func (s *WordsBotStorage) SetTest(userID int64, test string) error {
	_, err := s.db.Exec("UPDATE users SET current_test = ?, current_word = (SELECT word_id FROM words WHERE word_id>0 AND test=? LIMIT 1) WHERE telegram_id = ?", test, test, userID)
	return err

}

func (s *WordsBotStorage) CurrentTest(userID int64) (string, error) {
	var currentTest string
	row := s.db.QueryRow("SELECT current_test FROM users WHERE telegram_id = ?", userID)
	err := row.Scan(&currentTest)
	return currentTest, err
}

func (s *WordsBotStorage) EndTest(userID int64) error {
	_, err := s.db.Exec("UPDATE users SET current_test = default WHERE telegram_id =?", userID)
	return err
}

func (s *WordsBotStorage) SetPosition(userID int64, pos int) error {
	_, err := s.db.Exec("UPDATE users SET position = ? WHERE telegram_id = ?", pos, userID)
	return err
}

func (s *WordsBotStorage) CurrentPosition(userID int64) (int, error) {
	var currentPosition int
	row := s.db.QueryRow("SELECT position FROM users WHERE telegram_id = ?", userID)
	err := row.Scan(&currentPosition)
	return currentPosition, err
}

func (s *WordsBotStorage) MakeTestsList(userID int64) ([]string, error) {
	rows, err := s.db.Query("SELECT DISTINCT test FROM words WHERE owner = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testNames []string

	for rows.Next() {
		var testName string
		if err := rows.Scan(&testName); err != nil {
			return nil, err
		}
		testNames = append(testNames, testName)
	}
	return testNames, nil
}

func (s *WordsBotStorage) AddNewPair(userID int64, pair []string) error {
	_, err := s.db.Exec("INSERT INTO words (word, translation, owner) VALUES(?, ?, ?)", pair[0], pair[1], userID)
	return err
}

func (s *WordsBotStorage) AddNewTestName(userID int64, name string) error {
	_, err := s.db.Exec("UPDATE words SET test = ? WHERE owner = ? AND test IS NULL", name, userID)
	return err
}
