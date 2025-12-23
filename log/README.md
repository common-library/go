# Log

Logging utilities for Go applications.

## Overview

The log package provides high-level logging abstractions with support for multiple logging backends. It includes wrappers for both structured logging (slog) and Kubernetes-style logging (klog).

## Subpackages

### slog

Structured logging with asynchronous processing, multiple outputs, and daily rotation.

[📖 Documentation](slog/README.md)

**Features:**
- Asynchronous log processing with queue-based buffering
- Multiple output destinations (stdout, stderr, file)
- Daily log file rotation
- Configurable log levels (Trace, Debug, Info, Warn, Error, Fatal)
- Caller information tracking

**Quick Example:**
```go
import "github.com/common-library/go/log/slog"

slog.SetOutputToFile("/var/log/app.log", 7, 100, 100*1024*1024)
slog.Info("Application started")
defer slog.Flush()
```

### klog

Kubernetes-style logging wrapper with structured and formatted logging support.

[📖 Documentation](klog/README.md)

**Features:**
- Kubernetes ecosystem compatibility
- Structured logging (InfoS, ErrorS)
- Formatted logging (Infof, Errorf)
- Caller tracking for debugging
- Fatal logging with application exit

**Quick Example:**
```go
import "github.com/common-library/go/log/klog"

klog.Info("Server started")
klog.InfoS("Request processed", "method", "GET", "path", "/api/users")
defer klog.Flush()
```

## Choosing a Logger

| Feature | slog | klog |
|---------|------|------|
| Async Processing | ✅ | ❌ |
| File Rotation | ✅ | ❌ |
| Kubernetes Integration | ❌ | ✅ |
| Structured Logging | ✅ | ✅ |
| Performance | High (async) | Moderate (sync) |
| Use Case | General apps | K8s controllers |

## Installation

```bash
go get -u github.com/common-library/go/log/slog
go get -u github.com/common-library/go/log/klog
```

## Best Practices

1. **Always Flush** - Call `Flush()` before application exit
2. **Set Appropriate Levels** - Use Debug/Trace for development, Info+ for production
3. **Structured Data** - Prefer structured logging for machine parsing
4. **Caller Info** - Enable in development, consider disabling in production for performance
5. **File Rotation** - Configure appropriate retention for disk space management

## Further Reading

- [slog Package Documentation](slog/README.md)
- [klog Package Documentation](klog/README.md)
