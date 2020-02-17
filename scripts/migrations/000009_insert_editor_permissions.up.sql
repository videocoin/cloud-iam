INSERT INTO roles_permissions
  (role_name, permission_name)
VALUES
  ('roles/editor', 'iam.roles.get'),
  ('roles/editor', 'iam.roles.list'),
  ('roles/editor', 'iam.serviceAccountKeys.create'),
  ('roles/editor', 'iam.serviceAccountKeys.get'),
  ('roles/editor', 'iam.serviceAccountKeys.list'),
  ('roles/editor', 'iam.serviceAccountKeys.delete'),
  ('roles/editor', 'iam.serviceAccounts.create'),
  ('roles/editor', 'iam.serviceAccounts.get'),
  ('roles/editor', 'iam.serviceAccounts.list'),
  ('roles/editor', 'iam.serviceAccounts.delete'),
  ('roles/editor', 'symphony.blockchain.use');