use push_service
-- 먼저 기존 테이블이 있다면 삭제
DROP TABLE IF EXISTS messages;

DROP TABLE IF EXISTS users;

-- users 테이블 (소프트 삭제 포함)
CREATE TABLE
    users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE,
        deleted_at DATETIME DEFAULT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );

-- messages 테이블
CREATE TABLE
    messages (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT NOT NULL,
        content TEXT,
        status VARCHAR(20) NOT NULL DEFAULT 'pending',
        sent_at DATETIME DEFAULT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_user_id (user_id),
        INDEX idx_status (status),
        CONSTRAINT fk_messages_users FOREIGN KEY (user_id) REFERENCES users (id)
    );