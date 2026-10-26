package domain

type RoleName string

const (
	RoleSuperAdmin RoleName = "super_admin"
	RoleStudent    RoleName = "student"
	RoleClassRep   RoleName = "class_rep"
	RoleCourseRep  RoleName = "course_rep"
)

type PermissionName string

const (
	PermMembersList       PermissionName = "list_members"
	PermMembersManage     PermissionName = "manage_members"
	PermFilesUpload       PermissionName = "upload_files"
	PermFilesApprove      PermissionName = "approve_files"
	PermCoursesManage     PermissionName = "manage_courses"
	PermLecturesManage    PermissionName = "manage_lectures"
	PermManageClasses     PermissionName = "manage_classes"
	PermManageDepartments PermissionName = "manage_departments"
	PermManageSemesters   PermissionName = "manage_semesters"
)

type ScopeType string

const (
	ScopeGlobal ScopeType = "global"
	ScopeClass  ScopeType = "class"
)

type Provider string

const (
	ProviderGoogle Provider = "google"
)

type FileCategory string

const (
	FileCategoryLectureNote  FileCategory = "lecture_note"
	FileCategoryPastQuestion FileCategory = "past_question"
	FileCategoryMaterial     FileCategory = "material"
	FileCategoryAssignment   FileCategory = "assignment"
)

type ReviewAction string

const (
	ReviewActionApprove ReviewAction = "approve"
	ReviewActionReject  ReviewAction = "reject"
)

type FileStatus string

const (
	FileStatusPending  FileStatus = "pending"
	FileStatusApproved FileStatus = "approved"
	FileStatusRejected FileStatus = "rejected"
	FileStatusArchived FileStatus = "archived"
)

// PointsFileApproval is awarded once per approved file.
const PointsFileApproval = 10

type SemesterKind string

const (
	SemesterFirst  SemesterKind = "first"
	SemesterSecond SemesterKind = "second"
)

type LectureKind string

const (
	LectureKindLecture LectureKind = "lecture"
	LectureKindTest    LectureKind = "test"
	LectureKindExam    LectureKind = "exam"
	LectureKindLab     LectureKind = "lab"
)

type LectureStatus string

const (
	LectureStatusPublished LectureStatus = "published"
	LectureStatusArchived  LectureStatus = "archived"
)

type JoinRequestStatus string

const (
	JoinStatusPending  JoinRequestStatus = "pending"
	JoinStatusApproved JoinRequestStatus = "approved"
	JoinStatusRejected JoinRequestStatus = "rejected"
)

type PointSource string

const (
	PointSourceFileApproval PointSource = "file_approval"
	PointSourceManual       PointSource = "manual"
	PointSourceBonus        PointSource = "bonus"
)
