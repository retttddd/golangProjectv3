CREATE TABLE IF NOT EXISTS secret_data  (
    key_data VARCHAR(5000) UNIQUE NOT NULL,
    value_data VARCHAR(5000) NOT NULL
    );