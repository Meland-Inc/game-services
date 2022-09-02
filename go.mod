module github.com/Meland-Inc/game-services

go 1.17

// 使用本地go代码仓库方式: https://zhuanlan.zhihu.com/p/109828249
require game-message-core v0.0.0

replace game-message-core => ./src/game-message-core/messageGo
