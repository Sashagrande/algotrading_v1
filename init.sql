CREATE TABLE IF NOT EXISTS trades (
                                      id SERIAL PRIMARY KEY,
                                      trade_id VARCHAR(50) NOT NULL,
    trade_price DECIMAL(18, 8) NOT NULL,
    quantity DECIMAL(18, 8) NOT NULL,
    trade_time TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL
    );