\c scheduler

INSERT INTO filials (name, start_time, end_time) VALUES
('Filial 1', '06:00:00', '22:00:00'),
('Filial 2', '07:00:00', '23:00:00');

INSERT INTO management (user_id, filial_id) VALUES
(1, 1),
(1, 1),
(1, 2),
(2, 2),
(1, 3),
(2, 3),
(1, 4),
(2, 4),
(1, 5),
(2, 5),
(1, 6),
(2, 6),
(1, 7),
(2, 7),
(1, 8),
(2, 8)
(1, 9),
(2, 9),
(1, 10),
(2, 10);

INSERT INTO workoutFormat (name, duration, filial_id, max_count_clients, price, description) VALUES
('Yoga Class', '00:60', 1, 20, 1500, 'A relaxing yoga session.'),
('Pilates Class', '00:45', 1, 15, 1300, 'A strength-focused Pilates class.'),
('Cardio Blast', '00:30', 1, 25, 1000, 'High-intensity cardio workout.'),
('Strength Training', '0:60', 1, 20, 1800, 'Build strength and endurance.'),
('HIIT', '00:30', 1, 25, 1200, 'High-Intensity Interval Training.'),
('Zumba', '00:45', 2, 30, 1100, 'Dance your way to fitness.'),
('Spinning', '00:60', 2, 20, 1600, 'An intense cycling session.'),
('CrossFit', '00:60', 2, 15, 2000, 'Full-body CrossFit workout.'),
('Boxing', '00:60', 2, 15, 1700, 'Learn boxing techniques.'),
('Mindfulness Meditation', '00:45', 2, 20, 1400, 'Meditation and mindfulness practice.');

INSERT INTO workouts (coach_id, format_id, date, start_time, end_time, status) VALUES
(2, 1, '2024-08-15', '07:00:00', '08:00:00', 'Scheduled'),
(2, 2, '2024-08-16', '09:00:00', '09:45:00', 'Scheduled'),
(3, 3, '2024-08-17', '10:00:00', '10:30:00', 'Scheduled'),
(3, 4, '2024-08-18', '11:00:00', '12:00:00', 'Scheduled'),
(2, 5, '2024-08-19', '12:30:00', '13:00:00', 'Scheduled'),
(2, 6, '2024-08-20', '14:00:00', '14:45:00', 'Scheduled'),
(3, 7, '2024-08-21', '15:00:00', '16:00:00', 'Scheduled'),
(3, 8, '2024-08-22', '17:00:00', '18:00:00', 'Scheduled'),
(2, 9, '2024-08-23', '18:30:00', '19:30:00', 'Scheduled'),
(3, 10, '2024-08-24', '20:00:00', '20:45:00', 'Scheduled');

INSERT INTO coachWorkouts (coach_id, filial_id, date, start_time, end_time) VALUES
(2, 1, '2024-08-15', '07:00:00', '08:00:00'),
(2, 1, '2024-08-16', '07:00:00', '09:45:00'),
(3, 1, '2024-08-17', '07:00:00', '10:30:00'),
(3, 1, '2024-08-18', '07:00:00', '12:00:00'),
(2, 1, '2024-08-19', '07:30:00', '13:00:00'),
(2, 2, '2024-08-20', '07:00:00', '14:45:00'),
(3, 2, '2024-08-21', '07:00:00', '16:00:00'),
(3, 2, '2024-08-22', '07:00:00', '18:00:00'),
(2, 2, '2024-08-23', '07:30:00', '19:30:00'),
(3, 2, '2024-08-24', '07:00:00', '20:45:00');

INSERT INTO clientsToWorkout (client_id, workout_id) VALUES
(4, 1),
(5, 2),
(6, 3),
(7, 4),
(8, 5),
(9, 6),
(4, 7),
(5, 8),
(6, 9),
(7, 10);
