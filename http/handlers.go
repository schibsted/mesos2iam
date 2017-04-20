package http

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/amazon-ecs-agent/agent/credentials"
	"github.com/docker/distribution/uuid"
	"github.schibsted.io/spt-infrastructure/mesos2iam.git/pkg"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func NewSecurityRequestHandler(finder pkg.JobFinder, httpClient *http.Client, smaugUrl string) *SecurityRequestHandler {
	return &SecurityRequestHandler{
		finder,
		httpClient,
		smaugUrl,
	}
}

type SecurityRequestHandler struct {
	JobFinder pkg.JobFinder
	netClient *http.Client
	smaugUrl  string
}

func (h *SecurityRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jobId, err := h.JobFinder.FindJobIdFromRequest(r)

	if err != nil {
		errorMessage := fmt.Sprintf("Error getting JobId from http request: %s", err)
		w.WriteHeader(400)
		writeErrorResponse(errorMessage, 400, w)
		log.Error(errorMessage)
		return
	}

	_, err = uuid.Parse(jobId)
	if err != nil {
		errorMessage := "Invalid JobId in http request: " + jobId
		writeErrorResponse(errorMessage, 400, w)
		log.Error(errorMessage)
		return
	}

	log.Debug("JobId found: " + jobId)

	response, err := h.netClient.Get(fmt.Sprintf("%s/credentials/%s", h.smaugUrl, jobId))

	if err != nil {
		errorMessage := fmt.Sprintf("Couldn't get credentials from Smaug: %s", err.Error())
		writeErrorResponse(errorMessage, 500, w)
		log.Error(errorMessage)
		return
	}

	buf, _ := ioutil.ReadAll(response.Body)

	var creds = credentials.IAMRoleCredentials{}
	json.Unmarshal(buf, &creds)

	log.Debug(string(buf[:]))

	w.Header().Add("Content-Type", "application/json")
	w.Write(buf)
}

func writeErrorResponse(errorMessage string, returnCode int, writer http.ResponseWriter) {
	log.Error(errorMessage)
	writer.WriteHeader(returnCode)
	writer.Write([]byte(errorMessage))
}

func LogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logWriter := &logResponseWriter{w, 200}

		defer func() {
			if e := recover(); e != nil {
				log.Panic("Panic in request handler: ", e)
				logWriter.WriteHeader(http.StatusInternalServerError)
			}

			elapsed := time.Since(start)
			log.Infof("%s \"%s %s %s\" %d %s", remoteIP(r.RemoteAddr), r.Method, r.URL.Path, r.Proto, logWriter.Status, elapsed)
		}()

		handler.ServeHTTP(logWriter, r)
	})
}

type logResponseWriter struct {
	Wrapped http.ResponseWriter
	Status  int
}

func (t *logResponseWriter) Header() http.Header {
	return t.Wrapped.Header()
}

func (t *logResponseWriter) Write(d []byte) (int, error) {
	return t.Wrapped.Write(d)
}

func (t *logResponseWriter) WriteHeader(s int) {
	t.Wrapped.WriteHeader(s)
	t.Status = s
}

func remoteIP(addr string) string {
	index := strings.Index(addr, ":")

	if index < 0 {
		return addr
	}

	return addr[:index]
}
