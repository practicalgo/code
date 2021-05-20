use package_server;

CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(30) NOT NULL
);

CREATE TABLE packages(
    owner_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    version VARCHAR(50) NOT NULL,
    object_store_id VARCHAR(300) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (owner_id, name, version),
    FOREIGN KEY (owner_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);