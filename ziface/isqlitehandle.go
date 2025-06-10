package ziface

import "database/sql"

type ISqliteHandle interface {


	GetDB() *sql.DB	//


}
