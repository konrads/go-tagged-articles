CREATE SEQUENCE article_seq;

CREATE TABLE IF NOT EXISTS article (
    id TEXT NOT NULL DEFAULT ''||nextval('article_seq'::regclass)::TEXT,
    title text not null,
    date  date not null,
    body  text not null,
    tags  text[],
    CONSTRAINT id_pk PRIMARY KEY(id)
);
