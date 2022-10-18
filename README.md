# game-services

game used services

## game services dapr config  criterion

    ServicesManager     appPort:(5100~5129)  grpc:(5130~5169) serviceHttp(5170~5199)
    accountService      appPort:(5200~5249)  grpc:(5250~5299)
    mainService         appPort:(5300~5349)  grpc:(5350~5399) 
    taskService         appPort:(5400~5449)  grpc:(5450~5499) 
    chatService         appPort:(5500~5549)  grpc:(5550~5599) 
    agentService        appPort:(5600~5649)  grpc:(5650~5699) serviceSocket(5700~5799)
    sceneService        appPort:(5800~5899)  http:(5900~5999)
