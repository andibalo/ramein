BEGIN;

CREATE TABLE if not exists users (
    id varchar(255) PRIMARY KEY NOT NULL,
    email varchar(255) not null,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    phone varchar(255) not null,
    role varchar(255) not null,
    password varchar(255) not null,
    is_super_user boolean not null,
    is_verified boolean not null,
    is_email_verified boolean not null,
    profile_summary varchar(255),
    last_login timestamptz default now(),
    created_by varchar(50) not null,
    created_at timestamptz not null default now(),
    updated_by varchar(50),
    updated_at timestamptz,
    deleted_by varchar(50),
    deleted_at timestamptz,
    UNIQUE(email, phone)
    );

CREATE TABLE if not exists user_images (
    id varchar(255) PRIMARY KEY NOT NULL,
    user_id varchar(255) not null,
    image_url varchar(255) not null,
    image_order int not null,
    created_by varchar(50) not null,
    created_at timestamptz not null default now(),
    updated_by varchar(50),
    updated_at timestamptz,
    deleted_by varchar(50),
    deleted_at timestamptz
    );

CREATE TABLE if not exists user_verify_emails (
    id varchar(255) PRIMARY KEY NOT NULL,
    user_id varchar(255) not null,
    email varchar(255) not null,
    secret_code varchar(255) not null,
    is_used boolean not null,
    expired_at timestamptz not null,
    created_by varchar(50) not null,
    created_at timestamptz not null default now(),
    updated_by varchar(50),
    updated_at timestamptz,
    deleted_by varchar(50),
    deleted_at timestamptz
    );

ALTER TABLE user_images ADD CONSTRAINT user_images_fk_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE user_verify_emails ADD CONSTRAINT user_verify_emails_fk_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION DEFERRABLE INITIALLY DEFERRED;


COMMIT;
