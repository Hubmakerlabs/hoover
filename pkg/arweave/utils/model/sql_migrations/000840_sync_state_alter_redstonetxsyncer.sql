-- +migrate Down

-- +migrate Up
ALTER TYPE synced_component RENAME VALUE 'WarpySyncer' TO 'WarpySyncerAvax';
