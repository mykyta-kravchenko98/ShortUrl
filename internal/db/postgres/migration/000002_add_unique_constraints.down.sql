DROP INDEX IF EXISTS short_url_unique_index;
DROP INDEX IF EXISTS long_url_unique_index;

CREATE INDEX short_url_index
ON shorten_url (short_url);

CREATE INDEX long_url_index
ON shorten_url (long_url);
