# Event

Event-driven architecture utilities and implementations for Go applications.

## Overview

This package provides tools and libraries for building event-driven systems in Go. It includes implementations of industry-standard event formats and protocols to enable loosely-coupled, scalable microservices architectures.

## Packages

### [CloudEvents](cloudevents/)

CloudEvents client and server implementation for standardized event messaging.

**What is CloudEvents?**

CloudEvents is a specification for describing event data in a common format, enabling interoperability across different services, platforms, and cloud providers. It provides a consistent way to produce and consume events regardless of the underlying transport or programming language.

**Features:**
- HTTP client and server for CloudEvents
- Request-response pattern support
- Asynchronous event receiver
- Result types for delivery status
- CloudEvents v1.0 specification compliance

**Quick Example:**
```go
import "github.com/common-library/go/event/cloudevents"

// Send event
client, _ := cloudevents.NewHttp("http://localhost:8080", nil, nil)
event := cloudevents.NewEvent()
event.SetType("com.example.user.created")
event.SetSource("users-service")
event.SetData("application/json", userData)
result := client.Send(event)

// Receive events
receiver.StartReceiver(handler, failureFunc)
```

**[Full Documentation →](cloudevents/)**

---

## Event-Driven Architecture Concepts

### What is Event-Driven Architecture?

Event-driven architecture (EDA) is a software design pattern where components communicate by producing and consuming events. Events represent state changes or significant occurrences in a system.

**Key Characteristics:**
- **Loose Coupling** - Components don't need to know about each other
- **Asynchronous** - Event producers don't wait for consumers
- **Scalable** - Easy to add new consumers without changing producers
- **Resilient** - Failures in one component don't cascade

### Event Patterns

#### 1. Fire-and-Forget (Event Notification)

Producer sends an event without expecting a response.

```go
// Producer
event := cloudevents.NewEvent()
event.SetType("order.placed")
event.SetData("application/json", orderData)
client.Send(event)

// Consumer
handler := func(ctx context.Context, event cloudevents.Event) {
    if event.Type() == "order.placed" {
        processOrder(event.Data())
    }
}
```

**Use Cases:**
- Logging and audit trails
- Analytics and metrics
- Notifications (email, SMS)
- Cache invalidation

#### 2. Request-Response (Event-Carried State Transfer)

Producer sends an event and waits for a response.

```go
// Request
request := cloudevents.NewEvent()
request.SetType("order.status.query")
request.SetData("application/json", map[string]string{"order_id": "123"})

response, result := client.Request(request)
// Process response
```

**Use Cases:**
- Synchronous queries
- Validation requests
- Command execution with feedback

#### 3. Publish-Subscribe

Multiple consumers receive the same event.

```go
// Publisher
client.Send(event)

// Multiple subscribers
receiver1.StartReceiver(analyticsHandler, failureFunc)
receiver2.StartReceiver(emailHandler, failureFunc)
receiver3.StartReceiver(auditHandler, failureFunc)
```

**Use Cases:**
- Broadcast notifications
- Multi-step workflows
- Event sourcing
- CQRS (Command Query Responsibility Segregation)

## Choosing the Right Event Package

### Use CloudEvents When:

- Building microservices that need to interoperate
- Working with cloud-native applications
- Need standardized event format across teams/organizations
- Integrating with third-party services
- Building event-driven APIs

### Event Format Comparison

| Format | Standardization | Interoperability | Flexibility | Use Case |
|--------|----------------|------------------|-------------|----------|
| **CloudEvents** | ✅ Industry standard | ✅ High | ✅ High | Microservices, APIs |
| Custom JSON | ❌ None | ⚠️ Low | ✅ Very High | Internal services |
| Protocol Buffers | ⚠️ Google standard | ⚠️ Medium | ⚠️ Medium | High performance |

## Quick Start Guide

### 1. Install Package

```bash
go get -u github.com/common-library/go/event/cloudevents
```

### 2. Create Event Producer

```go
package main

import (
    "log"
    "github.com/common-library/go/event/cloudevents"
)

func main() {
    client, err := cloudevents.NewHttp("http://event-gateway:8080", nil, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    event := cloudevents.NewEvent()
    event.SetType("com.example.event.created")
    event.SetSource("example-service")
    event.SetData("application/json", map[string]string{
        "message": "Hello, Events!",
    })
    
    result := client.Send(event)
    if result.IsUndelivered() {
        log.Fatal(result.Error())
    }
    
    log.Println("Event sent successfully")
}
```

