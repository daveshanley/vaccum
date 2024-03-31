package statistics

import (
	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/model/reports"
	"github.com/pb33f/libopenapi/datamodel"
	"github.com/pb33f/libopenapi/index"
)

// CreateReportStatistics generates a ready to render breakdown of the document's statistics. A convenience function
// that reduces churn on building stats over and over for different interfaces.
func CreateReportStatistics(index *index.SpecIndex, info *datamodel.SpecInfo, results *model.RuleResultSet) *reports.ReportStatistics {

	// don't go looking for stats if we don't have the necessary data
	if index == nil || info == nil || results == nil {
		return nil
	}

	opPCount := index.GetOperationsParameterCount()
	cPCount := index.GetComponentParameterCount()

	var catStats []*reports.CategoryStatistic
	for _, cat := range model.RuleCategoriesOrdered {
		var numIssues, numWarnings, numErrors, numInfo, numHints int
		numIssues = len(results.GetResultsByRuleCategory(cat.Id))
		numWarnings = len(results.GetWarningsByRuleCategory(cat.Id))
		numErrors = len(results.GetErrorsByRuleCategory(cat.Id))
		numInfo = len(results.GetInfoByRuleCategory(cat.Id))
		numHints = len(results.GetHintByRuleCategory(cat.Id))
		numResults := len(results.Results)
		var score int
		if numResults == 0 && numIssues == 0 {
			score = 100 // perfect
		} else {
			score = numIssues / numResults * 100
		}
		catStats = append(catStats, &reports.CategoryStatistic{
			CategoryName: cat.Name,
			CategoryId:   cat.Id,
			NumIssues:    numIssues,
			Warnings:     numWarnings,
			Errors:       numErrors,
			Info:         numInfo,
			Hints:        numHints,
			Score:        score,
		})
	}

	total := 100.0
	score := total - float64(results.GetInfoCount())*0.1
	score = score - (0.4 * float64(results.GetWarnCount()))
	score = score - (15.0 * float64(results.GetErrorCount())) // errors are failures they should be judged harshly.

	if results.GetErrorCount() <= 0 && score < 0 {
		// floor at 25% if there are no errors, but a ton of warnings lowering the score
		score = 25.0
	}

	// if there are any oas-schema rile violations, bottom out the score, an invalid schema is a big deal.
	for _, result := range results.Results {
		if result.Rule.Id == "oas3-schema" {
			score = score - 90
		}
	}

	if score < 0 {
		score = 10 // the lowest score we want to present can't be 0, there has to be some hope!
	}

	stats := &reports.ReportStatistics{
		FilesizeBytes:      len(*info.SpecBytes),
		FilesizeKB:         len(*info.SpecBytes) / 1024,
		SpecType:           info.SpecType,
		SpecFormat:         info.SpecFormat,
		Version:            info.Version,
		References:         len(index.GetMappedReferences()),
		ExternalDocs:       len(index.GetAllExternalDocuments()),
		Schemas:            len(index.GetAllSchemas()),
		Parameters:         opPCount + cPCount,
		Links:              len(index.GetAllLinks()),
		Paths:              index.GetPathCount(),
		Operations:         index.GetOperationCount(),
		Tags:               index.GetTotalTagsCount(),
		Examples:           len(index.GetAllExamples()),
		Enums:              len(index.GetAllEnums()),
		Security:           len(index.GetAllSecuritySchemes()),
		OverallScore:       int(score),
		TotalErrors:        results.GetErrorCount(),
		TotalWarnings:      results.GetWarnCount(),
		TotalInfo:          results.GetInfoCount(),
		CategoryStatistics: catStats,
	}
	return stats
}
