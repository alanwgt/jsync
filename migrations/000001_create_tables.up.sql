CREATE TABLE IF NOT EXISTS brokers
(
  id                         INT PRIMARY KEY,
  tenant_id                  INT    NULL,
  "name"                     TEXT   NOT NULL,
  avatar                     TEXT   NULL,
  creci                      TEXT   NOT NULL,
  biography                  TEXT   NULL,
  job_position               TEXT   NULL,
  birthday                   DATE   NULL,
  city_name                  TEXT   NULL,
  state_abbrev               TEXT   NULL,
  is_broker                  BOOL   NOT NULL,
  email                      TEXT   NULL,
  primary_phone_number       TEXT   NULL,
  alt_phone_number           TEXT   NULL,
  primary_phone_has_whatsapp TEXT   NULL,
  alt_phone_has_whatsapp     TEXT   NULL,
  twitter_url                TEXT   NULL,
  facebook_url               TEXT   NULL,
  linkedin_url               TEXT   NULL,
  skype_url                  TEXT   NULL,
  instagram_url              TEXT   NULL,
  personal_website_url       TEXT   NULL,
  show_personal_website      BOOL   NOT NULL,
  testimonials               TEXT[] NOT NULL
);

CREATE INDEX brokers_tenant_id_idx ON brokers (tenant_id) WHERE tenant_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS banners
(
  id          INT PRIMARY KEY,
  image_id    INT  NOT NULL,
  tenant_id   INT  NULL,
  url         TEXT NULL,
  "order"     INT  NOT NULL,
  title       TEXT NULL,
  description TEXT NULL,
  href_target TEXT,
  is_video    BOOL NOT NULL,
  video_url   TEXT NULL,
  image_url   TEXT NULL
);

CREATE INDEX banners_tenant_id_idx ON banners (tenant_id) WHERE tenant_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS condominiums
(
  id                         INT PRIMARY KEY,
  type                       TEXT           NOT NULL,
  cover_image_id             INT            NULL,
  tenant_id                  INT            NULL,
  name                       TEXT           NOT NULL,
  featured                   BOOL           NOT NULL,
  launch                     BOOL           NOT NULL,
  gated                      BOOL           NOT NULL,
  construction_percentage    INT            NULL,
  brick_percentage           INT            NULL,
  structure_percentage       INT            NULL,
  foundation_percentage      INT            NULL,
  installations_percentage   INT            NULL,
  landscaping_percentage     INT            NULL,
  project_percentage         INT            NULL,
  ground_leveling_percentage INT            NULL,
  latitude                   NUMERIC        NULL,
  longitude                  NUMERIC        NULL,
  notes                      TEXT           NOT NULL,
  incorporation_record       TEXT           NULL,
  neighborhood_id            INT            NOT NULL,
  city_id                    INT            NOT NULL,
  state_id                   INT            NOT NULL,
  address_zipcode            TEXT           NULL,
  address_street             TEXT           NOT NULL,
  address_neighborhood       TEXT           NOT NULL,
  address_number             TEXT           NOT NULL,
  city_name                  TEXT           NOT NULL,
  state_name                 TEXT           NOT NULL,
  status                     TEXT           NOT NULL,
  delivery_month             INT            NULL,
  delivery_year              INT            NULL,
  administering_company_name TEXT           NULL,
  building_company_name      TEXT           NULL,
  real_state_developer_name  TEXT           NULL,
  architect_name             TEXT           NULL,
  landscaper_name            TEXT           NULL,
  decorator_name             TEXT           NULL,
  cover_image                TEXT           NOT NULL,
  available_properties       INT            NOT NULL,
  infrastructures            TEXT[]         NOT NULL,
  labels                     TEXT[]         NOT NULL,
  videos                     JSONB          NOT NULL,
  images                     JSONB          NOT NULL,
  blueprints                 JSONB          NOT NULL,
  ar_tour                    TEXT[]         NOT NULL,
  created_at                 TIMESTAMPTZ(3) NOT NULL,
  updated_at                 TIMESTAMPTZ(3) NOT NULL
  );

