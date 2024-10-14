CREATE TABLE IF NOT EXISTS certificate_monitor_sites(
   id serial PRIMARY KEY,
   url VARCHAR(512) NOT NULL,
   data VARCHAR(4096),
   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
   user_id INTEGER NOT NULL REFERENCES users
);
