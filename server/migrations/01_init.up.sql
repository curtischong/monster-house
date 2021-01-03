CREATE TABLE photo (
    id        uuid PRIMARY KEY NOT NULL,
    name      text             NOT NULL,
    extension text             NOT NULL
);

CREATE TABLE tag (
    id uuid PRIMARY KEY NOT NULL,
    name text UNIQUE NOT NULL
);
CREATE INDEX tag_name_index ON tag USING btree (name);

-- This table is used to link associate photos with tags
CREATE TABLE photo_tag (
    photo_id uuid NOT NULL,
    tag_id uuid NOT NULL,
    is_auto_generated bool NOT NULL -- is false when the user specifies this tag for this photo
);