### 3. Create Event Consumer

```go
package main

import (
    "context"
    "log"
    "github.com/common-library/go/event/cloudevents"
    "github.com/cloudevents/sdk-go/v2/protocol/http"
)

func main() {
    httpOpts := []http.Option{http.WithPort(8080)}
    receiver, err := cloudevents.NewHttp("", httpOpts, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    handler := func(ctx context.Context, event cloudevents.Event) {
        log.Printf("Received: %s", event.Type())
        log.Printf("Data: %v", event.Data())
    }
    
    receiver.StartReceiver(handler, func(err error) {
        log.Fatal(err)
    })
    defer receiver.StopReceiver()
    
    // Keep running...
    select {}
}
```

## Common Use Cases

### 1. Microservices Communication

```go
// Order Service → Payment Service → Fulfillment Service

// Order Service (Producer)
event := cloudevents.NewEvent()
event.SetType("com.example.order.created")
event.SetSource("order-service")
event.SetData("application/json", orderDetails)
client.Send(event)

// Payment Service (Consumer & Producer)
handler := func(ctx context.Context, event cloudevents.Event) {
    if event.Type() == "com.example.order.created" {
        processPayment(event.Data())
        
        // Emit next event
        paymentEvent := cloudevents.NewEvent()
        paymentEvent.SetType("com.example.payment.completed")
        client.Send(paymentEvent)
    }
}
```

### 2. Event Sourcing

```go
// Store events as source of truth
events := []cloudevents.Event{
    createEvent("account.created", accountData),
    createEvent("funds.deposited", depositData),
    createEvent("funds.withdrawn", withdrawalData),
}

// Rebuild state from events
currentState := replayEvents(events)
```

### 3. CQRS (Command Query Responsibility Segregation)

```go
// Command (Write Model)
commandEvent := cloudevents.NewEvent()
commandEvent.SetType("command.user.create")
client.Send(commandEvent)

// Query (Read Model)
queryEvent := cloudevents.NewEvent()
queryEvent.SetType("query.user.get")
response, _ := client.Request(queryEvent)
```

### 4. Audit Logging

```go
// Every action produces an audit event
auditEvent := cloudevents.NewEvent()
auditEvent.SetType("audit.user.login")
auditEvent.SetSource("auth-service")
auditEvent.SetData("application/json", map[string]interface{}{
    "user_id":   "123",
    "timestamp": time.Now(),
    "ip":        "192.168.1.1",
})
client.Send(auditEvent)
```

### 5. Saga Pattern

```go
// Distributed transaction coordination
// 1. Start saga
startSaga := createEvent("saga.order.started", sagaData)
client.Send(startSaga)

// 2. Each service processes and emits next step
// Reserve inventory → Process payment → Ship order

// 3. Compensate on failure
if err != nil {
    compensateEvent := createEvent("saga.order.compensate", sagaID)
    client.Send(compensateEvent)
}
```

## Best Practices

### 1. Event Naming Conventions

```go
// Use reverse DNS notation
// Format: <domain>.<entity>.<action>

// Good
"com.example.user.created"
"com.example.order.shipped"
"com.example.payment.failed"

// Avoid
"user_created"
"ORDER_SHIPPED"
"paymentFailed"
```

### 2. Event Versioning

```go
// Include version in event type or data schema
event.SetType("com.example.user.created.v2")

// Or in data schema
event.SetDataSchema("https://example.com/schemas/user-v2.json")
```

### 3. Idempotency

```go
// Include unique event ID for deduplication
event.SetID(uuid.New().String())

// Consumer tracks processed events
func handler(ctx context.Context, event cloudevents.Event) {
    if alreadyProcessed(event.ID()) {
        return // Skip duplicate
    }
    
    processEvent(event)
    markProcessed(event.ID())
}
```

### 4. Error Handling

```go
// Implement retry logic with backoff
func sendWithRetry(event cloudevents.Event, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        result := client.Send(event)
        if result.IsACK() {
            return nil
        }
        
        if result.IsNACK() {
            return fmt.Errorf("rejected: %s", result.Error())
        }
        
        // Exponential backoff
        time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)
    }
    
    return errors.New("max retries exceeded")
}
```

