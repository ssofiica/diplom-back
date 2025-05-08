CREATE TABLE orders (
    order_id UInt64,
    user_id UInt32,
    status Enum('finished' = 1, 'canceled' = 2),
    type Enum('pickup' = 1, 'delivery' = 2),
    sum Float32,
    restaurant_id UInt32,
    created_at DateTime,
    accepted_at DateTime,
    ready_at DateTime,
    finished_at DateTime,
    canceled_at DateTime
)
ENGINE = MergeTree()
ORDER BY (toDate(created_at), restaurant_id);

CREATE TABLE order_food (
    order_id UInt64,
    food_id UInt32,
    count UInt16,
    food_name String,
    food_price Float32,
    food_weight UInt16,
    restaurant_id UInt32,
    category String,
    category_id UInt32,
    ordered_at DateTime,
    order_status Enum('finished' = 1, 'canceled' = 2)
)
ENGINE = MergeTree()
ORDER BY (ordered_at, restaurant_id);