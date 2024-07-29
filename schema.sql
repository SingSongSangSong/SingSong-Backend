-- 테이블이 존재할 경우 삭제
DROP TABLE IF EXISTS keepSong;
DROP TABLE IF EXISTS keepList;
DROP TABLE IF EXISTS member;
DROP TABLE IF EXISTS songInfo;
DROP TABLE IF EXISTS artistInfo;

-- member 테이블 생성
create table if not exists member (
  id BIGINT AUTO_INCREMENT primary key,
  nickname varchar(255),
    email varchar(50) not null,
    gender varchar(20),
    birthday date,
    provider varchar(20) not null
    );

-- keepList 테이블 생성
create table if not exists keepList (
    keepId BIGINT AUTO_INCREMENT primary key,
    memberId BIGINT not null,
    keepName varchar(255),
    foreign key (memberId) references member(id)
    );

-- artistInfo 테이블 생성
create table if not exists artistInfo (
  artistId BIGINT AUTO_INCREMENT primary key,
  artistName varchar(255) not null,
    artistType varchar(100),
    relatedArtists varchar(255),
    country varchar(255)
    );

-- songInfo 테이블 생성
create table if not exists songInfo (
    songId BIGINT AUTO_INCREMENT primary key,
    songName varchar(255) not null,
    artistId BIGINT not null,
    album varchar(255),
    songNumber int not null,
    octave varchar(10),
    tjLink varchar(255),
    tags varchar(255),
    foreign key (artistId) references artistInfo(artistId)
    );

-- keepSong 테이블 생성
create table if not exists keepSong (
    keepSongId BIGINT AUTO_INCREMENT primary key,
    keepId BIGINT not null,
    songId BIGINT not null,
    foreign key (keepId) references keepList(keepId),
    foreign key (songId) references songInfo(songId)
    );