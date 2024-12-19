BEGIN;

-- Создание таблицы tokens
CREATE TABLE tokens (
    user_id INT PRIMARY KEY,
    value   VARCHAR(128) UNIQUE NOT NULL,
    expires TIMESTAMP NOT NULL,
    ip      VARCHAR(15) NOT NULL
);

-- Создание таблицы ips
CREATE TABLE ips (
    id      SERIAL PRIMARY KEY,
    value   VARCHAR(15) NOT NULL,
    user_id INT NOT NULL
);

-- Создание индекса на user_id в таблице ips
CREATE INDEX "user_id_idx" ON ips (user_id);

COMMIT;
