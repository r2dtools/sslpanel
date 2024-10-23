CREATE TABLE IF NOT EXISTS accounts(
   id INT NOT NULL AUTO_INCREMENT,
   confirmed TINYINT NOT NULL DEFAULT 0,
   confirmation_code INT NOT NULL DEFAULT 0,
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

   PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS users(
   id INT NOT NULL AUTO_INCREMENT,
   email VARCHAR(64) NOT NULL,
   `password` VARCHAR(512) NOT NULL,
   active INT NOT NULL DEFAULT 0,
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
   account_id INT NOT NULL,
   account_owner INT NOT NULL DEFAULT 0,
   confirmation_code INT NOT NULL DEFAULT 0,

   PRIMARY KEY(id),
   INDEX account_id_index (account_id),
   UNIQUE KEY (email),

   FOREIGN KEY (account_id) REFERENCES accounts(id)
);
