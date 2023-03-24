CREATE TABLE author (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE book (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    count int NOT NULL,
    income_price NUMERIC NOT NULL,
    profit_status VARCHAR NOT NULL,
    profit_price NUMERIC NOT NULL,
    sell_price NUMERIC DEFAULT 0,
    author_id VARCHAR REFERENCES author(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE "user" (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    balance NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE OR REPLACE FUNCTION sell_price() RETURNS TRIGGER LANGUAGE PLPGSQL
AS
$$
BEGIN
    if new.profit_status = 'fixed' THEN
        new.sell_price = new.income_price + new.profit_price;
    ELSIF new.profit_status = 'percent' THEN
        new.sell_price = new.income_price + (new.income_price * new.profit_price / 100);
    END IF;
    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER sell_price_tg
BEFORE INSERT OR UPDATE ON book
FOR EACH ROW EXECUTE PROCEDURE sell_price();