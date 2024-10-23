CREATE TABLE IF NOT EXISTS servers(
   id INT NOT NULL AUTO_INCREMENT,
   name VARCHAR(64) NOT NULL,
   os_code VARCHAR(64) DEFAULT NULL,
   os_version VARCHAR(64) DEFAULT NULL,
   ipv4_address VARCHAR(64) DEFAULT NULL,
   ipv6_address VARCHAR(256) DEFAULT NULL,
   agent_version VARCHAR(64) DEFAULT NULL,
   agent_port INT DEFAULT NULL,
   token VARCHAR(512) NOT NULL,
   is_registered TINYINT NOT NULL DEFAULT 0,
   is_active TINYINT NOT NULL DEFAULT 0,
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
   account_id INT NOT NULL,

   PRIMARY KEY(id),
   INDEX token_index (token),
   INDEX account_id_index (account_id),

   FOREIGN KEY (account_id) REFERENCES accounts(id)
);
