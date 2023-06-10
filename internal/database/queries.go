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
		name TEXT NOT NULL,
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
		cardname TEXT NOT NULL,
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
		name TEXT NOT NULL,
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
		name TEXT NOT NULL,
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
	INSERT INTO gk_cards(cardname, data, user_id)
	VALUES ($1,$2,(SELECT id FROM gk_users WHERE username=$3));
`

const getCardQuery = `
	SELECT gk_cards.cardname, gk_cards.data
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
	SELECT gk_cards.cardname
	FROM gk_cards
	JOIN gk_users ON gk_cards.user_id = gk_users.id
	WHERE gk_cards.cardname=$1 AND gk_users.username=$2;
`

const setLoginCredsQuery = `
	INSERT INTO gk_logincreds(name, data, user_id)
	VALUES ($1,$2,(SELECT id FROM gk_users WHERE username=$3));
`

const getLoginCredsQuery = `
	SELECT gk_logincreds.name, gk_logincreds.data 
	FROM gk_logincreds
	JOIN gk_users ON gk_logincreds.user_id = gk_users.id
	WHERE gk_logincreds.name=$1 AND gk_users.username=$2;
`

const listLoginCredsQuery = `
	SELECT gk_logincreds.name
	FROM gk_logincreds
	JOIN gk_users ON gk_logincreds.user_id = gk_users.id
	WHERE gk_users.username = $1;
`

const setNoteQuery = `
	INSERT INTO gk_notes(name, data, user_id)
	VALUES ($1,$2,(SELECT id FROM gk_users WHERE username=$3));
`

const getNoteQuery = `
	SELECT gk_notes.name, gk_notes.data 
	FROM gk_notes
	JOIN gk_users ON gk_notes.user_id = gk_users.id
	WHERE gk_notes.name=$1 AND gk_users.username=$2;
`

const listNotesQuery = `
	SELECT gk_notes.name
	FROM gk_notes
	JOIN gk_users ON gk_notes.user_id = gk_users.id
	WHERE gk_users.username = $1;
`

const setBinaryQuery = `
	INSERT INTO gk_binaries(name, data, user_id)
	VALUES ($1,$2,(SELECT id FROM gk_users WHERE username=$3));
`

const getBinaryQuery = `
	SELECT gk_binaries.name, gk_binaries.data 
	FROM gk_binaries
	JOIN gk_users ON gk_binaries.user_id = gk_users.id
	WHERE gk_binaries.name=$1 AND gk_users.username=$2;
`

const listBinariesQuery = `
	SELECT gk_binaries.name
	FROM gk_binaries
	JOIN gk_users ON gk_binaries.user_id = gk_users.id
	WHERE gk_users.username = $1;
`
