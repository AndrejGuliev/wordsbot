package storage

import (
	"database/sql"
	"fmt"
	"strings"

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

func (s *WordsBotStorage) GetCurrentWord(userID int64) (int, string, string, error) {
	var wordID int
	var word, translation string
	row := s.db.QueryRow("SELECT u.current_word, w.translation, w.word FROM users u inner join words w on u.current_word = w.word_id where telegram_id =?", userID)
	err := row.Scan(&wordID, &word, &translation)
	if err == sql.ErrNoRows {
		return 0, "", "", nil
	}
	return wordID, word, translation, err
}

func (s *WordsBotStorage) EncCurrentWordNum(userID int64) error {
	_, err := s.db.Exec("UPDATE users as u1 JOIN users as u2 ON u2.telegram_id = ? SET u1.current_word = (SELECT word_id from words WHERE word_id >u2.current_word AND test =u2.current_test LIMIT 1)WHERE u1.telegram_id =?", userID, userID)
	return err
}

func (s *WordsBotStorage) GetCurrentAnswNum(userID int64) (int, error) {
	var CurrentAnswNum int
	row := s.db.QueryRow("SELECT current_answ_num FROM users WHERE telegram_id = ?", userID)
	err := row.Scan(&CurrentAnswNum)
	return CurrentAnswNum, err
}

func (s *WordsBotStorage) EncCurrentAnswNum(userID int64) error {
	_, err := s.db.Exec("UPDATE users SET current_answ_num = current_answ_num + 1 WHERE telegram_id = ?", userID)
	return err
}

func (s *WordsBotStorage) SetTest(userID int64, test string) error {
	_, err := s.db.Exec("UPDATE users SET current_test = ?, current_word = (SELECT word_id FROM words WHERE word_id>0 AND test=? LIMIT 1) WHERE telegram_id = ?", test, test, userID)
	fmt.Println(err)
	return err
}

func (s *WordsBotStorage) CurrentTest(userID int64) (string, error) {
	var currentTest string
	row := s.db.QueryRow("SELECT current_test FROM users WHERE telegram_id = ?", userID)
	err := row.Scan(&currentTest)
	return currentTest, err
}

func (s *WordsBotStorage) EndTest(userID int64) error {
	_, err := s.db.Exec("UPDATE users SET current_test = default, current_answ_num = default WHERE telegram_id =?", userID)
	return err
}

func (s *WordsBotStorage) SetPosition(userID int64, pos int) error {
	_, err := s.db.Exec("UPDATE users SET position = ? WHERE telegram_id = ?", pos, userID)
	return err
}

func (s *WordsBotStorage) GetCurrentPosition(userID int64) (int, error) {
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
	_, err := s.db.Exec("INSERT INTO words (word, translation, owner) VALUES(?, ?, ?)", strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1]), userID)
	return err
}

func (s *WordsBotStorage) AddNewTestName(userID int64, testName string) error {
	_, err := s.db.Exec("UPDATE words SET test = ? WHERE owner = ? AND test IS NULL", testName, userID)
	return err
}

func (s *WordsBotStorage) ValidateName(userID int64, testName string) (bool, error) {
	var tmp interface{}
	row := s.db.QueryRow("SELECT test FROM words WHERE owner =? AND test =? LIMIT 1", userID, testName)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	fmt.Println(err)
	return false, err
}

func (s *WordsBotStorage) TestIdRange(userID int64, testName string) (int, int, error) {
	var start, end int
	row := s.db.QueryRow("SELECT (SELECT MIN(word_id) FROM words WHERE owner =? AND test =?), (SELECT MAX(word_id) FROM words WHERE owner =? AND test =?)", userID, testName, userID, testName)
	err := row.Scan(&start, &end)
	return start, end, err
}

func (s *WordsBotStorage) SetRandomWord(userID int64, wordID int) error {
	_, err := s.db.Exec("UPDATE users SET current_word = (SELECT word_id FROM words WHERE word_id>=? LIMIT 1) WHERE telegram_id = ?", wordID, userID)
	return err
}

func (s *WordsBotStorage) DeletePocket(userID int64, testName string) error {
	_, err := s.db.Exec("DELETE FROM words WHERE owner =? AND test=? ", userID, testName)
	return err
}

func (s *WordsBotStorage) GetMenuMessageID(userID int64) (int, error) {
	var MenuMessageID int
	row := s.db.QueryRow("SELECT menu_message_id FROM users WHERE telegram_id = ?", userID)
	err := row.Scan(&MenuMessageID)
	return MenuMessageID, err
}

func (s *WordsBotStorage) SetMenuMessageID(userID int64, messageID int) error {
	_, err := s.db.Exec("UPDATE users SET menu_message_id=? WHERE telegram_id = ?", messageID, userID)
	return err
}
