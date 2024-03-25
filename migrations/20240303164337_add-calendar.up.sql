CREATE TABLE calendar_user (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    surname VARCHAR NOT NULL,
    telegram_id BIGSERIAL
);

CREATE INDEX telegram_id_idx ON calendar_user (telegram_id);

CREATE TABLE event (
    id SERIAL PRIMARY KEY,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    name VARCHAR,
    description VARCHAR,
    creator_id INTEGER
);

ALTER TABLE event
ADD CONSTRAINT event_fk_creator_id
FOREIGN KEY (creator_id) REFERENCES calendar_user(id);

CREATE TABLE notification (
    user_id INTEGER,
    event_id INTEGER,
    notify_time TIMESTAMPTZ,
    CONSTRAINT notification_pk PRIMARY KEY (user_id, event_id)
);

ALTER TABLE notification
ADD CONSTRAINT notification_fk_user_id
FOREIGN KEY (user_id) REFERENCES calendar_user(id);

ALTER TABLE notification
ADD CONSTRAINT notification_fk_event_id
FOREIGN KEY (event_id) REFERENCES event(id);

CREATE TABLE event_to_user (
    user_id INTEGER,
    event_id INTEGER,
    CONSTRAINT event_to_user_pk PRIMARY KEY (user_id, event_id)
);

ALTER TABLE event_to_user
ADD CONSTRAINT event_to_user_fk_user_id
FOREIGN KEY (user_id) REFERENCES calendar_user(id);

ALTER TABLE event_to_user
ADD CONSTRAINT event_to_user_fk_event_id
FOREIGN KEY (event_id) REFERENCES event(id);
