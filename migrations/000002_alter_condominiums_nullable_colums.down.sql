ALTER TABLE condominiums
  ALTER COLUMN notes SET NOT NULL;

ALTER TABLE condominiums
  ALTER COLUMN cover_image SET NOT NULL;

UPDATE condominiums SET notes = '' WHERE notes IS NULL;
UPDATE condominiums SET cover_image = '' WHERE cover_image IS NULL;
