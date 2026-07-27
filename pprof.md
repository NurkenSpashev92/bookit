 В Flame Graph ищешь толстые горизонтальные блоки — это hotspots.

  Дополнительные endpoint'ы pprof

  # CPU 30s
  http://localhost:6060/debug/pprof/profile?seconds=30

  # Heap (память)
  http://localhost:6060/debug/pprof/heap

  # Goroutines (deadlocks/leaks)
  http://localhost:6060/debug/pprof/goroutine

  # Mutex contention (блокировки)
  http://localhost:6060/debug/pprof/mutex

  # Block (sync блокировки I/O)
  http://localhost:6060/debug/pprof/block

  # Trace (детальная трассировка событий 5s)
  curl -o trace.out 'http://localhost:6060/debug/pprof/trace?seconds=5'
  go tool trace trace.out

  Если хочешь без UI

  # текстовый top-20 функций
  go tool pprof -top -cum 'http://localhost:6060/debug/pprof/profile?seconds=30'

  # или сохранить и потом смотреть
  curl -o cpu.pprof 'http://localhost:6060/debug/pprof/profile?seconds=30'
  go tool pprof -top -cum cpu.pprof
