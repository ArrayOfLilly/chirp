-- +goose Up
-- Step 1: Add the new column as nullable
ALTER TABLE users ADD COLUMN hashed_password TEXT;

-- Step 2: Update existing rows with the default value
UPDATE users SET hashed_password = 'unset' WHERE hashed_password IS NULL;

-- Step 3: Make the column non-null
ALTER TABLE users ALTER COLUMN hashed_password SET NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN hashed_password;