-- Add capabilities column for bb-sites compatibility (e.g. ["network"])
ALTER TABLE scripts ADD COLUMN capabilities TEXT NOT NULL DEFAULT '[]';
