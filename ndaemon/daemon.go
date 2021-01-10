package ndaemon

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron"

	"github.com/influx6/npkg/nerror"
	"github.com/takama/daemon"

	"github.com/influx6/npkg/njson"
)

var (
	defaultLocation         = time.Now().Location()
	defaultCronParserOption = cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor
)

type Logger interface {
	Log(json *njson.JSON)
}

type DaemonJob func(ctx context.Context, logger Logger)

type ctxJob struct {
	ctx    context.Context
	logger Logger
	job    DaemonJob
}

func (c *ctxJob) Run() {
	c.job(c.ctx, c.logger)
}

// CronDaemon returns a new cron based daemon which will execute the function
// based on a cron schedule.
//
// See for marker format: http://www.quartz-scheduler.org/documentation/quartz-2.x/tutorials/crontrigger.html
// CRON Expression Format
// A cron expression represents a set of times, using 5 space-separated fields.
// Field name   | Mandatory? | Allowed values  | Allowed special characters
// ----------   | ---------- | --------------  | --------------------------
// Minutes      | Yes        | 0-59            | * / , -
// Hours        | Yes        | 0-23            | * / , -
// Day of month | Yes        | 1-31            | * / , - ?
// Month        | Yes        | 1-12 or JAN-DEC | * / , -
// Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
// Month and Day-of-week field values are case insensitive.  "SUN", "Sun", and "sun" are equally accepted.
// The specific interpretation of the format is based on the Cron Wikipedia page:
// https://en.wikipedia.org/wiki/Cron
//
func CronDaemon(
	ctx context.Context,
	canceler context.CancelFunc,
	cronExpression string,
	location *time.Location,
	overrideParser cron.ParseOption,
	name string,
	desc string,
	logger Logger,
	kind daemon.Kind,
	templateOverride string,
	depends []string,
	do DaemonJob,
) *ServiceDaemon {
	var daemonService = ServiceDaemon{
		Cancel:           canceler,
		Ctx:              ctx,
		Name:             name,
		Desc:             desc,
		Logger:           logger,
		Kind:             kind,
		DependOnServices: depends,
		DaemonTemplate:   templateOverride,
	}

	daemonService.Job = func(ctx context.Context, logger Logger) {
		var logStack = njson.Log(logger)

		var waiter sync.WaitGroup
		waiter.Add(1)
		go func() {
			<-ctx.Done()
			waiter.Done()
		}()

		var parser = cron.NewParser(overrideParser)
		var cronManager = cron.NewWithLocation(location)

		var schedule, err = parser.Parse(cronExpression)
		if err != nil {
			var wrapErr = nerror.WrapOnly(err)
			logStack.New().
				LError().
				Message("failed to parse cron expression").
				String("error", wrapErr.Error()).
				End()
			return
		}

		cronManager.Schedule(schedule, &ctxJob{
			ctx:    ctx,
			logger: logger,
			job:    do,
		})

		cronManager.Start()

		waiter.Wait()
	}

	return &daemonService
}

// Cron returns a cron ServiceDaemon which uses default location and
// default parser options (e.g "* * * * * *") to define the cron expression.
func Cron(
	ctx context.Context,
	canceler context.CancelFunc,
	cronExpression string,
	name string,
	desc string,
	logger Logger,
	kind daemon.Kind,
	do DaemonJob,
	depends ...string,
) *ServiceDaemon {
	return CronDaemon(
		ctx,
		canceler,
		cronExpression,
		defaultLocation,
		defaultCronParserOption,
		name,
		desc,
		logger,
		kind,
		"",
		depends,
		do,
	)
}

