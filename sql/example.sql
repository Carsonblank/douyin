DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name`           varchar(128)        NOT NULL DEFAULT '' COMMENT '用户昵称',
    `password`       varchar(128)        NOT NULL DEFAULT '' COMMENT '用户密码',
    `follow_count`   bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '关注数量',
    `follower_count` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '粉丝数量',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';

INSERT INTO `user`
VALUES (1, 'Jerry', '123456' ,0, 0),
       (2, 'Tom', '123456' ,0, 0);

DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`
(
    `id`              bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `userID`          bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '作者id',
    `play_url`        varchar(128)        NOT NULL DEFAULT '' COMMENT '视频地址',
    `cover_url`       varchar(128)        NOT NULL DEFAULT '' COMMENT '封面地址',
    `favorite_count`  bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '点赞数量',
    `comment_count`   bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '评论数量',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='视频表';
INSERT INTO `video`
VALUES (1, '1', 'http://175.178.98.212:20000/root/test/simple-demo-main/public/bear.mp4','https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg', 0,0),
       (2, '1', 'http://175.178.98.212:20000/root/test/simple-demo-main/public/望梅止渴.mp4','http://175.178.98.212:20000/root/test/simple-demo-main/public/1.png', 0,0);
