-- migrate:up
BEGIN;

alter table directory_entries 
	add column if not exists entry_textsearchable_index_col tsvector generated always as (to_tsvector('english', entry_first_name || ' ' || entry_last_name || ' ' || entry_email)) stored;

COMMIT;

create index if not exists entries_textsearch_idx on directory_entries using gin (entry_textsearchable_index_col);

-- migrate:down
BEGIN;
COMMIT;
