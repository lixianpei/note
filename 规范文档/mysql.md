# 库表相关规范
- 库表统一使用InnoDB引擎
- 优先使用 utf8mb4 格式
- 表名称推荐：库名缩写_模块名称_业务名称 例如：tp_system_admin
- 添加合适的索引：唯一索引：uni_xxx、普通索引：idx_xxx

# 表字段命名规范
- 数据表名称全部小写、单词之间下划线拼接；数据表名称根据业务名称定义
- 数据表字段按照下划线格式，例如，创建时间：created_time、注册用户手机号：register_user_phone
- 数据表必须包含的字段：id-主键自增、create_time-创建时间、update_time-更新时间、delete_time-删除时间、is_delete-软删除
- 数据表中存在后台管理员修改数据必含字段：create_admin_id-创建人、update_admin_id-更新人、delete_admin_id-删除人
- 使用 tinyint 代替 enum
- 字符串长度固定的使用 char 类型