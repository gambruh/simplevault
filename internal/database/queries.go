package database

// SQL queries
// database init queries
const checkTableExistsQuery = `
	SELECT EXISTS (
		SELECT 	1 
		FROM 	information_schema.tables
		WHERE 	table_name = $1
	);
`

const createLoginCredsTableQuery = `
	CREATE TABLE gk_logincreds (
		name TEXT,
		login TEXT,
		password TEXT,
		site TEXT,
		user_id integer NOT NULL,
		CONSTRAINT fk_gk_users
			FOREIGN KEY (user_id)
				REFERENCES gk_users(id)
				ON DELETE CASCADE
	);
`

const createCardsTableQuery = `
	CREATE TABLE gk_cards (
		number TEXT,
		name TEXT,
		surname TEXT,
		code integer,
		valid_till TEXT,
		bank TEXT,
		user_id integer NOT NULL,
		CONSTRAINT fk_gk_users
			FOREIGN KEY (user_id)
				REFERENCES gk_users(id)
				ON DELETE CASCADE
	);
`

const createNotesTableQuery = `
	CREATE TABLE gk_notes (
		id SERIAL,
		user_id integer NOT NULL,
		name TEXT,
		note TEXT,
		PRIMARY KEY (id),
		CONSTRAINT fk_gk_users
			FOREIGN KEY (user_id)
				REFERENCES gk_users(id)
				ON DELETE CASCADE
	)
`

const createBinariesTableQuery = `
	CREATE TABLE gk_binaries (
		id SERIAL,
		user_id integer NOT NULL,
		name TEXT,
		data BYTEA,
		PRIMARY KEY (id),
		CONSTRAINT fk_gk_users
			FOREIGN KEY (user_id)
				REFERENCES gk_users(id)
				ON DELETE CASCADE
	)
`

// set/get queries

const listCardsQuery = `
	SELECT gk_cards.bank, gk_cards.number, gk_cards.name, gk_cards.surname, gk_cards.valid_till, gk_cards.code
	FROM gk_cards
	JOIN gk_users ON cards.user_id = gk_users.id
	WHERE gk_users.username = $1;
`

const setCardQuery = `
	INSERT INTO gk_cards(number, name, surname, code, bank, valid_till, user_id)
	VALUES ($1,$2,$3,$4,$5,$6);
`

const CheckIDbyUsernameQuery = `
	SELECT id 
	FROM gk_users 
	WHERE username = $1;
`
