package mutation

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mjpitz/cluster-preset/internal/config"
	"github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/pkg/apis/core/v1"
)

var (
	errNoBody = fmt.Errorf("empty body")
	errInvalidBody = fmt.Errorf("invalid json payload")

	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

	// (https://github.com/kubernetes/kubernetes/issues/57982)
	defaulter = runtime.ObjectDefaulter(runtimeScheme)
)

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1beta1.AddToScheme(runtimeScheme)
	// defaulting with webhooks:
	// https://github.com/kubernetes/kubernetes/issues/57982
	_ = v1.AddToScheme(runtimeScheme)
}

// RegisterMutateWebhook manages binding endpoint invocations to underlying business logic.
func RegisterMutateWebhook(server *gin.Engine, holder *config.Holder) {
	webhook := &mutateWebhook{holder }

	server.POST("/mutate", func(ctx *gin.Context) {
		webhook.mutationHandler(ctx)
	})
}

type mutateWebhook struct {
	Holder *config.Holder
}

func (w *mutateWebhook) mutationHandler(ctx *gin.Context) {
	var body []byte
	if ctx.Request.Body != nil {
		if data, err := ioutil.ReadAll(ctx.Request.Body); err == nil {
			body = data
		}
	}

	if len(body) == 0 {
		if err := ctx.AbortWithError(http.StatusBadRequest, errNoBody); err != nil {
			logrus.Errorf("failed to abort request: %s", err.Error())
		}
		return
	}

	review := &v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, review); err != nil {
		if err := ctx.AbortWithError(http.StatusBadRequest, errInvalidBody); err != nil {
			logrus.Errorf("failed to abort request: %s", err.Error())
		}
		return
	}

	admissionResponse := mutate(w.Holder.Get(), review)

	if admissionResponse != nil && review.Request != nil {
		admissionResponse.UID = review.Request.UID
	}

	ctx.JSON(http.StatusOK, &v1beta1.AdmissionReview{
		Response: admissionResponse,
	})
}
