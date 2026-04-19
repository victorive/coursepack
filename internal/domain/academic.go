package domain

import "time"

type University struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string    `gorm:"column:name;not null;uniqueIndex"`
	Acronym   string    `gorm:"column:acronym;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type Department struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UniversityID int64     `gorm:"column:university_id;not null"`
	Name         string    `gorm:"column:name;not null"`
	Code         string    `gorm:"column:code;not null"`
	Faculty      *string   `gorm:"column:faculty"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type Semester struct {
	ID        int64      `gorm:"column:id;primaryKey;autoIncrement"`
	ClassID   int64      `gorm:"column:class_id;not null"`
	Session   string     `gorm:"column:session;not null"`
	Kind      string     `gorm:"column:kind;not null"`
	IsCurrent bool       `gorm:"column:is_current;not null;default:false"`
	StartsAt  *time.Time `gorm:"column:starts_at"`
	EndsAt    *time.Time `gorm:"column:ends_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

type Level struct {
	ID    int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Value int16  `gorm:"column:value;not null;uniqueIndex"`
	Label string `gorm:"column:label;not null"`
}

type Class struct {
	ID                  int64     `gorm:"column:id;primaryKey;autoIncrement"`
	DepartmentID     int64     `gorm:"column:department_id;not null"`
	AdmissionSession string    `gorm:"column:admission_session;not null"`
	CurrentLevelID   int64     `gorm:"column:current_level_id;not null"`
	Name                string    `gorm:"column:name;not null"`
	InviteCode          string    `gorm:"column:invite_code;uniqueIndex;not null"`
	CreatedBy           *int64    `gorm:"column:created_by"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type JoinRequest struct {
	ID          int64             `gorm:"column:id;primaryKey;autoIncrement"`
	ClassID     int64             `gorm:"column:class_id;not null"`
	UserID      int64             `gorm:"column:user_id;not null"`
	Status      JoinRequestStatus `gorm:"column:status;not null;default:pending"`
	ReviewedBy  *int64            `gorm:"column:reviewed_by"`
	RequestedAt time.Time         `gorm:"column:requested_at;autoCreateTime"`
	ResolvedAt  *time.Time        `gorm:"column:resolved_at"`
}

type ClassMember struct {
	ID       int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ClassID  int64     `gorm:"column:class_id;not null"`
	UserID   int64     `gorm:"column:user_id;not null"`
	User     *User     `gorm:"foreignKey:UserID;references:ID"`
	JoinedAt time.Time `gorm:"column:joined_at;autoCreateTime"`
}

type Course struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	ClassID    int64     `gorm:"column:class_id;not null"`
	SemesterID int64     `gorm:"column:semester_id;not null"`
	Code       string    `gorm:"column:code;not null"`
	Title      string    `gorm:"column:title;not null"`
	Units      int16     `gorm:"column:units;not null;default:2"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

type Lecture struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`
	CourseID     int64     `gorm:"column:course_id;not null"`
	Topic        string    `gorm:"column:topic;not null"`
	Description  *string   `gorm:"column:description"`
	LecturerName *string   `gorm:"column:lecturer_name"`
	WeekNumber   *int16    `gorm:"column:week_number"`
	Kind         string    `gorm:"column:kind;not null;default:lecture"`
	Status       string    `gorm:"column:status;not null;default:published"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
