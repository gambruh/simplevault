package database

const encryptedCardsTable = `
	ALTER TABLE gk_cards
	DROP COLUMN number,
	DROP COLUMN name,
	DROP COLUMN surname,
	DROP COLUMN valid_till,
	DROP COLUMN code
`

const addColumnData = `
	ALTER TABLE gk_cards
	ADD COLUMN data TEXT
`

const alterTableLC = `
	ALTER TABLE gk_logincreds
	ADD COLUMN data TEXT
`

const alterTableLC2 = `
	ALTER TABLE gk_logincreds
	DROP COLUMN login,
	DROP COLUMN password,
	DROP COLUMN site
`

const alterTableNotes = `
	ALTER TABLE gk_notes
	RENAME COLUMN note TO data
`

// UNIQUE constrains
const createUniqueCardConstraint = `
	ALTER TABLE gk_cards
	ADD CONSTRAINT gk_unique_cardname UNIQUE (cardname, user_id)
`

const createUniqueNoteConstraint = `
	ALTER TABLE gk_notes
	ADD CONSTRAINT gk_unique_notename UNIQUE (name, user_id)
`

const createUniqueLoginCredsConstraint = `
	ALTER TABLE gk_logincreds
	ADD CONSTRAINT gk_unique_name UNIQUE (name, user_id)
`

const createUniqueBinaryConstraint = `
	ALTER TABLE gk_binaries
	ADD CONSTRAINT gk_unique_binaryname UNIQUE (name, user_id)
`
