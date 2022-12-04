UPDATE properties SET building_status = '' WHERE building_status IS NULL;
UPDATE properties SET notes = '' WHERE notes IS NULL;
UPDATE properties SET meta_title = '' WHERE meta_title IS NULL;
UPDATE properties SET meta_description = '' WHERE meta_description IS NULL;
UPDATE properties SET address_reference = '' WHERE address_reference IS NULL;
UPDATE properties SET ad_title = '' WHERE ad_title IS NULL;
UPDATE properties SET ad_description = '' WHERE ad_description IS NULL;

ALTER TABLE properties
  ALTER COLUMN building_status SET NOT NULL;

ALTER TABLE properties
  ALTER COLUMN notes SET NOT NULL;

ALTER TABLE properties
  ALTER COLUMN meta_title SET NOT NULL;

ALTER TABLE properties
  ALTER COLUMN meta_description SET NOT NULL;

ALTER TABLE properties
  ALTER COLUMN address_reference SET NOT NULL;

ALTER TABLE properties
  ALTER COLUMN ad_title SET NOT NULL;

ALTER TABLE properties
  ALTER COLUMN ad_description SET NOT NULL;
