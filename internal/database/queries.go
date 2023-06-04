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
		name TEXT NOT NULL UNIQUE,
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
		cardname TEXT NOT NULL UNIQUE,
		number TEXT,
		name TEXT,
		surname TEXT,
		valid_till TEXT,
		code TEXT,
		user_id integer NOT NULL,
		CONSTRAINT fk_gk_users
			FOREIGN KEY (user_id)
				REFERENCES gk_users(id)
				ON DELETE CASCADE
	);
`

const createNotesTableQuery = `
	CREATE TABLE gk_notes (
		name TEXT NOT NULL UNIQUE,
		note TEXT,
		user_id integer NOT NULL,
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
		name TEXT NOT NULL UNIQUE,
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
	SELECT gk_cards.cardname
	FROM gk_cards
	JOIN gk_users ON gk_cards.user_id = gk_users.id
	WHERE gk_users.username = $1;
`

const setCardQuery = `
	INSERT INTO gk_cards(cardname, number, name, surname, code, valid_till, user_id)
	VALUES ($1,$2,$3,$4,$5,$6,(SELECT id FROM gk_users WHERE username=$7));
`

const setLoginCredsQuery = `
	INSERT INTO gk_logincreds(name, login, password, site, user_id)
	VALUES ($1,$2,$3,$4,(SELECT id FROM gk_users WHERE username=$5));
`

const getCardQuery = `
	SELECT gk_cards.cardname, gk_cards.number, gk_cards.name, gk_cards.surname, gk_cards.valid_till, gk_cards.code
	FROM gk_cards
	JOIN gk_users ON gk_cards.user_id = gk_users.id
	WHERE gk_cards.cardname=$1 AND gk_users.username=$2;
`

const CheckIDbyUsernameQuery = `
	SELECT id 
	FROM gk_users 
	WHERE username = $1;
`

const checkCardNameQuery = `
	SELECT cardname
	FROM gk_cards
	WHERE cardname = $1;
`
