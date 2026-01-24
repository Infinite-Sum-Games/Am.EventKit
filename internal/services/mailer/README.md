# 📬 Mailer Service — Persistent, Concurrent Email Dispatcher in Go

This package implements a production-grade mailer system that is:

- ✅ Concurrent (via worker pool)
- ✅ Reliable (via persistent disk queue)
- ✅ Gracefully shutdown-safe
- ✅ Templated and extensible
- ✅ Idiomatic and modular in Go

---

## ✨ Why This Design?

### 1. **Concurrency via Worker Pool**
To send multiple emails without overwhelming resources, a pool of fixed workers are used. Each worker:

- Listens to the persistent queue
- Sends email using a persistent SMTP connection
- Handles errors individually
- Terminates cleanly when the app shuts down

> This prevents spawning unbounded goroutines and allows safe parallelism.

---

### 2. **Persistence with `dque`**

In-memory queues would lose jobs on crash.  
We use [`dque`](https://github.com/joncrlsn/dque), a disk-backed, durable queue to ensure **email jobs survive restarts**.

- No external DB required
- Fast on-disk segment file format
- Custom serialization supported

> Durable queues help when the machine crashes, or SMTP servers are temporarily down.

---

### 3. **One Mailer per Worker**

Each worker holds its own `Mailer` instance (with an SMTP connection). This allows:

- Avoiding connection throttling
- Simplifying sync (no mutex)
- Clean resource release during shutdown

> Instead of one shared SMTP dialer with mutexes, this is cleaner and performs better under load.

---

### 4. **Graceful Shutdown**

We use:

- `context.WithCancel()` for signaling all workers to stop
- `sync.WaitGroup` to wait for in-flight emails to finish
- `os.Signal` in main to trigger shutdown

> Ensures the app doesn't exit while an email is being sent.

---

## 🧱 Architecture Overview

![Mailer architecture diagram](../../public/Mailer%20architecture.png)