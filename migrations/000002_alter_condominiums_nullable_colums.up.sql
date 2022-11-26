ALTER TABLE condominiums
  ALTER COLUMN notes DROP NOT NULL;

ALTER TABLE condominiums
  ALTER COLUMN cover_image DROP NOT NULL;

UPDATE condominiums SET notes = NULL WHERE notes = '';
UPDATE condominiums SET cover_image = NULL WHERE cover_image = '';
