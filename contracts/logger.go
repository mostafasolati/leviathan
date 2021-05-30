package contracts

// LogFields is the date-type for the collection of logging fields.
type LogFields map[string]interface{}

// ILogger is an abstraction over all loggings.
type ILogger interface {
	// Trace designates finer-grained informational events than the Debug.
	Trace(message string)
	// Debug is very verbose logging. Usually only enabled when debugging.
	Debug(message string)
	// Info is the general operational entries about what's going on inside the
	// application.
	Info(message string)
	// Warn is non-critical entries that deserve eyes.
	Warn(message string)
	// Error is used for errors that should definitely be noted. Commonly used
	// for hooks to send errors to an error tracking service.
	Error(message string)
	// Fatal logs and then calls `os.Exit(1)`.
	Fatal(message string)

	// WithFields provides additional information for the log entries. It may be
	// called before any of the above printing methods, e.g:
	//
	//      logger.WithFields(log.Fields{
	//          "transaction_id": transactionID,
	//          "gateway_id":     gatewayID,
	//      }).Info("payment succeeded.")
	//
	// Structured logging through logging fields is preferred upon long
	// unparseable error messages. For example, instead of:
	//
	//      logger.Error(
	// 	        fmt.Sprintf("cannot send %s to %s with %d", event, topic, key),
	//      )
	//
	// you should log the much more discoverable:
	//
	//      log.WithFields(log.Fields{
	//          "event": event,
	//          "topic": topic,
	//          "key":   key,
	//      }).Fatal("cannot send event")
	//
	// Also, in data visualization softwares like Kibana, it is much easier to
	// filter and group log entries using fields rather than text searches.
	WithFields(fields LogFields) ILogger
}
