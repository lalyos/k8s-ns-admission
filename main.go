package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/tidwall/gjson"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config contains the server (the webhook) cert and key.
type Config struct {
	CertFile string
	KeyFile  string
}

func (c *Config) addFlags() {
	flag.StringVar(&c.CertFile, "tls-cert-file", c.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")
	flag.StringVar(&c.KeyFile, "tls-private-key-file", c.KeyFile, ""+
		"File containing the default x509 private key matching --tls-cert-file.")
}

type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

var lastInput string

func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	glog.Info("---> Webhook::serve")
	glog.Error("---> Webhook::serve")
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	lastInput = string(body)
	glog.Info("===>INPUT: ", lastInput)

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	var reviewResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		glog.Error(err)
		reviewResponse = toAdmissionResponse(err)
	} else {
		reviewResponse = admit(ar)
	}

	response := v1beta1.AdmissionReview{}
	if reviewResponse != nil {
		response.Response = reviewResponse
		response.Response.UID = ar.Request.UID
	}
	// reset the Object and OldObject, they are not needed in a response.
	//ar.Request.Object = runtime.RawExtension{}
	//ar.Request.OldObject = runtime.RawExtension{}

	resp, err := json.Marshal(response)
	if err != nil {
		glog.Error(err)
	}
	if _, err := w.Write(resp); err != nil {
		glog.Error(err)
	}
}

func serveNamespaces(w http.ResponseWriter, r *http.Request) {
	serve(w, r, admitNamespace)
}

func serveLastInput(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, lastInput)
}
func admitNamespace(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {

	//glog.Infof("---> AdmissinReview.req: %#v", ar.Request)
	//glog.Infof("---> Kind: %#v, Name: %#v, Resource: %#v, UserInfo: %#v",
	//	ar.Request.Kind, ar.Request.Name, ar.Request.Resource, ar.Request.UserInfo)

	resp := &v1beta1.AdmissionResponse{
		Allowed: true,
	}

	user := ar.Request.UserInfo.Username
	resource := ar.Request.Resource.Resource
	glog.Infof("###### Object: %s, ResType: %#v, Owner: %#v", ar.Request.Object.Raw, resource, user)

	objName := ar.Request.Name
	glog.Infof("### Reques.Name: %q", objName)

	if objName == "" {
		objName = gjson.Get(string(ar.Request.Object.Raw), "metadata.name").Str
		glog.Infof("### objname (Request.Object.Raw): %#v", objName)
	}

	if resource == "namespaces" && strings.HasPrefix(user, "system:serviceaccount:default:sa-user") {
		username := strings.TrimPrefix(user, "system:serviceaccount:default:sa-")
		prefix := username + "-"
		allowed := strings.HasPrefix(objName, prefix)
		glog.Infof("user: %s, ns:%s allowed:%v", username, objName, allowed)
		if !allowed {
			resp.Allowed = false
			resp.Result = &metav1.Status{Message: "namespace prefix must match 'userX-'"}
		}
	}
	return resp

}

func main() {
	var config Config
	config.addFlags()
	flag.Parse()

	http.HandleFunc("/ns", serveNamespaces)
	http.HandleFunc("/last", serveLastInput)
	clientset := getClient()
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: configTLS(config, clientset),
	}
	server.ListenAndServeTLS("", "")
}
