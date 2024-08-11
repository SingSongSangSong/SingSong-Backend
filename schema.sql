-- 테이블이 존재할 경우 삭제
DROP TABLE IF EXISTS app_version;
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

-- appVersion 테이블 생성
CREATE TABLE IF NOT EXISTS app_version (
    app_version_id BIGINT AUTO_INCREMENT PRIMARY KEY,
    platform VARCHAR(10) NOT NULL,
    version VARCHAR(20) NOT NULL,
    force_update BOOLEAN NOT NULL DEFAULT FALSE,
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
    deleted_at TIMESTAMP NULL DEFAULT NULL,
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