-- 테이블이 존재할 경우 삭제
DROP TABLE IF EXISTS keep_song;
DROP TABLE IF EXISTS song_info;
DROP TABLE IF EXISTS artist;
DROP TABLE IF EXISTS keep_list;
DROP TABLE IF EXISTS member;
DROP TABLE IF EXISTS comment;
DROP TABLE IF EXISTS comment_like;
DROP TABLE IF EXISTS report;
DROP TABLE IF EXISTS song_review;
DROP TABLE IF EXISTS song_review_option;

-- member 테이블 생성
CREATE TABLE IF NOT EXISTS member (
    member_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    nickname VARCHAR(255),
    email VARCHAR(50) NOT NULL,
    gender VARCHAR(20),
    birthyear INT,
    provider VARCHAR(20) NOT NULL,
    UNIQUE(email, provider), -- 이메일과 제공자는 유니크하도록 설정
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE member MODIFY email VARCHAR(150);

-- keepList 테이블 생성
CREATE TABLE IF NOT EXISTS keep_list (
    keep_list_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT NOT NULL,
    keep_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- artistInfo 테이블 생성
CREATE TABLE IF NOT EXISTS artist (
    artist_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    artist_name VARCHAR(255) NOT NULL,
    artist_type VARCHAR(100),
    related_artists TEXT,
    country VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);


-- songInfo 테이블 생성
CREATE TABLE IF NOT EXISTS song_info (
    song_info_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    song_name VARCHAR(255) NOT NULL,
    artist_id BIGINT,
    artist_name VARCHAR(255) NOT NULL,
    artist_type VARCHAR(100),
    is_mr BOOLEAN DEFAULT FALSE,
    is_chosen_22000 BOOLEAN DEFAULT FALSE,
    related_artists TEXT,
    country VARCHAR(255),
    album VARCHAR(255),
    song_number INT NOT NULL,
    octave VARCHAR(10),
    tj_link VARCHAR(255),
    tags VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE song_info ADD COLUMN video_link TEXT;

-- keepSong 테이블 생성
CREATE TABLE IF NOT EXISTS keep_song (
    keep_song_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    keep_list_id BIGINT NOT NULL,
    song_info_id BIGINT NOT NULL,
    song_number INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS song_review_option (
    song_review_option_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS song_review (
    song_review_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    song_info_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    song_review_option_id BIGINT NOT NULL,
    gender VARCHAR(20),
    birthyear INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS comment (
    comment_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    parent_comment_id BIGINT,
    song_info_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    content TEXT,
    is_recomment BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE comment
    ADD CONSTRAINT fk_member_id
        FOREIGN KEY (member_id) REFERENCES member(member_id)
            ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE IF NOT EXISTS comment_like (
    comment_like_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    comment_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE comment_like
    ADD CONSTRAINT fk_comment_id
        FOREIGN KEY (comment_id) REFERENCES comment(comment_id)
            ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE IF NOT EXISTS report (
    report_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    comment_id BIGINT NOT NULL,
    reporter_member_id BIGINT NOT NULL,
    subject_member_id BIGINT NOT NULL,
    report_reason VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE comment ADD COLUMN likes int DEFAULT NULL;
ALTER TABLE comment MODIFY COLUMN likes int DEFAULT 0;
UPDATE comment SET likes = 0 WHERE likes IS NULL;


-- 현재 설정된 인덱스들
CREATE INDEX idx_song_info_song_number ON song_info(song_number);
CREATE INDEX idx_keep_list_member_id ON keep_list(member_id);
CREATE INDEX idx_keep_song_keep_list_id ON keep_song(keep_list_id);
CREATE INDEX idx_member_email_provider ON member(email, provider);

CREATE TABLE IF NOT EXISTS member_action (
    member_action_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT NOT NULL,
    gender VARCHAR(20),
    birthyear INT,
    song_info_id BIGINT NOT NULL,
    action_type VARCHAR(20) NOT NULL,
    action_score FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- soft delete시 재가입을 고려하여 기존 (email, provider) 유니크 인덱스 삭제
ALTER TABLE member DROP INDEX email;
ALTER TABLE member
    ADD not_archived BOOLEAN
        GENERATED ALWAYS AS (IF(deleted_at IS NULL, 1, NULL)) VIRTUAL;
ALTER TABLE member
    ADD CONSTRAINT UNIQUE (email, provider, not_archived);

-- song_review_option 테이블에 영어 enum 추가
ALTER TABLE song_review_option ADD COLUMN enum VARCHAR(20);

CREATE TABLE IF NOT EXISTS blacklist (
    blacklist_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    blocker_member_id BIGINT NOT NULL,
    blocked_member_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE song_info ADD COLUMN melon_song_id VARCHAR(255);
UPDATE song_info SI
    JOIN raw_song_info RSI ON SI.song_info_id = RSI.song_info_id
    SET SI.melon_song_id = RSI.melon_song_id;

ALTER TABLE song_info ADD COLUMN is_live BOOLEAN DEFAULT FALSE;
ALTER TABLE song_info ADD UNIQUE INDEX (song_number, song_name, artist_name);


-- 커뮤니티 기능 추가
CREATE TABLE IF NOT EXISTS board (
    board_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

INSERT INTO board (name) VALUES ('자유 게시판');

CREATE TABLE IF NOT EXISTS post (
    post_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    board_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    title VARCHAR(100) NOT NULL,
    content TEXT,
    likes INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE post_comment (
    post_comment_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    post_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    content TEXT,
    likes INT NOT NULL,
    is_recomment BOOLEAN DEFAULT FALSE,
    parent_post_comment_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE post_comment
    ADD CONSTRAINT fk_post_comment_post_id
        FOREIGN KEY (post_id) REFERENCES post(post_id)
            ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE post_comment
    ADD CONSTRAINT fk_post_comment_member_id
    FOREIGN KEY (member_id) REFERENCES member(member_id)
    ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE post
    ADD CONSTRAINT fk_post_member_id
    FOREIGN KEY (member_id) REFERENCES member(member_id)
    ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE IF NOT EXISTS post_song(
    post_song_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    post_id BIGINT NOT NULL,
    song_info_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE post_song
    ADD CONSTRAINT fk_post_song_post_id
    FOREIGN KEY (post_id) REFERENCES post(post_id)
    ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE IF NOT EXISTS post_comment_song (
    post_comment_song_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    post_comment_id BIGINT NOT NULL,
    song_info_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE post_comment_song
    ADD CONSTRAINT fk_post_comment_song_post_comment_id
    FOREIGN KEY (post_comment_id) REFERENCES post_comment(post_comment_id)
    ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE post_like (
    post_like_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE post_comment_like (
    post_comment_like_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT NOT NULL,
    post_comment_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE post_report (
    post_report_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    reporter_member_id BIGINT NOT NULL,
    subject_member_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    report_reason VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE post_comment_report (
    post_comment_report_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    reporter_member_id BIGINT NOT NULL,
    subject_member_id BIGINT NOT NULL,
    post_comment_id BIGINT NOT NULL,
    report_reason VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE `member`
    MODIFY `email` varchar(255) NOT NULL;

CREATE TABLE llm_search_log (
    llm_search_log_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT NOT NULL,
    search_text VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE llm_search_log
    ADD CONSTRAINT fk_llm_search_log_member_id
    FOREIGN KEY (member_id) REFERENCES member(member_id);

ALTER TABLE song_info
DROP COLUMN tags,
DROP COLUMN artist_id,
DROP COLUMN octave;

ALTER TABLE song_info
ADD COLUMN tj_score FLOAT,
ADD COLUMN genre VARCHAR(100),
ADD COLUMN year INT,
ADD COLUMN lyrics_video_link VARCHAR(255),
ADD COLUMN artist_gender VARCHAR(50),
ADD COLUMN octave VARCHAR(50),
ADD COLUMN tj_youtube_link VARCHAR(255);

ALTER TABLE song_info
ADD COLUMN classics BOOL DEFAULT FALSE,
ADD COLUMN finale BOOL DEFAULT FALSE,
ADD COLUMN high BOOL DEFAULT FALSE,
ADD COLUMN low BOOL DEFAULT FALSE,
ADD COLUMN rnb BOOL DEFAULT FALSE,
ADD COLUMN breakup BOOL DEFAULT FALSE,
ADD COLUMN ballads BOOL DEFAULT FALSE,
ADD COLUMN dance BOOL DEFAULT FALSE,
ADD COLUMN duet BOOL DEFAULT FALSE,
ADD COLUMN ssum BOOL DEFAULT FALSE,
ADD COLUMN carol BOOL DEFAULT FALSE,
ADD COLUMN rainy BOOL DEFAULT FALSE,
ADD COLUMN pop BOOL DEFAULT FALSE,
ADD COLUMN office BOOL DEFAULT FALSE,
ADD COLUMN wedding BOOL DEFAULT FALSE,
ADD COLUMN military BOOL DEFAULT FALSE;

ALTER TABLE keep_list ADD COLUMN likes INT DEFAULT 0;

CREATE TABLE keep_list_like (
    keep_list_like_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT NOT NULL,
    keep_list_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (member_id) REFERENCES member(member_id) ON DELETE CASCADE,
    FOREIGN KEY (keep_list_id) REFERENCES keep_list(keep_list_id) ON DELETE CASCADE
);

CREATE TABLE keep_list_subscribe (
    keep_list_subscribe_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT NOT NULL,
    keep_list_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (member_id) REFERENCES member(member_id) ON DELETE CASCADE,
    FOREIGN KEY (keep_list_id) REFERENCES keep_list(keep_list_id) ON DELETE CASCADE
);

ALTER TABLE song_info
    ADD COLUMN hiphop BOOL DEFAULT FALSE,
    ADD COLUMN jpop BOOL DEFAULT FALSE,
    ADD COLUMN musical BOOL DEFAULT FALSE,
    ADD COLUMN band BOOL DEFAULT FALSE;

-- appVersion 테이블 생성
drop table app_version;
CREATE TABLE IF NOT EXISTS app_version (
       app_version_id BIGINT AUTO_INCREMENT PRIMARY KEY,
       platform VARCHAR(10) NOT NULL,
       latest_version VARCHAR(20) NOT NULL,
       force_update_version VARCHAR(20) NOT NULL,
       update_url VARCHAR(255) NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
       deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE TABLE search_log (
    search_log_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT NOT NULL,
    search_text VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (member_id) REFERENCES member(member_id) ON DELETE CASCADE
);


CREATE TABLE song_recording (
    song_recording_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    song_info_id BIGINT NOT NULL,
    member_id BIGINT NOT NULL,
    recording_link VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    is_public BOOL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (song_info_id) REFERENCES song_info(song_info_id) ON DELETE CASCADE,
    FOREIGN KEY (member_id) REFERENCES member(member_id) ON DELETE CASCADE
);


CREATE TABLE member_device_token (
     member_device_token_id BIGINT AUTO_INCREMENT PRIMARY KEY,
     member_id BIGINT NOT NULL,
     device_token VARCHAR(255) NOT NULL,
     is_activate BOOL DEFAULT TRUE,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     deleted_at TIMESTAMP NULL DEFAULT NULL
);

ALTER TABLE member_device_token
    ADD CONSTRAINT fk_member_device_token_member_id
        FOREIGN KEY (member_id) REFERENCES member(member_id);

CREATE TABLE notification_history (
    notification_history_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    member_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    body VARCHAR(255) NOT NULL,
    screen_type VARCHAR(10),
    screen_type_id BIGINT,
    is_read BOOL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
)