import { getDueCards, getSRSStatistics, submitReview as submitReviewWails } from '@/api/wails/review'
import type { ApiResult } from '@/api/core/api-result'
import type { SRSStatisticsDTO } from '@/api/types/dto'
import type { ReviewItem } from '../types'
import type { ReviewGrade } from '../types'
import { mapCardWithSrsDtoToReviewItem, mapReviewGradeToRating } from './review.mapper'

export async function fetchDueReviewItems(
  knowledgeId: string | null,
  limit = 80,
): Promise<ApiResult<ReviewItem[]>> {
  const res = await getDueCards(knowledgeId, limit)
  if (!res.ok) {
    return res
  }
  return { ok: true, data: res.data.map((x) => mapCardWithSrsDtoToReviewItem(x)) }
}

export async function submitHostReview(cardId: string, grade: ReviewGrade): Promise<ApiResult<void>> {
  const rating = mapReviewGradeToRating(grade)
  return submitReviewWails({ cardId, rating })
}

export async function fetchSRSStats(knowledgeId: string | null): Promise<ApiResult<SRSStatisticsDTO>> {
  return getSRSStatistics(knowledgeId)
}
