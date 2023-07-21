CREATE TABLE "shorten_url" (
  "id" bigint PRIMARY KEY,
  "short_url" varchar NOT NULL,
  "long_url" varchar NOT NULL
);

CREATE INDEX short_url_index
ON shorten_url (short_url);

CREATE INDEX long_url_index
ON shorten_url (long_url);
