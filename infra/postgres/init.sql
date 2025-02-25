CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS restaurant (
    id            INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name          TEXT
        CONSTRAINT rest_name_length CHECK (LENGTH(name) BETWEEN 2 AND 100)  NOT NULL,
    address TEXT DEFAULT '',
    logo_url TEXT DEFAULT '',
    description_array TEXT[] DEFAULT ARRAY[]::TEXT[],
    img_urls TEXT[] DEFAULT ARRAY[]::TEXT[],
    phone TEXT DEFAULT '',
    email TEXT NOT NULL,
    media_links JSONB
);

CREATE TABLE IF NOT EXISTS schedule (
    id            INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restaurant_id INTEGER REFERENCES restaurant (id) ON DELETE CASCADE,
    day varchar(20) NOT NULL,
    open_time TIME NOT NULL,
    close_time TIME NOT NULL
);

CREATE TABLE IF NOT EXISTS category (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    restaurant_id INTEGER
        CONSTRAINT foreign_key_rest CHECK (restaurant_id > 0) REFERENCES restaurant (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS food (
    id            INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name          TEXT
        CONSTRAINT food_name_length CHECK (LENGTH(name) BETWEEN 2 AND 150)  NOT NULL,
    restaurant_id INTEGER
        CONSTRAINT foreign_key_rest CHECK (restaurant_id > 0) REFERENCES restaurant (id) ON DELETE CASCADE,
    category_id   INTEGER
        CONSTRAINT foreign_key_cat CHECK (category_id > 0) REFERENCES category (id) ON DELETE SET NULL,
    weight        INTEGER
        CONSTRAINT positive_weight CHECK (weight > 0) NOT NULL,
    price         INTEGER
        CONSTRAINT positive_price CHECK (price > 0) NOT NULL,
    img_url       TEXT
        CONSTRAINT restaurant_img_url CHECK (LENGTH(img_url) <= 60) NOT NULL,
    status        TEXT 
        CONSTRAINT food_status_length CHECK (LENGTH(status) <= 20) NOT NULL
);


insert into restaurant(name, email) values ('Mates', 'mates@yandex.ru')