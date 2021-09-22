create table agency
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6)  null,
    updated_at datetime(6)  null,
    deleted_at datetime(6)  null,
    name       varchar(64)  null,
    tel        varchar(64)  null,
    addr       varchar(128) null,
    email      varchar(128) null
)
    charset = utf8;

create table agency_acquirer
(
    id                   bigint auto_increment
        primary key,
    created_at           datetime(6)  null,
    updated_at           datetime(6)  null,
    deleted_at           datetime(6)  null,
    name                 varchar(32)  null,
    impl_name            varchar(32)  null,
    addition             varchar(255) null,
    config_file_url      varchar(512) null,
    agency_id            bigint       null,
    auto_settlement_time varchar(16)  null
)
    charset = utf8;

create table agency_entry_type
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6) null,
    updated_at datetime(6) null,
    deleted_at datetime(6) null,
    name       varchar(32) null
)
    charset = utf8;

create table agency_payment_method
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6) null,
    updated_at datetime(6) null,
    deleted_at datetime(6) null,
    name       varchar(32) null
)
    charset = utf8;

create table agency_payment_type
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6) null,
    updated_at datetime(6) null,
    deleted_at datetime(6) null,
    name       varchar(32) null
)
    charset = utf8;

create table agency_user_associate
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6)     null,
    updated_at datetime(6)     null,
    deleted_at datetime(6)     null,
    agency_id  bigint unsigned null,
    user_id    bigint          null,
    role       varchar(32)     null
)
    charset = utf8;

create index user_id
    on agency_user_associate (user_id);

create table merchant
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6)  null,
    updated_at datetime(6)  null,
    deleted_at datetime(6)  null,
    name       varchar(64)  null,
    tel        varchar(64)  null,
    addr       varchar(128) null,
    agency_id  bigint       null,
    email      varchar(128) null
)
    charset = utf8;

create table merchant_device
(
    id          bigint auto_increment
        primary key,
    created_at  datetime(6) null,
    updated_at  datetime(6) null,
    deleted_at  datetime(6) null,
    device_id   bigint      null,
    merchant_id bigint      null
)
    charset = utf8;

create index merchant_id
    on merchant_device (merchant_id);

create table merchant_payment_setting_in_device
(
    id                 bigint auto_increment
        primary key,
    created_at         datetime(6)                  null,
    updated_at         datetime(6)                  null,
    deleted_at         datetime(6)                  null,
    merchant_device_id bigint                       null,
    payment_methods    longtext collate utf8mb4_bin null,
    entry_types        longtext collate utf8mb4_bin null,
    payment_types      longtext collate utf8mb4_bin null,
    acquirer_id        bigint                       null,
    mid                varchar(32)                  null,
    tid                varchar(32)                  null,
    addition           varchar(512)                 null,
    constraint entry_types
        check (json_valid(`entry_types`)),
    constraint payment_methods
        check (json_valid(`payment_methods`)),
    constraint payment_types
        check (json_valid(`payment_types`))
)
    charset = utf8;

create index merchant_device_id
    on merchant_payment_setting_in_device (merchant_device_id);

create table merchant_user_associate
(
    id          bigint auto_increment
        primary key,
    created_at  datetime(6) null,
    updated_at  datetime(6) null,
    deleted_at  datetime(6) null,
    merchant_id bigint      null,
    user_id     bigint      null,
    role        varchar(16) null
)
    charset = utf8;

create index user_id
    on merchant_user_associate (user_id, merchant_id);

create table tms_app
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime(6)  null,
    updated_at  datetime(6)  null,
    deleted_at  datetime(6)  null,
    agency_id   bigint       null,
    name        varchar(32)  null,
    package_id  varchar(64)  null,
    description varchar(255) null
)
    collate = utf8_bin;

create table tms_app_file
(
    id                 bigint unsigned auto_increment
        primary key,
    deleted_at         datetime     null,
    updated_at         datetime     null,
    created_at         datetime     null,
    version_name       varchar(32)  null,
    version_code       int          null,
    update_description varchar(255) null,
    file_name          varchar(128) null,
    file_url           varchar(255) not null,
    decode_status      varchar(16)  null,
    decode_fail_msg    varchar(255) null,
    app_id             bigint       not null
)
    collate = utf8_bin;

create index appid
    on tms_app_file (app_id);

create table tms_app_in_device
(
    id               bigint unsigned auto_increment
        primary key,
    created_at       datetime(6)             null,
    updated_at       datetime(6)             null,
    deleted_at       datetime(6)             null,
    external_id      bigint                  null,
    external_id_type varchar(16) default '0' null,
    name             varchar(126)            null,
    package_id       varchar(64)             null,
    version_name     varchar(64)             null,
    version_code     int                     null,
    status           varchar(24)             null,
    app_id           bigint                  null,
    app_file_id      bigint                  null
)
    collate = utf8_bin;

