CREATE TABLE IF NOT EXISTS users(id SERIAL PRIMARY KEY, name TEXT, password TEXT, rank INTEGER);
INSERT INTO users(name, password, rank) VALUES('Bodya', '12345', 100);
INSERT INTO users(name, password, rank) VALUES('Anatoliy', 'asdad', 10);
