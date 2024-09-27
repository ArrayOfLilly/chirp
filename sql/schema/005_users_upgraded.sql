-- +goose Up
-- Step 1: Add the new column as nullable
ALTER TABLE users 
    ADD COLUMN is_chirpy_red BOOLEAN;

-- Step 2: Update existing rows with the default value
UPDATE users 
    SET is_chirpy_red = false 
    WHERE is_chirpy_red IS NULL;

-- Step 3: Make the column non-null
ALTER TABLE users 
    ALTER COLUMN is_chirpy_red 
        SET NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN is_chirpy_red;

