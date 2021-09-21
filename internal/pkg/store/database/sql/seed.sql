INSERT INTO PROJECT_TYPE (project_type_id, name)
VALUES ('f4a284cd-5a65-4468-8800-e0d8762933a9', 'Public'),
       ('0894f08f-1f45-4662-84bb-ac1d4f70b95d', 'Team'),
       ('b84bd65c-4fd5-4b67-abd1-b342fd67ee17', 'Private')
ON CONFLICT DO NOTHING;

INSERT INTO PROJECT (project_id, project_type_id, name, description, date_created, created_by_user_id, date_updated,
                     updated_by_user_id)
VALUES ('5b3ea10c-f6c6-4931-bbfc-ec20b190cca4', 'f4a284cd-5a65-4468-8800-e0d8762933a9', 'Test Project',
        'This project generated from migration seed preset.',
        '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6',
        '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6')
ON CONFLICT DO NOTHING;

INSERT INTO ROLE (role_id, name)
VALUES ('915c4e7e-a7fa-459d-9931-79de4b01621c', 'ProjectWorkspaceListRead'),
       ('5152caca-b43d-4b0b-8309-ac40a894eefc', 'ProjectReadAll'),
       ('16ab20b6-2016-4923-b14e-743b516efcf7', 'ProjectReadWriteAll')
ON CONFLICT DO NOTHING;

INSERT INTO GROUP (group_id, name, date_created, created_by_user_id, date_updated, updated_by_user_id)
VALUES ('f253b618-83c6-407b-af85-a1994e1e818c', 'Test Developer Group',
        '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6',
        '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6'),
       ('4b532822-59c8-4c67-941a-4b1704abad5f', 'Test Stakeholder Group',
        '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6',
        '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6')
ON CONFLICT DO NOTHING;

INSERT INTO GROUP_ROLE (group_role_id, group_id, role_id, date_created, created_by_user_id)
VALUES ('87c666e1-4141-4ca3-9632-2a8f934277c3', 'f253b618-83c6-407b-af85-a1994e1e818c',
        '16ab20b6-2016-4923-b14e-743b516efcf7', '2021-01-01 00:00:01.000001+00',
        '92eded9e-979c-4e94-afc5-2333fcc920f6'),
       ('461a010b-8791-4240-92d3-dbcfe8fa9b5d', '4b532822-59c8-4c67-941a-4b1704abad5f',
        '5152caca-b43d-4b0b-8309-ac40a894eefc', '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6')
ON CONFLICT DO NOTHING;

INSERT INTO GROUP_USER (group_user_id, group_id, user_id, date_created, created_by_user_id)
VALUES ('246f0914-7394-4341-96cd-0d27a66b1c37', 'f253b618-83c6-407b-af85-a1994e1e818c',
        '92eded9e-979c-4e94-afc5-2333fcc920f6', '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6')
ON CONFLICT DO NOTHING;

INSERT INTO PROJECT_GROUP_ACCESS (project_group_access_id, project_id, group_id, date_created, created_by_user_id)
VALUES ('949d2cd8-f0e1-4563-a29d-fbc3241b8d5c', '5b3ea10c-f6c6-4931-bbfc-ec20b190cca4',
        'f253b618-83c6-407b-af85-a1994e1e818c', '2021-01-01 00:00:01.000001+00',
        '92eded9e-979c-4e94-afc5-2333fcc920f6'),
       ('c19dae3f-0d97-45b7-b735-e0b9917f8067', '5b3ea10c-f6c6-4931-bbfc-ec20b190cca4',
        '4b532822-59c8-4c67-941a-4b1704abad5f', '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6')
ON CONFLICT DO NOTHING;

INSERT INTO STEM (stem_id, name)
VALUES ('2fdf996e-2372-4f3c-bccf-d8efcca8bd49', 'Sticker Pane'),
       ('78e95523-4ed2-49e6-8b1a-b8c073daab41', '2D Environment'),
       ('b8d78dda-027c-498e-8609-33cc6f4a6dbe', '3D Environment')
ON CONFLICT DO NOTHING;

INSERT INTO WORKSPACE (workspace_id, project_id, stem_id, name, description, asset_amount_limit,
                       x_max, y_max, z_max, date_created, created_by_user_id, date_updated, updated_by_user_id)
VALUES ('23b10a77-c45a-4bfc-a6c7-84cf8c6ab24e', '5b3ea10c-f6c6-4931-bbfc-ec20b190cca4',
        '2fdf996e-2372-4f3c-bccf-d8efcca8bd49', 'Sample Test Sticker Area', 'This workspace created by seed preset.',
        10, 1000, 1000, 1, '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6',
        '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6'),
       ('c89d7686-7b31-4818-93ee-ff146b79ae62', '5b3ea10c-f6c6-4931-bbfc-ec20b190cca4',
        '78e95523-4ed2-49e6-8b1a-b8c073daab41', 'Sample Collage Area', 'This workspace created by seed preset.',
        10, 1000, 1000, 1, '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6',
        '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6'),
       ('e3cd56d6-acf6-460f-8f0d-4674d9fab0a4', '5b3ea10c-f6c6-4931-bbfc-ec20b190cca4',
        'b8d78dda-027c-498e-8609-33cc6f4a6dbe', 'Sample 3D Scene', 'This workspace created by seed preset.',
        10, 1000, 1000, 1000, '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6',
        '2021-01-01 00:00:01.000001+00', '92eded9e-979c-4e94-afc5-2333fcc920f6')
ON CONFLICT DO NOTHING;
