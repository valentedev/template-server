Create image
docker build -t books_img .

Run the docker container
docker run -d --name books_ctr -p 5432:5432 books_img

Open new terminal to access the container
docker exec -it books_ctr bash

Once inside the container, enter the USER and DATABASE provided on dockerfile
psql -U admin -d books

CREATE TABLE books (
   id BIGSERIAL NOT NULL PRIMARY KEY,
   created timestamp(0) with time zone NOT NULL DEFAULT NOW(),
   title VARCHAR(100) NOT NULL,
   author VARCHAR(100) NOT NULL
);

INSERT INTO books (title, author) VALUES ('Harry Potter and the Chamber of Secrets', 'J.K.Rowling'); 

INSERT INTO books (title, author) VALUES ('The Lost Symbol', 'Dan Brown'); 

INSERT INTO books (title, author) VALUES ('Les Miserables', 'Victor Hugo'); 

CREATE TABLE sessions (
   token VARCHAR(43) PRIMARY KEY,
   data bytea NOT NULL,
   expiry timestamp(6) NOT NULL
)

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password VARCHAR(60) NOT NULL,
    created timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

ALTER TABLE users
ADD CONSTRAINT users_uc_email UNIQUE (email);
