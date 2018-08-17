% CREATE entropy DATABASE
CREATE DATABASE entropy;

% Create Migration User
CREATE USER 'migrate-api'@'%' IDENTIFIED BY '1e958331-b0b7-4464-b5e2-79f92042cdae';
GRANT ALTER, CREATE, DROP ON entropy.* TO 'migrate-api'@'%';
SHOW GRANTS FOR 'migrate-api'@'%';

% Create Entropy User
CREATE USER 'entropy-api'@'%' IDENTIFIED BY '3ae651ac-490f-4e7f-a693-f558648e1135';
GRANT SELECT, INSERT, UPDATE ON entropy.* TO 'entropy-api'@'%';
SHOW GRANTS FOR 'entropy-api'@'%';
