CREATE TABLE IF NOT EXISTS certificate_monitor_sites(
   id INT NOT NULL AUTO_INCREMENT,
   url VARCHAR(512) NOT NULL,
   data VARCHAR(4096),
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
   user_id INTEGER NOT NULL,

   PRIMARY KEY(id),
   INDEX user_id_index (user_id),

   FOREIGN KEY (user_id) REFERENCES users(id)
);
