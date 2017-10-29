package command

import (
	"fmt"
	"log"
	"os/exec"
)

type Executer interface {
	ExecuteRpc(cmd *exec.Cmd) (result []byte, err error)
	Execute(cmd *exec.Cmd, verbose bool) (err error)
}

type CommandExecuter struct {
	errLogger *log.Logger
	infLogger *log.Logger
}


func New(errLogger, infLogger *log.Logger) *CommandExecuter {
	return &CommandExecuter{
		errLogger: errLogger,
		infLogger: infLogger,
	}
}

func (me CommandExecuter) Execute(cmd *exec.Cmd, verbose bool) (err error) {
	me.infLogger.Println("Processing message...")

	if verbose {
		cmd.Stdout = NewLogWriter(me.infLogger)
		cmd.Stderr = NewLogWriter(me.errLogger)
		err = cmd.Run()
	} else if out, outErr := cmd.CombinedOutput(); outErr != nil {
		me.errLogger.Printf("Failed: %s\n", string(out[:]))
		err = outErr
	}

	if err != nil {
		me.infLogger.Println("Failed. Check error log for details.")
		me.errLogger.Printf("Error: %s\n", err)

		return fmt.Errorf("Error occured during execution of command: %s", err)
	}

	me.infLogger.Println("Processed!")

	return nil
}

func (me CommandExecuter) ExecuteRpc(cmd *exec.Cmd) (result []byte, err error) {
	me.infLogger.Println("Processing message...")

	out, err := cmd.Output()

	if err != nil {
		me.infLogger.Println("Failed. Check error log for details.")
		me.errLogger.Printf("Failed: %s\n", string(out[:]))
		me.errLogger.Printf("Error: %s\n", err)

		return out, fmt.Errorf("Error occured during execution of command: %s", err)
	}

	me.infLogger.Println("Processed!")

	return out, nil
}

type LogWriter struct {
	logger *log.Logger
}

func NewLogWriter(l *log.Logger) *LogWriter {
	lw := &LogWriter{}
	lw.logger = l
	return lw
}

func (lw LogWriter) Write (p []byte) (n int, err error) {
	lw.logger.SetFlags(0)
	lw.logger.Printf("%s", p)
	lw.logger.SetFlags(log.Ldate|log.Ltime)
	return len(p), nil
}