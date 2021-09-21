CREATE TABLE PROJECT_TYPE
(
    project_type_id UUID,
    name            varchar(255) UNIQUE,

    PRIMARY KEY (project_type_id)
);

CREATE TABLE PROJECT
(
    project_id         UUID,
    project_type_id    UUID,
    name               varchar(255) UNIQUE,
    description        varchar(500),
    date_created       timestamp,
    created_by_user_id UUID,
    date_updated       timestamp,
    updated_by_user_id UUID,

    PRIMARY KEY (project_id),
    FOREIGN KEY (project_type_id) REFERENCES PROJECT_TYPE (project_type_id)
);

CREATE TABLE ROLE
(
    role_id UUID,
    name    varchar(255) UNIQUE,

    PRIMARY KEY (role_id)
);

CREATE TABLE "GROUP"
(
    group_id           UUID,
    name               varchar(255) UNIQUE,
    date_created       timestamp,
    created_by_user_id UUID,
    date_updated       timestamp,
    updated_by_user_id UUID,

    PRIMARY KEY (group_id)
);

CREATE TABLE GROUP_ROLE
(
    group_role_id      UUID,
    group_id           UUID,
    role_id            UUID,
    date_created       timestamp,
    created_by_user_id UUID,

    PRIMARY KEY (group_role_id),
    FOREIGN KEY (group_id) REFERENCES "GROUP" (group_id),
    FOREIGN KEY (role_id) REFERENCES ROLE (role_id)
);

CREATE TABLE GROUP_USER
(
    group_user_id      UUID,
    group_id           UUID,
    user_id            UUID,
    date_created       timestamp,
    created_by_user_id UUID,

    PRIMARY KEY (group_user_id),
    FOREIGN KEY (group_id) REFERENCES "GROUP" (group_id)
);

CREATE TABLE PROJECT_GROUP_ACCESS
(
    project_group_access_id UUID,
    project_id              UUID,
    group_id                UUID,
    date_created            timestamp,
    created_by_user_id      UUID,

    PRIMARY KEY (project_group_access_id),
    FOREIGN KEY (project_id) REFERENCES PROJECT (project_id),
    FOREIGN KEY (group_id) REFERENCES "GROUP" (group_id)
);

CREATE TABLE STEM
(
    stem_id UUID,
    name    varchar(255) UNIQUE,

    PRIMARY KEY (stem_id)
);

CREATE TABLE WORKSPACE
(
    workspace_id       UUID,
    project_id         UUID,
    stem_id            UUID,
    name               varchar(255) UNIQUE,
    description        varchar(500),
    asset_amount_limit numeric,
    x_max              numeric,
    y_max              numeric,
    z_max              numeric,
    date_created       timestamp,
    created_by_user_id UUID,
    date_updated       timestamp,
    updated_by_user_id UUID,

    PRIMARY KEY (workspace_id),
    FOREIGN KEY (project_id) REFERENCES PROJECT (project_id),
    FOREIGN KEY (stem_id) REFERENCES STEM (stem_id)
);

CREATE TABLE ASSET
(
    asset_id              UUID,
    workspace_id          UUID,
    asset_external_ref_id UUID,
    position_x            numeric,
    position_y            numeric,
    position_z            numeric,
    scale                 numeric,
    height_by_y           numeric,
    width_by_x            numeric,
    length_by_z           numeric,
    date_created          timestamp,
    created_by_user_id    UUID,
    date_updated          timestamp,
    updated_by_user_id    UUID,

    PRIMARY KEY (asset_id),
    FOREIGN KEY (workspace_id) REFERENCES WORKSPACE (workspace_id)
);
