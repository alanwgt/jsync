ALTER TABLE properties
  ALTER COLUMN building_status DROP NOT NULL;

ALTER TABLE properties
  ALTER COLUMN notes DROP NOT NULL;

ALTER TABLE properties
  ALTER COLUMN meta_title DROP NOT NULL;

ALTER TABLE properties
  ALTER COLUMN meta_description DROP NOT NULL;

ALTER TABLE properties
  ALTER COLUMN address_reference DROP NOT NULL;

ALTER TABLE properties
  ALTER COLUMN ad_title DROP NOT NULL;

ALTER TABLE properties
  ALTER COLUMN ad_description DROP NOT NULL;

UPDATE properties SET building_status = NULL WHERE building_status = '';
UPDATE properties SET notes = NULL WHERE notes = '';
UPDATE properties SET meta_title = NULL WHERE meta_title = '';
UPDATE properties SET meta_description = NULL WHERE meta_description = '';
UPDATE properties SET address_reference = NULL WHERE address_reference = '';
UPDATE properties SET ad_title = NULL WHERE ad_title = '';
UPDATE properties SET ad_description = NULL WHERE ad_description = '';
