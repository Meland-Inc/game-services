package configData

import (
	"fmt"
	"os"

	"game-message-core/proto"
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

	taskCnf      map[int32]xlsxTable.TaskTableRow
	taskListCnf  map[int32]xlsxTable.TaskListTableRow
	roleLvCnf    map[int32]xlsxTable.RoleLvTableRow
	slotLvCnf    map[int32][]xlsxTable.SlotLvTableRow
	rewardCnf    map[int32]xlsxTable.RewardTableRow
	itemCnf      map[int32]xlsxTable.ItemTable
	dropCnf      map[int32]xlsxTable.DropTableRow
	chatCnf      map[proto.ChatChannelType]xlsxTable.ChatTableRow
	gameValueCnf map[int32]xlsxTable.GameValueTable
	sceneAreaCnf map[int32]xlsxTable.SceneAreaRow
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
		LoadFunc{"reward", mgr.initReward},
		LoadFunc{"item", mgr.initItem},
		LoadFunc{"drop", mgr.initDrop},
		LoadFunc{"chat", mgr.initChatCnf},
		LoadFunc{"gameValue", mgr.initGameValue},
		LoadFunc{"SceneArea", mgr.initSceneArea},
	}
}

func initDB(mgr *ConfigDataManager) (err error) {
	host := os.Getenv("GAME_CONFIG_DB_HOST")
	port := os.Getenv("GAME_CONFIG_DB_PORT")
	user := os.Getenv("GAME_CONFIG_DB_USER")
	password := os.Getenv("GAME_CONFIG_DB_PASS")
	dbName := os.Getenv("GAME_CONFIG_DB_DATABASE")

	serviceLog.Info("game config DB url:[%+s] DbName:[%s]", host, dbName)

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
