package configModule

import (
	"fmt"
	"os"

	xlsxTable "game-message-core/xlsxTableData"

	"github.com/Meland-Inc/game-services/src/common/gormDB"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"gorm.io/gorm"
)

var configMgr *ConfigDataManager

func ConfigMgr() *ConfigDataManager { return configMgr }

type LoadFunc struct {
	Name string
	F    func() error
}

type ConfigDataManager struct {
	configDb  *gorm.DB
	loadFuncs []LoadFunc

	taskCnf     map[int32]*xlsxTable.TaskTableRow
	taskListCnf map[int32]*xlsxTable.TaskListTableRow
	roleLvCnf   map[int32]*xlsxTable.RoleLvTableRow
	slotLvCnf   map[int32][]*xlsxTable.SlotLvTableRow
}

func Init() error {
	mgr := &ConfigDataManager{}

	if err := initDB(mgr); err != nil {
		serviceLog.Error("init config data DB fail err: %v", err)
		return err
	}

	mgr.registerLoadFunctions()

	if err := mgr.load(); err != nil {
		return err
	}

	configMgr = mgr
	return nil
}

func (mgr *ConfigDataManager) registerLoadFunctions() {
	mgr.loadFuncs = []LoadFunc{
		LoadFunc{"task", mgr.initTask},
		LoadFunc{"taskList", mgr.initTaskList},
		LoadFunc{"roleLv", mgr.initRoleLv},
		LoadFunc{"slotLv", mgr.initSlotLv},
	}
}

func initDB(mgr *ConfigDataManager) (err error) {
	host := os.Getenv("MELAND_CONFIG_DB_HOST")
	port := os.Getenv("MELAND_CONFIG_DB_PORT")
	user := os.Getenv("MELAND_CONFIG_DB_USER")
	password := os.Getenv("MELAND_CONFIG_DB_PASS")
	dbName := os.Getenv("MELAND_CONFIG_DB_DATABASE")

	mgr.configDb, err = gormDB.InitGormDB(host, port, user, password, dbName, xlsxTable.TableModels())
	return err
}

func (mgr *ConfigDataManager) load() error {
	errorExist := false
	for _, lf := range mgr.loadFuncs {
		if err := lf.F(); err != nil {
			serviceLog.Error("load config table [%s] failed err: %v", lf.Name, err)
			errorExist = true
		}
	}
	if errorExist {
		return fmt.Errorf("load config tables fail")
	}

	return nil
}
