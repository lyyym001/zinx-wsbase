package core

//交互资源对象
type JHObject struct {

	//下面是new
	ObjId       int					//资源ID
	InteractiveType 		int		//交互类型(0-不可交互，1-可交互),
	Tb 			int					//同步类型(0-不可同步，1-可同步)
	Visiable	int 				//显隐状态(0-隐藏 1-显示 2-销毁)
	X    float32            //坐标x
	Y    float32            //坐标y
	Z    float32            //坐标z
	RX    float32            //旋转x
	RY    float32            //旋转y
	RZ    float32            //旋转z

}
