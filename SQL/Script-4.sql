CREATE TABLE users (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    age  INT NOT NULL,
    city TEXT NOT NULL
);

CREATE TABLE orders (
    id      BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id),
    item    TEXT NOT NULL,
    price   INT NOT NULL
);
INSERT INTO users (name, age, city) VALUES
    ('Анна',  25, 'Москва'),
    ('Борис', 40, 'Москва'),
    ('Вера',  40, 'Москва'),
    ('Глеб',  31, 'Казань'),
    ('Дина',  19, 'Казань'),
    ('Егор',  27, 'СПб');

INSERT INTO orders (user_id, item, price) VALUES
    (1, 'телефон', 50000),
    (1, 'чехол',   1500),
    (2, 'ноутбук', 120000),
    (4, 'телефон', 45000),
    (4, 'наушники', 12000),
    (6, 'наушники', 20000),
    (6, 'наушники', 18000);
-- у Веры (3) и Дины (5) заказов нет
///////////////////////////////////////////////////////
1 задание
SELECT name,age
FROM users
WHERE (city = 'Москва') AND (age = (SELECT max(age) FROM users WHERE (city = 'Москва')));
2 задание
SELECT city, COUNT(*) AS PeapleCount
FROM users
GROUP BY city
HAVING COUNT(*) > 1
ORDER BY PeapleCount DESC, city;
3 задание
SELECT u.name, COALESCE(sum(o.price), 0) AS TotalPrice
FROM users u
LEFT JOIN orders o 
ON u.id = o.user_id
GROUP BY u.name
ORDER BY TotalPrice DESC;
4 задание
SELECT DISTINCT u.name
FROM users u
JOIN orders o
ON u.id = o.user_id
WHERE o.item = 'наушники'
ORDER BY u.name
5 задание
SELECT item, price
FROM orders
WHERE price = (SELECT min(price) FROM orders WHERE price > 10000) 