create index d
    on tms_app_in_device (external_id, external_id_type);

create table tms_batch_update
(
    id              bigint unsigned auto_increment
        primary key,
    created_at      datetime(6)                             null,
    updated_at      datetime(6)                             null,
    deleted_at      datetime(6)                             null,
    agency_id       int(20)                                 null,
    description     varchar(255) collate utf8mb4_unicode_ci null,
    status          varchar(11)                             null,
    update_fail_msg varchar(64) collate utf8mb4_unicode_ci  null,
    tags            longtext collate utf8mb4_bin            null,
    device_models   longtext collate utf8mb4_bin            null,
    constraint device_models
        check (json_valid(`device_models`)),
    constraint tags
        check (json_valid(`tags`))
)
    collate = utf8_bin;

create table tms_device
(
    id                  bigint unsigned auto_increment,
    created_at          datetime(6)       null,
    updated_at          datetime(6)       null,
    deleted_at          datetime(6)       null,
    agency_id           int unsigned      null,
    device_sn           varchar(126)      not null,
    device_csn          varchar(126)      null comment '第三方ID',
    device_model        int(16) default 0 null,
    alias               varchar(16)       null,
    reboot_mode         varchar(16)       not null,
    reboot_time         varchar(6)        not null,
    reboot_day_in_week  int               null,
    reboot_day_in_month int               null,
    location_lat        varchar(32)       null,
    location_lon        varchar(32)       null,
    push_token          varchar(64)       null comment '推送ID',
    battery             int               null,
    primary key (id, device_sn),
    constraint device_sn
        unique (device_sn)
)
    collate = utf8_bin;

create table tms_device_and_tag_mid
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(6) null,
    updated_at datetime(6) null,
    deleted_at datetime(6) null,
    device_id  bigint      not null,
    tag_id     bigint      not null
)
    collate = utf8_bin;

create index device_id
    on tms_device_and_tag_mid (device_id, tag_id);

create table tms_device_batch_update
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6) null,
    updated_at datetime(6) null,
    deleted_at datetime(6) null,
    batch_id   bigint      null,
    device_id  bigint      null,
    status     varchar(16) null
);

create index tms_device_batch_update_batch_id_status_index
    on tms_device_batch_update (batch_id, status);

create index tms_device_batch_update_device_id_index
    on tms_device_batch_update (device_id);

create table tms_model
(
    id         bigint unsigned auto_increment
        primary key,
    name       varchar(32) not null,
    created_at datetime(6) null,
    updated_at datetime(6) null,
    deleted_at datetime(6) null
)
    collate = utf8_bin;

create table tms_tags
(
    id          bigint unsigned auto_increment
        primary key,
    created_at  datetime(6)     null,
    updated_at  datetime(6)     null,
    deleted_at  datetime(6)     null,
    agency_id   bigint unsigned null,
    name        varchar(64)     not null,
    description varchar(256)    null
)
    collate = utf8_bin;

create table tms_upload_file
(
    id         bigint unsigned auto_increment
        primary key,
    created_at datetime(6)                             null,
    updated_at datetime(6)                             null,
    deleted_at datetime(6)                             null,
    device_sn  varchar(64) collate utf8mb4_unicode_ci  null,
    file_name  varchar(128) collate utf8mb4_unicode_ci null,
    file_url   varchar(256) collate utf8mb4_unicode_ci null,
    agency_id  bigint                                  null
)
    collate = utf8_bin;

create index device_sn
    on tms_upload_file (device_sn);

create table user
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6)     null,
    updated_at datetime(6)     null,
    deleted_at datetime(6)     null,
    agency_id  bigint unsigned null,
    email      varchar(64)     not null,
    pwd        varchar(64)     null,
    name       varchar(64)     null,
    role       varchar(16)     null,
    active     tinyint(1)      null,
    constraint email
        unique (email) using hash
)
    charset = utf8;

create table user_app_id
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6)  null,
    updated_at datetime(6)  null,
    deleted_at datetime(6)  null,
    app_id     varchar(64)  null,
    app_secret varchar(255) null
)
    charset = utf8;

create table user_role
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6)  null,
    updated_at datetime(6)  null,
    deleted_at datetime(6)  null,
    store_id   int unsigned null,
    name       varchar(16)  null
)
    charset = utf8;

create table user_token
(
    id         bigint auto_increment
        primary key,
    created_at datetime(6) null,
    updated_at datetime(6) null,
    deleted_at datetime(6) null,
    user_id    bigint      null,
    app_id     bigint      null,
    token      varchar(64) null
)
    charset = utf8;

create index token
    on user_token (token);

create index user_id
    on user_token (user_id, app_id);

