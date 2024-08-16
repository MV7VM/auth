\c content

INSERT INTO profileInfo (id, fullname, phone, mail, photoUrl, description) VALUES
(1, 'John Doe', '+10', 'user1@example.com', 'http://example.com/photo1.jpg', 'Admin'),
(2, 'Jane Smith', '+11', 'user2@example.com', 'http://example.com/photo2.jpg', 'Couch'),
(3, 'Alice Brown', '+12', 'user3@example.com', 'http://example.com/photo3.jpg', 'Couch'),
(4, 'Bob Johnson', '+13', 'user4@example.com', 'http://example.com/photo4.jpg', 'Client'),
(5, 'Charlie Black', '+14', 'user5@example.com', 'http://example.com/photo5.jpg', 'Client'),
(6, 'Diana White', '+15', 'user6@example.com', 'http://example.com/photo6.jpg', 'Client'),
(7, 'Eve Green', '+16', 'user7@example.com', 'http://example.com/photo7.jpg', 'Client'),
(8, 'Frank Blue', '+17', 'user8@example.com', 'http://example.com/photo8.jpg', 'Client'),
(9, 'Grace Yellow', '+18', 'user9@example.com', 'http://example.com/photo9.jpg', 'Client'),
(10, 'Hank Red', '+19', 'ser10@example.com', 'http://example.com/photo10.jpg', 'Client');

INSERT INTO certificates (name, description, userID) VALUES
('Fitness Certification', 'Certified fitness trainer.', 1),
('Nutrition Certification', 'Expert in nutrition and dietetics.', 2),
('Yoga Certification', 'Certified yoga instructor.', 3),
('Therapy License', 'Completed customer service training.', 4),
('Management Training', 'Completed customer service training.', 6),
('Mindfulness Training', 'Completed customer service training.', 7),
('Plant-based Nutrition', 'Completed customer service training.', 8),
('Pilates Certification', 'Expert in nutrition and dietetics.', 3),
('Mental Health Certification', 'Completed customer service training.', 9),
('Customer Service Training', 'StrongMan', 5);

INSERT INTO content (name, type, landingPoint, content, userID) VALUES
('Workout Plan', 'pdf', 'http://example.com/workoutplan1.pdf', NULL, NULL),
('Nutrition Guide', 'pdf', 'http://example.com/nutritionguide.pdf', NULL, 2),
('Yoga Routine', 'video', 'http://example.com/yogaroutine.mp4', NULL, 3),
('Therapy Session', 'audio', 'http://example.com/therapysession.mp3', NULL, 4),
('Management Report', 'pdf', 'http://example.com/managementreport.pdf', NULL, 6),
('Mindfulness Guide', 'pdf', 'http://example.com/mindfulnessguide.pdf', NULL, 7),
('Plant-Based Recipes', 'pdf', 'http://example.com/recipes.pdf', NULL, 8),
('Pilates Video', 'video', 'http://example.com/pilatesvideo.mp4', NULL, 3),
('Mental Health Tips', 'article', 'http://example.com/mentalhealthtips.html', NULL, 9),
('Customer Service Manual', 'pdf', 'http://example.com/customerservice.pdf', NULL, 5);
