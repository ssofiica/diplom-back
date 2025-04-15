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

--TODO: доделать
CREATE TABLE IF NOT EXISTS "user" (
    id            INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT,
    phone TEXT NOT NULL,
    password TEXT
);

CREATE TABLE IF NOT EXISTS "order" (
    id            INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER
        CONSTRAINT foreign_key_user CHECK (user_id > 0) REFERENCES "user" (id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    address TEXT CONSTRAINT order_address_length CHECK (LENGTH(address) BETWEEN 2 AND 200),
    sum INTEGER NOT NULL DEFAULT 0,
    restaurant_id INTEGER
        CONSTRAINT foreign_key_rest CHECK (restaurant_id > 0) REFERENCES restaurant (id) ON DELETE CASCADE,
    comment TEXT CONSTRAINT order_comment_length CHECK (LENGTH(comment) BETWEEN 2 AND 256),
    "type" TEXT,
    created_at TIME NOT NULL DEFAULT NOW(),
    accepted_at TIME,
    ready_at TIME,
    finished_at TIME,
    canceled_at TIME
);

CREATE TABLE IF NOT EXISTS order_food (
    order_id INTEGER
        CONSTRAINT foreign_key_order CHECK (order_id > 0) REFERENCES "order" (id) ON DELETE CASCADE,
    food_id INTEGER
        CONSTRAINT foreign_key_food CHECK (food_id > 0) REFERENCES food (id) ON DELETE CASCADE,
    count      INTEGER
        CONSTRAINT food_count_in_order CHECK (count > 0) NOT NULL,
    PRIMARY KEY (food_id, order_id)
);

insert into restaurant(name, email) values ('Mates', 'mates@yandex.ru');
insert into category(name, restaurant_id) values ('Пицца', 1), ('Салаты', 1), ('Паста',1);
insert into food(name, weight, category_id, price, restaurant_id, img_url, status) values 
('Пицца Маргарита', 370, 1, 500, 1, 'food/1/1.jpg', 'in'),
('Буррата со страчателой', 400, 2, 500, 1, 'food/1/2.jpg', 'in'),
('С креветками', 340, 3, 600, 1, 'food/1/3.jpg', 'in'),
('Карбонара', 400, 3, 700, 1, 'food/1/4.jpg', 'in'),
('Лазанья', 380, 3, 670, 1, 'food/1/5.jpg', 'in');
insert into "user"(name, phone) values ('sofia', '89009009090');