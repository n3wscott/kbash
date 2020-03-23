package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	// Port the http server on.
	Port string `envconfig:"PORT" default:"8080" required:"true"`

	// FilePath is the base location of the script folder.
	FilePath string `envconfig:"FILE_PATH" default:"/var/run/ko/" required:"true"`

	// Script is the name of the script to run, relative to FilePath.
	Script string `envconfig:"SCRIPT" default:"now.sh" required:"true"`
}

// RunCmd takes in a string and tries its best to run it as a command line
// returning the results.
func RunCmd(cmdLine string) ([]byte, error) {
	cmdSplit := strings.Split(cmdLine, " ")
	cmd := cmdSplit[0]
	args := cmdSplit[1:] // TODO: This will not be correct for passing quoted strings.

	cmdOut, err := exec.Command(cmd, args...).Output() // TODO: could get stderr too in the future.
	return cmdOut, err
}

// Handler will handle the http requests.
type Handler struct{ cmd string }

// ServeHTTP handles an incoming http request.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Run the command.
	out, err := RunCmd(h.cmd)
	// If it ended in an error, return the output, error and 500.
	// TODO: might also want the stderr.
	if err != nil {
		_, _ = fmt.Fprintf(w, "%s\n%s", out, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If out had any length, then 200 or 204 based on output length.
	if len(out) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Send the output of the command as the response of the http request.
	_, err = w.Write(out)
	if err != nil {
		fmt.Printf("[ERROR] unable to write out to ResponseWriter, %s\n", err)
		return
	}
}

func main() {
	// Process the env vars.
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	// Fix up the file path so we can just slap the script after it.
	if !strings.HasSuffix(env.FilePath, "/") {
		env.FilePath = env.FilePath + "/"
	}

	// Create the handler with the full path of the script.
	h := &Handler{cmd: fmt.Sprintf("%s%s", env.FilePath, env.Script)}

	// Start the http server.
	log.Printf("kbash: listening on port %s", env.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", env.Port), h); err != nil {
		log.Fatal(err)
	}
}
