-- Migration 019: Seed Test Menghafal module and questions

-- Insert Module
INSERT INTO modules (id, name, description, total_weight, created_by, memorization_content, memorization_duration) VALUES
  ('307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'Test Menghafal', 'Test untuk mengukur daya ingat kandidat', 20.0, '11111111-1111-1111-1111-111111111111', 
   'BUNGA : SOKA, LARAT, FLAMBOYAN, YASMIN, DAHLIA\nPERKAKAS : WAJAN, JARUM, KIKIR, CANGKUL, PALU\nBURUNG : ITIK, ELANG, WALET, TEKUKUR, NURI\nKESENIAN : QUINTET, ARCA, OPERA, UKIRAN, GAMELAN\nBINATANG : RUSA, MUSANG, BERUANG, HARIMAU, ZEBRA', 200)
ON CONFLICT (id) DO NOTHING;

-- Insert Questions and Options
-- Question: Kata yang mempunyai huruf permulaan – M – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('11c8a571-e8b3-45e2-8569-3cfceebb3180', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – M – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('99974f94-fd95-4f7f-af1b-ed2ddd9a59b3', '11c8a571-e8b3-45e2-8569-3cfceebb3180', 'kesenian', FALSE),
  ('e932ca88-abd9-46fb-a05e-c1e71ed0b22e', '11c8a571-e8b3-45e2-8569-3cfceebb3180', 'bunga', TRUE),
  ('8d692635-20e7-4c67-804f-2ddee1088e76', '11c8a571-e8b3-45e2-8569-3cfceebb3180', 'perkakas', FALSE),
  ('1c6abaa2-d127-4ba4-b28e-405d297796fc', '11c8a571-e8b3-45e2-8569-3cfceebb3180', 'burung', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – C – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('140288c1-ffb9-408e-831d-99603d36a2a6', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – C – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('a4d1cae5-f342-4fc0-a3a7-b40fd33b77c1', '140288c1-ffb9-408e-831d-99603d36a2a6', 'bunga', FALSE),
  ('291b0dd6-793d-4daa-9c4b-1f9a6a0e2edd', '140288c1-ffb9-408e-831d-99603d36a2a6', 'perkakas', TRUE),
  ('cd31a7d4-4dbb-478a-903f-4d869fc315eb', '140288c1-ffb9-408e-831d-99603d36a2a6', 'burung', FALSE),
  ('35fff50d-393a-4496-afc0-3a743ce9c59c', '140288c1-ffb9-408e-831d-99603d36a2a6', 'kesenian', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – K – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('1dc73240-b1a4-4f6a-b41b-7f2b8af43fcc', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – K – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('4698912a-da24-4606-8ebc-5f03a079176d', '1dc73240-b1a4-4f6a-b41b-7f2b8af43fcc', 'burung', FALSE),
  ('eb8fd113-5a7e-4add-97c8-87cd4e9aeb58', '1dc73240-b1a4-4f6a-b41b-7f2b8af43fcc', 'kesenian', FALSE),
  ('2083a3d1-0cca-4610-baa1-1932351dda3d', '1dc73240-b1a4-4f6a-b41b-7f2b8af43fcc', 'bunga', FALSE),
  ('6564a2c6-2795-4153-8870-37252bf59e01', '1dc73240-b1a4-4f6a-b41b-7f2b8af43fcc', 'perkakas', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – L – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('1ed46ac9-8b3f-451b-8ce9-7115f385b75d', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – L – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('3f2c0c91-8829-4d5f-8ea4-a11c35ac846a', '1ed46ac9-8b3f-451b-8ce9-7115f385b75d', 'kesenian', FALSE),
  ('b84754ba-ba57-4bb7-a807-f3f724eb7daf', '1ed46ac9-8b3f-451b-8ce9-7115f385b75d', 'bunga', FALSE),
  ('bfcfff43-3f00-4c2e-8f5f-aeb49f0dad67', '1ed46ac9-8b3f-451b-8ce9-7115f385b75d', 'perkakas', TRUE),
  ('b2e490f2-33f9-4f40-9acc-a948a0d123d7', '1ed46ac9-8b3f-451b-8ce9-7115f385b75d', 'burung', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – U – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('3165bb31-e2e9-4129-a9e1-d4b0a3fbce03', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – U – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('9980d7df-d5e3-4882-b5b5-e1f640856d6c', '3165bb31-e2e9-4129-a9e1-d4b0a3fbce03', 'perkakas', FALSE),
  ('fa7052fa-dff5-47fa-881d-9d5f662b10db', '3165bb31-e2e9-4129-a9e1-d4b0a3fbce03', 'kesenian', FALSE),
  ('c8be6a27-455a-4e14-9b18-efc0f472d5ac', '3165bb31-e2e9-4129-a9e1-d4b0a3fbce03', 'bunga', FALSE),
  ('9ae3d12d-24fd-49a6-8815-7fbe872800de', '3165bb31-e2e9-4129-a9e1-d4b0a3fbce03', 'burung', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – T – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('3fa3db20-ac09-4ff4-93ab-773a5c1456b6', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – T – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('8262db45-b998-4e7e-b972-c3b724c73eca', '3fa3db20-ac09-4ff4-93ab-773a5c1456b6', 'kesenian', FALSE),
  ('27ed71dc-eb2c-4acd-8ff2-30e0de0b991a', '3fa3db20-ac09-4ff4-93ab-773a5c1456b6', 'burung', FALSE),
  ('8bca49a6-56c5-4fbd-a9e2-262e93879278', '3fa3db20-ac09-4ff4-93ab-773a5c1456b6', 'perkakas', FALSE),
  ('a745f407-645a-4a6d-aaec-c83e044216c7', '3fa3db20-ac09-4ff4-93ab-773a5c1456b6', 'bunga', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – E – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('40d38061-8c1c-45c3-bd8d-08210c461fa0', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – E – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('e2872a27-56fd-4e00-9c30-499bbce3509b', '40d38061-8c1c-45c3-bd8d-08210c461fa0', 'bunga', FALSE),
  ('e1565ae3-f229-4216-a895-55462465f27c', '40d38061-8c1c-45c3-bd8d-08210c461fa0', 'perkakas', FALSE),
  ('44be2eea-6c73-4fcb-8f84-b9ecca5ad5e9', '40d38061-8c1c-45c3-bd8d-08210c461fa0', 'burung', TRUE),
  ('7fd73e3f-3276-406c-a570-762042fb4d8d', '40d38061-8c1c-45c3-bd8d-08210c461fa0', 'kesenian', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – O – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('44424d3c-6268-44a3-b356-5ff3b03ecb7e', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – O – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('f04c81b4-4c54-49df-8a86-99238833314a', '44424d3c-6268-44a3-b356-5ff3b03ecb7e', 'burung', TRUE),
  ('28881b25-b98d-4073-afde-194a34ca3b27', '44424d3c-6268-44a3-b356-5ff3b03ecb7e', 'kesenian', FALSE),
  ('1edfa0f3-39e3-4149-9432-dd291ba04739', '44424d3c-6268-44a3-b356-5ff3b03ecb7e', 'perkakas', FALSE),
  ('dd9e5046-2d38-4259-b919-05fdee8522be', '44424d3c-6268-44a3-b356-5ff3b03ecb7e', 'bunga', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – S – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('50861603-3605-4c55-931d-fc439f27b4f7', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – S – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('d6098145-2fb6-4921-b9d9-e4892d29fe8a', '50861603-3605-4c55-931d-fc439f27b4f7', 'kesenian', FALSE),
  ('d16911a3-6224-4f1d-9bb8-6b9eb7823697', '50861603-3605-4c55-931d-fc439f27b4f7', 'bunga', FALSE),
  ('86cafc2e-a5b2-41ed-a9da-d5c98708b52b', '50861603-3605-4c55-931d-fc439f27b4f7', 'perkakas', FALSE),
  ('90e1841d-9873-4ab9-86a7-4594ffb50201', '50861603-3605-4c55-931d-fc439f27b4f7', 'binatang', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – B – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('73a247f1-0807-45f7-91b7-e110c7248eb9', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – B – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('1382d8fa-0b0f-441e-bea6-2dc13c9b3501', '73a247f1-0807-45f7-91b7-e110c7248eb9', 'bunga', FALSE),
  ('b39daa2b-00b9-4d17-8f45-7a9d423ea650', '73a247f1-0807-45f7-91b7-e110c7248eb9', 'binatang', TRUE),
  ('708a619b-1623-40e4-9751-ddc5d6577e03', '73a247f1-0807-45f7-91b7-e110c7248eb9', 'burung', FALSE),
  ('51abd6fb-bf4e-4f3c-91c3-5a46c5a4564d', '73a247f1-0807-45f7-91b7-e110c7248eb9', 'perkakas', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – A – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('772e267d-6405-4465-9cd3-1a261ee0eec8', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – A – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('4ba697cc-2964-44f3-a56d-91f5fddb8717', '772e267d-6405-4465-9cd3-1a261ee0eec8', 'bunga', FALSE),
  ('115a0330-bc30-482d-b6b4-69eed98c1e45', '772e267d-6405-4465-9cd3-1a261ee0eec8', 'perkakas', FALSE),
  ('59bf4a96-a0cb-46cd-9008-8cc3ad59b28b', '772e267d-6405-4465-9cd3-1a261ee0eec8', 'burung', FALSE),
  ('5b50252f-90cf-4d5e-9ed2-4d94a205eafd', '772e267d-6405-4465-9cd3-1a261ee0eec8', 'kesenian', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – D – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('8309ccbf-d98c-40fc-8eef-963221d3d723', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – D – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('aeb30283-10b0-4df1-ba64-4c13b62ba294', '8309ccbf-d98c-40fc-8eef-963221d3d723', 'kesenian', FALSE),
  ('902b0216-9d2b-4e12-9daf-2c24ac08606e', '8309ccbf-d98c-40fc-8eef-963221d3d723', 'bunga', TRUE),
  ('d0a95e22-c0e5-429f-9196-1c28a78edf2c', '8309ccbf-d98c-40fc-8eef-963221d3d723', 'perkakas', FALSE),
  ('d9c937ef-4481-4393-a364-406d327b396b', '8309ccbf-d98c-40fc-8eef-963221d3d723', 'burung', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – I – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('8c5e5c77-fcd7-4f03-95f8-6cf768d856b1', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – I – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c6f0d3e6-0046-4924-b304-92c27858c5c8', '8c5e5c77-fcd7-4f03-95f8-6cf768d856b1', 'perkakas', FALSE),
  ('bf75d9d9-6a36-4519-be48-cdbd94e2c3fc', '8c5e5c77-fcd7-4f03-95f8-6cf768d856b1', 'binatang', TRUE),
  ('dcf874d5-5d68-4f64-8eb4-2d80602cb296', '8c5e5c77-fcd7-4f03-95f8-6cf768d856b1', 'burung', FALSE),
  ('7a5a530b-0df3-4803-bd2a-5b41c1902c91', '8c5e5c77-fcd7-4f03-95f8-6cf768d856b1', 'bunga', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – F – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('9a60f310-9f88-40fa-98b6-225ef327b600', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – F – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('ceb9ff47-0f04-4ce5-adc7-2e5baf9b8d11', '9a60f310-9f88-40fa-98b6-225ef327b600', 'kesenian', FALSE),
  ('2b5edd75-68a2-4931-b597-11b7fb924088', '9a60f310-9f88-40fa-98b6-225ef327b600', 'burung', FALSE),
  ('ac73af6c-1044-4446-b664-10fcb6bc7385', '9a60f310-9f88-40fa-98b6-225ef327b600', 'perkakas', FALSE),
  ('a1ced027-c29e-4726-b184-4cd08e9fdcb8', '9a60f310-9f88-40fa-98b6-225ef327b600', 'bunga', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – G – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('c0287ce5-f71e-4476-939a-d8e3b667e7a7', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – G – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('d2087252-4b4d-4221-a17a-21113659ba97', 'c0287ce5-f71e-4476-939a-d8e3b667e7a7', 'bunga', FALSE),
  ('1963bb6e-be45-4e0c-84fa-f85f4b3f8f4e', 'c0287ce5-f71e-4476-939a-d8e3b667e7a7', 'kesenian', TRUE),
  ('cf1a0554-5290-4ced-a0cf-460dc89f8d82', 'c0287ce5-f71e-4476-939a-d8e3b667e7a7', 'burung', FALSE),
  ('7e5a4f51-fcb5-4e90-9866-1e2874e8b41f', 'c0287ce5-f71e-4476-939a-d8e3b667e7a7', 'perkakas', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – R – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('c29f5b08-dcb8-47d4-b7a1-68b0a51d8963', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – R – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('66444bea-046e-43e6-8699-239928ad9029', 'c29f5b08-dcb8-47d4-b7a1-68b0a51d8963', 'bunga', FALSE),
  ('00b6de65-db05-4f87-8e05-8fce8544cfc5', 'c29f5b08-dcb8-47d4-b7a1-68b0a51d8963', 'perkakas', TRUE),
  ('659f6a29-8d01-4fd2-ac27-340e0faa5770', 'c29f5b08-dcb8-47d4-b7a1-68b0a51d8963', 'burung', FALSE),
  ('11cdecb2-ff61-496f-bb43-58c21790d0d5', 'c29f5b08-dcb8-47d4-b7a1-68b0a51d8963', 'binatang', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – P – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('d93f96c7-1a6c-439e-8aba-0ab3b73b4c43', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – P – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('79c69374-c40c-4f8e-92e3-e1500fae9cf2', 'd93f96c7-1a6c-439e-8aba-0ab3b73b4c43', 'bunga', FALSE),
  ('f44d2e08-b075-4187-8d05-32d1b2fec671', 'd93f96c7-1a6c-439e-8aba-0ab3b73b4c43', 'perkakas', FALSE),
  ('7286ea5c-1585-4c5d-95c4-a1e77f940d66', 'd93f96c7-1a6c-439e-8aba-0ab3b73b4c43', 'burung', FALSE),
  ('d5a46882-ef9f-4a25-8afb-d0b41d370429', 'd93f96c7-1a6c-439e-8aba-0ab3b73b4c43', 'kesenian', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – N – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('e1a7d481-2001-4fc4-8b09-1e3c4ebd1f72', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – N – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('a1a6751c-5ff6-4858-a0ce-481716bcc808', 'e1a7d481-2001-4fc4-8b09-1e3c4ebd1f72', 'binatang', TRUE),
  ('5f9b5587-4441-4c89-9c65-193d1c587379', 'e1a7d481-2001-4fc4-8b09-1e3c4ebd1f72', 'bunga', FALSE),
  ('7531372b-e651-4766-9a60-f612c589ba37', 'e1a7d481-2001-4fc4-8b09-1e3c4ebd1f72', 'burung', FALSE),
  ('6cc1f722-b38a-46f6-98e9-718e292c1366', 'e1a7d481-2001-4fc4-8b09-1e3c4ebd1f72', 'kesenian', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – H – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('ea6cbdb2-9bb0-4109-9c14-6ad99e65e1ac', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – H – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('4faafbaf-1e28-484e-9c0d-24d0a6672631', 'ea6cbdb2-9bb0-4109-9c14-6ad99e65e1ac', 'perkakas', FALSE),
  ('d882dad9-c14e-4524-bd78-fa2d186965a0', 'ea6cbdb2-9bb0-4109-9c14-6ad99e65e1ac', 'binatang', FALSE),
  ('359b4b55-cd77-446b-bcc8-7026ff512c09', 'ea6cbdb2-9bb0-4109-9c14-6ad99e65e1ac', 'bunga', FALSE),
  ('486eb044-0821-4ae6-a10c-3662621489ba', 'ea6cbdb2-9bb0-4109-9c14-6ad99e65e1ac', 'kesenian', TRUE)
ON CONFLICT (id) DO NOTHING;

-- Question: Kata yang mempunyai huruf permulaan – J – adalah …….
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('f1e823bb-7e1c-4bd9-b286-1d0cd206dbe0', '307006f7-6baa-4c31-a8dd-5dab9b4d82b6', 'multiple_choice', 'Kata yang mempunyai huruf permulaan – J – adalah …….', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('f81271e0-fe79-4ad7-8b00-af6e6dae43af', 'f1e823bb-7e1c-4bd9-b286-1d0cd206dbe0', 'binatang', FALSE),
  ('27a3640c-e3ae-4165-a96b-ad794b554713', 'f1e823bb-7e1c-4bd9-b286-1d0cd206dbe0', 'burung', TRUE),
  ('f1bde4ac-5719-4c37-bcfa-1098ec501d15', 'f1e823bb-7e1c-4bd9-b286-1d0cd206dbe0', 'perkakas', FALSE),
  ('b5916443-34c7-470c-a307-d9a6f69632ab', 'f1e823bb-7e1c-4bd9-b286-1d0cd206dbe0', 'bunga', FALSE)
ON CONFLICT (id) DO NOTHING;
