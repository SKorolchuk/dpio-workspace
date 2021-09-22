CREATE TABLE PROJECT_COLLABORATION_TYPE
(
    project_collaboration_type_id UUID,
    name                          varchar(255) UNIQUE,

    PRIMARY KEY (project_collaboration_type_id)
);

CREATE TABLE PROJECT
(
    project_id                    UUID,
    project_collaboration_type_id UUID,
    name                          varchar(255) UNIQUE,
    description                   varchar(500),
    date_created                  timestamp,
    created_by_user_id            UUID,
    date_updated                  timestamp,
    updated_by_user_id            UUID,

    PRIMARY KEY (project_id),
    FOREIGN KEY (project_collaboration_type_id) REFERENCES PROJECT_COLLABORATION_TYPE (project_collaboration_type_id)
);

CREATE TABLE PROJECT_ROLE
(
    project_role_id UUID,
    name            varchar(255) UNIQUE,

    PRIMARY KEY (project_role_id)
);

CREATE TABLE PROJECT_GROUP
(
    project_group_id   UUID,
    name               varchar(255) UNIQUE,
    date_created       timestamp,
    created_by_user_id UUID,
    date_updated       timestamp,
    updated_by_user_id UUID,

    PRIMARY KEY (project_group_id)
);

CREATE TABLE PROJECT_GROUP_ROLE
(
    project_group_role_id UUID,
    project_group_id      UUID,
    project_role_id       UUID,
    date_created          timestamp,
    created_by_user_id    UUID,

    PRIMARY KEY (project_group_role_id),
    FOREIGN KEY (project_group_id) REFERENCES PROJECT_GROUP (project_group_id),
    FOREIGN KEY (project_role_id) REFERENCES PROJECT_ROLE (project_role_id)
);

CREATE TABLE PROJECT_GROUP_USER
(
    project_group_user_id UUID,
    project_group_id      UUID,
    user_id               UUID,
    date_created          timestamp,
    created_by_user_id    UUID,

    PRIMARY KEY (project_group_user_id),
    FOREIGN KEY (project_group_id) REFERENCES PROJECT_GROUP (project_group_id)
);

CREATE TABLE PROJECT_GROUP_ACCESS
(
    project_group_access_id UUID,
    project_id              UUID,
    project_group_id        UUID,
    date_created            timestamp,
    created_by_user_id      UUID,

    PRIMARY KEY (project_group_access_id),
    FOREIGN KEY (project_id) REFERENCES PROJECT (project_id),
    FOREIGN KEY (project_group_id) REFERENCES PROJECT_GROUP (project_group_id)
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
