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


