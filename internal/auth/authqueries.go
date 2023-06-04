package auth

// auth queries
const CheckUsernameQuery = `
	SELECT id
	FROM gk_users 
	WHERE username = $1;
`

const CheckPasswordQuery = `
	SELECT password
	FROM gk_passwords 
	WHERE id = $1;
`

const AddNewUserQuery = `
	WITH new_user AS (
		INSERT INTO gk_users (username)
		VALUES ($1)
		RETURNING id
	)
	INSERT INTO gk_passwords (id, password)
	VALUES ((SELECT id FROM new_user), $2);
`

const checkTableExistsQuery = `
	SELECT EXISTS (
		SELECT 	1 
		FROM 	information_schema.tables
		WHERE 	table_name = $1
	);
`

const createUsersTableQuery = `
	CREATE TABLE gk_users (
		id SERIAL,
		username text NOT NULL UNIQUE,
		PRIMARY KEY (id)
	);
`

const createPasswordsTableQuery = `
	CREATE TABLE gk_passwords (
		id integer PRIMARY KEY,
		password TEXT NOT NULL,
		CONSTRAINT fk_gk_users
			FOREIGN KEY (id) 
				REFERENCES gk_users(id)
				ON DELETE CASCADE
	);
`

//const dropPasswordsTableQuery = `
//	DROP TABLE gk_passwords CASCADE;
//`

//const dropUsersTableQuery = `
//	DROP TABLE gk_users CASCADE;
//`
