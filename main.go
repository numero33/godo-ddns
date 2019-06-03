package main

import (
	"context"
	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/digitalocean/godo"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

// --------------------------------------------------------------------------------------------------------------------

type TokenSource struct {
	AccessToken string
}

// --------------------------------------------------------------------------------------------------------------------

func main() {
	if len(os.Getenv("GODOHOST")) == 0 {
		if err := os.Setenv("GODOHOST", ":80"); err != nil {
			logrus.WithError(err).Error("setenv default GODOHOST")
		}
	}

	http.HandleFunc("/update", update)
	logrus.Info("Start.. listen on " + os.Getenv("GODOHOST"))
	logrus.WithError(http.ListenAndServe(os.Getenv("GODOHOST"), nil)).Error("listenAndServe")
}

// --------------------------------------------------------------------------------------------------------------------

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

// --------------------------------------------------------------------------------------------------------------------

func update(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("h")

	if !domainutil.HasSubdomain(host) || len(domainutil.Domain(host)) <= 3 {
		logrus.WithField("host", host).Error("host error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`notfqdn`))
		return
	}

	ip := r.URL.Query().Get("ip")
	if len(r.URL.Query().Get("ip")) < 7 {
		ip = r.RemoteAddr
	}

	tokenSource := &TokenSource{
		AccessToken: r.URL.Query().Get("p"),
	}

	client := godo.NewClient(oauth2.NewClient(context.Background(), tokenSource))
	ctx := context.TODO()

	logrus.WithFields(logrus.Fields{"host": host, "ip": ip}).Info("new update")

	domains, err := domainList(ctx, client, domainutil.Domain(host))
	if err != nil {
		logrus.WithError(err).Error("get domains")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`badauth`))
		return
	}

	for _, d := range domains {
		if d.Name == domainutil.Subdomain(host) {

			if d.Data == ip {
				w.Write([]byte(`nochg ` + ip))
				return
			}

			editRequest := &godo.DomainRecordEditRequest{
				Data: ip,
			}

			domainRecord, _, err := client.Domains.EditRecord(ctx, domainutil.Domain(host), d.ID, editRequest)
			if err != nil {
				logrus.WithError(err).Error("get domains")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if domainRecord.Data == ip {
				w.Write([]byte(`good ` + domainRecord.Data))
				return
			}

			w.Write([]byte(`dnserr`))
			return
		}
	}

	w.Write([]byte(`nohost`))
}

// --------------------------------------------------------------------------------------------------------------------

func domainList(ctx context.Context, client *godo.Client, domain string) ([]godo.DomainRecord, error) {
	list := []godo.DomainRecord{}

	opt := &godo.ListOptions{}
	for {
		domains, resp, err := client.Domains.Records(ctx, domain, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range domains {
			list = append(list, d)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		opt.Page = page + 1
	}

	return list, nil
}
