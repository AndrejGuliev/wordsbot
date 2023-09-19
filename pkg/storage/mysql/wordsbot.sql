CREATE DATABASE IF NOT EXISTS wordsbot;

USE wordsbot;

-- Таблица "users"
CREATE TABLE IF NOT EXISTS users (
  telegram_id BIGINT NOT NULL,
  current_word INT DEFAULT '0',
  current_answ_num INT NOT NULL DEFAULT '0',
  current_test VARCHAR(60) DEFAULT '',
  position SMALLINT DEFAULT '0',
  PRIMARY KEY (telegram_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Таблица "words"
CREATE TABLE IF NOT EXISTS words (
  word_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  word VARCHAR(100) DEFAULT NULL,
  translation VARCHAR(120) DEFAULT NULL,
  owner BIGINT NOT NULL,
  test VARCHAR(60) DEFAULT NULL,
  PRIMARY KEY (word_id),
  KEY fk_owner_id (owner),
  CONSTRAINT fk_owner_id FOREIGN KEY (owner) REFERENCES users (telegram_id) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=237 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
