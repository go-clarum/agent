package control

import "sync"

var RunningActions sync.WaitGroup
var ShutdownHook sync.WaitGroup
