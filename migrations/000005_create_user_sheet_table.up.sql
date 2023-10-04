CREATE TABLE IF NOT EXISTS sheet_user (
    sheet_id uuid,
    user_id uuid,
    PRIMARY KEY (sheet_id, user_id),
    owner boolean,
    permission smallint,
    CONSTRAINT fk_sheet FOREIGN KEY(sheet_id) REFERENCES sheets(id),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);