type ServiceDaemon struct {
	Name             string
	Desc             string
	Logger           Logger
	Kind             daemon.Kind
	DependOnServices []string

	// Job function which will be lunched into a goroutine of it's own.
	// If the job immediately returns then the context attached will be
	// cancelled and the daemon killed.
	Job DaemonJob

	/***
	DaemonTemplate is a go template string to be used to create the service.
	This is optional, but can be modified to suit usage.

		[Unit]
		Description={{.Description}}
		Requires={{.Dependencies}}
		After={{.Dependencies}}

		[Service]
		PIDFile=/var/run/{{.Name}}.pid
		ExecStartPre=/bin/rm -f /var/run/{{.Name}}.pid
		ExecStart={{.Path}} {{.Args}}
		Restart=on-failure

		[Install]
		WantedBy=multi-user.target
	*/
	DaemonTemplate string

	// Context to use for delivering to operation.
	Ctx context.Context

	// Cancel function to cancel context when application stops
	// or is removed/uninstalled.
	Cancel context.CancelFunc
}

func (c *ServiceDaemon) Run(args []string) (string, error) {
	var logStack = njson.Log(c.Logger)

	var usageHelp = "Usage: " + c.Name + " install | remove | start | stop | status"

	var serviceDaemon, err = daemon.New(c.Name, c.Desc, c.Kind, c.DependOnServices...)
	if err != nil {
		return "", nerror.WrapOnly(err)
	}

	// start service details
	if len(args) == 0 {
		var waiter = WaiterForCtxSignal(c.Ctx, c.Cancel)

		go func() {
			defer c.Cancel()

			c.Job(c.Ctx, c.Logger)
		}()

		waiter.Wait()
		return "", nil
	}

	if len(c.DaemonTemplate) != 0 {
		if err := serviceDaemon.SetTemplate(c.DaemonTemplate); err != nil {
			var wrapErr = nerror.WrapOnly(err)
			logStack.New().
				LError().
				Message("failed to set daemon template").
				String("error", wrapErr.Error()).
				End()
			return "", wrapErr
		}
	}

	switch strings.ToLower(args[0]) {
	case "install":
		var message, opErr = serviceDaemon.Install(args[1:]...)
		if opErr != nil {
			var wrapErr = nerror.WrapOnly(opErr)
			logStack.New().
				LError().
				Message("failed to complete install operation").
				String("message", message).
				String("error", wrapErr.Error()).
				End()
			return "", wrapErr
		}
		logStack.New().
			LInfo().
			Message("installed daemon").
			String("message", message).
			End()
		return message, nil
	case "remove":
		var message, opErr = serviceDaemon.Remove()
		if opErr != nil {
			var wrapErr = nerror.WrapOnly(opErr)
			logStack.New().
				LError().
				Message("failed to complete remove operation").
				String("message", message).
				String("error", wrapErr.Error()).
				End()
			return "", wrapErr
		}
		logStack.New().
			LInfo().
			Message("removed daemon").
			String("message", message).
			End()
		return message, nil
	case "start":
		var message, opErr = serviceDaemon.Start()
		if opErr != nil {
			var wrapErr = nerror.WrapOnly(opErr)
			logStack.New().
				LError().
				Message("failed to complete start operation").
				String("message", message).
				String("error", wrapErr.Error()).
				End()
			return "", wrapErr
		}
		logStack.New().
			LInfo().
			Message("started daemon").
			String("message", message).
			End()
		return message, nil
	case "stop":
		var message, opErr = serviceDaemon.Stop()
		if opErr != nil {
			var wrapErr = nerror.WrapOnly(opErr)
			logStack.New().
				LError().
				Message("failed to complete start operation").
				String("message", message).
				String("error", wrapErr.Error()).
				End()
			return "", wrapErr
		}
		logStack.New().
			LInfo().
			Message("started daemon").
			String("message", message).
			End()
		return message, nil
	case "status":
		var message, opErr = serviceDaemon.Status()
		if opErr != nil {
			var wrapErr = nerror.WrapOnly(opErr)
			logStack.New().
				LError().
				Message("failed to complete start operation").
				String("message", message).
				String("error", wrapErr.Error()).
				End()
			return "", wrapErr
		}
		logStack.New().
			LInfo().
			Message("started daemon").
			String("message", message).
			End()
		return message, nil
	}
	return usageHelp, nil
}
