package spider

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/muniere/screw/internal/pkg/sys"
)

type DownloadOptions struct {
	Prefix      string
	Concurrency int
	Blocking    bool
	Overwrite   bool
	DryRun      bool
}

var SkipDownload = errors.New("skip download")

type command struct {
	value *url.URL
	index int
	total int
}

func (c command) format(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return fmt.Sprintf("[%03d/%03d]%s", c.index+1, c.total, msg)
}

func Download(uris []*url.URL, options DownloadOptions) error {
	ch := make(chan command, len(uris))
	wg := &sync.WaitGroup{}

	// workers
	for i := 0; i < options.Concurrency; i++ {
		wg.Add(1)
		go work(ch, options, wg)
	}

	// enqueue
	for i, u := range uris {
		ch <- command{
			value: u,
			index: i,
			total: len(uris),
		}
	}
	close(ch)

	// join
	wg.Wait()

	return nil
}

func work(ch chan command, options DownloadOptions, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		cmd, ok := <-ch
		if !ok {
			return
		}

		err := download(cmd, options)
		if err == SkipDownload {
			continue
		}
		if err != nil {
			log.Warn(err)
		}

		log.Debug(cmd.format("[ACTION] Sleep 500 milliseconds for next job"))
		time.Sleep(500 * time.Millisecond)
	}
}

func download(cmd command, options DownloadOptions) error {
	log.Debug(cmd.format("[ACTION] Try to download from URL: %s", cmd.value.String()))

	base := filepath.Base(cmd.value.String())

	var path string

	if options.Prefix != "" {
		path = filepath.Join(options.Prefix, base)
	} else {
		path = base
	}

	if !options.Overwrite && sys.Exists(path) {
		log.Info(cmd.format("[SKIP  ] File already exists: %s", path))
		return SkipDownload
	}

	log.Info(cmd.format("[START ] %s => %s", cmd.value.String(), path))

	if options.DryRun {
		log.Info(cmd.format("[FINISH] %s => %s", cmd.value.String(), path))
		return SkipDownload
	}

	log.Debug(cmd.format("[ACTION] Get contents of URI: %s", cmd.value.String()))

	res, err := http.Get(cmd.value.String())
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	log.Debug(cmd.format("[ACTION] Create a file at path: %s", cmd.value.String()))

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	log.Info(cmd.format("[FINISH] %s => %s", cmd.value.String(), path))

	return nil
}