### 5. Dead Letter Queue

```go
// Handle permanently failed events
func handler(ctx context.Context, event cloudevents.Event) {
    if err := processEvent(event); err != nil {
        sendToDeadLetterQueue(event, err)
    }
}
```

### 6. Monitoring and Observability

```go
// Add tracing and metrics
func handler(ctx context.Context, event cloudevents.Event) {
    span := trace.StartSpan(ctx, "handle_event")
    defer span.End()
    
    metrics.EventsReceived.Inc()
    start := time.Now()
    
    processEvent(event)
    
    metrics.EventProcessingDuration.Observe(time.Since(start).Seconds())
}
```

## Architecture Patterns

### Microservices Event Flow

```
┌─────────────┐      ┌──────────────┐      ┌───────────────┐
│   Service A │─────>│ Event Gateway│─────>│   Service B   │
│  (Producer) │      │   (Broker)   │      │  (Consumer)   │
└─────────────┘      └──────────────┘      └───────────────┘
                            │
                            ├──────────────>┌───────────────┐
                            │               │   Service C   │
                            │               │  (Consumer)   │
                            │               └───────────────┘
                            │
                            └──────────────>┌───────────────┐
                                           │   Service D   │
                                           │  (Consumer)   │
                                           └───────────────┘
```

### Event Sourcing Pattern

```
┌──────────┐    ┌────────────────┐    ┌─────────────┐
│ Commands │───>│  Event Store   │───>│ Event Stream│
└──────────┘    │ (Append Only)  │    └─────────────┘
                └────────────────┘           │
                        │                    │
                        ▼                    ▼
                ┌──────────────┐     ┌──────────────┐
                │  Read Model  │     │ Projections  │
                │   (Cache)    │     │   (Views)    │
                └──────────────┘     └──────────────┘
```

## Performance Considerations

### Throughput

| Pattern | Throughput | Latency | Use Case |
|---------|-----------|---------|----------|
| Fire-and-Forget | High | Low | Async notifications |
| Request-Response | Medium | Medium | Sync operations |
| Batch Processing | Very High | High | Analytics |

### Optimization Tips

1. **Connection Pooling** - Reuse HTTP connections
2. **Async Processing** - Use StartReceiver for non-blocking receipt
3. **Batch Events** - Group related events when possible
4. **Event Filtering** - Filter events early to reduce processing
5. **Compression** - Use gzip for large event payloads

## Troubleshooting

### Event Not Received

```go
// Check network connectivity
result := client.Send(event)
if result.IsUndelivered() {
    log.Printf("Network error: %s", result.Error())
}

// Verify receiver is running
// Check firewall/port settings
// Verify event gateway address
```

### Event Rejected

```go
// Check event structure
if result.IsNACK() {
    statusCode, _ := result.GetHttpStatusCode()
    log.Printf("Rejected: %d - %s", statusCode, result.Error())
}

// Validate required CloudEvents fields
// Check data schema
// Verify authentication/authorization
```

### High Latency

```go
// Add timeout configuration
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Monitor event processing time
// Check network latency
// Optimize event handler
```

## Migration Guide

### From Custom Events to CloudEvents

```go
// Before: Custom event format
type CustomEvent struct {
    Type      string
    Timestamp time.Time
    Data      interface{}
}

// After: CloudEvents
event := cloudevents.NewEvent()
event.SetType("com.example.custom.event")
event.SetTime(time.Now())
event.SetData("application/json", data)
```

### From Message Queue to CloudEvents

```go
// Before: Direct message queue
mqClient.Publish("events", message)

// After: CloudEvents over HTTP
event := cloudevents.NewEvent()
event.SetType("com.example.message")
event.SetData("application/json", message)
client.Send(event)
```

## Dependencies

- `github.com/cloudevents/sdk-go/v2` - Official CloudEvents SDK

## Further Reading

- [CloudEvents Specification](https://cloudevents.io/)
- [Event-Driven Architecture](https://martinfowler.com/articles/201701-event-driven.html)
- [Microservices Patterns](https://microservices.io/patterns/index.html)
- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)
- [Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html)
- [Saga Pattern](https://microservices.io/patterns/data/saga.html)
