/**
 * Clamp a point so a box of size (w, h) stays inside the viewport with margin.
 */
export function clampMenuPosition(
  x: number,
  y: number,
  width: number,
  height: number,
  margin = 8,
): { x: number; y: number } {
  const vw = typeof window !== 'undefined' ? window.innerWidth : 800
  const vh = typeof window !== 'undefined' ? window.innerHeight : 600
  const maxX = Math.max(margin, vw - width - margin)
  const maxY = Math.max(margin, vh - height - margin)
  return {
    x: Math.min(maxX, Math.max(margin, x)),
    y: Math.min(maxY, Math.max(margin, y)),
  }
}
