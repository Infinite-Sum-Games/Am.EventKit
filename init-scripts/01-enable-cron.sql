-- Connect to the postgres database 
\c postgres

-- Create the extension inside this database
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Verify it is working by logging a message
SELECT cron.schedule('startup-check', '* * * * *', 'SELECT 1');
