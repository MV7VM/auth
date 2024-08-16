\с users

INSERT INTO roles (role) VALUES
('admin'),
('coach'),
('client');

INSERT INTO users (login, password, token, mail, roleID) VALUES
('+10', 'pass123', NULL, 'user1@example.com', 1),
('+11', 'pass123', NULL, 'user2@example.com', 2),
('+12', 'pass123', NULL, 'user3@example.com', 2),
('+13', 'pass123', NULL, 'user4@example.com', 3),
('+14', 'pass123', NULL, 'user5@example.com', 3),
('+15', 'pass123', NULL, 'user6@example.com', 3),
('+16', 'pass123', NULL, 'user7@example.com', 3),
('+17', 'pass123', NULL, 'user8@example.com', 3),
('+18', 'pass123', NULL, 'user9@example.com', 3),
('+19', 'pass123', NULL, 'user10@example.com', 3);
