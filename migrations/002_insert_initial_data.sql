-- 初始数据插入脚本

-- 插入默认主题
INSERT INTO themes (name, description, primary_color, logo_url, active)
VALUES ('Default', '默认主题', '#1890ff', '', TRUE);

-- 插入管理员角色
INSERT INTO roles (name, description)
VALUES ('admin', '系统管理员');

-- 插入基础权限
INSERT INTO permissions (name, description, resource, action)
VALUES 
('manage_users', '管理用户', 'user', 'manage'),
('view_users', '查看用户', 'user', 'view'),
('manage_roles', '管理角色', 'role', 'manage'),
('view_roles', '查看角色', 'role', 'view'),
('manage_permissions', '管理权限', 'permission', 'manage'),
('view_permissions', '查看权限', 'permission', 'view'),
('manage_applications', '管理应用', 'application', 'manage'),
('view_applications', '查看应用', 'application', 'view'),
('manage_themes', '管理主题', 'theme', 'manage'),
('view_themes', '查看主题', 'theme', 'view');

-- 为管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'admin'), 
    id 
FROM permissions;

-- 插入默认管理员用户（密码需要在应用中使用bcrypt加密）
-- 这里密码是占位符，实际实现时应该通过应用程序生成加密后的密码
INSERT INTO users (username, email, password, full_name, active, theme_id)
VALUES ('admin', 'admin@example.com', '$2a$10$PLACEHOLDER_HASH', '系统管理员', TRUE, 1);

-- 为默认管理员分配管理员角色
INSERT INTO user_roles (user_id, role_id)
VALUES (
    (SELECT id FROM users WHERE username = 'admin'),
    (SELECT id FROM roles WHERE name = 'admin')
);