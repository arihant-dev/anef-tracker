package report

func (r *EvidenceReport) RenderPDFText() string {
	return r.RenderMarkdown()
}
