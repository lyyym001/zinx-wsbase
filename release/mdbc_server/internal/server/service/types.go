package service

type UserLoginRequest struct {
	Username    string                    `json:"username,omitempty"`
	AccountType uint32                    `json:"accountType"` //0-学生 1-老师
	NickName    string                    `json:"nickname"`
	Ip          string                    `json:"ip"`
	Members     string                    `json:"members"`
	Dirs        map[string]DirInfoRequest `json:"dirs"`
}

type DirInfoRequest struct {
	Did   int    `json:"did"`
	Sort  int    `json:"sort"`
	DName string `json:"dname"`
}

type UserRegisterRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
	NickName string `json:"nickname"`
	Role     uint32 `json:"role"`
}

type DeviceListResponse struct {
	Ds map[string]DeviceInfoResponse `json:"ds"`
}

type DeviceInfoResponse struct {
	Username string `json:"username,omitempty"`
	Status   int32  `json:"status"`
	Battery  int32  `json:"battery"`
	Ip       string `json:"ip"`
	Free     int32  `json:"free"`
}

type DeviceDetailRequest struct {
	Username string `json:"username,omitempty"`
}

type DeviceDetailResponse struct {
	Username string `json:"username,omitempty"`
	Ip       string `json:"ip"`
	Dir      int32  `json:"dir"`
}

type DirListResponse struct {
	DirVersion int64                      `json:"dirVersion"`
	Ds         map[string]DirInfoResponse `json:"ds"`
}

type DirInfoResponse struct {
	Did   string                        `json:"did,omitempty"`
	Sort  int                           `json:"sort"`
	DName string                        `json:"dname"`
	Cs    map[string]CourseInfoResponse `json:"cs"`
}

type CourseInfoResponse struct {
	Rid   string `json:"rid,omitempty"`
	RName string `json:"rname"`
}

type DirNameSetRequest struct {
	Did   int    `json:"did,omitempty"`
	DName string `json:"dname"`
}

type DirCourseRequest struct {
	Cid string `json:"cid,omitempty"`
	Did int    `json:"did"`
}

type ScCourseRequest struct {
	Cid    string `json:"cid,omitempty"`
	Status int    `json:"status"`
	Date   string `json:"date"`
	CName  string `json:"cname"`
}

type ScCourseListRequest struct {
	Scs []ScCourseRequest `json:"scs"`
}

type UserListResponse struct {
	Us []string `json:"us"`
}

type Tcp_Tj1 struct {
	Cid   string `json:"Cid,omitempty"`
	Mode2 int32  `json:"Mode2,omitempty"`
	Mode3 int32  `json:"Mode3,omitempty"`
	Page  int32  `json:"Page,omitempty"`
}

type Tcp_Tj1Data struct {
	MaxPage int32          `json:"MaxPage,omitempty"`
	Data    []*Tcp_Tj1Info `json:"Data,omitempty"`
}

type Tcp_Tj1Info struct {
	Date     string  `json:"Date,omitempty"`
	Mode     int32   `json:"Mode,omitempty"`
	Number   int32   `json:"Number,omitempty"`
	MaxScore float32 `json:"MaxScore,omitempty"`
	Score    float32 `json:"Score,omitempty"`
	UniqueId int32   `json:"UniqueId,omitempty"`
}

type Tcp_Tj2 struct {
	UniqueId int32 `json:"UniqueId,omitempty"`
	Page     int32 `json:"Page,omitempty"`
}

type Tcp_Tj2Data struct {
	MaxPage   int32          `json:"MaxPage,omitempty"`
	MaxNumber int32          `json:"MaxNumber,omitempty"`
	Data      []*Tcp_Tj2Info `json:"Data,omitempty"`
}

type Tcp_Tj2Info struct {
	UName string  `json:"UName,omitempty"`
	Score float32 `json:"Score,omitempty"`
}

type Tcp_Tj3 struct {
	UniqueId int32  `json:"UniqueId,omitempty"`
	UName    string `json:"UName,omitempty"`
	WorkName string `json:"WorkName,omitempty"`
	Mode     string `json:"Mode,omitempty"`
	MaxScore int32  `json:"MaxScore,omitempty"`
	Number   int32  `json:"Number,omitempty"`
}

type Tcp_Tj3Data struct {
	WorkName string         `json:"WorkName,omitempty"`
	Mode     string         `json:"Mode,omitempty"`
	MaxScore float32        `json:"MaxScore,omitempty"`
	Number   int32          `json:"Number,omitempty"`
	Data     []*Tcp_Tj3Info `json:"Data,omitempty"`
}

type Tcp_Tj3Info struct {
	UName   string  `json:"UName,omitempty"`
	Date    string  `json:"Date,omitempty"`
	State   int32   `json:"State,omitempty"`
	Type    int32   `json:"Type,omitempty"`
	Score   float32 `json:"Score,omitempty"`
	Content string  `json:"Content,omitempty"`
}

type Tcp_WorkInfoRecord struct {
	Username string `json:"Username,omitempty"`
	Uname    string `json:"Uname,omitempty"`
	Date     string `json:"Date,omitempty"`
	Type     int32  `json:"Type,omitempty"`
	SetId    int32  `json:"SetId,omitempty"`
	State    int32  `json:"State,omitempty"`
	Score    int32  `json:"Score,omitempty"`
	Content  string `json:"Content,omitempty"`
	UniqueId int32  `json:"UniqueId,omitempty"`
}

type Tcp_ResponseRecordInfo struct {
	Code     int32 `json:"Code,omitempty"`
	RecordId int32 `json:"RecordId,omitempty"`
}

type Tcp_WorkStatus struct {
	TeacherStatus int32 `json:"TeacherStatus,omitempty"` //0-老师不在线 1-老师在线
	AutoStatus    int32 `json:"AutoStatus,omitempty"`
}

//type Tcp_Step struct {
//	StepId    int32  `json:"StepId,omitempty"`
//	StepDate  string `json:"StepDate,omitempty"`
//	StepState int32  `json:"StepState,omitempty"`
//	UName     string `json:"UName,omitempty"`
//}
//
//type Tcp_QuestionInfo struct {
//	Qid       int32  `json:"Qid,omitempty"`
//	Code      int32  `json:"Code,omitempty"`
//	StepDate  string `json:"StepDate,omitempty"`
//	StepState int32  `json:"StepState,omitempty"`
//	UName     string `json:"UName,omitempty"`
//}
