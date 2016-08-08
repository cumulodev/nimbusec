package nimbusec

import "io/ioutil"

type Agent struct {
	OS      string `json:"os"`
	Arch    string `json:"arch"`
	Version int    `json:"version"`
	Md5     string `json:"md5"`
	Sha1    string `json:"sha1"`
	Format  string `json:"format"`
	URL     string `json:"url"`
}

func (a *API) DownloadAgent(agent Agent) ([]byte, error) {
	url := a.geturl("/v2/agent/download/nimbusagent-%s-%s-v%d.%s", agent.OS, agent.Arch, agent.Version, agent.Format)
	res, err := a.client.Get(url, params{}, a.token)
	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)

}

func (a *API) FindAgents(filter string) ([]Agent, error) {
	params := params{}
	if filter != EmptyFilter {
		params["q"] = filter
	}

	dst := make([]Agent, 0)
	url := a.geturl("/v2/agent/download")
	err := a.get(url, params, &dst)
	return dst, err
}
