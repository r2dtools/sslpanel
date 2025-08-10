CREATE TABLE IF NOT EXISTS renewal_logs(
   id INT NOT NULL AUTO_INCREMENT,
   server_id INT NOT NULL,
   domain_name VARCHAR(64) NOT NULL,
   error VARCHAR(1024) NOT NULL DEFAULT '',
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),

   PRIMARY KEY(id),
   INDEX server_id_index (server_id),

   FOREIGN KEY (server_id) REFERENCES servers(id)
      ON DELETE CASCADE
);
