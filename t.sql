-- Active: 1727140703105@@vividly-climbing-mallard.data-1.use1.tembo.io@5432@armazem_DB@public
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE books(  
    book_id CHAR(8) PRIMARY KEY DEFAULT substring(md5(random()::text) FROM 1 FOR 8),
    title VARCHAR(200) NOT NULL,
    synopsis VARCHAR(500) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    amount INT DEFAULT 0,
    genre_id INT REFERENCES generos(gene_id),
    author_id INT REFERENCES author(author_id)
);

CREATE TABLE authors(
  author_id CHAR(8) PRIMARY KEY DEFAULT substring(md5(random()::text) FROM 1 FOR 8),
  author_name VARCHAR(100) NOT NULL,
  author_nationality VARCHAR(100)
);

CREATE TABLE genres(
  genre_id CHAR(8) PRIMARY KEY DEFAULT substring(md5(random()::text) FROM 1 FOR 8),
  genre_name VARCHAR(100) NOT NULL
);

CREATE TABLE genre_book (
  fk_book_id INT REFERENCES book(book_id),
  fk_genre_id INT REFERENCES genres(genre_id),
  PRIMARY KEY (fk_book_id,fk_genre_id)
);

CREATE VIEW book_genre_view AS
SELECT 
  b.book_id,
  b.title,
  b.synopsis,
  g.genre_name
FROM 
  book b
JOIN
  genre_book gb ON b.book_id = gb.fk_book_id
JOIN 
  genres g ON gb.fk_genre_id = g.genre_id;

CREATE OR REPLACE PROCEDURE add_book_with_genre(
  book_title VARCHAR,
  book_synopsis VARCHAR,
  book_price DECIMAL,
  book_amount INT,
  genre_id INT 
)
LANGUAGE plpgsql
AS $$
DECLARE
  new_book_id CHAR(8);
BEGIN 

  new_book_id := substring(md5(random()::text)FROM 1 for 8)

  INSERT INTO book (book_id, book_title, book_price, book_synopsis,)
  VALUES (new_book_id, book_title, book_price, book_amount, genre_id)
  RETURNING book_id INTO new_book_id;

  INSERT INTO genre_book (fk_book_id, fk_genre_id)
  VALUES (new_book_id, genre_id);
END;
$$;