-- migrate:up
alter table contacts
	add column if not exists contact_textsearchable_index_col tsvector generated always as (to_tsvector('english', contact_first_name || ' ' || contact_last_name || ' ' || contact_email)) stored;

create index if not exists contacts_textsearch_idx on contacts using gin (contact_textsearchable_index_col);

-- migrate:down
BEGIN;
COMMIT;
