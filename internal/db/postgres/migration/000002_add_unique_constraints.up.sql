DROP INDEX IF EXISTS short_url_index;
DROP INDEX IF EXISTS long_url_index;

CREATE UNIQUE INDEX short_url_unique_index
ON shorten_url (short_url);

CREATE UNIQUE INDEX long_url_unique_index
ON shorten_url (long_url);
