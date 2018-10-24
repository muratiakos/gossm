package printer

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/glassechidna/gossm/pkg/gossm"
	"github.com/mitchellh/go-wordwrap"
	"github.com/nsf/termbox-go"
	"io"
	"os"
	"strings"
)

type Printer struct {
	Out       io.Writer
	outColors []*color.Color
	Err       io.Writer
	errColors []*color.Color
	Quiet     bool
}

func New() *Printer {
	return &Printer{
		Out:       os.Stdout,
		outColors: []*color.Color{color.New(color.FgGreen)},
		Err:       os.Stderr,
		errColors: []*color.Color{color.New(color.FgRed)},
		Quiet:     false,
	}
}

func (p *Printer) PrintInfo(command string, resp *gossm.DoitResponse) {
	p.printInfo("Command: ", command)
	p.printInfo("Command ID: ", resp.CommandId)

	instanceIds := resp.InstanceIds.InstanceIds
	prefix := fmt.Sprintf("Running command on %d instances: ", len(instanceIds))
	p.printInfo(prefix, fmt.Sprintf("%+v", instanceIds))
}

func (p *Printer) printInfo(prefix, info string) {
	if !p.Quiet {
		faint := color.New(color.Faint)
		blue := color.New(color.FgBlue)
		_, _ = fmt.Fprintf(os.Stderr, "%s%s\n", blue.Sprintf("%s", prefix), faint.Sprintf("%s", info))
	}
}

func (p *Printer) Print(msg gossm.SsmMessage) {
	if len(msg.StdoutChunk) > 0 {
		p.print(p.Out, p.outColors[0], msg, msg.StdoutChunk)
	}

	if len(msg.StderrChunk) > 0 {
		p.print(p.Err, p.errColors[0], msg, msg.StderrChunk)
	}

	if !p.Quiet {
		_, _ = fmt.Fprintln(p.Out) // split em out
	}

	if msg.Error != nil {
		panic(msg.Error)
	}
}

func (p *Printer) print(w io.Writer, prefixColor *color.Color, msg gossm.SsmMessage, payload string) {
	if p.Quiet {
		_, _ = fmt.Fprintln(w, payload)
		return
	}

	windowWidth := 80

	if err := termbox.Init(); err == nil {
		windowWidth, _ = termbox.Size()
		termbox.Close()
	}

	prefix := prefixColor.Sprintf("[%s] ", msg.InstanceId)

	outputWidth := windowWidth - len(prefix)
	wrapped := wordwrap.WrapString(payload, uint(outputWidth))
	lines := strings.Split(wrapped, "\n")

	for _, line := range lines {
		_, _ = fmt.Fprintf(w, "%s%s\n", prefix, line)
	}
}
