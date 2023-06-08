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
	DROP COLUMN password
`