CREATE INDEX condominiums_tenant_id_idx ON condominiums (tenant_id) WHERE tenant_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS properties
(
  id                          INT PRIMARY KEY,
  broker_id                   INT            NOT NULL,
  condominium_id              INT            NULL,
  state_id                    INT            NOT NULL,
  city_id                     INT            NOT NULL,
  neighborhood_id             INT            NOT NULL,
  tenant_id                   INT            NULL,
  identifier_code             TEXT           NOT NULL,
  contracts                   TEXT[]         NOT NULL,
  type                        TEXT           NOT NULL,
  active                      BOOL           NOT NULL DEFAULT false,
  subtype                     TEXT           NOT NULL,
  notes                       TEXT           NOT NULL,
  building_type               TEXT           NULL,
  delivery_year               INT            NULL,
  delivery_month              INT            NULL,
  furnished                   INT            NOT NULL,
  suites                      INT            NOT NULL,
  bathrooms                   INT            NOT NULL,
  bedrooms                    INT            NOT NULL,
  garages                     INT            NOT NULL,
  financeable                 INT            NOT NULL,
  has_exclusivity             BOOL           NOT NULL,
  total_area                  NUMERIC(30, 2),
  private_area                NUMERIC(30, 2),
  useful_area                 NUMERIC(30, 2),
  measurement_type            TEXT           NOT NULL,
  floor_types                 TEXT[]         NOT NULL,
  terrain_width_front         NUMERIC(30, 2) NULL,
  terrain_width_back          NUMERIC(30, 2) NULL,
  terrain_length_left         NUMERIC(30, 2) NULL,
  terrain_length_right        NUMERIC(30, 2) NULL,
  terrain_area                NUMERIC(30, 2) NULL,
  created_at                  TIMESTAMPTZ(3),
  building_status             TEXT           NOT NULL,
  show_condominium_value      BOOL           NOT NULL,
  condominium_value           NUMERIC        NULL,
  show_sale_value             BOOLEAN        NOT NULL,
  sale_value                  NUMERIC        NULL,
  show_rental_value           BOOLEAN        NOT NULL,
  rental_value                NUMERIC        NULL,
  show_seasonal_value         BOOLEAN        NOT NULL,
  seasonal_value              NUMERIC        NULL,
  iptu_frequency              TEXT           NULL,
  iptu_value_exempt           TEXT,
  show_iptu_value             BOOLEAN        NOT NULL,
  iptu_value                  NUMERIC        NULL,
  season_calendar             JSONB          NOT NULL,
  rural_activities            TEXT[]         NOT NULL,
  rural_headquarters          INT            NULL,
  arable_area                 TEXT           NULL,
  allowed_guests              INT            NULL,
  meta_title                  TEXT           NOT NULL,
  meta_description            TEXT           NOT NULL,
  fire_insurance_value        NUMERIC        NULL,
  cleaning_fee_value          NUMERIC        NULL,
  position                    TEXT           NULL,
  solar_positions             TEXT[]         NOT NULL,
  sea_distance                INT            NULL,
  accept_exchange             BOOL           NOT NULL,
  latitude                    NUMERIC        NULL,
  longitude                   NUMERIC        NULL,
  occupancy_status            TEXT           NOT NULL,
  featured                    TEXT           NOT NULL,
  feature_until               TIMESTAMPTZ(3) NULL,
  condominium_type            TEXT           NULL,
  condominium_name            TEXT           NULL,
  gated_condominium           BOOL           NULL,
  show_full_address           BOOL           NOT NULL,
  show_address_state          BOOL           NOT NULL,
  show_address_city           BOOL           NOT NULL,
  show_address_neighborhood   BOOL           NOT NULL,
  show_address_street         BOOL           NOT NULL,
  show_address_reference      BOOL           NOT NULL,
  show_address_number         BOOL           NOT NULL,
  show_address_floor          BOOL           NOT NULL,
  address_state               TEXT           NOT NULL,
  address_city                TEXT           NOT NULL,
  address_neighborhood        TEXT           NOT NULL,
  address_street              TEXT           NOT NULL,
  address_zipcode             TEXT           NULL,
  address_reference           TEXT           NULL,
  address_number              TEXT           NOT NULL,
  address_floor               INT            NULL,
  geoposition_visibility      INT            NOT NULL,
  ad_title                    TEXT           NOT NULL,
  ad_description              TEXT           NOT NULL,
  labels                      TEXT[]         NOT NULL,
  surety_insurance            BOOL           NULL,
  property_infrastructures    TEXT[]         NOT NULL,
  condominium_infrastructures TEXT[]         NOT NULL,
  updated_at                  TIMESTAMPTZ(3) NOT NULL,
  validated_at                TIMESTAMPTZ(3) NOT NULL,
  videos                      JSONB          NOT NULL,
  blueprints                  JSONB          NOT NULL,
  images                      JSONB          NOT NULL,
  ar_tour                     TEXT[]         NOT NULL
  );

CREATE UNIQUE INDEX properties_unique_identifier_code_idx ON properties (identifier_code);
CREATE INDEX properties_tenant_id_idx ON properties (tenant_id) WHERE tenant_id IS NOT NULL;
