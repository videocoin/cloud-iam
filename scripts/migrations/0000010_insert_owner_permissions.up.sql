INSERT INTO roles_permissions
  (role_name, permission_name)
VALUES
  ('roles/owner', 'iam.roles.get'),
  ('roles/owner', 'iam.roles.list'),
  ('roles/owner', 'iam.serviceAccountKeys.create'),
  ('roles/owner', 'iam.serviceAccountKeys.get'),
  ('roles/owner', 'iam.serviceAccountKeys.list'),
  ('roles/owner', 'iam.serviceAccountKeys.delete'),
  ('roles/owner', 'iam.serviceAccounts.create'),
  ('roles/owner', 'iam.serviceAccounts.get'),
  ('roles/owner', 'iam.serviceAccounts.list'),
  ('roles/owner', 'iam.serviceAccounts.delete'),
  ('roles/owner', 'iam.serviceAccounts.setIamPolicy'),
  ('roles/owner', 'iam.serviceAccounts.getIamPolicy'),
  ('roles/owner', 'symphony.blockchain.use');