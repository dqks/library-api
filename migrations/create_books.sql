DROP TABLE IF EXISTS books
CREATE TABLE IF NOT EXISTS books (
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	author TEXT NOT NULL,
	available BOOL NOT NULL,
	year SMALLINT NOT NULL,
	CONSTRAINT ck_title_not_empty CHECK (TRIM(title) <> ''),
	CONSTRAINT ck_author_not_empty CHECK (TRIM(author) <> ''),
	CONSTRAINT ck_year_not_greater_than_current CHECK (year <= EXTRACT(YEAR FROM CURRENT_DATE))
)
