-- Create additional databases for FarmAgent services
-- This script runs on first postgres container startup

CREATE DATABASE farmagent_crops;
CREATE DATABASE farmagent_payments;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE farmagent_crops TO farmagent;
GRANT ALL PRIVILEGES ON DATABASE farmagent_payments TO farmagent;
