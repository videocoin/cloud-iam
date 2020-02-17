INSERT INTO roles
  (name, title, created_at, updated_at)
VALUES
  ('roles/editor', 'Editor', now(), now()),
  ('roles/viewer', 'Viewer', now(), now()),
  ('roles/owner', 'Owner', now(), now());