package database

type Users struct {
	User_id       int
	Name          string
	Email         string
	Password      string
	Profile_image string
	Deactive      int
	User_level    string
}
type Posts struct {
	Post_id      int
	User_id      int
	Heading      string
	Body         string
	Closed_user  int
	Closed_admin int
	Closed_date  string
	Insert_time  string
	Update_time  string
	Image        string
}

type Comments struct {
	Comment_id  int
	Post_id     int
	User_id     int
	Body        string
	Insert_time string
}

type Categories struct {
	Category_id   int
	Category_Name string
	Closed        int
}

type Reaction struct {
	User_id     int
	Post_id     int
	Comment_id  int
	Reaction_id string
}

type PostCategory struct {
	Category_id int
	Post_id     int
}
type UserLevel struct {
	User_level string
	value      int
}
