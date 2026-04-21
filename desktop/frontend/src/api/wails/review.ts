import { invokeApp, safeWailsCall } from '@/api/core/api-guard'
import type { ApiResult } from '@/api/core/api-result'
import type { CardWithSRSDTO, ReviewLogDTO, SRSStatisticsDTO, ReviewStatisticsDTO } from '@/api/types/dto'
import type { ReviewRequest } from '@/api/types/requests'

export function getDueCards(
  knowledgeID: string | null,
  limit: number,
): Promise<ApiResult<CardWithSRSDTO[]>> {
  return safeWailsCall(() => invokeApp<CardWithSRSDTO[]>('GetDueCards', knowledgeID, limit))
}

export function submitReview(req: ReviewRequest): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('SubmitReview', req)
  })
}

export function undoLastReview(cardID: string): Promise<ApiResult<void>> {
  return safeWailsCall(async () => {
    await invokeApp<void>('UndoLastReview', cardID)
  })
}

export function getSRSStatistics(knowledgeID: string | null): Promise<ApiResult<SRSStatisticsDTO>> {
  return safeWailsCall(() => invokeApp<SRSStatisticsDTO>('GetSRSStatistics', knowledgeID))
}

export function getReviewHistory(cardID: string, limit: number): Promise<ApiResult<ReviewLogDTO[]>> {
  return safeWailsCall(() => invokeApp<ReviewLogDTO[]>('GetReviewHistory', cardID, limit))
}

export function getReviewStats(startDate: Date, endDate: Date): Promise<ApiResult<ReviewStatisticsDTO>> {
  return safeWailsCall(() => invokeApp<ReviewStatisticsDTO>('GetReviewStats', startDate, endDate))
}
