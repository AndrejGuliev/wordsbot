package mysql

import (
	"database/sql"
	"log"
	"strings"

	//"github.com/AndrejGuliev/wordsbot/pkg/storage"
	_ "github.com/go-sql-driver/mysql"
)

type WordsBotStorage struct {
	db *sql.DB
}

// NewWordsBotStorage creates a new WordsBotStorage instance with a reference to the given database.
func NewWordsBotStorage(db *sql.DB) *WordsBotStorage {
	return &WordsBotStorage{db: db}
}

// AddUser inserts a new user into the database with the given userID.
func (s *WordsBotStorage) AddUser(userID int64) error {
	_, err := s.db.Exec("INSERT INTO users (telegram_id) VALUES(?)", userID)
	return err
}

// GetCurrentWord retrieves the current word, its translation, and its ID for a given userID.
// It performs a SQL query to fetch this information.
func (s *WordsBotStorage) GetCurrentWord(userID int64) (int, string, string, error) {
	var wordID int
	var word, translation string
	row := s.db.QueryRow("SELECT u.current_word, w.word, w.translation FROM users u INNER JOIN words w ON u.current_word = w.word_id WHERE telegram_id = ?", userID)
	err := row.Scan(&wordID, &word, &translation)
	if err == sql.ErrNoRows {
		return 0, "", "", nil
	}
	return wordID, word, translation, err
}

// EncCurrentWordNum updates the current word number for a user, based on certain conditions.
func (s *WordsBotStorage) EncCurrentWordNum(userID int64) error {
	_, err := s.db.Exec(`UPDATE users AS u1 JOIN users AS u2 ON u2.telegram_id = ? SET u1.current_word = (SELECT word_id FROM words WHERE word_id > u2.current_word AND test = u2.current_test LIMIT 1) WHERE u1.telegram_id = ?`, userID, userID)
	return err
}

// GetCurrentAnswNum retrieves the current answer number for a user.
func (s *WordsBotStorage) GetCurrentAnswNum(userID int64) (int, error) {
	var CurrentAnswNum int
	row := s.db.QueryRow("SELECT current_answ_num FROM users WHERE telegram_id = ?", userID)
	err := row.Scan(&CurrentAnswNum)
	return CurrentAnswNum, err
}

// EncCurrentAnswNum increments the current answer number for a user.
func (s *WordsBotStorage) EncCurrentAnswNum(userID int64) error {
	_, err := s.db.Exec("UPDATE users SET current_answ_num = current_answ_num + 1 WHERE telegram_id = ?", userID)
	return err
}

// SetTest sets the current test for a user and updates their current word accordingly.
func (s *WordsBotStorage) SetTest(userID int64, test string) error {
	_, err := s.db.Exec("UPDATE users SET current_test = ?, current_word = (SELECT word_id FROM words WHERE word_id > 0 AND test = ? LIMIT 1) WHERE telegram_id = ?", test, test, userID)
	return err
}

// GetCurrentTest retrieves the current test name for a user.
func (s *WordsBotStorage) GetCurrentTest(userID int64) (string, error) {
	var currentTest string
	row := s.db.QueryRow("SELECT current_test FROM users WHERE telegram_id = ?", userID)
	err := row.Scan(&currentTest)
	return currentTest, err
}

// EndTest resets the current test and answer number for a user.
func (s *WordsBotStorage) EndTest(userID int64) error {
	_, err := s.db.Exec("UPDATE users SET current_test = DEFAULT, current_answ_num = DEFAULT, position = DEFAULT WHERE telegram_id = ?", userID)
	return err
}

// SetPosition updates the position for a user.
func (s *WordsBotStorage) SetPosition(userID int64, pos int) error {
	_, err := s.db.Exec("UPDATE users SET position = ? WHERE telegram_id = ?", pos, userID)
	return err
}

// GetCurrentPosition retrieves the current position for a user.
func (s *WordsBotStorage) GetCurrentPosition(userID int64) (int, error) {
	var currentPosition int
	row := s.db.QueryRow("SELECT position FROM users WHERE telegram_id = ?", userID)
	err := row.Scan(&currentPosition)
	return currentPosition, err
}

// MakeTestsList retrieves a list of distinct test names associated with a user.
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

// AddNewPair adds a new word pair (word and translation) for a user into the database.
func (s *WordsBotStorage) AddNewPair(userID int64, words [][]string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	for _, pair := range words {
		if len(pair) != 2 {
			log.Println(len(pair))
			tx.Rollback()
			return err
		}
		_, err := tx.Exec("INSERT INTO words (word, translation, owner) VALUES(?, ?, ?)", strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1]), userID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return err
}

// AddNewTestName updates the test name for words owned by a user in the database.
func (s *WordsBotStorage) AddNewTestName(userID int64, testName string) error {
	_, err := s.db.Exec("UPDATE words SET test = ? WHERE owner = ? AND test IS NULL", testName, userID)
	return err
}

// ValidateName checks if a test name is valid for a user and returns true if it exists, false otherwise.
func (s *WordsBotStorage) ValidateName(userID int64, testName string) (bool, error) {
	var tmp interface{}
	row := s.db.QueryRow("SELECT test FROM words WHERE owner = ? AND test = ? LIMIT 1", userID, testName)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

// GetTestIdRange retrieves the minimum and maximum word IDs associated with a specific test for a user.
func (s *WordsBotStorage) GetTestIdRange(userID int64, testName string) (int, int, error) {
	var start, end int
	row := s.db.QueryRow("SELECT (SELECT MIN(word_id) FROM words WHERE owner = ? AND test = ?), (SELECT MAX(word_id) FROM words WHERE owner = ? AND test = ?)", userID, testName, userID, testName)
	err := row.Scan(&start, &end)
	return start, end, err
}

// SetRandomWord updates the current word for a user to a random word ID within a specified range.
func (s *WordsBotStorage) SetRandomWord(userID int64, wordID int) error {
	_, err := s.db.Exec("UPDATE users SET current_word = (SELECT word_id FROM words WHERE word_id >= ? LIMIT 1) WHERE telegram_id = ?", wordID, userID)
	return err
}

// DeletePocket deletes all words associated with a specific test for a user.
func (s *WordsBotStorage) DeletePackage(userID int64, testName string) error {
	_, err := s.db.Exec("DELETE FROM words WHERE owner = ? AND test = ? ", userID, testName)
	return err
}
