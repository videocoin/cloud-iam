INSERT INTO permissions
  (name, created_at, updated_at)
VALUES
  ('iam.roles.get', now(), now()),
  ('iam.roles.list', now(), now()),
  ('iam.serviceAccountKeys.create', now(), now()),
  ('iam.serviceAccountKeys.get', now(), now()),
  ('iam.serviceAccountKeys.list', now(), now()),
  ('iam.serviceAccountKeys.delete', now(), now()),
  ('iam.serviceAccounts.create', now(), now()),
  ('iam.serviceAccounts.get', now(), now()),
  ('iam.serviceAccounts.list', now(), now()),
  ('iam.serviceAccounts.delete', now(), now()),
  ('iam.serviceAccounts.setIamPolicy', now(), now()),
  ('iam.serviceAccounts.getIamPolicy', now(), now()),
  ('symphony.blockchain.use', now(), now());
  