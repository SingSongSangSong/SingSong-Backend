-- 테이블이 존재할 경우 삭제
DROP TABLE IF EXISTS keepSong;
DROP TABLE IF EXISTS songTempInfo;
DROP TABLE IF EXISTS songInfo;
DROP TABLE IF EXISTS artistInfo;
DROP TABLE IF EXISTS keepList;
DROP TABLE IF EXISTS member;

-- member 테이블 생성
CREATE TABLE IF NOT EXISTS member (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  nickname VARCHAR(255),
    email VARCHAR(50) NOT NULL,
    gender VARCHAR(20),
    birthday DATE,
    provider VARCHAR(20) NOT NULL,
    UNIQUE(email, provider) -- 이메일과 제공자는 유니크하도록 설정
    );

-- keepList 테이블 생성
CREATE TABLE IF NOT EXISTS keepList (
    keepId BIGINT AUTO_INCREMENT PRIMARY KEY,
    memberId BIGINT NOT NULL,
    keepName VARCHAR(255)
    );

-- artistInfo 테이블 생성
CREATE TABLE IF NOT EXISTS artistInfo (
  artistId BIGINT AUTO_INCREMENT PRIMARY KEY,
  artistName VARCHAR(255) NOT NULL,
    artistType VARCHAR(100),
    relatedArtists VARCHAR(255),
    country VARCHAR(255)
    );

-- songInfo 테이블 생성
CREATE TABLE IF NOT EXISTS songInfo (
    songId BIGINT AUTO_INCREMENT PRIMARY KEY,
    songName VARCHAR(255) NOT NULL,
    artistId BIGINT NOT NULL,
    album VARCHAR(255),
    songNumber INT NOT NULL,
    octave VARCHAR(10),
    tjLink VARCHAR(255),
    tags VARCHAR(255)
    );

-- songTempInfo 테이블 생성
CREATE TABLE IF NOT EXISTS songTempInfo (
    songTempId BIGINT AUTO_INCREMENT PRIMARY KEY,
    songName VARCHAR(255) NOT NULL,
    artistName VARCHAR(255) NOT NULL,
    album VARCHAR(255),
    songNumber INT NOT NULL,
    octave VARCHAR(10),
    tjLink VARCHAR(255),
    tags VARCHAR(255)
    );

-- keepSong 테이블 생성
CREATE TABLE IF NOT EXISTS keepSong (
    keepSongId BIGINT AUTO_INCREMENT PRIMARY KEY,
    keepId BIGINT NOT NULL,
    songTempId BIGINT NOT NULL,
    songNumber INT NOT NULL,
);