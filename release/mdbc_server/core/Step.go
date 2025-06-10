package core

//关键步骤得分
type Step struct {

	//下面是new
	StepId int				//步骤ID
	StepDate string			//日期
	StepState int			//操作状态 0-失败 1-成功
	UName string			//操作人名称

}
