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
     last_login timestamp default now(),
     created_by varchar(50) not null,
     created_at timestamp not null default now(),
     updated_by varchar(50),
     updated_at timestamp
     deleted_by varchar(50),
     deleted_at timestamp,
     UNIQUE(email, phone)
);

CREATE TABLE if not exists user_images (
    id varchar(255) PRIMARY KEY NOT NULL,
    user_id varchar(255) not null,
    image_url varchar(255) not null,
    order int not null,
    created_by varchar(50) not null,
    created_at timestamp not null default now(),
    updated_by varchar(50),
    updated_at timestamp
    deleted_by varchar(50),
    deleted_at timestamp
);


ALTER TABLE user_image ADD CONSTRAINT user_image_fk_user FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE NO ACTION ON UPDATE NO ACTION DEFERRABLE INITIALLY DEFERRED;

COMMIT;
