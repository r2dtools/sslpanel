CREATE TABLE IF NOT EXISTS domain_settings(
   id INT NOT NULL AUTO_INCREMENT,
   server_id INT NOT NULL,
   domain_name VARCHAR(64) NOT NULL,
   setting_name VARCHAR(64) NOT NULL,
   setting_value VARCHAR(64) NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

   PRIMARY KEY(id),
   UNIQUE (server_id, domain_name, setting_name),

   FOREIGN KEY (server_id) REFERENCES servers(id)
      ON DELETE CASCADE
);
