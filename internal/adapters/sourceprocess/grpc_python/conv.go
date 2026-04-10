package grpcpython

import (
	kmemov1 "kmemo/gen/kmemo/v1"
	"kmemo/internal/contracts/sourceprocess"
)

func cleanHTMLRequest(in sourceprocess.CleanHTMLInput) *kmemov1.CleanHtmlRequest {
	return &kmemov1.CleanHtmlRequest{
		SourceId: in.SourceID,
		Html:     in.HTML,
	}
}

func submitImportJobRequest(in sourceprocess.SubmitImportJobInput) *kmemov1.SubmitImportJobRequest {
	req := &kmemov1.SubmitImportJobRequest{
		JobId:        in.JobID,
		SourceType:   in.SourceType,
		WorkspaceDir: in.WorkspaceDir,
		OutputDir:    in.OutputDir,
		TempDir:      in.TempDir,
		Options: &kmemov1.ImportOptions{
			ConversionMode:       in.Options.ConversionMode,
			FallbackModes:        append([]string(nil), in.Options.FallbackModes...),
			ExtractMainContent:   in.Options.ExtractMainContent,
			SanitizeHtml:         in.Options.SanitizeHTML,
			PreserveSemanticTags: in.Options.PreserveSemanticTags,
			DownloadRemoteAssets: in.Options.DownloadRemoteAssets,
			InlineSmallImages:    in.Options.InlineSmallImages,
			GenerateToc:          in.Options.GenerateTOC,
			AnalyzeStructure:     in.Options.AnalyzeStructure,
			KeepSourceCopy:       in.Options.KeepSourceCopy,
			EnabledCleaners:      append([]string(nil), in.Options.EnabledCleaners...),
		},
		Metadata: copyStringMap(in.Metadata),
	}
	if in.SourcePath != nil {
		req.SourcePath = in.SourcePath
	}
	if in.SourceURI != nil {
		req.SourceUri = in.SourceURI
	}
	if in.SourceURL != nil {
		req.SourceUrl = in.SourceURL
	}
	if in.RawHTML != nil {
		req.RawHtml = in.RawHTML
	}
	if in.Options.ConverterParamsJSON != nil {
		req.Options.ConverterParamsJson = in.Options.ConverterParamsJSON
	}
	if in.IdempotencyKey != nil {
		req.IdempotencyKey = in.IdempotencyKey
	}
	return req
}

func listJobEventsRequest(in sourceprocess.ListJobEventsInput) *kmemov1.ListJobEventsRequest {
	req := &kmemov1.ListJobEventsRequest{JobId: in.JobID}
	if in.AfterSequence != nil {
		req.AfterSequence = in.AfterSequence
	}
	return req
}

func cleanHTMLOutput(resp *kmemov1.CleanHtmlResponse) *sourceprocess.CleanHTMLOutput {
	if resp == nil {
		return nil
	}
	return &sourceprocess.CleanHTMLOutput{
		OK:          resp.GetOk(),
		Message:     resp.GetMessage(),
		CleanedHTML: resp.GetCleanedHtml(),
	}
}

func submitImportJobOutput(resp *kmemov1.SubmitImportJobResponse) *sourceprocess.SubmitImportJobOutput {
	if resp == nil {
		return nil
	}
	return &sourceprocess.SubmitImportJobOutput{
		JobID:    resp.GetJobId(),
		Status:   resp.GetStatus(),
		Accepted: resp.GetAccepted(),
	}
}

func jobFromProto(protoJob *kmemov1.SourceProcessJob) *sourceprocess.Job {
	if protoJob == nil {
		return nil
	}
	job := &sourceprocess.Job{
		JobID:    protoJob.GetJobId(),
		Status:   protoJob.GetStatus(),
		Stage:    protoJob.GetStage(),
		Progress: protoJob.GetProgress(),
		Result:   importResultFromProto(protoJob.GetResult()),
	}
	if protoJob.ResultPath != nil {
		v := protoJob.GetResultPath()
		job.ResultPath = &v
	}
	if protoJob.ErrorCode != nil {
		v := protoJob.GetErrorCode()
		job.ErrorCode = &v
	}
	if protoJob.ErrorMessage != nil {
		v := protoJob.GetErrorMessage()
		job.ErrorMessage = &v
	}
	return job
}

func importResultFromProto(result *kmemov1.ImportResult) *sourceprocess.ImportResult {
	if result == nil {
		return nil
	}
	return &sourceprocess.ImportResult{
		EntryHTMLPath:           result.GetEntryHtmlPath(),
		CleanedHTMLPath:         result.GetCleanedHtmlPath(),
		RawTextPath:             result.GetRawTextPath(),
		Assets:                  append([]string(nil), result.GetAssets()...),
		ExtractedMetadata:       copyStringMap(result.GetExtractedMetadata()),
		ManifestPath:            result.GetManifestPath(),
		ContentHash:             result.GetContentHash(),
		EffectiveConversionMode: result.GetEffectiveConversionMode(),
		ConverterName:           result.GetConverterName(),
		ConverterVersion:        result.GetConverterVersion(),
		CleanerVersion:          result.GetCleanerVersion(),
	}
}

func jobEventsFromProto(events []*kmemov1.SourceProcessJobEvent) []sourceprocess.JobEvent {
	out := make([]sourceprocess.JobEvent, 0, len(events))
	for _, event := range events {
		if event == nil {
			continue
		}
		out = append(out, sourceprocess.JobEvent{
			JobID:         event.GetJobId(),
			Sequence:      event.GetSequence(),
			Stage:         event.GetStage(),
			Message:       event.GetMessage(),
			CreatedAtUnix: event.GetCreatedAtUnix(),
		})
	}
	return out
}

func cancelJobOutput(resp *kmemov1.CancelJobResponse) *sourceprocess.CancelJobOutput {
	if resp == nil {
		return nil
	}
	return &sourceprocess.CancelJobOutput{
		JobID:  resp.GetJobId(),
		Status: resp.GetStatus(),
	}
}

func capabilitiesFromProto(resp *kmemov1.GetCapabilitiesResponse) *sourceprocess.Capabilities {
	if resp == nil {
		return nil
	}
	caps := &sourceprocess.Capabilities{
		SourceTypes:     append([]string(nil), resp.GetSourceTypes()...),
		ConversionModes: append([]string(nil), resp.GetConversionModes()...),
		Converters:      make([]sourceprocess.ConverterCapability, 0, len(resp.GetConverters())),
		Cleaners:        make([]sourceprocess.CleanerCapability, 0, len(resp.GetCleaners())),
	}
	for _, item := range resp.GetConverters() {
		if item == nil {
			continue
		}
		caps.Converters = append(caps.Converters, sourceprocess.ConverterCapability{
			SourceType:       item.GetSourceType(),
			ConversionMode:   item.GetConversionMode(),
			ConverterName:    item.GetConverterName(),
			ConverterVersion: item.GetConverterVersion(),
			Description:      item.GetDescription(),
		})
	}
	for _, item := range resp.GetCleaners() {
		if item == nil {
			continue
		}
		caps.Cleaners = append(caps.Cleaners, sourceprocess.CleanerCapability{
			CleanerName:    item.GetCleanerName(),
			CleanerVersion: item.GetCleanerVersion(),
			Description:    item.GetDescription(),
		})
	}
	return caps
}

func copyStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
