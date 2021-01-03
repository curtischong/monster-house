CREATE TABLE tag (
    id uuid PRIMARY KEY NOT NULL,
    photo_id uuid NOT NULL,
    is_auto_generated bool NOT NULL,
    name text NOT NULL
) ;
CREATE INDEX tag_name_index ON tag USING btree (name);